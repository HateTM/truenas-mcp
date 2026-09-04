package tools

import (
	"testing"
)

func TestPriority2DeleteTools(t *testing.T) {
	tools := []string{
		"app_registry_delete",
		"alertservice_delete",
		"cloud_backup_delete",
	}

	for _, tool := range tools {
		t.Run(tool+"/requires confirmed", func(t *testing.T) {
			cs, cleanup := connectTestServer(t, &mockTruenasClient{})
			defer cleanup()

			res := callTool(t, cs, tool, map[string]any{"id": 1, "confirmed": false})
			assertError(t, res, "confirmed must be true")
		})

		t.Run(tool+"/invalid id", func(t *testing.T) {
			cs, cleanup := connectTestServer(t, &mockTruenasClient{})
			defer cleanup()

			res := callTool(t, cs, tool, map[string]any{"id": 0, "confirmed": true})
			assertError(t, res, "id must be a positive integer")
		})

		t.Run(tool+"/succeeds", func(t *testing.T) {
			cs, cleanup := connectTestServer(t, &mockTruenasClient{})
			defer cleanup()

			res := callTool(t, cs, tool, map[string]any{"id": 1, "confirmed": true})
			assertResultJSON(t, res)
		})
	}
}

func TestCloudBackupDeleteSnapshot(t *testing.T) {
	t.Run("requires confirmed", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "cloud_backup_delete_snapshot", map[string]any{
			"id": 1, "snapshot_id": "snap-1", "confirmed": false,
		})
		assertError(t, res, "confirmed must be true")
	})

	t.Run("invalid id", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "cloud_backup_delete_snapshot", map[string]any{
			"id": 0, "snapshot_id": "snap-1", "confirmed": true,
		})
		assertError(t, res, "id must be a positive integer")
	})

	t.Run("requires snapshot_id", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "cloud_backup_delete_snapshot", map[string]any{
			"id": 1, "snapshot_id": "", "confirmed": true,
		})
		assertError(t, res, "snapshot_id must not be empty")
	})

	t.Run("succeeds", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "cloud_backup_delete_snapshot", map[string]any{
			"id": 1, "snapshot_id": "snap-1", "confirmed": true,
		})
		assertResultJSON(t, res)
	})
}
