package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCloudBackupDestructiveTools adds opt-in destructive Cloud Backup tools to the server.
// These tools are only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerCloudBackupDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type cloudBackupIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric Cloud Backup task ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_delete",
		Description: "Permanently delete a Cloud Backup task configuration. Does not delete the underlying restic repository or its snapshots. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p cloudBackupIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("cloud_backup_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_delete: id must be a positive integer"))
		}
		if err := client.DeleteCloudBackup(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("cloud_backup_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})

	type deleteCloudBackupSnapshotInput struct {
		ID         int    `json:"id"          jsonschema:"Numeric Cloud Backup task ID"`
		SnapshotID string `json:"snapshot_id" jsonschema:"Snapshot ID to delete (from cloud_backup_list_snapshots)"`
		Confirmed  bool   `json:"confirmed"   jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "cloud_backup_delete_snapshot",
		Description: "Permanently delete a single snapshot from a Cloud Backup task's repository. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p deleteCloudBackupSnapshotInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("cloud_backup_delete_snapshot: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("cloud_backup_delete_snapshot: id must be a positive integer"))
		}
		if p.SnapshotID == "" {
			return errorResult(errors.New("cloud_backup_delete_snapshot: snapshot_id must not be empty"))
		}
		if err := client.DeleteCloudBackupSnapshot(ctx, p.ID, p.SnapshotID); err != nil {
			return errorResult(fmt.Errorf("cloud_backup_delete_snapshot: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID, "snapshot_id": p.SnapshotID})
	})
}
