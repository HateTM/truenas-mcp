package truenas

import (
	"context"
	"errors"
	"fmt"
)

// NFSConfig represents the global NFS service configuration (distinct from
// individual NFS exports managed by sharing.nfs.*).
type NFSConfig struct {
	Servers      int      `json:"servers,omitempty"`
	BindIP       []string `json:"bindip,omitempty"`
	AllowNonroot bool     `json:"allow_nonroot,omitempty"`
	V4           bool     `json:"v4,omitempty"`
}

// UpdateNFSConfigParams holds fields for updating the global NFS service configuration.
type UpdateNFSConfigParams struct {
	Servers      int      `json:"servers,omitempty"`
	BindIP       []string `json:"bindip,omitempty"`
	AllowNonroot bool     `json:"allow_nonroot,omitempty"`
	V4           bool     `json:"v4,omitempty"`
}

// NFSClient represents a client currently connected to the NFS service.
type NFSClient struct {
	ID     int    `json:"id,omitempty"`
	Client string `json:"client,omitempty"`
}

// NFSBindIPChoices returns the IP addresses available for the NFS service to bind to.
func (c *Client) NFSBindIPChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "nfs.bindip_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing nfs bindip choices: %w", err)
	}
	return choices, nil
}

// NFSClientCount returns the number of clients currently connected to the NFS service.
func (c *Client) NFSClientCount(ctx context.Context) (int, error) {
	var count int
	if err := c.call(ctx, "nfs.client_count", nil, &count); err != nil {
		return 0, fmt.Errorf("getting nfs client_count: %w", err)
	}
	return count, nil
}

// NFSConfigGet returns the global NFS service configuration.
func (c *Client) NFSConfigGet(ctx context.Context) (*NFSConfig, error) {
	var cfg NFSConfig
	if err := c.call(ctx, "nfs.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting nfs config: %w", err)
	}
	return &cfg, nil
}

// GetNFS3Clients returns clients currently connected via NFSv3.
func (c *Client) GetNFS3Clients(ctx context.Context) ([]NFSClient, error) {
	var clients []NFSClient
	if err := c.call(ctx, "nfs.get_nfs3_clients", nil, &clients); err != nil {
		return nil, fmt.Errorf("listing nfsv3 clients: %w", err)
	}
	return clients, nil
}

// GetNFS4Clients returns clients currently connected via NFSv4.
func (c *Client) GetNFS4Clients(ctx context.Context) ([]NFSClient, error) {
	var clients []NFSClient
	if err := c.call(ctx, "nfs.get_nfs4_clients", nil, &clients); err != nil {
		return nil, fmt.Errorf("listing nfsv4 clients: %w", err)
	}
	return clients, nil
}

// UpdateNFSConfig updates the global NFS service configuration.
func (c *Client) UpdateNFSConfig(ctx context.Context, p *UpdateNFSConfigParams) (*NFSConfig, error) {
	if p == nil {
		return nil, errors.New("update nfs config: params required")
	}
	var cfg NFSConfig
	if err := c.call(ctx, "nfs.update", []any{p}, &cfg); err != nil {
		return nil, fmt.Errorf("updating nfs config: %w", err)
	}
	return &cfg, nil
}
