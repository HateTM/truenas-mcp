package tools

import (
	"context"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNFSConfigTools registers nfs.* global service configuration MCP tools.
// These are distinct from sharing.nfs.* (individual exports), registered in sharing.go.
func registerNFSConfigTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nfs_bindip_choices",
		Description: "List the IP addresses available for the NFS service to bind to.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.NFSBindIPChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nfs_bindip_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nfs_client_count",
		Description: "Get the number of clients currently connected to the NFS service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		count, err := client.NFSClientCount(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nfs_client_count: %w", err))
		}
		return jsonResult(map[string]int{"client_count": count})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nfs_config",
		Description: "Get the global NFS service configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.NFSConfigGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nfs_config: %w", err))
		}
		return jsonResult(cfg)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nfs_get_nfs3_clients",
		Description: "List clients currently connected via NFSv3.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		clients, err := client.GetNFS3Clients(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nfs_get_nfs3_clients: %w", err))
		}
		return jsonResult(clients)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nfs_get_nfs4_clients",
		Description: "List clients currently connected via NFSv4.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		clients, err := client.GetNFS4Clients(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nfs_get_nfs4_clients: %w", err))
		}
		return jsonResult(clients)
	})

	type updateNFSConfigInput struct {
		Servers      int      `json:"servers,omitempty"       jsonschema:"Number of NFS server threads"`
		BindIP       []string `json:"bindip,omitempty"        jsonschema:"IP addresses to bind the NFS service to; empty means all interfaces"`
		AllowNonroot bool     `json:"allow_nonroot,omitempty" jsonschema:"Allow non-root mount requests"`
		V4           bool     `json:"v4,omitempty"            jsonschema:"Enable NFSv4"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nfs_update",
		Description: "Update the global NFS service configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNFSConfigInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.UpdateNFSConfig(ctx, &truenas.UpdateNFSConfigParams{
			Servers: p.Servers, BindIP: p.BindIP, AllowNonroot: p.AllowNonroot, V4: p.V4,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nfs_update: %w", err))
		}
		return jsonResult(cfg)
	})
}
