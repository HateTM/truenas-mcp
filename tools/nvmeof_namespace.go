package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNVMeOFNamespaceTools registers nvmet.namespace.* and nvmet.port_subsys.* MCP tools.
func registerNVMeOFNamespaceTools(s *mcp.Server, client truenasClient) {
	type listNVMeOFNamespacesInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of namespaces to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of namespaces to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_namespace_list",
		Description: "List configured NVMe-oF namespaces.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listNVMeOFNamespacesInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListNVMetNamespaces(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_namespace_list: %w", err))
		}
		return jsonResult(result)
	})

	type getNVMeOFNamespaceInput struct {
		ID int `json:"id" jsonschema:"Numeric NVMe-oF namespace ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_namespace_get",
		Description: "Get a single NVMe-oF namespace by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getNVMeOFNamespaceInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_namespace_get: id must be a positive integer"))
		}
		result, err := client.GetNVMetNamespace(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_namespace_get: %w", err))
		}
		return jsonResult(result)
	})

	type createNVMeOFNamespaceInput struct {
		Subsys     int    `json:"subsys"               jsonschema:"NVMe-oF subsystem ID (from nvmeof_subsys_create)"`
		NSID       int    `json:"nsid,omitempty"       jsonschema:"Namespace ID within the subsystem; auto-assigned if omitted"`
		DeviceType string `json:"device_type"          jsonschema:"ZVOL (backed by a zvol) or FILE (backed by a file on a dataset)"`
		DevicePath string `json:"device_path"          jsonschema:"zvol device path or file path, matching device_type"`
		Filesize   int64  `json:"filesize,omitempty"   jsonschema:"File size in bytes, required when device_type=FILE"`
		Enabled    bool   `json:"enabled,omitempty"    jsonschema:"Whether the namespace is active"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_namespace_create",
		Description: "Create a new NVMe-oF namespace backed by a zvol (ZVOL) or a file (FILE), within a subsystem.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createNVMeOFNamespaceInput) (*mcp.CallToolResult, any, error) {
		if p.Subsys <= 0 {
			return errorResult(errors.New("nvmeof_namespace_create: subsys is required"))
		}
		if p.DeviceType == "" || p.DevicePath == "" {
			return errorResult(errors.New("nvmeof_namespace_create: device_type and device_path are required"))
		}
		if p.DeviceType == "FILE" && p.Filesize <= 0 {
			return errorResult(errors.New("nvmeof_namespace_create: filesize is required when device_type=FILE"))
		}
		result, err := client.CreateNVMetNamespace(ctx, &truenas.CreateNVMetNamespaceParams{
			Subsys: p.Subsys, NSID: p.NSID, DeviceType: p.DeviceType, DevicePath: p.DevicePath,
			Filesize: p.Filesize, Enabled: p.Enabled,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_namespace_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateNVMeOFNamespaceInput struct {
		ID         int    `json:"id"                   jsonschema:"Numeric NVMe-oF namespace ID"`
		Subsys     int    `json:"subsys"               jsonschema:"NVMe-oF subsystem ID"`
		NSID       int    `json:"nsid,omitempty"       jsonschema:"Namespace ID within the subsystem"`
		DeviceType string `json:"device_type"          jsonschema:"ZVOL or FILE"`
		DevicePath string `json:"device_path"          jsonschema:"zvol device path or file path"`
		Filesize   int64  `json:"filesize,omitempty"   jsonschema:"File size in bytes"`
		Enabled    bool   `json:"enabled,omitempty"    jsonschema:"Whether the namespace is active"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_namespace_update",
		Description: "Update an existing NVMe-oF namespace.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNVMeOFNamespaceInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_namespace_update: id must be a positive integer"))
		}
		if p.Subsys <= 0 {
			return errorResult(errors.New("nvmeof_namespace_update: subsys is required"))
		}
		result, err := client.UpdateNVMetNamespace(ctx, p.ID, &truenas.CreateNVMetNamespaceParams{
			Subsys: p.Subsys, NSID: p.NSID, DeviceType: p.DeviceType, DevicePath: p.DevicePath,
			Filesize: p.Filesize, Enabled: p.Enabled,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_namespace_update: %w", err))
		}
		return jsonResult(result)
	})

	// nvmeof_namespace_delete is destructive and lives in destructive.go,
	// gated behind Config.AllowDestructive.

	type listNVMeOFPortSubsysInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of bindings to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of bindings to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_subsys_list",
		Description: "List configured NVMe-oF port/subsystem bindings (which subsystems are exposed on which ports).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listNVMeOFPortSubsysInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListNVMetPortSubsys(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_subsys_list: %w", err))
		}
		return jsonResult(result)
	})

	type getNVMeOFPortSubsysInput struct {
		ID int `json:"id" jsonschema:"Numeric NVMe-oF port/subsystem binding ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_subsys_get",
		Description: "Get a single NVMe-oF port/subsystem binding by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getNVMeOFPortSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_port_subsys_get: id must be a positive integer"))
		}
		result, err := client.GetNVMetPortSubsys(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_subsys_get: %w", err))
		}
		return jsonResult(result)
	})

	type createNVMeOFPortSubsysInput struct {
		Port   int `json:"port"   jsonschema:"NVMe-oF port ID (from nvmeof_port_create)"`
		Subsys int `json:"subsys" jsonschema:"NVMe-oF subsystem ID (from nvmeof_subsys_create)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_subsys_create",
		Description: "Expose an NVMe-oF subsystem on a port.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createNVMeOFPortSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.Port <= 0 || p.Subsys <= 0 {
			return errorResult(errors.New("nvmeof_port_subsys_create: port and subsys are required"))
		}
		result, err := client.CreateNVMetPortSubsys(ctx, &truenas.CreateNVMetPortSubsysParams{Port: p.Port, Subsys: p.Subsys})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_subsys_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateNVMeOFPortSubsysInput struct {
		ID     int `json:"id"     jsonschema:"Numeric NVMe-oF port/subsystem binding ID"`
		Port   int `json:"port"   jsonschema:"NVMe-oF port ID"`
		Subsys int `json:"subsys" jsonschema:"NVMe-oF subsystem ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_port_subsys_update",
		Description: "Update an existing NVMe-oF port/subsystem binding.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNVMeOFPortSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_port_subsys_update: id must be a positive integer"))
		}
		if p.Port <= 0 || p.Subsys <= 0 {
			return errorResult(errors.New("nvmeof_port_subsys_update: port and subsys are required"))
		}
		result, err := client.UpdateNVMetPortSubsys(ctx, p.ID, &truenas.CreateNVMetPortSubsysParams{Port: p.Port, Subsys: p.Subsys})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_port_subsys_update: %w", err))
		}
		return jsonResult(result)
	})

	// nvmeof_port_subsys_delete is destructive and lives in destructive.go,
	// gated behind Config.AllowDestructive.
}
