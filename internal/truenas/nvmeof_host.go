package truenas

import (
	"context"
	"errors"
	"fmt"
)

// NVMetGlobalConfig represents the global NVMe-oF (nvmet) service configuration.
type NVMetGlobalConfig struct {
	ANA  bool `json:"ana,omitempty"`
	RDMA bool `json:"rdma,omitempty"`
}

// UpdateNVMetGlobalParams holds fields for updating the global NVMe-oF service configuration.
type UpdateNVMetGlobalParams struct {
	ANA  bool `json:"ana,omitempty"`
	RDMA bool `json:"rdma,omitempty"`
}

// NVMetHost represents an NVMe-oF host allowed to connect to subsystems (via nvmet.host_subsys.*).
type NVMetHost struct {
	ID            int    `json:"id"`
	HostNQN       string `json:"hostnqn"`
	DHCHAPKey     string `json:"dhchap_key,omitempty"`
	DHCHAPDHGroup string `json:"dhchap_dhgroup,omitempty"`
	DHCHAPHash    string `json:"dhchap_hash,omitempty"`
}

// CreateNVMetHostParams holds fields for creating or updating an NVMe-oF host.
type CreateNVMetHostParams struct {
	HostNQN       string `json:"hostnqn"`
	DHCHAPKey     string `json:"dhchap_key,omitempty"`
	DHCHAPDHGroup string `json:"dhchap_dhgroup,omitempty"`
	DHCHAPHash    string `json:"dhchap_hash,omitempty"`
}

// NVMetGlobalConfigGet returns the global NVMe-oF service configuration.
func (c *Client) NVMetGlobalConfigGet(ctx context.Context) (*NVMetGlobalConfig, error) {
	var cfg NVMetGlobalConfig
	if err := c.call(ctx, "nvmet.global.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting nvmet global config: %w", err)
	}
	return &cfg, nil
}

// UpdateNVMetGlobal updates the global NVMe-oF service configuration.
func (c *Client) UpdateNVMetGlobal(ctx context.Context, p *UpdateNVMetGlobalParams) (*NVMetGlobalConfig, error) {
	if p == nil {
		return nil, errors.New("update nvmet global config: params required")
	}
	var cfg NVMetGlobalConfig
	if err := c.call(ctx, "nvmet.global.update", []any{p}, &cfg); err != nil {
		return nil, fmt.Errorf("updating nvmet global config: %w", err)
	}
	return &cfg, nil
}

// ListNVMetHosts lists configured NVMe-oF hosts.
func (c *Client) ListNVMetHosts(ctx context.Context, opts ...ListOptions) ([]NVMetHost, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []NVMetHost
	if err := c.call(ctx, "nvmet.host.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing nvmet hosts: %w", err)
	}
	return result, nil
}

// GetNVMetHost returns a single NVMe-oF host by ID.
func (c *Client) GetNVMetHost(ctx context.Context, id int) (*NVMetHost, error) {
	var result NVMetHost
	if err := c.call(ctx, "nvmet.host.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting nvmet host %d: %w", id, err)
	}
	return &result, nil
}

// CreateNVMetHost creates a new NVMe-oF host.
func (c *Client) CreateNVMetHost(ctx context.Context, p *CreateNVMetHostParams) (*NVMetHost, error) {
	if p == nil {
		return nil, errors.New("create nvmet host: params required")
	}
	var result NVMetHost
	if err := c.call(ctx, "nvmet.host.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating nvmet host: %w", err)
	}
	return &result, nil
}

// UpdateNVMetHost updates an existing NVMe-oF host.
func (c *Client) UpdateNVMetHost(ctx context.Context, id int, p *CreateNVMetHostParams) (*NVMetHost, error) {
	if p == nil {
		return nil, errors.New("update nvmet host: params required")
	}
	var result NVMetHost
	if err := c.call(ctx, "nvmet.host.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating nvmet host %d: %w", id, err)
	}
	return &result, nil
}

// DeleteNVMetHost deletes an NVMe-oF host.
func (c *Client) DeleteNVMetHost(ctx context.Context, id int) error {
	if err := c.call(ctx, "nvmet.host.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting nvmet host %d: %w", id, err)
	}
	return nil
}

// NVMetHostDHCHAPDHGroupChoices returns the DH groups available for host DH-HMAC-CHAP authentication.
func (c *Client) NVMetHostDHCHAPDHGroupChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "nvmet.host.dhchap_dhgroup_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing nvmet host dhchap dhgroup choices: %w", err)
	}
	return choices, nil
}

// NVMetHostDHCHAPHashChoices returns the hash algorithms available for host DH-HMAC-CHAP authentication.
func (c *Client) NVMetHostDHCHAPHashChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "nvmet.host.dhchap_hash_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing nvmet host dhchap hash choices: %w", err)
	}
	return choices, nil
}

// NVMetHostGenerateKey generates a new random DH-HMAC-CHAP key suitable for use as a host's dhchap_key.
func (c *Client) NVMetHostGenerateKey(ctx context.Context) (string, error) {
	var key string
	if err := c.call(ctx, "nvmet.host.generate_key", nil, &key); err != nil {
		return "", fmt.Errorf("generating nvmet host dhchap key: %w", err)
	}
	return key, nil
}
