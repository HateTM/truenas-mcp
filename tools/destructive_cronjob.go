package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCronJobDestructiveTools adds the opt-in destructive cronjob.delete tool to the
// server. Only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerCronJobDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type cronJobIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric cron job ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "cronjob_delete",
		Description: "Permanently delete a scheduled cron job. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p cronJobIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("cronjob_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("cronjob_delete: id must be a positive integer"))
		}
		if err := client.DeleteCronJob(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("cronjob_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
