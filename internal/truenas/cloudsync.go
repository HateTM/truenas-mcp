package truenas

import (
	"context"
	"fmt"
	"maps"
)

// CloudProvider describes a backend supported by TrueNAS Cloud Sync Tasks
type CloudProvider struct {
	Name string `json:"name"`
}

// CloudCredential is a stored Cloud Sync credential
type CloudCredential struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Provider map[string]any `json:"provider"`
}

// CloudSyncTask is a configured Cloud Sync Task
type CloudSyncTask struct {
	ID           int            `json:"id"`
	Description  string         `json:"description"`
	Path         string         `json:"path"`
	Credentials  map[string]any `json:"credentials"`
	Direction    string         `json:"direction"`
	TransferMode string         `json:"transfer_mode"`
	Attributes   map[string]any `json:"attributes,omitempty"`
	Enabled      bool           `json:"enabled"`
}

// CreateCloudCredentialParams holds fields for creating a Cloud Sync credential.
type CreateCloudCredentialParams struct {
	Name       string
	Provider   string
	Attributes map[string]any
}

// CreateCloudSyncTaskParams holds fields for creating a Cloud Sync Task.
type CreateCloudSyncTaskParams struct {
	Description  string
	Path         string
	Credentials  int
	Direction    string
	TransferMode string
	Attributes   map[string]any
	Extra        map[string]any
	Enabled      bool
}

// ListCloudProviders lists supported providers.
func (c *Client) ListCloudProviders(ctx context.Context) ([]CloudProvider, error) {
	var providers []CloudProvider
	if err := c.call(ctx, "cloudsync.providers", nil, &providers); err != nil {
		return nil, fmt.Errorf("listing cloud providers: %w", err)
	}
	return providers, nil
}

// ListCloudCredentials returns all stored Cloud Sync credentials.
func (c *Client) ListCloudCredentials(ctx context.Context, opts ...ListOptions) ([]CloudCredential, error) {
	var opt ListOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	var creds []CloudCredential
	if err := c.call(ctx, "cloudsync.credentials.query", []any{nil, map[string]any{"limit": opt.Limit, "offset": opt.Offset}}, &creds); err != nil {
		return nil, fmt.Errorf("listing cloud credentials: %w", err)
	}
	return creds, nil
}

// CreateCloudCredential creates a new cloud credential.
func (c *Client) CreateCloudCredential(ctx context.Context, p *CreateCloudCredentialParams) (*CloudCredential, error) {
	if p == nil {
		return nil, fmt.Errorf("creating cloud credential: params must not be nil")
	}
	params := map[string]any{
		"name": p.Name,
		"config": map[string]any{
			"type": p.Provider,
		},
	}
	// Flatten attributes directly into config
	maps.Copy(params["config"].(map[string]any), p.Attributes)

	var cred CloudCredential
	if err := c.call(ctx, "cloudsync.credentials.create", []any{params}, &cred); err != nil {
		return nil, fmt.Errorf("creating cloud credential: %w", err)
	}
	return &cred, nil
}

// ListCloudSyncTasks lists configured tasks.
func (c *Client) ListCloudSyncTasks(ctx context.Context, opts ...ListOptions) ([]CloudSyncTask, error) {
	var opt ListOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	var tasks []CloudSyncTask
	if err := c.call(ctx, "cloudsync.query", []any{nil, map[string]any{"limit": opt.Limit, "offset": opt.Offset}}, &tasks); err != nil {
		return nil, fmt.Errorf("listing cloud sync tasks: %w", err)
	}
	return tasks, nil
}

// CreateCloudSyncTask creates a task.
func (c *Client) CreateCloudSyncTask(ctx context.Context, p *CreateCloudSyncTaskParams) (*CloudSyncTask, error) {
	if p == nil {
		return nil, fmt.Errorf("creating cloud sync task: params must not be nil")
	}
	params := map[string]any{
		"description":   p.Description,
		"path":          p.Path,
		"credentials":   p.Credentials,
		"direction":     p.Direction,
		"transfer_mode": p.TransferMode,
		"attributes":    p.Attributes,
		"enabled":       p.Enabled,
	}
	maps.Copy(params, p.Extra)
	var task CloudSyncTask
	if err := c.call(ctx, "cloudsync.create", []any{params}, &task); err != nil {
		return nil, fmt.Errorf("creating cloud sync task: %w", err)
	}
	return &task, nil
}

// RunCloudSyncTask runs a task.
func (c *Client) RunCloudSyncTask(ctx context.Context, id int, dryRun bool) (int, error) {
	var jobID int
	if err := c.call(ctx, "cloudsync.sync", []any{id, map[string]any{"dry_run": dryRun}}, &jobID); err != nil {
		return 0, fmt.Errorf("running cloud sync task %d: %w", id, err)
	}
	return jobID, nil
}

// AbortCloudSyncTask aborts a task job.
func (c *Client) AbortCloudSyncTask(ctx context.Context, id int) error {
	if err := c.call(ctx, "cloudsync.abort", []any{id}, nil); err != nil {
		return fmt.Errorf("aborting cloud sync task %d: %w", id, err)
	}
	return nil
}

// DeleteCloudSyncTask deletes a task.
func (c *Client) DeleteCloudSyncTask(ctx context.Context, id int) error {
	if err := c.call(ctx, "cloudsync.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting cloud sync task %d: %w", id, err)
	}
	return nil
}

// DeleteCloudCredential deletes a credential.
func (c *Client) DeleteCloudCredential(ctx context.Context, id int) error {
	if err := c.call(ctx, "cloudsync.credentials.destroy", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting cloud sync credential %d: %w", id, err)
	}
	return nil
}
