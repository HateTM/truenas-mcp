package truenas

import (
	"context"
	"errors"
	"fmt"
)

// DirectoryServicesConfig represents the directory service (AD/LDAP) configuration.
type DirectoryServicesConfig struct {
	Service       string         `json:"service,omitempty"` // e.g. ACTIVEDIRECTORY, LDAP
	Enable        bool           `json:"enable,omitempty"`
	Configuration map[string]any `json:"configuration,omitempty"`
}

// UpdateDirectoryServicesParams holds fields for updating the directory service configuration.
type UpdateDirectoryServicesParams struct {
	Service       string         `json:"service,omitempty"`
	Enable        bool           `json:"enable,omitempty"`
	Configuration map[string]any `json:"configuration,omitempty"`
}

// DirectoryServicesStatus reports the current connection status of the directory service.
type DirectoryServicesStatus struct {
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
}

// DirectoryServicesCacheRefresh refreshes the cached identities (users/groups) from the directory service.
func (c *Client) DirectoryServicesCacheRefresh(ctx context.Context) error {
	if err := c.call(ctx, "directoryservices.cache_refresh", nil, nil); err != nil {
		return fmt.Errorf("refreshing directory services cache: %w", err)
	}
	return nil
}

// DirectoryServicesCertificateChoices returns the certificates available for the directory service connection.
func (c *Client) DirectoryServicesCertificateChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "directoryservices.certificate_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing directory services certificate choices: %w", err)
	}
	return choices, nil
}

// DirectoryServicesConfigGet returns the directory service (AD/LDAP) configuration.
func (c *Client) DirectoryServicesConfigGet(ctx context.Context) (*DirectoryServicesConfig, error) {
	var cfg DirectoryServicesConfig
	if err := c.call(ctx, "directoryservices.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting directory services config: %w", err)
	}
	return &cfg, nil
}

// DirectoryServicesLeave leaves the currently joined domain, using admin credentials passed in params.
func (c *Client) DirectoryServicesLeave(ctx context.Context, params map[string]any) error {
	if err := c.call(ctx, "directoryservices.leave", []any{params}, nil); err != nil {
		return fmt.Errorf("leaving directory service domain: %w", err)
	}
	return nil
}

// DirectoryServicesStatusGet returns the current connection status of the directory service.
func (c *Client) DirectoryServicesStatusGet(ctx context.Context) (*DirectoryServicesStatus, error) {
	var status DirectoryServicesStatus
	if err := c.call(ctx, "directoryservices.status", nil, &status); err != nil {
		return nil, fmt.Errorf("getting directory services status: %w", err)
	}
	return &status, nil
}

// DirectoryServicesSyncKeytab synchronizes the Kerberos keytab with the directory service.
func (c *Client) DirectoryServicesSyncKeytab(ctx context.Context) error {
	if err := c.call(ctx, "directoryservices.sync_keytab", nil, nil); err != nil {
		return fmt.Errorf("syncing directory services keytab: %w", err)
	}
	return nil
}

// UpdateDirectoryServices updates the directory service (AD/LDAP) configuration.
func (c *Client) UpdateDirectoryServices(ctx context.Context, p *UpdateDirectoryServicesParams) (*DirectoryServicesConfig, error) {
	if p == nil {
		return nil, errors.New("update directory services config: params required")
	}
	var cfg DirectoryServicesConfig
	if err := c.call(ctx, "directoryservices.update", []any{p}, &cfg); err != nil {
		return nil, fmt.Errorf("updating directory services config: %w", err)
	}
	return &cfg, nil
}
