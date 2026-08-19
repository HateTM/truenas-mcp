package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCloudBackupTools registers cloud_backup.* (restic-based) MCP tools onto the server.
// Distinct from registerCloudSyncTools (rclone-based cloudsync.*), which this does not overlap with.
func registerCloudBackupTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listCloudBackupsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of tasks to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of tasks to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_list",
		Description: "List configured Cloud Backup (restic) tasks.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listCloudBackupsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListCloudBackups(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_list: %w", err))
		}
		return jsonResult(result)
	})

	type getCloudBackupInput struct {
		ID int `json:"id" jsonschema:"Numeric Cloud Backup task ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_get",
		Description: "Get a single Cloud Backup task by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getCloudBackupInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_get: id must be a positive integer"))
		}
		result, err := client.GetCloudBackup(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_get: %w", err))
		}
		return jsonResult(result)
	})

	type cloudBackupInput struct {
		Description string         `json:"description,omitempty" jsonschema:"Name of the task to display in the TrueNAS UI"`
		Path        string         `json:"path"                  jsonschema:"Local path to back up, beginning with /mnt"`
		Credentials int            `json:"credentials"            jsonschema:"ID of a Cloud Sync credential (from create_cloud_credential) for the restic repository backend"`
		Repository  string         `json:"repository,omitempty"   jsonschema:"restic repository URI, e.g. s3:https://bucket.example.com/backups"`
		Password    string         `json:"password"                jsonschema:"restic repository encryption password"`
		Schedule    map[string]any `json:"schedule,omitempty"      jsonschema:"Cron schedule fields, e.g. {\"minute\": \"0\", \"hour\": \"3\"}"`
		Enabled     bool           `json:"enabled,omitempty"       jsonschema:"Whether the task becomes active on its own schedule immediately. Defaults to false — set true only if it should run unattended."`
		KeepLast    int            `json:"keep_last,omitempty"     jsonschema:"Number of most recent snapshots to retain on prune"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_create",
		Description: "Create a Cloud Backup (restic) task. Created disabled (enabled=false) by default. Does not run the task itself — use cloud_backup_sync for that.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p cloudBackupInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("cloud_backup_create: path must not be empty"))
		}
		if p.Credentials <= 0 {
			return errorResult(errors.New("cloud_backup_create: credentials must be a positive integer"))
		}
		if p.Password == "" {
			return errorResult(errors.New("cloud_backup_create: password must not be empty"))
		}
		result, err := client.CreateCloudBackup(ctx, &truenas.CreateCloudBackupParams{
			Description: p.Description, Path: p.Path, Credentials: p.Credentials, Repository: p.Repository,
			Password: p.Password, Schedule: p.Schedule, Enabled: p.Enabled, KeepLast: p.KeepLast,
		})
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateCloudBackupInput struct {
		ID          int            `json:"id"                     jsonschema:"Numeric Cloud Backup task ID"`
		Description string         `json:"description,omitempty" jsonschema:"Name of the task to display in the TrueNAS UI"`
		Path        string         `json:"path"                  jsonschema:"Local path to back up"`
		Credentials int            `json:"credentials"            jsonschema:"ID of a Cloud Sync credential for the restic repository backend"`
		Repository  string         `json:"repository,omitempty"   jsonschema:"restic repository URI"`
		Password    string         `json:"password"                jsonschema:"restic repository encryption password"`
		Schedule    map[string]any `json:"schedule,omitempty"      jsonschema:"Cron schedule fields"`
		Enabled     bool           `json:"enabled,omitempty"       jsonschema:"Whether the task is active"`
		KeepLast    int            `json:"keep_last,omitempty"     jsonschema:"Number of most recent snapshots to retain on prune"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_update",
		Description: "Update an existing Cloud Backup task.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateCloudBackupInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_update: id must be a positive integer"))
		}
		if p.Path == "" || p.Credentials <= 0 || p.Password == "" {
			return errorResult(errors.New("cloud_backup_update: path, credentials, and password are required"))
		}
		result, err := client.UpdateCloudBackup(ctx, p.ID, &truenas.CreateCloudBackupParams{
			Description: p.Description, Path: p.Path, Credentials: p.Credentials, Repository: p.Repository,
			Password: p.Password, Schedule: p.Schedule, Enabled: p.Enabled, KeepLast: p.KeepLast,
		})
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_update: %w", err))
		}
		return jsonResult(result)
	})

	// cloud_backup_delete and cloud_backup_delete_snapshot are destructive and
	// live in destructive_cloud_backup.go, gated behind Config.AllowDestructive.

	type cloudBackupIDInput struct {
		ID int `json:"id" jsonschema:"Numeric Cloud Backup task ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_abort",
		Description: "Abort a currently running Cloud Backup job.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p cloudBackupIDInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_abort: id must be a positive integer"))
		}
		if err := client.AbortCloudBackup(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("cloud_backup_abort: %w", err))
		}
		return jsonResult(map[string]any{"aborted": true, "id": p.ID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_sync",
		Description: "Run a Cloud Backup task, creating a new snapshot. Returns the async job ID immediately.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p cloudBackupIDInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_sync: id must be a positive integer"))
		}
		jobID, err := client.SyncCloudBackup(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_sync: %w", err))
		}
		return jsonResult(map[string]int{"job_id": jobID})
	})

	type restoreCloudBackupInput struct {
		ID              int    `json:"id"                         jsonschema:"Numeric Cloud Backup task ID"`
		SnapshotID      string `json:"snapshot_id"                jsonschema:"Snapshot ID to restore from (from cloud_backup_list_snapshots)"`
		Subfolder       string `json:"subfolder,omitempty"        jsonschema:"Restore only this subfolder of the snapshot; omit to restore everything"`
		DestinationPath string `json:"destination_path"           jsonschema:"Local path to restore into, beginning with /mnt"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_restore",
		Description: "Restore files from a Cloud Backup snapshot to a local path. Returns the async job ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p restoreCloudBackupInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_restore: id must be a positive integer"))
		}
		if p.SnapshotID == "" {
			return errorResult(errors.New("cloud_backup_restore: snapshot_id must not be empty"))
		}
		if p.DestinationPath == "" {
			return errorResult(errors.New("cloud_backup_restore: destination_path must not be empty"))
		}
		jobID, err := client.RestoreCloudBackup(ctx, p.ID, &truenas.CloudBackupRestoreParams{
			SnapshotID: p.SnapshotID, Subfolder: p.Subfolder, DestinationPath: p.DestinationPath,
		})
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_restore: %w", err))
		}
		return jsonResult(map[string]int{"job_id": jobID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_list_snapshots",
		Description: "List the snapshots stored in a Cloud Backup task's repository.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p cloudBackupIDInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_list_snapshots: id must be a positive integer"))
		}
		snapshots, err := client.ListCloudBackupSnapshots(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_list_snapshots: %w", err))
		}
		return jsonResult(snapshots)
	})

	type listCloudBackupSnapshotDirInput struct {
		ID         int    `json:"id"          jsonschema:"Numeric Cloud Backup task ID"`
		SnapshotID string `json:"snapshot_id" jsonschema:"Snapshot ID to browse"`
		Path       string `json:"path"        jsonschema:"Directory path within the snapshot to list"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_list_snapshot_directory",
		Description: "List the contents of a directory within a Cloud Backup snapshot.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listCloudBackupSnapshotDirInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_list_snapshot_directory: id must be a positive integer"))
		}
		if p.SnapshotID == "" {
			return errorResult(errors.New("cloud_backup_list_snapshot_directory: snapshot_id must not be empty"))
		}
		entries, err := client.ListCloudBackupSnapshotDirectory(ctx, p.ID, p.SnapshotID, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_list_snapshot_directory: %w", err))
		}
		return jsonResult(entries)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_transfer_setting_choices",
		Description: "List the valid transfer-setting values for a Cloud Backup task.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.CloudBackupTransferSettingChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("cloud_backup_transfer_setting_choices: %w", err))
		}
		return jsonResult(choices)
	})
}
