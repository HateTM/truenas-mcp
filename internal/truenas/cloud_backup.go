// Package truenas — cloud_backup.go implements the restic-based cloud_backup.* domain.
// This is distinct from cloudsync.go, which implements the rclone-based cloudsync.* domain;
// the two do not overlap.
package truenas

import (
	"context"
	"errors"
	"fmt"
)

// CloudBackupTask represents a configured restic-based Cloud Backup task.
type CloudBackupTask struct {
	ID          int            `json:"id"`
	Description string         `json:"description,omitempty"`
	Path        string         `json:"path"`
	Credentials int            `json:"credentials"`
	Repository  string         `json:"repository,omitempty"`
	Schedule    map[string]any `json:"schedule,omitempty"`
	Enabled     bool           `json:"enabled,omitempty"`
	KeepLast    int            `json:"keep_last,omitempty"`
}

// CreateCloudBackupParams holds fields for creating or updating a Cloud Backup task.
type CreateCloudBackupParams struct {
	Description string         `json:"description,omitempty"`
	Path        string         `json:"path"`
	Credentials int            `json:"credentials"`
	Repository  string         `json:"repository,omitempty"`
	Password    string         `json:"password"`
	Schedule    map[string]any `json:"schedule,omitempty"`
	Enabled     bool           `json:"enabled,omitempty"`
	KeepLast    int            `json:"keep_last,omitempty"`
}

// CloudBackupSnapshot represents a single restic snapshot within a Cloud Backup task's repository.
type CloudBackupSnapshot struct {
	ID    string   `json:"id"`
	Time  string   `json:"time,omitempty"`
	Paths []string `json:"paths,omitempty"`
}

// CloudBackupRestoreParams holds fields for restoring a Cloud Backup snapshot.
type CloudBackupRestoreParams struct {
	SnapshotID      string `json:"snapshot_id"`
	Subfolder       string `json:"subfolder,omitempty"`
	DestinationPath string `json:"destination_path"`
}

// ListCloudBackups lists configured Cloud Backup tasks.
func (c *Client) ListCloudBackups(ctx context.Context, opts ...ListOptions) ([]CloudBackupTask, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []CloudBackupTask
	if err := c.call(ctx, "cloud_backup.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing cloud backup tasks: %w", err)
	}
	return result, nil
}

// GetCloudBackup returns a single Cloud Backup task by ID.
func (c *Client) GetCloudBackup(ctx context.Context, id int) (*CloudBackupTask, error) {
	var result CloudBackupTask
	if err := c.call(ctx, "cloud_backup.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting cloud backup task %d: %w", id, err)
	}
	return &result, nil
}

// CreateCloudBackup creates a new Cloud Backup task.
func (c *Client) CreateCloudBackup(ctx context.Context, p *CreateCloudBackupParams) (*CloudBackupTask, error) {
	if p == nil {
		return nil, errors.New("create cloud backup task: params required")
	}
	var result CloudBackupTask
	if err := c.call(ctx, "cloud_backup.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating cloud backup task: %w", err)
	}
	return &result, nil
}

// UpdateCloudBackup updates an existing Cloud Backup task.
func (c *Client) UpdateCloudBackup(ctx context.Context, id int, p *CreateCloudBackupParams) (*CloudBackupTask, error) {
	if p == nil {
		return nil, errors.New("update cloud backup task: params required")
	}
	var result CloudBackupTask
	if err := c.call(ctx, "cloud_backup.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating cloud backup task %d: %w", id, err)
	}
	return &result, nil
}

// DeleteCloudBackup deletes a Cloud Backup task.
func (c *Client) DeleteCloudBackup(ctx context.Context, id int) error {
	if err := c.call(ctx, "cloud_backup.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting cloud backup task %d: %w", id, err)
	}
	return nil
}

// AbortCloudBackup aborts a running Cloud Backup job.
func (c *Client) AbortCloudBackup(ctx context.Context, id int) error {
	if err := c.call(ctx, "cloud_backup.abort", []any{id}, nil); err != nil {
		return fmt.Errorf("aborting cloud backup task %d: %w", id, err)
	}
	return nil
}

// SyncCloudBackup runs a Cloud Backup task and returns the async job ID.
func (c *Client) SyncCloudBackup(ctx context.Context, id int) (int, error) {
	var jobID int
	if err := c.call(ctx, "cloud_backup.sync", []any{id}, &jobID); err != nil {
		return 0, fmt.Errorf("syncing cloud backup task %d: %w", id, err)
	}
	return jobID, nil
}

// RestoreCloudBackup restores files from a Cloud Backup snapshot and returns the async job ID.
func (c *Client) RestoreCloudBackup(ctx context.Context, id int, p *CloudBackupRestoreParams) (int, error) {
	if p == nil {
		return 0, errors.New("restore cloud backup: params required")
	}
	var jobID int
	if err := c.call(ctx, "cloud_backup.restore", []any{id, p}, &jobID); err != nil {
		return 0, fmt.Errorf("restoring cloud backup task %d snapshot %q: %w", id, p.SnapshotID, err)
	}
	return jobID, nil
}

// DeleteCloudBackupSnapshot deletes a single snapshot from a Cloud Backup task's repository.
func (c *Client) DeleteCloudBackupSnapshot(ctx context.Context, id int, snapshotID string) error {
	if snapshotID == "" {
		return errors.New("delete cloud backup snapshot: snapshot_id must not be empty")
	}
	if err := c.call(ctx, "cloud_backup.delete_snapshot", []any{id, snapshotID}, nil); err != nil {
		return fmt.Errorf("deleting cloud backup task %d snapshot %q: %w", id, snapshotID, err)
	}
	return nil
}

// ListCloudBackupSnapshots lists the snapshots stored in a Cloud Backup task's repository.
func (c *Client) ListCloudBackupSnapshots(ctx context.Context, id int) ([]CloudBackupSnapshot, error) {
	var snapshots []CloudBackupSnapshot
	if err := c.call(ctx, "cloud_backup.list_snapshots", []any{id}, &snapshots); err != nil {
		return nil, fmt.Errorf("listing cloud backup task %d snapshots: %w", id, err)
	}
	return snapshots, nil
}

// ListCloudBackupSnapshotDirectory lists the contents of a directory within a Cloud Backup snapshot.
func (c *Client) ListCloudBackupSnapshotDirectory(ctx context.Context, id int, snapshotID, path string) ([]DirEntry, error) {
	if snapshotID == "" {
		return nil, errors.New("list cloud backup snapshot directory: snapshot_id must not be empty")
	}
	var entries []DirEntry
	if err := c.call(ctx, "cloud_backup.list_snapshot_directory", []any{id, snapshotID, path}, &entries); err != nil {
		return nil, fmt.Errorf("listing cloud backup task %d snapshot %q directory %q: %w", id, snapshotID, path, err)
	}
	return entries, nil
}

// CloudBackupTransferSettingChoices returns the valid transfer-setting values for a Cloud Backup task.
func (c *Client) CloudBackupTransferSettingChoices(ctx context.Context) ([]string, error) {
	var choices []string
	if err := c.call(ctx, "cloud_backup.transfer_setting_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing cloud backup transfer setting choices: %w", err)
	}
	return choices, nil
}
