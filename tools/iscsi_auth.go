package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerISCSIAuthTools registers iscsi.auth.* and iscsi.global.* MCP tools.
func registerISCSIAuthTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listISCSIAuthInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of credentials to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of credentials to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_auth_list",
		Description: "List configured iSCSI CHAP authentication credentials.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listISCSIAuthInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListISCSIAuth(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_auth_list: %w", err))
		}
		return jsonResult(result)
	})

	type getISCSIAuthInput struct {
		ID int `json:"id" jsonschema:"Numeric iSCSI auth credential ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_auth_get",
		Description: "Get a single iSCSI CHAP authentication credential by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getISCSIAuthInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_auth_get: id must be a positive integer"))
		}
		result, err := client.GetISCSIAuth(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_auth_get: %w", err))
		}
		return jsonResult(result)
	})

	type createISCSIAuthInput struct {
		Tag        int    `json:"tag"                  jsonschema:"CHAP auth group tag, referenced from a portal's discovery_authgroup or a target group's auth"`
		User       string `json:"user"                 jsonschema:"CHAP username"`
		Secret     string `json:"secret"                jsonschema:"CHAP secret (12-16 characters)"`
		PeerUser   string `json:"peeruser,omitempty"   jsonschema:"Mutual CHAP peer username"`
		PeerSecret string `json:"peersecret,omitempty" jsonschema:"Mutual CHAP peer secret"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_auth_create",
		Description: "Create a new iSCSI CHAP authentication credential.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createISCSIAuthInput) (*mcp.CallToolResult, any, error) {
		if p.User == "" || p.Secret == "" {
			return errorResult(errors.New("iscsi_auth_create: user and secret are required"))
		}
		result, err := client.CreateISCSIAuth(ctx, &truenas.CreateISCSIAuthParams{
			Tag: p.Tag, User: p.User, Secret: p.Secret, PeerUser: p.PeerUser, PeerSecret: p.PeerSecret,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_auth_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateISCSIAuthInput struct {
		ID         int    `json:"id"                   jsonschema:"Numeric iSCSI auth credential ID"`
		Tag        int    `json:"tag"                  jsonschema:"CHAP auth group tag"`
		User       string `json:"user"                 jsonschema:"CHAP username"`
		Secret     string `json:"secret"                jsonschema:"CHAP secret (12-16 characters)"`
		PeerUser   string `json:"peeruser,omitempty"   jsonschema:"Mutual CHAP peer username"`
		PeerSecret string `json:"peersecret,omitempty" jsonschema:"Mutual CHAP peer secret"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_auth_update",
		Description: "Update an existing iSCSI CHAP authentication credential.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateISCSIAuthInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_auth_update: id must be a positive integer"))
		}
		if p.User == "" || p.Secret == "" {
			return errorResult(errors.New("iscsi_auth_update: user and secret are required"))
		}
		result, err := client.UpdateISCSIAuth(ctx, p.ID, &truenas.CreateISCSIAuthParams{
			Tag: p.Tag, User: p.User, Secret: p.Secret, PeerUser: p.PeerUser, PeerSecret: p.PeerSecret,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_auth_update: %w", err))
		}
		return jsonResult(result)
	})

	// iscsi_auth_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_global_alua_enabled",
		Description: "Check whether ALUA (Asymmetric Logical Unit Access) is enabled for the iSCSI service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		enabled, err := client.ISCSIGlobalALUAEnabled(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_global_alua_enabled: %w", err))
		}
		return jsonResult(map[string]bool{"alua_enabled": enabled})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_global_client_count",
		Description: "Get the number of clients currently connected to the iSCSI service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		count, err := client.ISCSIGlobalClientCount(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_global_client_count: %w", err))
		}
		return jsonResult(map[string]int{"client_count": count})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_global_config",
		Description: "Get the global iSCSI service configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.ISCSIGlobalConfigGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_global_config: %w", err))
		}
		return jsonResult(cfg)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_global_iser_enabled",
		Description: "Check whether iSER (iSCSI Extensions for RDMA) is enabled for the iSCSI service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		enabled, err := client.ISCSIGlobalISEREnabled(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_global_iser_enabled: %w", err))
		}
		return jsonResult(map[string]bool{"iser_enabled": enabled})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_global_sessions",
		Description: "List all active iSCSI sessions across all targets.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		sessions, err := client.ISCSIGlobalSessions(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_global_sessions: %w", err))
		}
		return jsonResult(sessions)
	})

	type updateISCSIGlobalInput struct {
		Basename           string   `json:"basename,omitempty"             jsonschema:"Base name (IQN prefix) for auto-generated target/extent names"`
		ISNSServers        []string `json:"isns_servers,omitempty"         jsonschema:"iSNS server hostnames/IPs for target discovery"`
		ListenPort         int      `json:"listen_port,omitempty"          jsonschema:"TCP port the iSCSI service listens on"`
		ALUA               bool     `json:"alua,omitempty"                 jsonschema:"Enable Asymmetric Logical Unit Access"`
		PoolAvailThreshold int      `json:"pool_avail_threshold,omitempty" jsonschema:"Warn when pool free space drops below this percentage"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_global_update",
		Description: "Update the global iSCSI service configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateISCSIGlobalInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.UpdateISCSIGlobal(ctx, &truenas.UpdateISCSIGlobalParams{
			Basename: p.Basename, ISNSServers: p.ISNSServers, ListenPort: p.ListenPort,
			ALUA: p.ALUA, PoolAvailThreshold: p.PoolAvailThreshold,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_global_update: %w", err))
		}
		return jsonResult(cfg)
	})
}
