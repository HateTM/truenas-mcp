package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNVMeOFSubsysTools registers nvmet.host_subsys.* and nvmet.subsys.* MCP tools.
func registerNVMeOFSubsysTools(s *mcp.Server, client truenasClient) {
	type listNVMeOFHostSubsysInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of bindings to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of bindings to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_subsys_list",
		Description: "List configured NVMe-oF host/subsystem bindings (which hosts may connect to which subsystems).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listNVMeOFHostSubsysInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListNVMetHostSubsys(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_subsys_list: %w", err))
		}
		return jsonResult(result)
	})

	type getNVMeOFHostSubsysInput struct {
		ID int `json:"id" jsonschema:"Numeric NVMe-oF host/subsystem binding ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_subsys_get",
		Description: "Get a single NVMe-oF host/subsystem binding by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getNVMeOFHostSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_host_subsys_get: id must be a positive integer"))
		}
		result, err := client.GetNVMetHostSubsys(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_subsys_get: %w", err))
		}
		return jsonResult(result)
	})

	type createNVMeOFHostSubsysInput struct {
		Host   int `json:"host"   jsonschema:"NVMe-oF host ID (from nvmeof_host_create)"`
		Subsys int `json:"subsys" jsonschema:"NVMe-oF subsystem ID (from nvmeof_subsys_create)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_subsys_create",
		Description: "Allow an NVMe-oF host to connect to a subsystem.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createNVMeOFHostSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.Host <= 0 || p.Subsys <= 0 {
			return errorResult(errors.New("nvmeof_host_subsys_create: host and subsys are required"))
		}
		result, err := client.CreateNVMetHostSubsys(ctx, &truenas.CreateNVMetHostSubsysParams{Host: p.Host, Subsys: p.Subsys})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_subsys_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateNVMeOFHostSubsysInput struct {
		ID     int `json:"id"     jsonschema:"Numeric NVMe-oF host/subsystem binding ID"`
		Host   int `json:"host"   jsonschema:"NVMe-oF host ID"`
		Subsys int `json:"subsys" jsonschema:"NVMe-oF subsystem ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_subsys_update",
		Description: "Update an existing NVMe-oF host/subsystem binding.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNVMeOFHostSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_host_subsys_update: id must be a positive integer"))
		}
		if p.Host <= 0 || p.Subsys <= 0 {
			return errorResult(errors.New("nvmeof_host_subsys_update: host and subsys are required"))
		}
		result, err := client.UpdateNVMetHostSubsys(ctx, p.ID, &truenas.CreateNVMetHostSubsysParams{Host: p.Host, Subsys: p.Subsys})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_subsys_update: %w", err))
		}
		return jsonResult(result)
	})

	// nvmeof_host_subsys_delete is destructive and lives in destructive.go,
	// gated behind Config.AllowDestructive.

	type listNVMeOFSubsysInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of subsystems to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of subsystems to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_subsys_list",
		Description: "List configured NVMe-oF subsystems.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listNVMeOFSubsysInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListNVMetSubsys(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_subsys_list: %w", err))
		}
		return jsonResult(result)
	})

	type getNVMeOFSubsysInput struct {
		ID int `json:"id" jsonschema:"Numeric NVMe-oF subsystem ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_subsys_get",
		Description: "Get a single NVMe-oF subsystem by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getNVMeOFSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_subsys_get: id must be a positive integer"))
		}
		result, err := client.GetNVMetSubsys(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_subsys_get: %w", err))
		}
		return jsonResult(result)
	})

	type createNVMeOFSubsysInput struct {
		Name         string `json:"name"                     jsonschema:"Subsystem name"`
		SubNQN       string `json:"subnqn,omitempty"         jsonschema:"Subsystem NQN; auto-generated from name if omitted"`
		AllowAnyHost bool   `json:"allow_any_host,omitempty" jsonschema:"Allow any host to connect without an explicit binding"`
		Serial       string `json:"serial,omitempty"         jsonschema:"Serial number reported to initiators"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_subsys_create",
		Description: "Create a new NVMe-oF subsystem.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createNVMeOFSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" {
			return errorResult(errors.New("nvmeof_subsys_create: name must not be empty"))
		}
		result, err := client.CreateNVMetSubsys(ctx, &truenas.CreateNVMetSubsysParams{
			Name: p.Name, SubNQN: p.SubNQN, AllowAnyHost: p.AllowAnyHost, Serial: p.Serial,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_subsys_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateNVMeOFSubsysInput struct {
		ID           int    `json:"id"                       jsonschema:"Numeric NVMe-oF subsystem ID"`
		Name         string `json:"name"                     jsonschema:"Subsystem name"`
		SubNQN       string `json:"subnqn,omitempty"         jsonschema:"Subsystem NQN"`
		AllowAnyHost bool   `json:"allow_any_host,omitempty" jsonschema:"Allow any host to connect without an explicit binding"`
		Serial       string `json:"serial,omitempty"         jsonschema:"Serial number reported to initiators"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_subsys_update",
		Description: "Update an existing NVMe-oF subsystem.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNVMeOFSubsysInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_subsys_update: id must be a positive integer"))
		}
		if p.Name == "" {
			return errorResult(errors.New("nvmeof_subsys_update: name must not be empty"))
		}
		result, err := client.UpdateNVMetSubsys(ctx, p.ID, &truenas.CreateNVMetSubsysParams{
			Name: p.Name, SubNQN: p.SubNQN, AllowAnyHost: p.AllowAnyHost, Serial: p.Serial,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_subsys_update: %w", err))
		}
		return jsonResult(result)
	})

	// nvmeof_subsys_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.
}
