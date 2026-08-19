package truenas

import (
	"context"
	"errors"
	"fmt"
)

// ISCSITargetGroup binds a portal, an initiator group, and an auth group to an iSCSI target.
type ISCSITargetGroup struct {
	Portal     int    `json:"portal"`
	Initiator  int    `json:"initiator,omitempty"`
	Auth       int    `json:"auth,omitempty"`
	AuthMethod string `json:"authmethod,omitempty"`
}

// ISCSITarget represents an iSCSI target.
type ISCSITarget struct {
	ID     int                `json:"id"`
	Name   string             `json:"name"`
	Alias  string             `json:"alias,omitempty"`
	Mode   string             `json:"mode,omitempty"`
	Groups []ISCSITargetGroup `json:"groups,omitempty"`
}

// CreateISCSITargetParams holds fields for creating or updating an iSCSI target.
type CreateISCSITargetParams struct {
	Name   string             `json:"name"`
	Alias  string             `json:"alias,omitempty"`
	Mode   string             `json:"mode,omitempty"`
	Groups []ISCSITargetGroup `json:"groups,omitempty"`
}

// ISCSITargetExtent binds an extent to a target at a given LUN ID.
type ISCSITargetExtent struct {
	ID     int `json:"id"`
	Target int `json:"target"`
	Extent int `json:"extent"`
	LUNID  int `json:"lunid,omitempty"`
}

// CreateISCSITargetExtentParams holds fields for creating or updating a target/extent mapping.
type CreateISCSITargetExtentParams struct {
	Target int `json:"target"`
	Extent int `json:"extent"`
	LUNID  int `json:"lunid,omitempty"`
}

// ListISCSITargets lists configured iSCSI targets.
func (c *Client) ListISCSITargets(ctx context.Context, opts ...ListOptions) ([]ISCSITarget, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ISCSITarget
	if err := c.call(ctx, "iscsi.target.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing iscsi targets: %w", err)
	}
	return result, nil
}

// GetISCSITarget returns a single iSCSI target by ID.
func (c *Client) GetISCSITarget(ctx context.Context, id int) (*ISCSITarget, error) {
	var result ISCSITarget
	if err := c.call(ctx, "iscsi.target.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting iscsi target %d: %w", id, err)
	}
	return &result, nil
}

// CreateISCSITarget creates a new iSCSI target.
func (c *Client) CreateISCSITarget(ctx context.Context, p *CreateISCSITargetParams) (*ISCSITarget, error) {
	if p == nil {
		return nil, errors.New("create iscsi target: params required")
	}
	var result ISCSITarget
	if err := c.call(ctx, "iscsi.target.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating iscsi target %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateISCSITarget updates an existing iSCSI target.
func (c *Client) UpdateISCSITarget(ctx context.Context, id int, p *CreateISCSITargetParams) (*ISCSITarget, error) {
	if p == nil {
		return nil, errors.New("update iscsi target: params required")
	}
	var result ISCSITarget
	if err := c.call(ctx, "iscsi.target.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating iscsi target %d: %w", id, err)
	}
	return &result, nil
}

// DeleteISCSITarget deletes an iSCSI target.
func (c *Client) DeleteISCSITarget(ctx context.Context, id int) error {
	if err := c.call(ctx, "iscsi.target.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting iscsi target %d: %w", id, err)
	}
	return nil
}

// ValidateISCSITargetName reports whether name is a valid, unused iSCSI target name.
func (c *Client) ValidateISCSITargetName(ctx context.Context, name string) (bool, error) {
	var valid bool
	if err := c.call(ctx, "iscsi.target.validate_name", []any{name}, &valid); err != nil {
		return false, fmt.Errorf("validating iscsi target name %q: %w", name, err)
	}
	return valid, nil
}

// ListISCSITargetExtents lists configured iSCSI target/extent mappings.
func (c *Client) ListISCSITargetExtents(ctx context.Context, opts ...ListOptions) ([]ISCSITargetExtent, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ISCSITargetExtent
	if err := c.call(ctx, "iscsi.targetextent.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing iscsi target/extent mappings: %w", err)
	}
	return result, nil
}

// GetISCSITargetExtent returns a single iSCSI target/extent mapping by ID.
func (c *Client) GetISCSITargetExtent(ctx context.Context, id int) (*ISCSITargetExtent, error) {
	var result ISCSITargetExtent
	if err := c.call(ctx, "iscsi.targetextent.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting iscsi target/extent mapping %d: %w", id, err)
	}
	return &result, nil
}

// CreateISCSITargetExtent creates a new iSCSI target/extent mapping.
func (c *Client) CreateISCSITargetExtent(ctx context.Context, p *CreateISCSITargetExtentParams) (*ISCSITargetExtent, error) {
	if p == nil {
		return nil, errors.New("create iscsi target/extent mapping: params required")
	}
	var result ISCSITargetExtent
	if err := c.call(ctx, "iscsi.targetextent.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating iscsi target/extent mapping: %w", err)
	}
	return &result, nil
}

// UpdateISCSITargetExtent updates an existing iSCSI target/extent mapping.
func (c *Client) UpdateISCSITargetExtent(ctx context.Context, id int, p *CreateISCSITargetExtentParams) (*ISCSITargetExtent, error) {
	if p == nil {
		return nil, errors.New("update iscsi target/extent mapping: params required")
	}
	var result ISCSITargetExtent
	if err := c.call(ctx, "iscsi.targetextent.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating iscsi target/extent mapping %d: %w", id, err)
	}
	return &result, nil
}

// DeleteISCSITargetExtent deletes an iSCSI target/extent mapping.
func (c *Client) DeleteISCSITargetExtent(ctx context.Context, id int) error {
	if err := c.call(ctx, "iscsi.targetextent.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting iscsi target/extent mapping %d: %w", id, err)
	}
	return nil
}
