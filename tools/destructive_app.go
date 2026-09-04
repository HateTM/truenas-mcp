package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerAppManagementDestructiveTools adds opt-in destructive app management tools to the server.
// These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerAppManagementDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type appRegistryIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric registry ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_registry_delete",
		Description: "Permanently delete a private container registry credential. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p appRegistryIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("app_registry_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("app_registry_delete: id must be a positive integer"))
		}
		if err := client.DeleteAppRegistry(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("app_registry_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
