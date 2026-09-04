package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerUserDestructiveTools adds opt-in destructive/sensitive user management tools to
// the server. These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerUserDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type deleteUserInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric user ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_delete",
		Description: "Permanently delete a local user account. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p deleteUserInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("user_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("user_delete: id must be a positive integer"))
		}
		if err := client.DeleteUser(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("user_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	type updateUserPasswordInput struct {
		ID        int    `json:"id"        jsonschema:"Numeric user ID"`
		Password  string `json:"password"  jsonschema:"New login password"`
		Confirmed bool   `json:"confirmed" jsonschema:"Must be set to true to confirm the password change"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_update_password",
		Description: "Set a new password for a local user account, invalidating the old one. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateUserPasswordInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("user_update_password: confirmed must be true to proceed"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("user_update_password: id must be a positive integer"))
		}
		if p.Password == "" {
			return errorResult(errors.New("user_update_password: password must not be empty"))
		}
		if err := client.UpdateUserPassword(ctx, p.ID, p.Password); err != nil {
			return errorResult(fmt.Errorf("user_update_password: %w", err))
		}
		return jsonResult(map[string]any{"updated": true, "id": p.ID})
	})
}
