package truenas

import (
	"context"
	"errors"
	"fmt"
)

// NFSShare represents an NFS export.
type NFSShare struct {
	ID           int      `json:"id"`
	Path         string   `json:"path"`
	Comment      string   `json:"comment,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
	ReadOnly     bool     `json:"ro,omitempty"`
	MaprootUser  string   `json:"maproot_user,omitempty"`
	MaprootGroup string   `json:"maproot_group,omitempty"`
	MapallUser   string   `json:"mapall_user,omitempty"`
	MapallGroup  string   `json:"mapall_group,omitempty"`
	Hosts        []string `json:"hosts,omitempty"`
	Networks     []string `json:"networks,omitempty"`
}

// CreateNFSShareParams holds fields for creating or updating an NFS export.
type CreateNFSShareParams struct {
	Path         string   `json:"path"`
	Comment      string   `json:"comment,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
	ReadOnly     bool     `json:"ro,omitempty"`
	MaprootUser  string   `json:"maproot_user,omitempty"`
	MaprootGroup string   `json:"maproot_group,omitempty"`
	MapallUser   string   `json:"mapall_user,omitempty"`
	MapallGroup  string   `json:"mapall_group,omitempty"`
	Hosts        []string `json:"hosts,omitempty"`
	Networks     []string `json:"networks,omitempty"`
}

// SMBShare represents an SMB share.
type SMBShare struct {
	ID         int      `json:"id"`
	Path       string   `json:"path"`
	Name       string   `json:"name,omitempty"`
	Comment    string   `json:"comment,omitempty"`
	Enabled    bool     `json:"enabled,omitempty"`
	ReadOnly   bool     `json:"ro,omitempty"`
	Browsable  bool     `json:"browsable,omitempty"`
	GuestOK    bool     `json:"guestok,omitempty"`
	HostsAllow []string `json:"hostsallow,omitempty"`
	HostsDeny  []string `json:"hostsdeny,omitempty"`
}

// CreateSMBShareParams holds fields for creating or updating an SMB share.
type CreateSMBShareParams struct {
	Path       string   `json:"path"`
	Name       string   `json:"name,omitempty"`
	Comment    string   `json:"comment,omitempty"`
	Enabled    bool     `json:"enabled,omitempty"`
	ReadOnly   bool     `json:"ro,omitempty"`
	Browsable  bool     `json:"browsable,omitempty"`
	GuestOK    bool     `json:"guestok,omitempty"`
	HostsAllow []string `json:"hostsallow,omitempty"`
	HostsDeny  []string `json:"hostsdeny,omitempty"`
}

// WebDAVShare represents a WebDAV share.
type WebDAVShare struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Comment  string `json:"comment,omitempty"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"ro,omitempty"`
	Perm     bool   `json:"perm,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

// CreateWebDAVShareParams holds fields for creating or updating a WebDAV share.
type CreateWebDAVShareParams struct {
	Name     string `json:"name"`
	Comment  string `json:"comment,omitempty"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"ro,omitempty"`
	Perm     bool   `json:"perm,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

// ListNFSShares lists configured NFS exports.
func (c *Client) ListNFSShares(ctx context.Context, opts ...ListOptions) ([]NFSShare, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []NFSShare
	if err := c.call(ctx, "sharing.nfs.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing NFS shares: %w", err)
	}
	return result, nil
}

// GetNFSShare returns a single NFS export by ID.
func (c *Client) GetNFSShare(ctx context.Context, id int) (*NFSShare, error) {
	var result NFSShare
	if err := c.call(ctx, "sharing.nfs.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting NFS share %d: %w", id, err)
	}
	return &result, nil
}

// CreateNFSShare creates a new NFS export.
func (c *Client) CreateNFSShare(ctx context.Context, p *CreateNFSShareParams) (*NFSShare, error) {
	if p == nil {
		return nil, errors.New("create NFS share: params required")
	}
	var result NFSShare
	if err := c.call(ctx, "sharing.nfs.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating NFS share %q: %w", p.Path, err)
	}
	return &result, nil
}

