package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerAlertManagementDestructiveTools adds opt-in destructive alert management tools to the server.
// These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerAlertManagementDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type alertServiceIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric alert service ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertservice_delete",
		Description: "Permanently delete an alert notification service. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p alertServiceIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("alertservice_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("alertservice_delete: id must be a positive integer"))
		}
		if err := client.DeleteAlertService(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("alertservice_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
