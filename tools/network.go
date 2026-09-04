package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNetworkTools registers network-related MCP tools onto the server.
func registerNetworkTools(s *mcp.Server, client truenasClient) {
	type listInterfacesInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_interfaces",
		Description: "List all network interfaces configured on the TrueNAS host, including bridges, physical ports, bonds, and VLANs. Useful for finding bridge names when attaching NICs to VMs.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ listInterfacesInput) (*mcp.CallToolResult, any, error) {
		ifaces, err := client.ListInterfaces(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("list_interfaces: %w", err))
		}
		return jsonResult(ifaces)
	})

	type getInterfaceInput struct {
		ID string `json:"id" jsonschema:"Interface ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_interface",
		Description: "Get detailed information about a specific network interface.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getInterfaceInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("get_interface: id is required"))
		}
		iface, err := client.GetInterface(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("get_interface: %w", err))
		}
		return jsonResult(iface)
	})

	type updateInterfaceInput struct {
		ID          string `json:"id"          jsonschema:"Interface ID"`
		Description string `json:"description" jsonschema:"New interface description"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_interface",
		Description: "Update a network interface description.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateInterfaceInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("update_interface: id is required"))
		}
		iface, err := client.UpdateInterface(ctx, p.ID, &truenas.UpdateInterfaceParams{Description: p.Description})
		if err != nil {
			return errorResult(fmt.Errorf("update_interface: %w", err))
		}
		return jsonResult(iface)
	})
}
