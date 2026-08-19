// Package truenas — replication.go implements the replication.* task domain (replication.create,
// replication.global.*, etc.). Not to be confused with pool.dataset.export_keys_for_replication(),
// a dataset key-export helper unrelated to this domain. Replication endpoints (replication.endpoint.*)
// live in replication_endpoint.go.
package truenas

import (
	"context"
	"errors"
	"fmt"
)

// ReplicationTask represents a configured ZFS replication task.
type ReplicationTask struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	SourceDatasets []string       `json:"source_datasets,omitempty"`
	TargetDataset  string         `json:"target_dataset,omitempty"`
	Direction      string         `json:"direction,omitempty"` // PUSH or PULL
	SSHCredentials int            `json:"ssh_credentials,omitempty"`
	Schedule       map[string]any `json:"schedule,omitempty"`
	Enabled        bool           `json:"enabled,omitempty"`
	ReadOnly       string         `json:"readonly,omitempty"` // SET, REQUIRE, or IGNORE
}

// CreateReplicationParams holds fields for creating or updating a replication task.
type CreateReplicationParams struct {
	Name           string         `json:"name"`
	SourceDatasets []string       `json:"source_datasets"`
	TargetDataset  string         `json:"target_dataset"`
	Direction      string         `json:"direction"`
	SSHCredentials int            `json:"ssh_credentials,omitempty"`
	Schedule       map[string]any `json:"schedule,omitempty"`
	Enabled        bool           `json:"enabled,omitempty"`
	ReadOnly       string         `json:"readonly,omitempty"`
}

// ReplicationGlobalConfig holds global replication settings.
type ReplicationGlobalConfig struct {
	MaxParallelReplicationTasks int `json:"max_parallel_replication_tasks,omitempty"`
}

// UpdateReplicationGlobalParams holds fields for updating global replication settings.
type UpdateReplicationGlobalParams struct {
	MaxParallelReplicationTasks int `json:"max_parallel_replication_tasks,omitempty"`
}

// ListReplicationTasks lists configured replication tasks.
func (c *Client) ListReplicationTasks(ctx context.Context, opts ...ListOptions) ([]ReplicationTask, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []ReplicationTask
	if err := c.call(ctx, "replication.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing replication tasks: %w", err)
	}
	return result, nil
}

// GetReplicationTask returns a single replication task by ID.
func (c *Client) GetReplicationTask(ctx context.Context, id int) (*ReplicationTask, error) {
	var result ReplicationTask
	if err := c.call(ctx, "replication.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting replication task %d: %w", id, err)
	}
	return &result, nil
}

// CreateReplicationTask creates a new replication task.
func (c *Client) CreateReplicationTask(ctx context.Context, p *CreateReplicationParams) (*ReplicationTask, error) {
	if p == nil {
		return nil, errors.New("create replication task: params required")
	}
	var result ReplicationTask
	if err := c.call(ctx, "replication.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating replication task %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateReplicationTask updates an existing replication task.
func (c *Client) UpdateReplicationTask(ctx context.Context, id int, p *CreateReplicationParams) (*ReplicationTask, error) {
	if p == nil {
		return nil, errors.New("update replication task: params required")
	}
	var result ReplicationTask
	if err := c.call(ctx, "replication.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating replication task %d: %w", id, err)
	}
	return &result, nil
}

// DeleteReplicationTask deletes a replication task.
func (c *Client) DeleteReplicationTask(ctx context.Context, id int) error {
	if err := c.call(ctx, "replication.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting replication task %d: %w", id, err)
	}
	return nil
}

// TestReplicationTask runs a dry-run validation of a replication task configuration.
func (c *Client) TestReplicationTask(ctx context.Context, id int) (map[string]any, error) {
	var result map[string]any
	if err := c.call(ctx, "replication.test", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("testing replication task %d: %w", id, err)
	}
	return result, nil
}

// ReplicationSchemas returns the JSON schemas describing valid replication task configurations.
func (c *Client) ReplicationSchemas(ctx context.Context) (map[string]any, error) {
	var schemas map[string]any
	if err := c.call(ctx, "replication.replication_schemas", nil, &schemas); err != nil {
		return nil, fmt.Errorf("getting replication schemas: %w", err)
	}
	return schemas, nil
}

// ReplicationGlobalConfigGet returns global replication settings.
func (c *Client) ReplicationGlobalConfigGet(ctx context.Context) (*ReplicationGlobalConfig, error) {
	var cfg ReplicationGlobalConfig
	if err := c.call(ctx, "replication.global.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting replication global config: %w", err)
	}
	return &cfg, nil
}

// UpdateReplicationGlobal updates global replication settings.
func (c *Client) UpdateReplicationGlobal(ctx context.Context, p *UpdateReplicationGlobalParams) (*ReplicationGlobalConfig, error) {
	if p == nil {
		return nil, errors.New("update replication global config: params required")
	}
	var cfg ReplicationGlobalConfig
	if err := c.call(ctx, "replication.global.update", []any{p}, &cfg); err != nil {
		return nil, fmt.Errorf("updating replication global config: %w", err)
	}
	return &cfg, nil
}
