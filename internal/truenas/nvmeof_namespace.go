package truenas

import (
	"context"
	"errors"
	"fmt"
)

// NVMetNamespace represents an NVMe-oF namespace — the backing storage (a
// zvol or a file) exposed as a namespace within a subsystem.
type NVMetNamespace struct {
	ID         int    `json:"id"`
	Subsys     int    `json:"subsys"`
	NSID       int    `json:"nsid,omitempty"`
	DeviceType string `json:"device_type"` // ZVOL or FILE
	DevicePath string `json:"device_path"`
	Filesize   int64  `json:"filesize,omitempty"`
	Enabled    bool   `json:"enabled,omitempty"`
}

// CreateNVMetNamespaceParams holds fields for creating or updating an NVMe-oF namespace.
type CreateNVMetNamespaceParams struct {
	Subsys     int    `json:"subsys"`
	NSID       int    `json:"nsid,omitempty"`
	DeviceType string `json:"device_type"`
	DevicePath string `json:"device_path"`
	Filesize   int64  `json:"filesize,omitempty"`
	Enabled    bool   `json:"enabled,omitempty"`
}

// NVMetPortSubsys binds an NVMe-oF port to a subsystem it exposes.
type NVMetPortSubsys struct {
	ID     int `json:"id"`
	Port   int `json:"port"`
	Subsys int `json:"subsys"`
}

// CreateNVMetPortSubsysParams holds fields for creating a port/subsystem binding.
type CreateNVMetPortSubsysParams struct {
	Port   int `json:"port"`
	Subsys int `json:"subsys"`
}

// ListNVMetNamespaces lists configured NVMe-oF namespaces.
func (c *Client) ListNVMetNamespaces(ctx context.Context, opts ...ListOptions) ([]NVMetNamespace, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []NVMetNamespace
	if err := c.call(ctx, "nvmet.namespace.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing nvmet namespaces: %w", err)
	}
	return result, nil
}

// GetNVMetNamespace returns a single NVMe-oF namespace by ID.
func (c *Client) GetNVMetNamespace(ctx context.Context, id int) (*NVMetNamespace, error) {
	var result NVMetNamespace
	if err := c.call(ctx, "nvmet.namespace.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting nvmet namespace %d: %w", id, err)
	}
	return &result, nil
}

// CreateNVMetNamespace creates a new NVMe-oF namespace.
func (c *Client) CreateNVMetNamespace(ctx context.Context, p *CreateNVMetNamespaceParams) (*NVMetNamespace, error) {
	if p == nil {
		return nil, errors.New("create nvmet namespace: params required")
	}
	var result NVMetNamespace
	if err := c.call(ctx, "nvmet.namespace.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating nvmet namespace: %w", err)
	}
	return &result, nil
}

// UpdateNVMetNamespace updates an existing NVMe-oF namespace.
func (c *Client) UpdateNVMetNamespace(ctx context.Context, id int, p *CreateNVMetNamespaceParams) (*NVMetNamespace, error) {
	if p == nil {
		return nil, errors.New("update nvmet namespace: params required")
	}
	var result NVMetNamespace
	if err := c.call(ctx, "nvmet.namespace.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating nvmet namespace %d: %w", id, err)
	}
	return &result, nil
}

// DeleteNVMetNamespace deletes an NVMe-oF namespace.
func (c *Client) DeleteNVMetNamespace(ctx context.Context, id int) error {
	if err := c.call(ctx, "nvmet.namespace.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting nvmet namespace %d: %w", id, err)
	}
	return nil
}

// ListNVMetPortSubsys lists configured NVMe-oF port/subsystem bindings.
func (c *Client) ListNVMetPortSubsys(ctx context.Context, opts ...ListOptions) ([]NVMetPortSubsys, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []NVMetPortSubsys
	if err := c.call(ctx, "nvmet.port_subsys.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing nvmet port/subsys bindings: %w", err)
	}
	return result, nil
}

// GetNVMetPortSubsys returns a single NVMe-oF port/subsystem binding by ID.
func (c *Client) GetNVMetPortSubsys(ctx context.Context, id int) (*NVMetPortSubsys, error) {
	var result NVMetPortSubsys
	if err := c.call(ctx, "nvmet.port_subsys.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting nvmet port/subsys binding %d: %w", id, err)
	}
	return &result, nil
}

// CreateNVMetPortSubsys creates a new NVMe-oF port/subsystem binding.
func (c *Client) CreateNVMetPortSubsys(ctx context.Context, p *CreateNVMetPortSubsysParams) (*NVMetPortSubsys, error) {
	if p == nil {
		return nil, errors.New("create nvmet port/subsys binding: params required")
	}
	var result NVMetPortSubsys
	if err := c.call(ctx, "nvmet.port_subsys.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating nvmet port/subsys binding: %w", err)
	}
	return &result, nil
}

// UpdateNVMetPortSubsys updates an existing NVMe-oF port/subsystem binding.
func (c *Client) UpdateNVMetPortSubsys(ctx context.Context, id int, p *CreateNVMetPortSubsysParams) (*NVMetPortSubsys, error) {
	if p == nil {
		return nil, errors.New("update nvmet port/subsys binding: params required")
	}
	var result NVMetPortSubsys
	if err := c.call(ctx, "nvmet.port_subsys.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating nvmet port/subsys binding %d: %w", id, err)
	}
	return &result, nil
}

// DeleteNVMetPortSubsys deletes an NVMe-oF port/subsystem binding.
func (c *Client) DeleteNVMetPortSubsys(ctx context.Context, id int) error {
	if err := c.call(ctx, "nvmet.port_subsys.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting nvmet port/subsys binding %d: %w", id, err)
	}
	return nil
}
