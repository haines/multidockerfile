package multistagedockerfile

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type MultiStageDockerfile struct {
	directivesList []*dockerfileDirective
	directivesMap  map[string]*dockerfileDirective
	globals        []string
	stagesList     []*dockerfileStage
	stagesMap      map[string]*dockerfileStage
	currentStage   *dockerfileStage
}

type dockerfileDirective struct {
	Location string
	Name     string
	Value    string
}

type dockerfileStage struct {
	Index        int
	Location     string
	Name         string
	Instructions []string
	Dependencies []string
}

func New() *MultiStageDockerfile {
	return &MultiStageDockerfile{
		directivesMap: make(map[string]*dockerfileDirective),
		stagesMap:     make(map[string]*dockerfileStage),
	}
}

func (m *MultiStageDockerfile) Read(path string) error {
	dockerfile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open Dockerfile %s: %w", path, err)
	}
	defer dockerfile.Close()

	directives := dockerfile2llb.ParseDirectives(dockerfile)
	err = m.addDirectives(directives, path)
	if err != nil {
		return err
	}

	_, err = dockerfile.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to rewind Dockerfile %s: %w", path, err)
	}

	result, err := parser.Parse(dockerfile)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile %s: %w", path, err)
	}

	err = m.addInstructions(result.AST.Children, path)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile %s: %w", path, err)
	}

	return nil
}

func (m *MultiStageDockerfile) addDirectives(directives map[string]dockerfile2llb.Directive, path string) error {
	for name, directive := range directives {
		newDirective := &dockerfileDirective{
			Location: location(path, directive.Location),
			Name:     name,
			Value:    directive.Value,
		}

		existingDirective, ok := m.directivesMap[name]
		if !ok {
			m.directivesList = append(m.directivesList, newDirective)
			m.directivesMap[name] = newDirective
		} else if newDirective.Value != existingDirective.Value {
			return fmt.Errorf(
				"incompatible %s directives:\n  %s : %q\n  %s : %q",
				name,
				existingDirective.Location,
				existingDirective.Value,
				newDirective.Location,
				newDirective.Value,
			)
		}
	}

	return nil
}

func (m *MultiStageDockerfile) addInstructions(nodes []*parser.Node, path string) error {
	m.currentStage = nil

	for _, node := range nodes {
		stageOrCommand, err := instructions.ParseInstruction(node)
		if err != nil {
			return err
		}

		stage, ok := stageOrCommand.(*instructions.Stage)
		if ok {
			err = m.addStage(node, stage, path)
			if err != nil {
				return err
			}
		} else {
			m.addCommand(node, stageOrCommand)
		}
	}

	return nil
}

func (m *MultiStageDockerfile) addStage(node *parser.Node, instruction *instructions.Stage, path string) error {
	m.currentStage = &dockerfileStage{
		Index:    len(m.stagesList),
		Location: location(path, instruction.Location),
		Name:     instruction.Name,
		Instructions: []string{
			node.Original,
		},
		Dependencies: []string{
			instruction.BaseName,
		},
	}

	if m.currentStage.Name == "" {
		m.currentStage.Name = strconv.Itoa(m.currentStage.Index)
	}

	existing, exists := m.stagesMap[m.currentStage.Name]
	if exists {
		return fmt.Errorf("found multiple definitions for stage %q:\n  %s\n  %s", m.currentStage.Name, existing.Location, m.currentStage.Location)
	}

	m.stagesList = append(m.stagesList, m.currentStage)
	m.stagesMap[m.currentStage.Name] = m.currentStage

	return nil
}

func location(path string, location []parser.Range) string {
	return fmt.Sprintf("%s:%d", path, location[0].Start.Line)
}

func (m *MultiStageDockerfile) addCommand(node *parser.Node, command interface{}) {
	if m.currentStage == nil {
		m.globals = append(m.globals, node.Original)
		return
	}

	m.currentStage.Instructions = append(m.currentStage.Instructions, node.Original)

	copy, ok := command.(*instructions.CopyCommand)
	if ok && copy.From != "" {
		m.currentStage.Dependencies = append(m.currentStage.Dependencies, copy.From)
	}
}

func (m *MultiStageDockerfile) Write(w io.Writer) (int, error) {
	written := 0
	for _, directive := range m.directivesList {
		n, err := fmt.Fprintf(w, "# %s = %s\n", directive.Name, directive.Value)
		written += n
		if err != nil {
			return written, err
		}
	}

	for _, global := range m.globals {
		n, err := fmt.Fprintln(w, global)
		written += n
		if err != nil {
			return written, err
		}
	}

	stages, err := m.sortedStages()
	if err != nil {
		return written, err
	}

	for _, stage := range stages {
		for _, instruction := range stage.Instructions {
			n, err := fmt.Fprintln(w, instruction)
			written += n
			if err != nil {
				return written, err
			}
		}
	}

	return written, nil
}

func (m *MultiStageDockerfile) sortedStages() ([]*dockerfileStage, error) {
	var result []*dockerfileStage
	edges := make(map[string][]*dockerfileStage)
	degree := make(map[string]int)

	for _, stage := range m.stagesList {
		for _, dependency := range stage.Dependencies {
			_, ok := m.stagesMap[dependency]
			if ok {
				edges[dependency] = append(edges[dependency], stage)
				degree[stage.Name]++
			}
		}

		if degree[stage.Name] == 0 {
			result = append(result, stage)
		}
	}

	for i := 0; i < len(result); i++ {
		vertex := result[i]

		for _, neighbour := range edges[vertex.Name] {
			if degree[neighbour.Name] == 1 {
				result = append(result, neighbour)
			} else {
				degree[neighbour.Name]--
			}
		}

		delete(edges, vertex.Name)
	}

	if len(edges) > 0 {
		var message strings.Builder

		fmt.Fprintf(&message, "cycle detected between stages:\n\n")

		vertices := make([]*dockerfileStage, 0, len(edges))
		for vertexName := range edges {
			vertices = append(vertices, m.stagesMap[vertexName])
		}

		sortStagesByIndex(vertices)

		for _, vertex := range vertices {
			fmt.Fprintf(&message, "  stage %q defined at %s is depended on by\n", vertex.Name, vertex.Location)

			neighbours := edges[vertex.Name]
			sortStagesByIndex(neighbours)

			for _, neighbour := range neighbours {
				fmt.Fprintf(&message, "    - stage %q defined at %s\n", neighbour.Name, neighbour.Location)
			}

			fmt.Fprintf(&message, "\n")
		}

		return nil, fmt.Errorf(strings.TrimSuffix(message.String(), "\n\n"))
	}

	return result, nil
}

func sortStagesByIndex(stages []*dockerfileStage) {
	sort.Slice(stages, func(i int, j int) bool {
		return stages[i].Index < stages[j].Index
	})
}
