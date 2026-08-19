package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerISCSIExtentTools registers iscsi.extent.* and iscsi.initiator.* MCP tools.
func registerISCSIExtentTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listISCSIExtentsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of extents to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of extents to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_extent_list",
		Description: "List configured iSCSI extents (the backing storage exposed through targets).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listISCSIExtentsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListISCSIExtents(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_extent_list: %w", err))
		}
		return jsonResult(result)
	})

	type getISCSIExtentInput struct {
		ID int `json:"id" jsonschema:"Numeric iSCSI extent ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_extent_get",
		Description: "Get a single iSCSI extent by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getISCSIExtentInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_extent_get: id must be a positive integer"))
		}
		result, err := client.GetISCSIExtent(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_extent_get: %w", err))
		}
		return jsonResult(result)
	})

	type createISCSIExtentInput struct {
		Name        string `json:"name"                  jsonschema:"Extent name"`
		Type        string `json:"type"                  jsonschema:"DISK (backed by a zvol) or FILE (backed by a file on a dataset)"`
		Disk        string `json:"disk,omitempty"        jsonschema:"zvol device path, required when type=DISK, e.g. zvol/Storage/vm-disk"`
		Path        string `json:"path,omitempty"        jsonschema:"File path, required when type=FILE"`
		Filesize    int64  `json:"filesize,omitempty"    jsonschema:"File size in bytes, required when type=FILE"`
		Blocksize   int    `json:"blocksize,omitempty"   jsonschema:"Logical block size in bytes, e.g. 512 or 4096"`
		Comment     string `json:"comment,omitempty"     jsonschema:"Free-text description"`
		InsecureTPC bool   `json:"insecure_tpc,omitempty" jsonschema:"Allow initiators to xcopy without proving access to source"`
		ReadOnly    bool   `json:"ro,omitempty"          jsonschema:"Expose the extent as read-only"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_extent_create",
		Description: "Create a new iSCSI extent backed by a zvol (DISK) or a file (FILE).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createISCSIExtentInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || p.Type == "" {
			return errorResult(errors.New("iscsi_extent_create: name and type are required"))
		}
		if p.Type == "DISK" && p.Disk == "" {
			return errorResult(errors.New("iscsi_extent_create: disk is required when type=DISK"))
		}
		if p.Type == "FILE" && (p.Path == "" || p.Filesize <= 0) {
			return errorResult(errors.New("iscsi_extent_create: path and filesize are required when type=FILE"))
		}
		result, err := client.CreateISCSIExtent(ctx, &truenas.CreateISCSIExtentParams{
			Name: p.Name, Type: p.Type, Disk: p.Disk, Path: p.Path, Filesize: p.Filesize,
			Blocksize: p.Blocksize, Comment: p.Comment, InsecureTPC: p.InsecureTPC, ReadOnly: p.ReadOnly,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_extent_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateISCSIExtentInput struct {
		ID          int    `json:"id"                    jsonschema:"Numeric iSCSI extent ID"`
		Name        string `json:"name"                  jsonschema:"Extent name"`
		Type        string `json:"type"                  jsonschema:"DISK or FILE"`
		Disk        string `json:"disk,omitempty"        jsonschema:"zvol device path, required when type=DISK"`
		Path        string `json:"path,omitempty"        jsonschema:"File path, required when type=FILE"`
		Filesize    int64  `json:"filesize,omitempty"    jsonschema:"File size in bytes"`
		Blocksize   int    `json:"blocksize,omitempty"   jsonschema:"Logical block size in bytes"`
		Comment     string `json:"comment,omitempty"     jsonschema:"Free-text description"`
		InsecureTPC bool   `json:"insecure_tpc,omitempty" jsonschema:"Allow initiators to xcopy without proving access to source"`
		ReadOnly    bool   `json:"ro,omitempty"          jsonschema:"Expose the extent as read-only"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_extent_update",
		Description: "Update an existing iSCSI extent.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateISCSIExtentInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_extent_update: id must be a positive integer"))
		}
		if p.Name == "" || p.Type == "" {
			return errorResult(errors.New("iscsi_extent_update: name and type are required"))
		}
		result, err := client.UpdateISCSIExtent(ctx, p.ID, &truenas.CreateISCSIExtentParams{
			Name: p.Name, Type: p.Type, Disk: p.Disk, Path: p.Path, Filesize: p.Filesize,
			Blocksize: p.Blocksize, Comment: p.Comment, InsecureTPC: p.InsecureTPC, ReadOnly: p.ReadOnly,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_extent_update: %w", err))
		}
		return jsonResult(result)
	})

	// iscsi_extent_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_extent_disk_choices",
		Description: "List zvol disks available for use as a new DISK-type iSCSI extent.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.ISCSIExtentDiskChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_extent_disk_choices: %w", err))
		}
		return jsonResult(choices)
	})

	type listISCSIInitiatorsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of initiator groups to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of initiator groups to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_initiator_list",
		Description: "List configured iSCSI initiator groups (allow-lists of initiator IQNs).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listISCSIInitiatorsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListISCSIInitiators(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_initiator_list: %w", err))
		}
		return jsonResult(result)
	})

	type getISCSIInitiatorInput struct {
		ID int `json:"id" jsonschema:"Numeric iSCSI initiator group ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_initiator_get",
		Description: "Get a single iSCSI initiator group by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getISCSIInitiatorInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_initiator_get: id must be a positive integer"))
		}
		result, err := client.GetISCSIInitiator(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_initiator_get: %w", err))
		}
		return jsonResult(result)
	})

	type createISCSIInitiatorInput struct {
		Initiators []string `json:"initiators,omitempty" jsonschema:"Allowed initiator IQNs; empty means allow all"`
		Comment    string   `json:"comment,omitempty"    jsonschema:"Free-text description"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_initiator_create",
		Description: "Create a new iSCSI initiator group.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createISCSIInitiatorInput) (*mcp.CallToolResult, any, error) {
		result, err := client.CreateISCSIInitiator(ctx, &truenas.CreateISCSIInitiatorParams{
			Initiators: p.Initiators, Comment: p.Comment,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_initiator_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateISCSIInitiatorInput struct {
		ID         int      `json:"id"                   jsonschema:"Numeric iSCSI initiator group ID"`
		Initiators []string `json:"initiators,omitempty" jsonschema:"Allowed initiator IQNs; empty means allow all"`
		Comment    string   `json:"comment,omitempty"    jsonschema:"Free-text description"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "iscsi_initiator_update",
		Description: "Update an existing iSCSI initiator group.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateISCSIInitiatorInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("iscsi_initiator_update: id must be a positive integer"))
		}
		result, err := client.UpdateISCSIInitiator(ctx, p.ID, &truenas.CreateISCSIInitiatorParams{
			Initiators: p.Initiators, Comment: p.Comment,
		})
		if err != nil {
			return errorResult(fmt.Errorf("iscsi_initiator_update: %w", err))
		}
		return jsonResult(result)
	})

	// iscsi_initiator_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.
}
