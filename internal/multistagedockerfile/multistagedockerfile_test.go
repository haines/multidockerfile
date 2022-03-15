package multistagedockerfile_test

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/haines/multidockerfile/internal/multistagedockerfile"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/success/expected.dockerfile
var expected string

func TestReadWrite(t *testing.T) {
	dockerfile := multistagedockerfile.New()

	err := dockerfile.Read("testdata/success/a.dockerfile")
	require.NoError(t, err)

	err = dockerfile.Read("testdata/success/b.dockerfile")
	require.NoError(t, err)

	err = dockerfile.Read("testdata/success/c.dockerfile")
	require.NoError(t, err)

	err = dockerfile.Read("testdata/success/d.dockerfile")
	require.NoError(t, err)

	err = dockerfile.Read("testdata/success/e.dockerfile")
	require.NoError(t, err)

	var result strings.Builder
	written, err := dockerfile.Write(&result)
	require.NoError(t, err)
	assert.Equal(t, len(expected), written, "Unexpected number of bytes written")
	assert.Equal(t, expected, result.String(), "Unexpected content written")
}

func TestReadFailsWhenFileCannotBeRead(t *testing.T) {
	err := multistagedockerfile.New().Read("not/a.dockerfile")
	require.Error(t, err, "Unexpected nil error reading non-existent Dockerfile")
	assert.Contains(t, err.Error(), "failed to open Dockerfile not/a.dockerfile", "Unexpected error")
	var pathError *os.PathError
	assert.ErrorAs(t, err, &pathError)
}

func TestReadFailsWhenFilesContainIncompatibleDirectives(t *testing.T) {
	dockerfile := multistagedockerfile.New()

	err := dockerfile.Read("testdata/incompatible_directives/a.dockerfile")
	require.NoError(t, err)

	err = dockerfile.Read("testdata/incompatible_directives/b.dockerfile")
	assert.EqualError(
		t, err,
		`incompatible syntax directives:
  testdata/incompatible_directives/a.dockerfile:1 : "docker/dockerfile:1.1"
  testdata/incompatible_directives/b.dockerfile:1 : "docker/dockerfile:1.2"`,
	)
}

func TestReadFailsWhenDockerfileCannotBeParsed(t *testing.T) {
	dockerfile := multistagedockerfile.New()

	err := dockerfile.Read("testdata/empty/Dockerfile")
	require.Error(t, err, "Unexpected nil error reading empty Dockerfile")
	assert.Contains(t, err.Error(), "failed to parse Dockerfile testdata/empty/Dockerfile", "Unexpected error")
	var locationError *parser.ErrorLocation
	assert.ErrorAs(t, err, &locationError)
}

func TestReadFailsWhenInstructionCannotBeParsed(t *testing.T) {
	dockerfile := multistagedockerfile.New()

	err := dockerfile.Read("testdata/unknown_instruction/Dockerfile")
	require.Error(t, err, "Unexpected nil error reading Dockerfile with unknown instruction")
	assert.Contains(t, err.Error(), "failed to parse Dockerfile testdata/unknown_instruction/Dockerfile", "Unexpected error")
	var unknownInstruction *instructions.UnknownInstructionError
	assert.ErrorAs(t, err, &unknownInstruction)
}

func TestReadFailsWhenStagesAreDuplicated(t *testing.T) {
	dockerfile := multistagedockerfile.New()

	err := dockerfile.Read("testdata/duplicate_stages/a.dockerfile")
	require.NoError(t, err)

	err = dockerfile.Read("testdata/duplicate_stages/b.dockerfile")
	assert.EqualError(
		t, err,
		`failed to parse Dockerfile testdata/duplicate_stages/b.dockerfile: found multiple definitions for stage "a":
  testdata/duplicate_stages/a.dockerfile:1
  testdata/duplicate_stages/b.dockerfile:1`,
	)
}

func TestWriteFailsWhenStagesAreCyclic(t *testing.T) {
	dockerfile := multistagedockerfile.New()

	err := dockerfile.Read("testdata/cyclic_stages/Dockerfile")
	require.NoError(t, err)

	var result strings.Builder
	_, err = dockerfile.Write(&result)
	assert.EqualError(
		t, err,
		`cycle detected between stages:

  stage "a" defined at testdata/cyclic_stages/Dockerfile:1 is depended on by
    - stage "b" defined at testdata/cyclic_stages/Dockerfile:11

  stage "b" defined at testdata/cyclic_stages/Dockerfile:11 is depended on by
    - stage "c" defined at testdata/cyclic_stages/Dockerfile:21

  stage "c" defined at testdata/cyclic_stages/Dockerfile:21 is depended on by
    - stage "a" defined at testdata/cyclic_stages/Dockerfile:1
    - stage "b" defined at testdata/cyclic_stages/Dockerfile:11`,
	)
}
