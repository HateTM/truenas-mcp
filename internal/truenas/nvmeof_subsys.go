package truenas

import (
	"context"
	"errors"
	"fmt"
)

// NVMetHostSubsys binds an NVMe-oF host to a subsystem it is allowed to connect to.
type NVMetHostSubsys struct {
	ID     int `json:"id"`
	Host   int `json:"host"`
	Subsys int `json:"subsys"`
}

// CreateNVMetHostSubsysParams holds fields for creating a host/subsystem binding.
type CreateNVMetHostSubsysParams struct {
	Host   int `json:"host"`
	Subsys int `json:"subsys"`
}

// NVMetSubsys represents an NVMe-oF subsystem.
type NVMetSubsys struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SubNQN       string `json:"subnqn,omitempty"`
	AllowAnyHost bool   `json:"allow_any_host,omitempty"`
	Serial       string `json:"serial,omitempty"`
}

// CreateNVMetSubsysParams holds fields for creating or updating an NVMe-oF subsystem.
type CreateNVMetSubsysParams struct {
	Name         string `json:"name"`
	SubNQN       string `json:"subnqn,omitempty"`
	AllowAnyHost bool   `json:"allow_any_host,omitempty"`
	Serial       string `json:"serial,omitempty"`
}

// ListNVMetHostSubsys lists configured NVMe-oF host/subsystem bindings.
func (c *Client) ListNVMetHostSubsys(ctx context.Context, opts ...ListOptions) ([]NVMetHostSubsys, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []NVMetHostSubsys
	if err := c.call(ctx, "nvmet.host_subsys.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing nvmet host/subsys bindings: %w", err)
	}
	return result, nil
}

// GetNVMetHostSubsys returns a single NVMe-oF host/subsystem binding by ID.
func (c *Client) GetNVMetHostSubsys(ctx context.Context, id int) (*NVMetHostSubsys, error) {
	var result NVMetHostSubsys
	if err := c.call(ctx, "nvmet.host_subsys.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting nvmet host/subsys binding %d: %w", id, err)
	}
	return &result, nil
}

// CreateNVMetHostSubsys creates a new NVMe-oF host/subsystem binding.
func (c *Client) CreateNVMetHostSubsys(ctx context.Context, p *CreateNVMetHostSubsysParams) (*NVMetHostSubsys, error) {
	if p == nil {
		return nil, errors.New("create nvmet host/subsys binding: params required")
	}
	var result NVMetHostSubsys
	if err := c.call(ctx, "nvmet.host_subsys.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating nvmet host/subsys binding: %w", err)
	}
	return &result, nil
}

// UpdateNVMetHostSubsys updates an existing NVMe-oF host/subsystem binding.
func (c *Client) UpdateNVMetHostSubsys(ctx context.Context, id int, p *CreateNVMetHostSubsysParams) (*NVMetHostSubsys, error) {
	if p == nil {
		return nil, errors.New("update nvmet host/subsys binding: params required")
	}
	var result NVMetHostSubsys
	if err := c.call(ctx, "nvmet.host_subsys.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating nvmet host/subsys binding %d: %w", id, err)
	}
	return &result, nil
}

// DeleteNVMetHostSubsys deletes an NVMe-oF host/subsystem binding.
func (c *Client) DeleteNVMetHostSubsys(ctx context.Context, id int) error {
	if err := c.call(ctx, "nvmet.host_subsys.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting nvmet host/subsys binding %d: %w", id, err)
	}
	return nil
}

// ListNVMetSubsys lists configured NVMe-oF subsystems.
func (c *Client) ListNVMetSubsys(ctx context.Context, opts ...ListOptions) ([]NVMetSubsys, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []NVMetSubsys
	if err := c.call(ctx, "nvmet.subsys.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing nvmet subsystems: %w", err)
	}
	return result, nil
}

// GetNVMetSubsys returns a single NVMe-oF subsystem by ID.
func (c *Client) GetNVMetSubsys(ctx context.Context, id int) (*NVMetSubsys, error) {
	var result NVMetSubsys
	if err := c.call(ctx, "nvmet.subsys.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting nvmet subsystem %d: %w", id, err)
	}
	return &result, nil
}

// CreateNVMetSubsys creates a new NVMe-oF subsystem.
func (c *Client) CreateNVMetSubsys(ctx context.Context, p *CreateNVMetSubsysParams) (*NVMetSubsys, error) {
	if p == nil {
		return nil, errors.New("create nvmet subsystem: params required")
	}
	var result NVMetSubsys
	if err := c.call(ctx, "nvmet.subsys.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating nvmet subsystem %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateNVMetSubsys updates an existing NVMe-oF subsystem.
func (c *Client) UpdateNVMetSubsys(ctx context.Context, id int, p *CreateNVMetSubsysParams) (*NVMetSubsys, error) {
	if p == nil {
		return nil, errors.New("update nvmet subsystem: params required")
	}
	var result NVMetSubsys
	if err := c.call(ctx, "nvmet.subsys.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating nvmet subsystem %d: %w", id, err)
	}
	return &result, nil
}

// DeleteNVMetSubsys deletes an NVMe-oF subsystem.
func (c *Client) DeleteNVMetSubsys(ctx context.Context, id int) error {
	if err := c.call(ctx, "nvmet.subsys.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting nvmet subsystem %d: %w", id, err)
	}
	return nil
}
