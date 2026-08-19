package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerReplicationDestructiveTools adds opt-in destructive replication tools to the server.
// These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerReplicationDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type replicationIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_delete",
		Description: "Permanently delete a replication task. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p replicationIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("replication_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("replication_delete: id must be a positive integer"))
		}
		if err := client.DeleteReplicationTask(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("replication_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_endpoint_delete",
		Description: "Permanently delete a remote replication endpoint. Fails if it is still referenced by a replication task. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p replicationIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("replication_endpoint_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("replication_endpoint_delete: id must be a positive integer"))
		}
		if err := client.DeleteReplicationEndpoint(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("replication_endpoint_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
