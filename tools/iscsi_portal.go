package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// portalListenInput is the wire shape for a single iSCSI portal listen IP/port pair.
type portalListenInput struct {
	IP   string `json:"ip"             jsonschema:"IP address to listen on, e.g. 0.0.0.0"`
	Port int    `json:"port,omitempty" jsonschema:"TCP port to listen on; defaults to 3260"`
}

// registerISCSIPortalTools registers iscsi.portal.* MCP tools.
func registerISCSIPortalTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listISCSIPortalsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of portals to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of portals to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_portal_list",
		Description: "List configured iSCSI portals (listen IP/port groups).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listISCSIPortalsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListISCSIPortals(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_portal_list: %w", err))
		}
		return jsonResult(result)
	})

	type getISCSIPortalInput struct {
		ID int `json:"id" jsonschema:"Numeric iSCSI portal ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_portal_get",
		Description: "Get a single iSCSI portal by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getISCSIPortalInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_portal_get: id must be a positive integer"))
		}
		result, err := client.GetISCSIPortal(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_portal_get: %w", err))
		}
		return jsonResult(result)
	})

	type createISCSIPortalInput struct {
		Listen              []portalListenInput `json:"listen"                          jsonschema:"IP/port pairs to listen on"`
		Comment             string              `json:"comment,omitempty"               jsonschema:"Free-text description"`
		DiscoveryAuthMethod string              `json:"discovery_authmethod,omitempty"  jsonschema:"NONE, CHAP, or CHAP_MUTUAL"`
		DiscoveryAuthGroup  int                 `json:"discovery_authgroup,omitempty"   jsonschema:"CHAP auth group tag (from iscsi_auth_create) used for discovery"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_portal_create",
		Description: "Create a new iSCSI portal.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createISCSIPortalInput) (*mcp.CallToolResult, any, error) {
		if len(p.Listen) == 0 {
			return errorResult(errors.New("iscsi_portal_create: listen must have at least one entry"))
		}
		result, err := client.CreateISCSIPortal(ctx, &truenas.CreateISCSIPortalParams{
			Listen:              toPortalListen(p.Listen),
			Comment:             p.Comment,
			DiscoveryAuthMethod: p.DiscoveryAuthMethod,
			DiscoveryAuthGroup:  p.DiscoveryAuthGroup,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_portal_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateISCSIPortalInput struct {
		ID                  int                 `json:"id"                              jsonschema:"Numeric iSCSI portal ID"`
		Listen              []portalListenInput `json:"listen"                          jsonschema:"IP/port pairs to listen on"`
		Comment             string              `json:"comment,omitempty"               jsonschema:"Free-text description"`
		DiscoveryAuthMethod string              `json:"discovery_authmethod,omitempty"  jsonschema:"NONE, CHAP, or CHAP_MUTUAL"`
		DiscoveryAuthGroup  int                 `json:"discovery_authgroup,omitempty"   jsonschema:"CHAP auth group tag used for discovery"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_portal_update",
		Description: "Update an existing iSCSI portal.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateISCSIPortalInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_portal_update: id must be a positive integer"))
		}
		if len(p.Listen) == 0 {
			return errorResult(errors.New("iscsi_portal_update: listen must have at least one entry"))
		}
		result, err := client.UpdateISCSIPortal(ctx, p.ID, &truenas.CreateISCSIPortalParams{
			Listen:              toPortalListen(p.Listen),
			Comment:             p.Comment,
			DiscoveryAuthMethod: p.DiscoveryAuthMethod,
			DiscoveryAuthGroup:  p.DiscoveryAuthGroup,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_portal_update: %w", err))
		}
		return jsonResult(result)
	})

	// iscsi_portal_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_portal_listen_ip_choices",
		Description: "List IP addresses available for an iSCSI portal to listen on.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.ISCSIPortalListenIPChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_portal_listen_ip_choices: %w", err))
		}
		return jsonResult(choices)
	})
}

func toPortalListen(in []portalListenInput) []truenas.ISCSIPortalListen {
	out := make([]truenas.ISCSIPortalListen, len(in))
	for i, l := range in {
		out[i] = truenas.ISCSIPortalListen{IP: l.IP, Port: l.Port}
	}
	return out
}
