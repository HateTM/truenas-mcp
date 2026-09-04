package truenas

import (
	"context"
	"errors"
	"fmt"
)

// ISCSIExtent represents an iSCSI extent — the backing storage (a zvol disk
// or a file) exposed to initiators through a target.
type ISCSIExtent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // DISK or FILE
	Disk        string `json:"disk,omitempty"`
	Path        string `json:"path,omitempty"`
	Filesize    int64  `json:"filesize,omitempty"`
	Blocksize   int    `json:"blocksize,omitempty"`
	Comment     string `json:"comment,omitempty"`
	InsecureTPC bool   `json:"insecure_tpc,omitempty"`
	ReadOnly    bool   `json:"ro,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
}

// CreateISCSIExtentParams holds fields for creating or updating an iSCSI extent.
type CreateISCSIExtentParams struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Disk        string `json:"disk,omitempty"`
	Path        string `json:"path,omitempty"`
	Filesize    int64  `json:"filesize,omitempty"`
	Blocksize   int    `json:"blocksize,omitempty"`
	Comment     string `json:"comment,omitempty"`
	InsecureTPC bool   `json:"insecure_tpc,omitempty"`
	ReadOnly    bool   `json:"ro,omitempty"`
}

// ISCSIInitiator represents an allowed set of iSCSI initiators (an "initiator group").
type ISCSIInitiator struct {
	ID         int      `json:"id"`
	Initiators []string `json:"initiators,omitempty"`
	Comment    string   `json:"comment,omitempty"`
}

// CreateISCSIInitiatorParams holds fields for creating or updating an initiator group.
type CreateISCSIInitiatorParams struct {
	Initiators []string `json:"initiators,omitempty"`
	Comment    string   `json:"comment,omitempty"`
}

// ListISCSIExtents lists configured iSCSI extents.
func (c *Client) ListISCSIExtents(ctx context.Context, opts ...ListOptions) ([]ISCSIExtent, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ISCSIExtent
	if err := c.call(ctx, "iscsi.extent.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing iscsi extents: %w", err)
	}
	return result, nil
}

// GetISCSIExtent returns a single iSCSI extent by ID.
func (c *Client) GetISCSIExtent(ctx context.Context, id int) (*ISCSIExtent, error) {
	var result ISCSIExtent
	if err := c.call(ctx, "iscsi.extent.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting iscsi extent %d: %w", id, err)
	}
	return &result, nil
}

// CreateISCSIExtent creates a new iSCSI extent.
func (c *Client) CreateISCSIExtent(ctx context.Context, p *CreateISCSIExtentParams) (*ISCSIExtent, error) {
	if p == nil {
		return nil, errors.New("create iscsi extent: params required")
	}
	var result ISCSIExtent
	if err := c.call(ctx, "iscsi.extent.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating iscsi extent %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateISCSIExtent updates an existing iSCSI extent.
func (c *Client) UpdateISCSIExtent(ctx context.Context, id int, p *CreateISCSIExtentParams) (*ISCSIExtent, error) {
	if p == nil {
		return nil, errors.New("update iscsi extent: params required")
	}
	var result ISCSIExtent
	if err := c.call(ctx, "iscsi.extent.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating iscsi extent %d: %w", id, err)
	}
	return &result, nil
}

// DeleteISCSIExtent deletes an iSCSI extent.
func (c *Client) DeleteISCSIExtent(ctx context.Context, id int) error {
	if err := c.call(ctx, "iscsi.extent.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting iscsi extent %d: %w", id, err)
	}
	return nil
}

// ISCSIExtentDiskChoices returns the zvol disks available for use as a new DISK-type extent.
func (c *Client) ISCSIExtentDiskChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "iscsi.extent.disk_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing iscsi extent disk choices: %w", err)
	}
	return choices, nil
}

// ListISCSIInitiators lists configured iSCSI initiator groups.
func (c *Client) ListISCSIInitiators(ctx context.Context, opts ...ListOptions) ([]ISCSIInitiator, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ISCSIInitiator
	if err := c.call(ctx, "iscsi.initiator.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing iscsi initiators: %w", err)
	}
	return result, nil
}

// GetISCSIInitiator returns a single iSCSI initiator group by ID.
func (c *Client) GetISCSIInitiator(ctx context.Context, id int) (*ISCSIInitiator, error) {
	var result ISCSIInitiator
	if err := c.call(ctx, "iscsi.initiator.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting iscsi initiator %d: %w", id, err)
	}
	return &result, nil
}

// CreateISCSIInitiator creates a new iSCSI initiator group.
func (c *Client) CreateISCSIInitiator(ctx context.Context, p *CreateISCSIInitiatorParams) (*ISCSIInitiator, error) {
	if p == nil {
		return nil, errors.New("create iscsi initiator: params required")
	}
	var result ISCSIInitiator
	if err := c.call(ctx, "iscsi.initiator.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating iscsi initiator: %w", err)
	}
	return &result, nil
}

// UpdateISCSIInitiator updates an existing iSCSI initiator group.
func (c *Client) UpdateISCSIInitiator(ctx context.Context, id int, p *CreateISCSIInitiatorParams) (*ISCSIInitiator, error) {
	if p == nil {
		return nil, errors.New("update iscsi initiator: params required")
	}
	var result ISCSIInitiator
	if err := c.call(ctx, "iscsi.initiator.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating iscsi initiator %d: %w", id, err)
	}
	return &result, nil
}

// DeleteISCSIInitiator deletes an iSCSI initiator group.
func (c *Client) DeleteISCSIInitiator(ctx context.Context, id int) error {
	if err := c.call(ctx, "iscsi.initiator.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting iscsi initiator %d: %w", id, err)
	}
	return nil
}
