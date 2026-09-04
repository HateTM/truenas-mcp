package truenas

import (
	"context"
	"errors"
	"fmt"
)

// ReplicationEndpoint represents a remote TrueNAS system replication tasks can target.
type ReplicationEndpoint struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URI   string `json:"uri,omitempty"`
	Token string `json:"token,omitempty"`
}

// CreateReplicationEndpointParams holds fields for creating or updating a replication endpoint.
type CreateReplicationEndpointParams struct {
	Name  string `json:"name"`
	URI   string `json:"uri"`
	Token string `json:"token,omitempty"`
}

// ListReplicationEndpoints lists configured remote replication endpoints.
func (c *Client) ListReplicationEndpoints(ctx context.Context, opts ...ListOptions) ([]ReplicationEndpoint, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ReplicationEndpoint
	if err := c.call(ctx, "replication.endpoint.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing replication endpoints: %w", err)
	}
	return result, nil
}

// GetReplicationEndpoint returns a single replication endpoint by ID.
func (c *Client) GetReplicationEndpoint(ctx context.Context, id int) (*ReplicationEndpoint, error) {
	var result ReplicationEndpoint
	if err := c.call(ctx, "replication.endpoint.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting replication endpoint %d: %w", id, err)
	}
	return &result, nil
}

// CreateReplicationEndpoint creates a new remote replication endpoint.
func (c *Client) CreateReplicationEndpoint(ctx context.Context, p *CreateReplicationEndpointParams) (*ReplicationEndpoint, error) {
	if p == nil {
		return nil, errors.New("create replication endpoint: params required")
	}
	var result ReplicationEndpoint
	if err := c.call(ctx, "replication.endpoint.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating replication endpoint %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateReplicationEndpoint updates an existing remote replication endpoint.
func (c *Client) UpdateReplicationEndpoint(ctx context.Context, id int, p *CreateReplicationEndpointParams) (*ReplicationEndpoint, error) {
	if p == nil {
		return nil, errors.New("update replication endpoint: params required")
	}
	var result ReplicationEndpoint
	if err := c.call(ctx, "replication.endpoint.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating replication endpoint %d: %w", id, err)
	}
	return &result, nil
}

// DeleteReplicationEndpoint deletes a remote replication endpoint.
func (c *Client) DeleteReplicationEndpoint(ctx context.Context, id int) error {
	if err := c.call(ctx, "replication.endpoint.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting replication endpoint %d: %w", id, err)
	}
	return nil
}
