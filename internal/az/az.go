// Package az provides the private bridge to the Azure CLI.
package az

import (
	"context"
	"os/exec"
)

// Run executes the az CLI with args and returns its combined stdout and stderr.
func Run(ctx context.Context, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, "az", args...).CombinedOutput()
}
