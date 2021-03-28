package command

import (
	"encoding/json"
	"fmt"

	"github.com/haines/multidockerfile/internal/version"
)

type versionCommand struct{}

func (versionCommand) Run(ctx *context) error {
	output, err := json.MarshalIndent(version.Get(), "", "  ")
	if err != nil {
		return fmt.Errorf("failed to JSON-encode version info: %w", err)
	}

	fmt.Fprintf(ctx.Stdout, "%s\n", output)

	return nil
}
