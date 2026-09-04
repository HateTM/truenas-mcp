package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestCloudBackupCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"cloud_backup_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.CloudBackupTask, error) {
				return []truenas.CloudBackupTask{{ID: 1, Path: "/mnt/Storage/docs"}}, nil
			},
			"cloud_backup_get": func(_ context.Context, id int) (*truenas.CloudBackupTask, error) {
				return &truenas.CloudBackupTask{ID: id}, nil
			},
			"cloud_backup_create": func(_ context.Context, p *truenas.CreateCloudBackupParams) (*truenas.CloudBackupTask, error) {
				return &truenas.CloudBackupTask{ID: 1, Path: p.Path, Credentials: p.Credentials}, nil
			},
			"cloud_backup_update": func(_ context.Context, id int, _ *truenas.CreateCloudBackupParams) (*truenas.CloudBackupTask, error) {
				return &truenas.CloudBackupTask{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "cloud_backup_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "cloud_backup_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_create", map[string]any{
			"path": "/mnt/Storage/docs", "credentials": 1, "password": "supersecret1",
		})
		assertResultJSON(t, res)
	})
	t.Run("create requires path/credentials/password", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_create", map[string]any{"path": "", "credentials": 0, "password": ""})
		assertError(t, res, "path must not be empty")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_update", map[string]any{
			"id": 1, "path": "/mnt/Storage/docs", "credentials": 1, "password": "supersecret1",
		})
		assertResultJSON(t, res)
	})
}

func TestCloudBackupOperations(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"cloud_backup_abort": func(_ context.Context, id int) error { return nil },
			"cloud_backup_sync": func(_ context.Context, id int) (int, error) {
				return 99, nil
			},
			"cloud_backup_restore": func(_ context.Context, id int, p *truenas.CloudBackupRestoreParams) (int, error) {
				if p.SnapshotID != "snap-1" {
					t.Errorf("SnapshotID = %q, want snap-1", p.SnapshotID)
				}
				return 100, nil
			},
			"cloud_backup_list_snapshots": func(_ context.Context, id int) ([]truenas.CloudBackupSnapshot, error) {
				return []truenas.CloudBackupSnapshot{{ID: "snap-1"}}, nil
			},
			"cloud_backup_list_snapshot_directory": func(_ context.Context, id int, snapshotID, path string) ([]truenas.DirEntry, error) {
				return []truenas.DirEntry{{Name: "file.txt"}}, nil
			},
			"cloud_backup_transfer_setting_choices": func(_ context.Context) ([]string, error) {
				return []string{"DEFAULT", "PERFORMANCE", "FAST_STORAGE"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("abort", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_abort", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})
	t.Run("sync", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_sync", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})
	t.Run("restore", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_restore", map[string]any{
			"id": 1, "snapshot_id": "snap-1", "destination_path": "/mnt/Storage/restore",
		})
		assertResultJSON(t, res)
	})
	t.Run("restore requires snapshot_id", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_restore", map[string]any{
			"id": 1, "snapshot_id": "", "destination_path": "/mnt/Storage/restore",
		})
		assertError(t, res, "snapshot_id must not be empty")
	})
	t.Run("list_snapshots", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_list_snapshots", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})
	t.Run("list_snapshot_directory", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_list_snapshot_directory", map[string]any{
			"id": 1, "snapshot_id": "snap-1", "path": "/",
		})
		assertResultJSON(t, res)
	})
	t.Run("transfer_setting_choices", func(t *testing.T) {
		res := callTool(t, cs, "cloud_backup_transfer_setting_choices", nil)
		assertResultJSON(t, res)
	})
}
