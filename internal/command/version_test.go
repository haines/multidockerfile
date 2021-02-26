package command_test

import (
	"encoding/json"
	"testing"

	"github.com/haines/multidockerfile/internal/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	result := command.Capture("version")

	assert.Equal(t, 0, result.ExitCode, "Unexpected exit code")
	assert.Empty(t, result.Stderr, "Unexpected output to stderr")

	var versionInfo map[string]interface{}
	err := json.Unmarshal([]byte(result.Stdout), &versionInfo)
	require.NoError(t, err, "Unable to parse stdout as JSON: %q", result.Stdout)

	assert.Equal(t, "unknown", versionInfo["Version"], "Unexpected version")
	assert.Equal(t, "unknown", versionInfo["GitCommit"], "Unexpected Git commit")
	assert.Equal(t, "unknown", versionInfo["Built"], "Unexpected built")
	assert.Contains(t, versionInfo, "GoVersion", "Missing Go version")
	assert.Contains(t, versionInfo, "OS", "Missing OS")
	assert.Contains(t, versionInfo, "Arch", "Missing arch")
}
