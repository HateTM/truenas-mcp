package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerSharingTools registers sharing.nfs.*, sharing.smb.*, and sharing.webdav.* MCP tools.
func registerSharingTools(s *mcp.Server, client truenasClient) {
	type listSharesInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of shares to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of shares to skip; 0 means start from the beginning"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_nfs_list",
		Description: "List configured NFS exports.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listSharesInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListNFSShares(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_nfs_list: %w", err))
		}
		return jsonResult(result)
	})

	type getShareInput struct {
		ID int `json:"id" jsonschema:"Numeric share ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_nfs_get",
		Description: "Get a single NFS export by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getShareInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_nfs_get: id must be a positive integer"))
		}
		result, err := client.GetNFSShare(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("sharing_nfs_get: %w", err))
		}
		return jsonResult(result)
	})

	type nfsShareInput struct {
		Path         string   `json:"path"                    jsonschema:"Dataset path to export, e.g. /mnt/Storage/backups"`
		Comment      string   `json:"comment,omitempty"       jsonschema:"Free-text description"`
		Enabled      bool     `json:"enabled,omitempty"       jsonschema:"Whether the export is active"`
		ReadOnly     bool     `json:"ro,omitempty"            jsonschema:"Export as read-only"`
		MaprootUser  string   `json:"maproot_user,omitempty"  jsonschema:"Map root user connections to this local user"`
		MaprootGroup string   `json:"maproot_group,omitempty" jsonschema:"Map root user connections to this local group"`
		MapallUser   string   `json:"mapall_user,omitempty"   jsonschema:"Map all client connections to this local user"`
		MapallGroup  string   `json:"mapall_group,omitempty"  jsonschema:"Map all client connections to this local group"`
		Hosts        []string `json:"hosts,omitempty"         jsonschema:"Hostnames/IPs allowed to mount; empty means allow all"`
		Networks     []string `json:"networks,omitempty"      jsonschema:"CIDR networks allowed to mount; empty means allow all"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_nfs_create",
		Description: "Create a new NFS export.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p nfsShareInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("sharing_nfs_create: path must not be empty"))
		}
		result, err := client.CreateNFSShare(ctx, &truenas.CreateNFSShareParams{
			Path: p.Path, Comment: p.Comment, Enabled: p.Enabled, ReadOnly: p.ReadOnly,
			MaprootUser: p.MaprootUser, MaprootGroup: p.MaprootGroup,
			MapallUser: p.MapallUser, MapallGroup: p.MapallGroup,
			Hosts: p.Hosts, Networks: p.Networks,
		})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_nfs_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateNFSShareInput struct {
		ID           int      `json:"id"                      jsonschema:"Numeric share ID"`
		Path         string   `json:"path"                    jsonschema:"Dataset path to export"`
		Comment      string   `json:"comment,omitempty"       jsonschema:"Free-text description"`
		Enabled      bool     `json:"enabled,omitempty"       jsonschema:"Whether the export is active"`
		ReadOnly     bool     `json:"ro,omitempty"            jsonschema:"Export as read-only"`
		MaprootUser  string   `json:"maproot_user,omitempty"  jsonschema:"Map root user connections to this local user"`
		MaprootGroup string   `json:"maproot_group,omitempty" jsonschema:"Map root user connections to this local group"`
		MapallUser   string   `json:"mapall_user,omitempty"   jsonschema:"Map all client connections to this local user"`
		MapallGroup  string   `json:"mapall_group,omitempty"  jsonschema:"Map all client connections to this local group"`
		Hosts        []string `json:"hosts,omitempty"         jsonschema:"Hostnames/IPs allowed to mount"`
		Networks     []string `json:"networks,omitempty"      jsonschema:"CIDR networks allowed to mount"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_nfs_update",
		Description: "Update an existing NFS export.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNFSShareInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_nfs_update: id must be a positive integer"))
		}
		if p.Path == "" {
			return errorResult(errors.New("sharing_nfs_update: path must not be empty"))
		}
		result, err := client.UpdateNFSShare(ctx, p.ID, &truenas.CreateNFSShareParams{
			Path: p.Path, Comment: p.Comment, Enabled: p.Enabled, ReadOnly: p.ReadOnly,
			MaprootUser: p.MaprootUser, MaprootGroup: p.MaprootGroup,
			MapallUser: p.MapallUser, MapallGroup: p.MapallGroup,
			Hosts: p.Hosts, Networks: p.Networks,
		})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_nfs_update: %w", err))
		}
		return jsonResult(result)
	})

	// sharing_nfs_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_smb_list",
		Description: "List configured SMB shares.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listSharesInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListSMBShares(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_smb_list: %w", err))
		}
		return jsonResult(result)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_smb_get",
		Description: "Get a single SMB share by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getShareInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_smb_get: id must be a positive integer"))
		}
		result, err := client.GetSMBShare(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("sharing_smb_get: %w", err))
		}
		return jsonResult(result)
	})

	type smbShareInput struct {
		Path       string   `json:"path"                 jsonschema:"Dataset path to share, e.g. /mnt/Storage/media"`
		Name       string   `json:"name,omitempty"       jsonschema:"Share name visible to clients; defaults to the last path component"`
		Comment    string   `json:"comment,omitempty"    jsonschema:"Free-text description"`
		Enabled    bool     `json:"enabled,omitempty"    jsonschema:"Whether the share is active"`
		ReadOnly   bool     `json:"ro,omitempty"         jsonschema:"Share as read-only"`
		Browsable  bool     `json:"browsable,omitempty"  jsonschema:"Show the share in network browse lists"`
		GuestOK    bool     `json:"guestok,omitempty"    jsonschema:"Allow guest (unauthenticated) access"`
		HostsAllow []string `json:"hostsallow,omitempty" jsonschema:"Hostnames/IPs explicitly allowed to connect"`
		HostsDeny  []string `json:"hostsdeny,omitempty"  jsonschema:"Hostnames/IPs explicitly denied"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_smb_create",
		Description: "Create a new SMB share.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p smbShareInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("sharing_smb_create: path must not be empty"))
		}
		result, err := client.CreateSMBShare(ctx, &truenas.CreateSMBShareParams{
			Path: p.Path, Name: p.Name, Comment: p.Comment, Enabled: p.Enabled, ReadOnly: p.ReadOnly,
			Browsable: p.Browsable, GuestOK: p.GuestOK, HostsAllow: p.HostsAllow, HostsDeny: p.HostsDeny,
		})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_smb_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateSMBShareInput struct {
		ID         int      `json:"id"                   jsonschema:"Numeric share ID"`
		Path       string   `json:"path"                 jsonschema:"Dataset path to share"`
		Name       string   `json:"name,omitempty"       jsonschema:"Share name visible to clients"`
		Comment    string   `json:"comment,omitempty"    jsonschema:"Free-text description"`
		Enabled    bool     `json:"enabled,omitempty"    jsonschema:"Whether the share is active"`
		ReadOnly   bool     `json:"ro,omitempty"         jsonschema:"Share as read-only"`
		Browsable  bool     `json:"browsable,omitempty"  jsonschema:"Show the share in network browse lists"`
		GuestOK    bool     `json:"guestok,omitempty"    jsonschema:"Allow guest (unauthenticated) access"`
		HostsAllow []string `json:"hostsallow,omitempty" jsonschema:"Hostnames/IPs explicitly allowed to connect"`
		HostsDeny  []string `json:"hostsdeny,omitempty"  jsonschema:"Hostnames/IPs explicitly denied"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_smb_update",
		Description: "Update an existing SMB share.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateSMBShareInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_smb_update: id must be a positive integer"))
		}
		if p.Path == "" {
			return errorResult(errors.New("sharing_smb_update: path must not be empty"))
		}
		result, err := client.UpdateSMBShare(ctx, p.ID, &truenas.CreateSMBShareParams{
			Path: p.Path, Name: p.Name, Comment: p.Comment, Enabled: p.Enabled, ReadOnly: p.ReadOnly,
			Browsable: p.Browsable, GuestOK: p.GuestOK, HostsAllow: p.HostsAllow, HostsDeny: p.HostsDeny,
		})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_smb_update: %w", err))
		}
		return jsonResult(result)
	})

	// sharing_smb_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_webdav_list",
		Description: "List configured WebDAV shares.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listSharesInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListWebDAVShares(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_webdav_list: %w", err))
		}
		return jsonResult(result)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_webdav_get",
		Description: "Get a single WebDAV share by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getShareInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_webdav_get: id must be a positive integer"))
		}
		result, err := client.GetWebDAVShare(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("sharing_webdav_get: %w", err))
		}
		return jsonResult(result)
	})

	type webdavShareInput struct {
		Name     string `json:"name"               jsonschema:"Share name, appears in the WebDAV URL"`
		Comment  string `json:"comment,omitempty"  jsonschema:"Free-text description"`
		Path     string `json:"path"               jsonschema:"Dataset path to share, e.g. /mnt/Storage/docs"`
		ReadOnly bool   `json:"ro,omitempty"       jsonschema:"Share as read-only"`
		Perm     bool   `json:"perm,omitempty"     jsonschema:"Recursively set default WebDAV permissions on the path"`
		Enabled  bool   `json:"enabled,omitempty"  jsonschema:"Whether the share is active"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_webdav_create",
		Description: "Create a new WebDAV share.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p webdavShareInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || p.Path == "" {
			return errorResult(errors.New("sharing_webdav_create: name and path are required"))
		}
		result, err := client.CreateWebDAVShare(ctx, &truenas.CreateWebDAVShareParams{
			Name: p.Name, Comment: p.Comment, Path: p.Path, ReadOnly: p.ReadOnly, Perm: p.Perm, Enabled: p.Enabled,
		})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_webdav_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateWebDAVShareInput struct {
		ID       int    `json:"id"                 jsonschema:"Numeric share ID"`
		Name     string `json:"name"               jsonschema:"Share name"`
		Comment  string `json:"comment,omitempty"  jsonschema:"Free-text description"`
		Path     string `json:"path"               jsonschema:"Dataset path to share"`
		ReadOnly bool   `json:"ro,omitempty"       jsonschema:"Share as read-only"`
		Perm     bool   `json:"perm,omitempty"     jsonschema:"Recursively set default WebDAV permissions on the path"`
		Enabled  bool   `json:"enabled,omitempty"  jsonschema:"Whether the share is active"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sharing_webdav_update",
		Description: "Update an existing WebDAV share.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateWebDAVShareInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("sharing_webdav_update: id must be a positive integer"))
		}
		if p.Name == "" || p.Path == "" {
			return errorResult(errors.New("sharing_webdav_update: name and path are required"))
		}
		result, err := client.UpdateWebDAVShare(ctx, p.ID, &truenas.CreateWebDAVShareParams{
			Name: p.Name, Comment: p.Comment, Path: p.Path, ReadOnly: p.ReadOnly, Perm: p.Perm, Enabled: p.Enabled,
		})
		if err != nil {
			return errorResult(fmt.Errorf("sharing_webdav_update: %w", err))
		}
		return jsonResult(result)
	})

	// sharing_webdav_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.
}
