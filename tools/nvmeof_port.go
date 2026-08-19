package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNVMeOFPortTools registers nvmet.port.* MCP tools.
func registerNVMeOFPortTools(s *mcp.Server, client truenasClient) {
	type listNVMeOFPortsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of ports to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of ports to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_list",
		Description: "List configured NVMe-oF ports (transport addresses subsystems are reachable on).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listNVMeOFPortsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListNVMetPorts(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_list: %w", err))
		}
		return jsonResult(result)
	})

	type getNVMeOFPortInput struct {
		ID int `json:"id" jsonschema:"Numeric NVMe-oF port ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_get",
		Description: "Get a single NVMe-oF port by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getNVMeOFPortInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_port_get: id must be a positive integer"))
		}
		result, err := client.GetNVMetPort(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_get: %w", err))
		}
		return jsonResult(result)
	})

	type createNVMeOFPortInput struct {
		AddrTrtype  string `json:"addr_trtype"            jsonschema:"Transport type: TCP, RDMA, or FC"`
		AddrTraddr  string `json:"addr_traddr"            jsonschema:"Local IP address to listen on, from nvmeof_port_transport_address_choices"`
		AddrTrsvcid string `json:"addr_trsvcid,omitempty" jsonschema:"Transport service ID (TCP/RDMA port number); defaults to 4420"`
		AddrAdrfam  string `json:"addr_adrfam,omitempty"  jsonschema:"Address family: ipv4 or ipv6"`
		Enabled     bool   `json:"enabled,omitempty"      jsonschema:"Whether the port is active"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_create",
		Description: "Create a new NVMe-oF port.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createNVMeOFPortInput) (*mcp.CallToolResult, any, error) {
		if p.AddrTrtype == "" || p.AddrTraddr == "" {
			return errorResult(errors.New("nvmeof_port_create: addr_trtype and addr_traddr are required"))
		}
		result, err := client.CreateNVMetPort(ctx, &truenas.CreateNVMetPortParams{
			AddrTrtype: p.AddrTrtype, AddrTraddr: p.AddrTraddr, AddrTrsvcid: p.AddrTrsvcid,
			AddrAdrfam: p.AddrAdrfam, Enabled: p.Enabled,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateNVMeOFPortInput struct {
		ID          int    `json:"id"                     jsonschema:"Numeric NVMe-oF port ID"`
		AddrTrtype  string `json:"addr_trtype"            jsonschema:"Transport type: TCP, RDMA, or FC"`
		AddrTraddr  string `json:"addr_traddr"            jsonschema:"Local IP address to listen on"`
		AddrTrsvcid string `json:"addr_trsvcid,omitempty" jsonschema:"Transport service ID"`
		AddrAdrfam  string `json:"addr_adrfam,omitempty"  jsonschema:"Address family: ipv4 or ipv6"`
		Enabled     bool   `json:"enabled,omitempty"      jsonschema:"Whether the port is active"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_update",
		Description: "Update an existing NVMe-oF port.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNVMeOFPortInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_port_update: id must be a positive integer"))
		}
		if p.AddrTrtype == "" || p.AddrTraddr == "" {
			return errorResult(errors.New("nvmeof_port_update: addr_trtype and addr_traddr are required"))
		}
		result, err := client.UpdateNVMetPort(ctx, p.ID, &truenas.CreateNVMetPortParams{
			AddrTrtype: p.AddrTrtype, AddrTraddr: p.AddrTraddr, AddrTrsvcid: p.AddrTrsvcid,
			AddrAdrfam: p.AddrAdrfam, Enabled: p.Enabled,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_update: %w", err))
		}
		return jsonResult(result)
	})

	// nvmeof_port_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	type nvmeofPortTransportAddressChoicesInput struct {
		Trtype string `json:"trtype,omitempty" jsonschema:"Transport type to filter addresses for: TCP, RDMA, or FC"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_transport_address_choices",
		Description: "List local IP addresses available for a new NVMe-oF port to listen on.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nvmeofPortTransportAddressChoicesInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.NVMetPortTransportAddressChoices(ctx, p.Trtype)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_transport_address_choices: %w", err))
		}
		return jsonResult(choices)
	})
}