// UpdateNFSShare updates an existing NFS export.
func (c *Client) UpdateNFSShare(ctx context.Context, id int, p *CreateNFSShareParams) (*NFSShare, error) {
	if p == nil {
		return nil, errors.New("update NFS share: params required")
	}
	var result NFSShare
	if err := c.call(ctx, "sharing.nfs.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating NFS share %d: %w", id, err)
	}
	return &result, nil
}

// DeleteNFSShare deletes an NFS export.
func (c *Client) DeleteNFSShare(ctx context.Context, id int) error {
	if err := c.call(ctx, "sharing.nfs.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting NFS share %d: %w", id, err)
	}
	return nil
}

// ListSMBShares lists configured SMB shares.
func (c *Client) ListSMBShares(ctx context.Context, opts ...ListOptions) ([]SMBShare, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []SMBShare
	if err := c.call(ctx, "sharing.smb.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing SMB shares: %w", err)
	}
	return result, nil
}

// GetSMBShare returns a single SMB share by ID.
func (c *Client) GetSMBShare(ctx context.Context, id int) (*SMBShare, error) {
	var result SMBShare
	if err := c.call(ctx, "sharing.smb.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting SMB share %d: %w", id, err)
	}
	return &result, nil
}

// CreateSMBShare creates a new SMB share.
func (c *Client) CreateSMBShare(ctx context.Context, p *CreateSMBShareParams) (*SMBShare, error) {
	if p == nil {
		return nil, errors.New("create SMB share: params required")
	}
	var result SMBShare
	if err := c.call(ctx, "sharing.smb.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating SMB share %q: %w", p.Path, err)
	}
	return &result, nil
}

// UpdateSMBShare updates an existing SMB share.
func (c *Client) UpdateSMBShare(ctx context.Context, id int, p *CreateSMBShareParams) (*SMBShare, error) {
	if p == nil {
		return nil, errors.New("update SMB share: params required")
	}
	var result SMBShare
	if err := c.call(ctx, "sharing.smb.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating SMB share %d: %w", id, err)
	}
	return &result, nil
}

// DeleteSMBShare deletes an SMB share.
func (c *Client) DeleteSMBShare(ctx context.Context, id int) error {
	if err := c.call(ctx, "sharing.smb.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting SMB share %d: %w", id, err)
	}
	return nil
}

// ListWebDAVShares lists configured WebDAV shares.
func (c *Client) ListWebDAVShares(ctx context.Context, opts ...ListOptions) ([]WebDAVShare, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []WebDAVShare
	if err := c.call(ctx, "sharing.webdav.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing WebDAV shares: %w", err)
	}
	return result, nil
}

// GetWebDAVShare returns a single WebDAV share by ID.
func (c *Client) GetWebDAVShare(ctx context.Context, id int) (*WebDAVShare, error) {
	var result WebDAVShare
	if err := c.call(ctx, "sharing.webdav.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting WebDAV share %d: %w", id, err)
	}
	return &result, nil
}

// CreateWebDAVShare creates a new WebDAV share.
func (c *Client) CreateWebDAVShare(ctx context.Context, p *CreateWebDAVShareParams) (*WebDAVShare, error) {
	if p == nil {
		return nil, errors.New("create WebDAV share: params required")
	}
	var result WebDAVShare
	if err := c.call(ctx, "sharing.webdav.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating WebDAV share %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateWebDAVShare updates an existing WebDAV share.
func (c *Client) UpdateWebDAVShare(ctx context.Context, id int, p *CreateWebDAVShareParams) (*WebDAVShare, error) {
	if p == nil {
		return nil, errors.New("update WebDAV share: params required")
	}
	var result WebDAVShare
	if err := c.call(ctx, "sharing.webdav.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating WebDAV share %d: %w", id, err)
	}
	return &result, nil
}

// DeleteWebDAVShare deletes a WebDAV share.
func (c *Client) DeleteWebDAVShare(ctx context.Context, id int) error {
	if err := c.call(ctx, "sharing.webdav.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting WebDAV share %d: %w", id, err)
	}
	return nil
}
