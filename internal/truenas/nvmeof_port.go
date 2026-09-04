package truenas

import (
	"context"
	"errors"
	"fmt"
)

// NVMetPort represents an NVMe-oF port — the transport address an NVMe-oF
// subsystem is reachable on.
type NVMetPort struct {
	ID          int    `json:"id"`
	AddrTrtype  string `json:"addr_trtype"` // TCP, RDMA, or FC
	AddrTraddr  string `json:"addr_traddr"`
	AddrTrsvcid string `json:"addr_trsvcid,omitempty"`
	AddrAdrfam  string `json:"addr_adrfam,omitempty"` // ipv4 or ipv6
	Enabled     bool   `json:"enabled,omitempty"`
}

// CreateNVMetPortParams holds fields for creating or updating an NVMe-oF port.
type CreateNVMetPortParams struct {
	AddrTrtype  string `json:"addr_trtype"`
	AddrTraddr  string `json:"addr_traddr"`
	AddrTrsvcid string `json:"addr_trsvcid,omitempty"`
	AddrAdrfam  string `json:"addr_adrfam,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
}

// ListNVMetPorts lists configured NVMe-oF ports.
func (c *Client) ListNVMetPorts(ctx context.Context, opts ...ListOptions) ([]NVMetPort, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []NVMetPort
	if err := c.call(ctx, "nvmet.port.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing nvmet ports: %w", err)
	}
	return result, nil
}

// GetNVMetPort returns a single NVMe-oF port by ID.
func (c *Client) GetNVMetPort(ctx context.Context, id int) (*NVMetPort, error) {
	var result NVMetPort
	if err := c.call(ctx, "nvmet.port.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting nvmet port %d: %w", id, err)
	}
	return &result, nil
}

// CreateNVMetPort creates a new NVMe-oF port.
func (c *Client) CreateNVMetPort(ctx context.Context, p *CreateNVMetPortParams) (*NVMetPort, error) {
	if p == nil {
		return nil, errors.New("create nvmet port: params required")
	}
	var result NVMetPort
	if err := c.call(ctx, "nvmet.port.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating nvmet port: %w", err)
	}
	return &result, nil
}

// UpdateNVMetPort updates an existing NVMe-oF port.
func (c *Client) UpdateNVMetPort(ctx context.Context, id int, p *CreateNVMetPortParams) (*NVMetPort, error) {
	if p == nil {
		return nil, errors.New("update nvmet port: params required")
	}
	var result NVMetPort
	if err := c.call(ctx, "nvmet.port.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating nvmet port %d: %w", id, err)
	}
	return &result, nil
}

// DeleteNVMetPort deletes an NVMe-oF port.
func (c *Client) DeleteNVMetPort(ctx context.Context, id int) error {
	if err := c.call(ctx, "nvmet.port.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting nvmet port %d: %w", id, err)
	}
	return nil
}

// NVMetPortTransportAddressChoices returns the local addresses available for a new NVMe-oF port to listen on.
func (c *Client) NVMetPortTransportAddressChoices(ctx context.Context, trtype string) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "nvmet.port.transport_address_choices", []any{trtype}, &choices); err != nil {
		return nil, fmt.Errorf("listing nvmet port transport address choices: %w", err)
	}
	return choices, nil
}
