package command_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/haines/multidockerfile/internal/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/join/success/expected.dockerfile
var expected string

func TestJoinToStdout(t *testing.T) {
	result := command.Capture(
		"join",
		"testdata/join/success/a.dockerfile",
		"testdata/join/success/b.dockerfile",
	)

	assert.Equal(t, 0, result.ExitCode, "Unexpected exit code")
	assert.Equal(t, expected, result.Stdout, "Unexpected output to stdout")
	assert.Empty(t, result.Stderr, "Unexpected output to stderr")
}

func TestJoinToFile(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "multidockerfile-test-join-to-file")
	require.NoError(t, err, "Failed to create output directory")

	outputFile := filepath.Join(outputDir, "Dockerfile")

	result := command.Capture(
		"join",
		"--output", outputFile,
		"testdata/join/success/a.dockerfile",
		"testdata/join/success/b.dockerfile",
	)

	assert.Equal(t, 0, result.ExitCode, "Unexpected exit code")
	assert.Empty(t, result.Stdout, "Unexpected output to stdout")
	assert.Empty(t, result.Stderr, "Unexpected output to stderr")

	output, err := os.ReadFile(outputFile)
	require.NoError(t, err, "Failed to read output file")
	assert.Equal(t, expected, string(output), "Unexpected output to file")
}

func TestJoinToFileFailsWhenOutputDockerfileCannotBeWritten(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "multidockerfile-test-join-to-file")
	require.NoError(t, err, "Failed to create output directory")

	result := command.Capture(
		"join",
		"--output", outputDir,
		"testdata/join/success/a.dockerfile",
		"testdata/join/success/b.dockerfile",
	)

	assert.Equal(t, 1, result.ExitCode, "Unexpected exit code")
	assert.Empty(t, result.Stdout, "Unexpected output to stdout")
	assert.Contains(t, result.Stderr, "failed to open output Dockerfile", "Unexpected output to stderr")
}

func TestJoinFailsWhenInputDockerfileCannotBeRead(t *testing.T) {
	result := command.Capture(
		"join",
		"testdata/join/unreadable/Dockerfile",
	)

	assert.Equal(t, 1, result.ExitCode, "Unexpected exit code")
	assert.Empty(t, result.Stdout, "Unexpected output to stdout")
	assert.Contains(t, result.Stderr, "failed to parse Dockerfile", "Unexpected output to stderr")
}

func TestJoinFailsWhenOutputDockerfileCannotBeGenerated(t *testing.T) {
	result := command.Capture(
		"join",
		"testdata/join/unwritable/Dockerfile",
	)

	assert.Equal(t, 1, result.ExitCode, "Unexpected exit code")
	assert.Empty(t, result.Stdout, "Unexpected output to stdout")
	assert.Contains(t, result.Stderr, "cycle detected", "Unexpected output to stderr")
}
