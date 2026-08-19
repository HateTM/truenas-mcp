package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerISCSIDestructiveTools adds opt-in destructive iSCSI tools to the server.
// These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerISCSIDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type iscsiIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_auth_delete",
		Description: "Permanently delete an iSCSI CHAP authentication credential. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p iscsiIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("iscsi_auth_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_auth_delete: id must be a positive integer"))
		}
		if err := client.DeleteISCSIAuth(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("iscsi_auth_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_extent_delete",
		Description: "Permanently delete an iSCSI extent. Fails if it is still mapped to a target. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p iscsiIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("iscsi_extent_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_extent_delete: id must be a positive integer"))
		}
		if err := client.DeleteISCSIExtent(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("iscsi_extent_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_initiator_delete",
		Description: "Permanently delete an iSCSI initiator group. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p iscsiIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("iscsi_initiator_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_initiator_delete: id must be a positive integer"))
		}
		if err := client.DeleteISCSIInitiator(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("iscsi_initiator_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_portal_delete",
		Description: "Permanently delete an iSCSI portal. Fails if it is still referenced by a target. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p iscsiIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("iscsi_portal_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_portal_delete: id must be a positive integer"))
		}
		if err := client.DeleteISCSIPortal(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("iscsi_portal_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_target_delete",
		Description: "Permanently delete an iSCSI target. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p iscsiIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("iscsi_target_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_target_delete: id must be a positive integer"))
		}
		if err := client.DeleteISCSITarget(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("iscsi_target_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_targetextent_delete",
		Description: "Permanently remove an iSCSI target/extent LUN mapping. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p iscsiIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("iscsi_targetextent_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_targetextent_delete: id must be a positive integer"))
		}
		if err := client.DeleteISCSITargetExtent(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("iscsi_targetextent_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
