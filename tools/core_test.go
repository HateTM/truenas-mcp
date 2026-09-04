package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestCoreConnectivityTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"core_ping":        func(_ context.Context) (string, error) { return "pong", nil },
			"core_ping_remote": func(_ context.Context, _ map[string]any) (string, error) { return "pong", nil },
			"core_get_methods": func(_ context.Context, _ string) (map[string]any, error) {
				return map[string]any{"pool.query": map[string]any{}}, nil
			},
			"core_get_services": func(_ context.Context) ([]map[string]any, error) {
				return []map[string]any{{"name": "pool"}}, nil
			},
			"core_arp": func(_ context.Context, _ string) (map[string]string, error) {
				return map[string]string{"192.168.1.1": "aa:bb:cc:dd:ee:ff"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("ping", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "core_ping", nil))
	})
	t.Run("ping_remote", func(t *testing.T) {
		res := callTool(t, cs, "core_ping_remote", map[string]any{"params": map[string]any{"type": "ICMP", "hosts": []string{"1.1.1.1"}}})
		assertResultJSON(t, res)
	})
	t.Run("get_methods", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "core_get_methods", nil))
	})
	t.Run("get_services", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "core_get_services", nil))
	})
	t.Run("arp", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "core_arp", nil))
	})
}

func TestCoreJobTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"core_job_download_logs": func(_ context.Context, jobID int) (string, error) {
				return "https://truenas.local/_download/1", nil
			},
			"core_job_wait": func(_ context.Context, jobID int) (*truenas.Job, error) {
				return &truenas.Job{ID: jobID, State: "SUCCESS"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("job_download_logs", func(t *testing.T) {
		res := callTool(t, cs, "core_job_download_logs", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})
	t.Run("job_download_logs invalid id", func(t *testing.T) {
		res := callTool(t, cs, "core_job_download_logs", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})
	t.Run("job_wait", func(t *testing.T) {
		res := callTool(t, cs, "core_job_wait", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})
}

func TestCoreBulkAndDownloadTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"core_bulk": func(_ context.Context, method string, _ [][]any) ([]any, error) {
				if method != "pool.dataset.update" {
					t.Errorf("method = %q, want pool.dataset.update", method)
				}
				return []any{map[string]any{"result": "ok"}}, nil
			},
			"core_download": func(_ context.Context, method string, _ []any, filename string) (*truenas.CoreDownloadResult, error) {
				return &truenas.CoreDownloadResult{JobID: 1, URL: "https://truenas.local/_download/1"}, nil
			},
			"core_resize_shell": func(_ context.Context, id string, cols, rows int) error { return nil },
			"core_subscribe": func(_ context.Context, event string) (string, error) {
				return "sub-1", nil
			},
			"core_unsubscribe": func(_ context.Context, subscriptionID string) error { return nil },
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("bulk", func(t *testing.T) {
		res := callTool(t, cs, "core_bulk", map[string]any{
			"method": "pool.dataset.update",
			"params": [][]any{{"tank/a", map[string]any{"comments": "x"}}},
		})
		assertResultJSON(t, res)
	})
	t.Run("bulk requires method", func(t *testing.T) {
		res := callTool(t, cs, "core_bulk", map[string]any{"method": "", "params": [][]any{}})
		assertError(t, res, "method must not be empty")
	})
	t.Run("download", func(t *testing.T) {
		res := callTool(t, cs, "core_download", map[string]any{"method": "config.save", "filename": "config.tar"})
		assertResultJSON(t, res)
	})
	t.Run("download requires filename", func(t *testing.T) {
		res := callTool(t, cs, "core_download", map[string]any{"method": "config.save", "filename": ""})
		assertError(t, res, "filename must not be empty")
	})
	t.Run("resize_shell", func(t *testing.T) {
		res := callTool(t, cs, "core_resize_shell", map[string]any{"id": "shell-1", "cols": 80, "rows": 24})
		assertResultJSON(t, res)
	})
	t.Run("resize_shell requires positive cols/rows", func(t *testing.T) {
		res := callTool(t, cs, "core_resize_shell", map[string]any{"id": "shell-1", "cols": 0, "rows": 24})
		assertError(t, res, "cols and rows must be positive integers")
	})
	t.Run("subscribe", func(t *testing.T) {
		res := callTool(t, cs, "core_subscribe", map[string]any{"event": "pool.query"})
		assertResultJSON(t, res)
	})
	t.Run("unsubscribe", func(t *testing.T) {
		res := callTool(t, cs, "core_unsubscribe", map[string]any{"subscription_id": "sub-1"})
		assertResultJSON(t, res)
	})
}
