package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNVMeOFDestructiveTools adds opt-in destructive NVMe-oF tools to the server.
// These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerNVMeOFDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type nvmeofIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_delete",
		Description: "Permanently delete an NVMe-oF host. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nvmeofIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("nvmeof_host_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_host_delete: id must be a positive integer"))
		}
		if err := client.DeleteNVMetHost(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_subsys_delete",
		Description: "Permanently remove an NVMe-oF host/subsystem binding. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nvmeofIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("nvmeof_host_subsys_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_host_subsys_delete: id must be a positive integer"))
		}
		if err := client.DeleteNVMetHostSubsys(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_subsys_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_subsys_delete",
		Description: "Permanently delete an NVMe-oF subsystem. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nvmeofIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("nvmeof_subsys_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_subsys_delete: id must be a positive integer"))
		}
		if err := client.DeleteNVMetSubsys(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("nvmeof_subsys_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_namespace_delete",
		Description: "Permanently delete an NVMe-oF namespace. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nvmeofIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("nvmeof_namespace_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_namespace_delete: id must be a positive integer"))
		}
		if err := client.DeleteNVMetNamespace(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("nvmeof_namespace_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_subsys_delete",
		Description: "Permanently remove an NVMe-oF port/subsystem binding. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nvmeofIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("nvmeof_port_subsys_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_port_subsys_delete: id must be a positive integer"))
		}
		if err := client.DeleteNVMetPortSubsys(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_subsys_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_delete",
		Description: "Permanently delete an NVMe-oF port. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nvmeofIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("nvmeof_port_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_port_delete: id must be a positive integer"))
		}
		if err := client.DeleteNVMetPort(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
