package truenas

import (
	"context"
	"errors"
	"fmt"
)

// ISCSIPortalListen is a single listen address/port pair for an iSCSI portal.
type ISCSIPortalListen struct {
	IP   string `json:"ip"`
	Port int    `json:"port,omitempty"`
}

// ISCSIPortal represents an iSCSI portal — the set of IP/port combinations
// an iSCSI target listens on.
type ISCSIPortal struct {
	ID                  int                 `json:"id"`
	Listen              []ISCSIPortalListen `json:"listen"`
	Comment             string              `json:"comment,omitempty"`
	DiscoveryAuthMethod string              `json:"discovery_authmethod,omitempty"`
	DiscoveryAuthGroup  int                 `json:"discovery_authgroup,omitempty"`
}

// CreateISCSIPortalParams holds fields for creating or updating an iSCSI portal.
type CreateISCSIPortalParams struct {
	Listen              []ISCSIPortalListen `json:"listen"`
	Comment             string              `json:"comment,omitempty"`
	DiscoveryAuthMethod string              `json:"discovery_authmethod,omitempty"`
	DiscoveryAuthGroup  int                 `json:"discovery_authgroup,omitempty"`
}

// ListISCSIPortals lists configured iSCSI portals.
func (c *Client) ListISCSIPortals(ctx context.Context, opts ...ListOptions) ([]ISCSIPortal, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ISCSIPortal
	if err := c.call(ctx, "iscsi.portal.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing iscsi portals: %w", err)
	}
	return result, nil
}

// GetISCSIPortal returns a single iSCSI portal by ID.
func (c *Client) GetISCSIPortal(ctx context.Context, id int) (*ISCSIPortal, error) {
	var result ISCSIPortal
	if err := c.call(ctx, "iscsi.portal.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting iscsi portal %d: %w", id, err)
	}
	return &result, nil
}

// CreateISCSIPortal creates a new iSCSI portal.
func (c *Client) CreateISCSIPortal(ctx context.Context, p *CreateISCSIPortalParams) (*ISCSIPortal, error) {
	if p == nil {
		return nil, errors.New("create iscsi portal: params required")
	}
	var result ISCSIPortal
	if err := c.call(ctx, "iscsi.portal.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating iscsi portal: %w", err)
	}
	return &result, nil
}

// UpdateISCSIPortal updates an existing iSCSI portal.
func (c *Client) UpdateISCSIPortal(ctx context.Context, id int, p *CreateISCSIPortalParams) (*ISCSIPortal, error) {
	if p == nil {
		return nil, errors.New("update iscsi portal: params required")
	}
	var result ISCSIPortal
	if err := c.call(ctx, "iscsi.portal.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating iscsi portal %d: %w", id, err)
	}
	return &result, nil
}

// DeleteISCSIPortal deletes an iSCSI portal.
func (c *Client) DeleteISCSIPortal(ctx context.Context, id int) error {
	if err := c.call(ctx, "iscsi.portal.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting iscsi portal %d: %w", id, err)
	}
	return nil
}

// ISCSIPortalListenIPChoices returns the IP addresses available for an iSCSI portal to listen on.
func (c *Client) ISCSIPortalListenIPChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "iscsi.portal.listen_ip_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing iscsi portal listen IP choices: %w", err)
	}
	return choices, nil
}
