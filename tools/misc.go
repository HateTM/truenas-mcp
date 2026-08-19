// Package tools — misc.go registers the device.get_info and dns.query MCP tools, the two
// smallest remaining tracked categories, sharing one file since each is a single method.
package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerMiscTools registers device.get_info and dns.query MCP tools onto the server.
func registerMiscTools(s *mcp.Server, client truenasClient) {
	type deviceGetInfoInput struct {
		Type string `json:"type" jsonschema:"Device category, e.g. SERIAL, DISK, or GPU"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "device_get_info",
		Description: "Get information about system devices of a given category (e.g. serial ports, disks, GPUs).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p deviceGetInfoInput) (*mcp.CallToolResult, any, error) {
		if p.Type == "" {
			return errorResult(errors.New("device_get_info: type must not be empty"))
		}
		info, err := client.DeviceGetInfo(ctx, p.Type)
		if err != nil {
			return errorResult(fmt.Errorf("device_get_info: %w", err))
		}
		return jsonResult(info)
	})

	type emptyInput struct{}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dns_query",
		Description: "Get the currently configured DNS resolvers.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.DNSQuery(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("dns_query: %w", err))
		}
		return jsonResult(cfg)
	})
}
