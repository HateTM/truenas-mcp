package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// targetGroupInput is the wire shape for a single iSCSI target group binding.
type targetGroupInput struct {
	Portal     int    `json:"portal"               jsonschema:"iSCSI portal ID (from iscsi_portal_create)"`
	Initiator  int    `json:"initiator,omitempty"  jsonschema:"iSCSI initiator group ID (from iscsi_initiator_create); omit to allow any initiator"`
	Auth       int    `json:"auth,omitempty"       jsonschema:"CHAP auth group tag (from iscsi_auth_create)"`
	AuthMethod string `json:"authmethod,omitempty" jsonschema:"NONE, CHAP, or CHAP_MUTUAL"`
}

func toTargetGroups(in []targetGroupInput) []truenas.ISCSITargetGroup {
	out := make([]truenas.ISCSITargetGroup, len(in))
	for i, g := range in {
		out[i] = truenas.ISCSITargetGroup{Portal: g.Portal, Initiator: g.Initiator, Auth: g.Auth, AuthMethod: g.AuthMethod}
	}
	return out
}

// registerISCSITargetTools registers iscsi.target.* and iscsi.targetextent.* MCP tools.
func registerISCSITargetTools(s *mcp.Server, client truenasClient) {
	type listISCSITargetsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of targets to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of targets to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_target_list",
		Description: "List configured iSCSI targets.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listISCSITargetsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListISCSITargets(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_target_list: %w", err))
		}
		return jsonResult(result)
	})

	type getISCSITargetInput struct {
		ID int `json:"id" jsonschema:"Numeric iSCSI target ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_target_get",
		Description: "Get a single iSCSI target by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getISCSITargetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_target_get: id must be a positive integer"))
		}
		result, err := client.GetISCSITarget(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_target_get: %w", err))
		}
		return jsonResult(result)
	})

	type createISCSITargetInput struct {
		Name   string             `json:"name"            jsonschema:"Target name (becomes part of the IQN)"`
		Alias  string             `json:"alias,omitempty" jsonschema:"Human-readable alias"`
		Mode   string             `json:"mode,omitempty"  jsonschema:"ISCSI, FC, or BOTH; defaults to ISCSI"`
		Groups []targetGroupInput `json:"groups,omitempty" jsonschema:"Portal/initiator/auth group bindings for this target"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_target_create",
		Description: "Create a new iSCSI target. Use iscsi_target_validate_name first to confirm the name is available.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createISCSITargetInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" {
			return errorResult(errors.New("iscsi_target_create: name must not be empty"))
		}
		result, err := client.CreateISCSITarget(ctx, &truenas.CreateISCSITargetParams{
			Name: p.Name, Alias: p.Alias, Mode: p.Mode, Groups: toTargetGroups(p.Groups),
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_target_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateISCSITargetInput struct {
		ID     int                `json:"id"               jsonschema:"Numeric iSCSI target ID"`
		Name   string             `json:"name"             jsonschema:"Target name"`
		Alias  string             `json:"alias,omitempty"  jsonschema:"Human-readable alias"`
		Mode   string             `json:"mode,omitempty"   jsonschema:"ISCSI, FC, or BOTH"`
		Groups []targetGroupInput `json:"groups,omitempty" jsonschema:"Portal/initiator/auth group bindings for this target"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_target_update",
		Description: "Update an existing iSCSI target.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateISCSITargetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_target_update: id must be a positive integer"))
		}
		if p.Name == "" {
			return errorResult(errors.New("iscsi_target_update: name must not be empty"))
		}
		result, err := client.UpdateISCSITarget(ctx, p.ID, &truenas.CreateISCSITargetParams{
			Name: p.Name, Alias: p.Alias, Mode: p.Mode, Groups: toTargetGroups(p.Groups),
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_target_update: %w", err))
		}
		return jsonResult(result)
	})

	// iscsi_target_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	type validateISCSITargetNameInput struct {
		Name string `json:"name" jsonschema:"Candidate iSCSI target name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_target_validate_name",
		Description: "Check whether a name is valid and not already in use for an iSCSI target.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p validateISCSITargetNameInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" {
			return errorResult(errors.New("iscsi_target_validate_name: name must not be empty"))
		}
		valid, err := client.ValidateISCSITargetName(ctx, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_target_validate_name: %w", err))
		}
		return jsonResult(map[string]bool{"valid": valid})
	})

	type listISCSITargetExtentsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of mappings to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of mappings to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_targetextent_list",
		Description: "List configured iSCSI target/extent mappings (which extent is exposed as which LUN on which target).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listISCSITargetExtentsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListISCSITargetExtents(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_targetextent_list: %w", err))
		}
		return jsonResult(result)
	})

	type getISCSITargetExtentInput struct {
		ID int `json:"id" jsonschema:"Numeric iSCSI target/extent mapping ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_targetextent_get",
		Description: "Get a single iSCSI target/extent mapping by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getISCSITargetExtentInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_targetextent_get: id must be a positive integer"))
		}
		result, err := client.GetISCSITargetExtent(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_targetextent_get: %w", err))
		}
		return jsonResult(result)
	})

	type createISCSITargetExtentInput struct {
		Target int `json:"target"          jsonschema:"iSCSI target ID (from iscsi_target_create)"`
		Extent int `json:"extent"          jsonschema:"iSCSI extent ID (from iscsi_extent_create)"`
		LUNID  int `json:"lunid,omitempty" jsonschema:"LUN ID to expose the extent as; auto-assigned if omitted"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_targetextent_create",
		Description: "Map an iSCSI extent onto a target as a LUN.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createISCSITargetExtentInput) (*mcp.CallToolResult, any, error) {
		if p.Target <= 0 || p.Extent <= 0 {
			return errorResult(errors.New("iscsi_targetextent_create: target and extent are required"))
		}
		result, err := client.CreateISCSITargetExtent(ctx, &truenas.CreateISCSITargetExtentParams{
			Target: p.Target, Extent: p.Extent, LUNID: p.LUNID,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_targetextent_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateISCSITargetExtentInput struct {
		ID     int `json:"id"              jsonschema:"Numeric iSCSI target/extent mapping ID"`
		Target int `json:"target"          jsonschema:"iSCSI target ID"`
		Extent int `json:"extent"          jsonschema:"iSCSI extent ID"`
		LUNID  int `json:"lunid,omitempty" jsonschema:"LUN ID to expose the extent as"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_targetextent_update",
		Description: "Update an existing iSCSI target/extent mapping.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateISCSITargetExtentInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_targetextent_update: id must be a positive integer"))
		}
		if p.Target <= 0 || p.Extent <= 0 {
			return errorResult(errors.New("iscsi_targetextent_update: target and extent are required"))
		}
		result, err := client.UpdateISCSITargetExtent(ctx, p.ID, &truenas.CreateISCSITargetExtentParams{
			Target: p.Target, Extent: p.Extent, LUNID: p.LUNID,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_targetextent_update: %w", err))
		}
		return jsonResult(result)
	})

	// iscsi_targetextent_delete is destructive and lives in destructive.go,
	// gated behind Config.AllowDestructive.
}
