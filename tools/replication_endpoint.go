package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerReplicationEndpointTools registers replication.endpoint.* MCP tools onto the server.
func registerReplicationEndpointTools(s *mcp.Server, client truenasClient) {
	type listReplicationEndpointsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of endpoints to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of endpoints to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_endpoint_list",
		Description: "List configured remote replication endpoints.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listReplicationEndpointsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListReplicationEndpoints(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("replication_endpoint_list: %w", err))
		}
		return jsonResult(result)
	})

	type getReplicationEndpointInput struct {
		ID int `json:"id" jsonschema:"Numeric replication endpoint ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_endpoint_get",
		Description: "Get a single remote replication endpoint by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getReplicationEndpointInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("replication_endpoint_get: id must be a positive integer"))
		}
		result, err := client.GetReplicationEndpoint(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("replication_endpoint_get: %w", err))
		}
		return jsonResult(result)
	})

	type replicationEndpointInput struct {
		Name  string `json:"name"           jsonschema:"Endpoint name"`
		URI   string `json:"uri"            jsonschema:"Remote system URI, e.g. ws://remote.example.com/websocket"`
		Token string `json:"token,omitempty" jsonschema:"API token or key for authenticating to the remote system"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_endpoint_create",
		Description: "Create a new remote replication endpoint, for use as a replication task's target.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p replicationEndpointInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || p.URI == "" {
			return errorResult(errors.New("replication_endpoint_create: name and uri are required"))
		}
		result, err := client.CreateReplicationEndpoint(ctx, &truenas.CreateReplicationEndpointParams{
			Name: p.Name, URI: p.URI, Token: p.Token,
		})
		if err != nil {
			return errorResult(fmt.Errorf("replication_endpoint_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateReplicationEndpointInput struct {
		ID    int    `json:"id"             jsonschema:"Numeric replication endpoint ID"`
		Name  string `json:"name"           jsonschema:"Endpoint name"`
		URI   string `json:"uri"            jsonschema:"Remote system URI"`
		Token string `json:"token,omitempty" jsonschema:"API token or key for authenticating to the remote system"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_endpoint_update",
		Description: "Update an existing remote replication endpoint.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateReplicationEndpointInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("replication_endpoint_update: id must be a positive integer"))
		}
		if p.Name == "" || p.URI == "" {
			return errorResult(errors.New("replication_endpoint_update: name and uri are required"))
		}
		result, err := client.UpdateReplicationEndpoint(ctx, p.ID, &truenas.CreateReplicationEndpointParams{
			Name: p.Name, URI: p.URI, Token: p.Token,
		})
		if err != nil {
			return errorResult(fmt.Errorf("replication_endpoint_update: %w", err))
		}
		return jsonResult(result)
	})

	// replication_endpoint_delete is destructive and lives in
	// destructive_replication.go, gated behind Config.AllowDestructive.
}
