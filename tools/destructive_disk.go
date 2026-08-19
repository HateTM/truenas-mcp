package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerDiskDestructiveTools adds the opt-in destructive disk.wipe tool to the server.
// Only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerDiskDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type wipeDiskInput struct {
		Identifier string `json:"identifier" jsonschema:"Disk identifier, e.g. {serial}sN12345"`
		Mode       string `json:"mode"       jsonschema:"QUICK, FULL, or FULL_WITH_ZEROS"`
		Confirmed  bool   `json:"confirmed"  jsonschema:"Must be set to true to confirm the wipe"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_wipe",
		Description: "Permanently erase all data on a disk. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p wipeDiskInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("disk_wipe: confirmed must be true to proceed with wiping"))
		}
		if p.Identifier == "" {
			return errorResult(errors.New("disk_wipe: identifier must not be empty"))
		}
		if p.Mode == "" {
			return errorResult(errors.New("disk_wipe: mode must not be empty"))
		}
		jobID, err := client.WipeDisk(ctx, p.Identifier, &truenas.WipeDiskParams{Mode: p.Mode})
		if err != nil {
			return errorResult(fmt.Errorf("disk_wipe: %w", err))
		}
		return jsonResult(map[string]any{"job_id": jobID, "identifier": p.Identifier})
	})
}
