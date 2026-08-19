package tools

import (
	"context"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerDirectoryServicesTools registers directoryservices.* MCP tools onto the server.
func registerDirectoryServicesTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "directoryservices_cache_refresh",
		Description: "Refresh the cached identities (users/groups) from the directory service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		if err := client.DirectoryServicesCacheRefresh(ctx); err != nil {
			return errorResult(fmt.Errorf("directoryservices_cache_refresh: %w", err))
		}
		return jsonResult(map[string]bool{"refreshed": true})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "directoryservices_certificate_choices",
		Description: "List the certificates available for the directory service connection.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.DirectoryServicesCertificateChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("directoryservices_certificate_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "directoryservices_config",
		Description: "Get the directory service (AD/LDAP) configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.DirectoryServicesConfigGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("directoryservices_config: %w", err))
		}
		return jsonResult(cfg)
	})

	type leaveInput struct {
		Params map[string]any `json:"params" jsonschema:"Admin credentials required to leave the domain, e.g. {\"username\": \"admin\", \"password\": \"...\"}"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "directoryservices_leave",
		Description: "Leave the currently joined Active Directory or LDAP domain.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p leaveInput) (*mcp.CallToolResult, any, error) {
		if err := client.DirectoryServicesLeave(ctx, p.Params); err != nil {
			return errorResult(fmt.Errorf("directoryservices_leave: %w", err))
		}
		return jsonResult(map[string]bool{"left": true})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "directoryservices_status",
		Description: "Get the current connection status of the directory service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		status, err := client.DirectoryServicesStatusGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("directoryservices_status: %w", err))
		}
		return jsonResult(status)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "directoryservices_sync_keytab",
		Description: "Synchronize the Kerberos keytab with the directory service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		if err := client.DirectoryServicesSyncKeytab(ctx); err != nil {
			return errorResult(fmt.Errorf("directoryservices_sync_keytab: %w", err))
		}
		return jsonResult(map[string]bool{"synced": true})
	})

	type updateDirectoryServicesInput struct {
		Service       string         `json:"service,omitempty"       jsonschema:"ACTIVEDIRECTORY or LDAP"`
		Enable        bool           `json:"enable,omitempty"        jsonschema:"Whether the directory service is active"`
		Configuration map[string]any `json:"configuration,omitempty" jsonschema:"Service-specific configuration fields, e.g. domain, bindname, bindpw for AD"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "directoryservices_update",
		Description: "Update the directory service (AD/LDAP) configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateDirectoryServicesInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.UpdateDirectoryServices(ctx, &truenas.UpdateDirectoryServicesParams{
			Service: p.Service, Enable: p.Enable, Configuration: p.Configuration,
		})
		if err != nil {
			return errorResult(fmt.Errorf("directoryservices_update: %w", err))
		}
		return jsonResult(cfg)
	})
}
