package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerSharingDestructiveTools adds opt-in destructive sharing (NFS/SMB/WebDAV) tools to the server.
// These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerSharingDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type sharingIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric share ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_nfs_delete",
		Description: "Permanently delete an NFS export. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p sharingIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("sharing_nfs_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_nfs_delete: id must be a positive integer"))
		}
		if err := client.DeleteNFSShare(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("sharing_nfs_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_smb_delete",
		Description: "Permanently delete an SMB share. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p sharingIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("sharing_smb_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_smb_delete: id must be a positive integer"))
		}
		if err := client.DeleteSMBShare(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("sharing_smb_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_webdav_delete",
		Description: "Permanently delete a WebDAV share. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p sharingIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("sharing_webdav_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_webdav_delete: id must be a positive integer"))
		}
		if err := client.DeleteWebDAVShare(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("sharing_webdav_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
