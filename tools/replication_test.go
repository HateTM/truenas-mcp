package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestReplicationTaskCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"replication_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ReplicationTask, error) {
				return []truenas.ReplicationTask{{ID: 1, Name: "nightly"}}, nil
			},
			"replication_get": func(_ context.Context, id int) (*truenas.ReplicationTask, error) {
				return &truenas.ReplicationTask{ID: id}, nil
			},
			"replication_create": func(_ context.Context, p *truenas.CreateReplicationParams) (*truenas.ReplicationTask, error) {
				return &truenas.ReplicationTask{ID: 1, Name: p.Name}, nil
			},
			"replication_update": func(_ context.Context, id int, _ *truenas.CreateReplicationParams) (*truenas.ReplicationTask, error) {
				return &truenas.ReplicationTask{ID: id}, nil
			},
			"replication_test": func(_ context.Context, id int) (map[string]any, error) {
				return map[string]any{"valid": true}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "replication_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "replication_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "replication_create", map[string]any{
			"name": "nightly", "source_datasets": []string{"Storage/backups"},
			"target_dataset": "Backup/Storage/backups", "direction": "PUSH",
		})
		assertResultJSON(t, res)
	})
	t.Run("create requires source_datasets", func(t *testing.T) {
		res := callTool(t, cs, "replication_create", map[string]any{
			"name": "nightly", "source_datasets": []string{},
			"target_dataset": "Backup/Storage/backups", "direction": "PUSH",
		})
		assertError(t, res, "source_datasets must have at least one entry")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "replication_update", map[string]any{
			"id": 1, "name": "nightly", "source_datasets": []string{"Storage/backups"},
			"target_dataset": "Backup/Storage/backups", "direction": "PUSH",
		})
		assertResultJSON(t, res)
	})
	t.Run("test", func(t *testing.T) {
		res := callTool(t, cs, "replication_test", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})
}

func TestReplicationGlobalTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"replication_schemas": func(_ context.Context) (map[string]any, error) {
				return map[string]any{"periodic_snapshot_task": map[string]any{}}, nil
			},
			"replication_global_config": func(_ context.Context) (*truenas.ReplicationGlobalConfig, error) {
				return &truenas.ReplicationGlobalConfig{MaxParallelReplicationTasks: 2}, nil
			},
			"replication_global_update": func(_ context.Context, p *truenas.UpdateReplicationGlobalParams) (*truenas.ReplicationGlobalConfig, error) {
				return &truenas.ReplicationGlobalConfig{MaxParallelReplicationTasks: p.MaxParallelReplicationTasks}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("schemas", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "replication_schemas", nil))
	})
	t.Run("global_config", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "replication_global_config", nil))
	})
	t.Run("global_update", func(t *testing.T) {
		res := callTool(t, cs, "replication_global_update", map[string]any{"max_parallel_replication_tasks": 3})
		assertResultJSON(t, res)
	})
}

func TestReplicationEndpointCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"replication_endpoint_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ReplicationEndpoint, error) {
				return []truenas.ReplicationEndpoint{{ID: 1, Name: "offsite"}}, nil
			},
			"replication_endpoint_get": func(_ context.Context, id int) (*truenas.ReplicationEndpoint, error) {
				return &truenas.ReplicationEndpoint{ID: id}, nil
			},
			"replication_endpoint_create": func(_ context.Context, p *truenas.CreateReplicationEndpointParams) (*truenas.ReplicationEndpoint, error) {
				return &truenas.ReplicationEndpoint{ID: 1, Name: p.Name, URI: p.URI}, nil
			},
			"replication_endpoint_update": func(_ context.Context, id int, _ *truenas.CreateReplicationEndpointParams) (*truenas.ReplicationEndpoint, error) {
				return &truenas.ReplicationEndpoint{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "replication_endpoint_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "replication_endpoint_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "replication_endpoint_create", map[string]any{
			"name": "offsite", "uri": "ws://remote.example.com/websocket",
		})
		assertResultJSON(t, res)
	})
	t.Run("create requires name and uri", func(t *testing.T) {
		res := callTool(t, cs, "replication_endpoint_create", map[string]any{"name": "", "uri": ""})
		assertError(t, res, "name and uri are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "replication_endpoint_update", map[string]any{
			"id": 1, "name": "offsite", "uri": "ws://remote.example.com/websocket",
		})
		assertResultJSON(t, res)
	})
}
