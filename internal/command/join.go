package command

import (
	"fmt"
	"io"
	"os"

	"github.com/haines/multidockerfile/internal/multistagedockerfile"
)

type joinCommand struct {
	Output string   `short:"o" type:"path" default:"-" help:"Where to write the multi-stage Dockerfile (- for stdout)."`
	Inputs []string `arg:"" type:"existingfile" help:"Paths to the Dockerfiles to be joined."`
}

func (j joinCommand) Run(ctx *context) error {
	var output io.Writer
	if j.Output == "-" {
		output = ctx.Stdout
	} else {
		file, err := os.Create(j.Output)
		if err != nil {
			return fmt.Errorf("failed to open output Dockerfile %q for writing: %w", j.Output, err)
		}
		defer file.Close()

		output = file
	}

	dockerfile := multistagedockerfile.New()

	for _, input := range j.Inputs {
		err := dockerfile.Read(input)
		if err != nil {
			return err
		}
	}

	_, err := dockerfile.Write(output)
	if err != nil {
		return err
	}

	return nil
}
