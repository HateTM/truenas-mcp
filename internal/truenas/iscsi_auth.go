package truenas

import (
	"context"
	"errors"
	"fmt"
)

// ISCSIAuth represents a CHAP authentication credential for iSCSI targets.
type ISCSIAuth struct {
	ID         int    `json:"id"`
	Tag        int    `json:"tag"`
	User       string `json:"user"`
	Secret     string `json:"secret,omitempty"`
	PeerUser   string `json:"peeruser,omitempty"`
	PeerSecret string `json:"peersecret,omitempty"`
}

// CreateISCSIAuthParams holds fields for creating or updating an iSCSI CHAP credential.
type CreateISCSIAuthParams struct {
	Tag        int    `json:"tag"`
	User       string `json:"user"`
	Secret     string `json:"secret"`
	PeerUser   string `json:"peeruser,omitempty"`
	PeerSecret string `json:"peersecret,omitempty"`
}

// ISCSIGlobalConfig represents the global iSCSI service configuration.
type ISCSIGlobalConfig struct {
	Basename           string   `json:"basename"`
	ISNSServers        []string `json:"isns_servers,omitempty"`
	ListenPort         int      `json:"listen_port,omitempty"`
	ALUA               bool     `json:"alua,omitempty"`
	PoolAvailThreshold int      `json:"pool_avail_threshold,omitempty"`
}

// UpdateISCSIGlobalParams holds fields for updating the global iSCSI service configuration.
type UpdateISCSIGlobalParams struct {
	Basename           string   `json:"basename,omitempty"`
	ISNSServers        []string `json:"isns_servers,omitempty"`
	ListenPort         int      `json:"listen_port,omitempty"`
	ALUA               bool     `json:"alua,omitempty"`
	PoolAvailThreshold int      `json:"pool_avail_threshold,omitempty"`
}

// ISCSISession represents an active iSCSI session.
type ISCSISession struct {
	ID        int    `json:"id"`
	Initiator string `json:"initiator"`
	Target    string `json:"target"`
}

// ListISCSIAuth lists configured iSCSI CHAP credentials.
func (c *Client) ListISCSIAuth(ctx context.Context, opts ...ListOptions) ([]ISCSIAuth, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ISCSIAuth
	if err := c.call(ctx, "iscsi.auth.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing iscsi auth credentials: %w", err)
	}
	return result, nil
}

// GetISCSIAuth returns a single iSCSI CHAP credential by ID.
func (c *Client) GetISCSIAuth(ctx context.Context, id int) (*ISCSIAuth, error) {
	var result ISCSIAuth
	if err := c.call(ctx, "iscsi.auth.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting iscsi auth credential %d: %w", id, err)
	}
	return &result, nil
}

// CreateISCSIAuth creates a new iSCSI CHAP credential.
func (c *Client) CreateISCSIAuth(ctx context.Context, p *CreateISCSIAuthParams) (*ISCSIAuth, error) {
	if p == nil {
		return nil, errors.New("create iscsi auth credential: params required")
	}
	var result ISCSIAuth
	if err := c.call(ctx, "iscsi.auth.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating iscsi auth credential: %w", err)
	}
	return &result, nil
}

// UpdateISCSIAuth updates an existing iSCSI CHAP credential.
func (c *Client) UpdateISCSIAuth(ctx context.Context, id int, p *CreateISCSIAuthParams) (*ISCSIAuth, error) {
	if p == nil {
		return nil, errors.New("update iscsi auth credential: params required")
	}
	var result ISCSIAuth
	if err := c.call(ctx, "iscsi.auth.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating iscsi auth credential %d: %w", id, err)
	}
	return &result, nil
}

// DeleteISCSIAuth deletes an iSCSI CHAP credential.
func (c *Client) DeleteISCSIAuth(ctx context.Context, id int) error {
	if err := c.call(ctx, "iscsi.auth.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting iscsi auth credential %d: %w", id, err)
	}
	return nil
}

// ISCSIGlobalALUAEnabled reports whether ALUA is enabled for the iSCSI service.
func (c *Client) ISCSIGlobalALUAEnabled(ctx context.Context) (bool, error) {
	var enabled bool
	if err := c.call(ctx, "iscsi.global.alua_enabled", nil, &enabled); err != nil {
		return false, fmt.Errorf("checking iscsi global alua_enabled: %w", err)
	}
	return enabled, nil
}

// ISCSIGlobalClientCount returns the number of connected iSCSI clients.
func (c *Client) ISCSIGlobalClientCount(ctx context.Context) (int, error) {
	var count int
	if err := c.call(ctx, "iscsi.global.client_count", nil, &count); err != nil {
		return 0, fmt.Errorf("getting iscsi global client_count: %w", err)
	}
	return count, nil
}

// ISCSIGlobalConfigGet returns the global iSCSI service configuration.
func (c *Client) ISCSIGlobalConfigGet(ctx context.Context) (*ISCSIGlobalConfig, error) {
	var cfg ISCSIGlobalConfig
	if err := c.call(ctx, "iscsi.global.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting iscsi global config: %w", err)
	}
	return &cfg, nil
}

// ISCSIGlobalISEREnabled reports whether iSER is enabled for the iSCSI service.
func (c *Client) ISCSIGlobalISEREnabled(ctx context.Context) (bool, error) {
	var enabled bool
	if err := c.call(ctx, "iscsi.global.iser_enabled", nil, &enabled); err != nil {
		return false, fmt.Errorf("checking iscsi global iser_enabled: %w", err)
	}
	return enabled, nil
}

// ISCSIGlobalSessions returns all active iSCSI sessions.
func (c *Client) ISCSIGlobalSessions(ctx context.Context) ([]ISCSISession, error) {
	var sessions []ISCSISession
	if err := c.call(ctx, "iscsi.global.sessions", nil, &sessions); err != nil {
		return nil, fmt.Errorf("listing iscsi global sessions: %w", err)
	}
	return sessions, nil
}

// UpdateISCSIGlobal updates the global iSCSI service configuration.
func (c *Client) UpdateISCSIGlobal(ctx context.Context, p *UpdateISCSIGlobalParams) (*ISCSIGlobalConfig, error) {
	if p == nil {
		return nil, errors.New("update iscsi global config: params required")
	}
	var cfg ISCSIGlobalConfig
	if err := c.call(ctx, "iscsi.global.update", []any{p}, &cfg); err != nil {
		return nil, fmt.Errorf("updating iscsi global config: %w", err)
	}
	return &cfg, nil
}
