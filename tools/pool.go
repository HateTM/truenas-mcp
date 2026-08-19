package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerPoolTools registers all pool-related MCP tools onto the server.
func registerPoolTools(s *mcp.Server, client truenasClient) {
	type listPoolsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of pools to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of pools to skip; 0 means start from the beginning"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pools",
		Description: "List all ZFS storage pools and their status, size, and health.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listPoolsInput) (*mcp.CallToolResult, any, error) {
		pools, err := client.ListPools(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("list_pools: %w", err))
		}
		return jsonResult(pools)
	})

	type getPoolInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_pool",
		Description: "Get detailed information about a specific ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getPoolInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("get_pool: id must be a positive integer"))
		}
		pool, err := client.GetPool(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("get_pool: %w", err))
		}
		return jsonResult(pool)
	})

	type createPoolInput struct {
		Name   string   `json:"name"   jsonschema:"Pool name"`
		Disks  []string `json:"disks"  jsonschema:"List of disk devices"`
		Layout string   `json:"layout" jsonschema:"ZFS layout, e.g. STRIPE, RAIDZ1"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_pool",
		Description: "Create a new ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createPoolInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || len(p.Disks) == 0 || p.Layout == "" {
			return errorResult(errors.New("create_pool: name, disks, and layout are required"))
		}
		pool, err := client.CreatePool(ctx, &truenas.CreatePoolParams{
			Name:   p.Name,
			Disks:  p.Disks,
			Layout: p.Layout,
		})
		if err != nil {
			return errorResult(fmt.Errorf("create_pool: %w", err))
		}
		return jsonResult(pool)
	})

	type updatePoolInput struct {
		ID   int    `json:"id"   jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"New pool name"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_pool",
		Description: "Update an existing ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updatePoolInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("update_pool: id and name are required"))
		}
		pool, err := client.UpdatePool(ctx, p.ID, &truenas.UpdatePoolParams{Name: p.Name})
		if err != nil {
			return errorResult(fmt.Errorf("update_pool: %w", err))
		}
		return jsonResult(pool)
	})
}
