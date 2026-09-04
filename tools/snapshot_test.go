package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestListSnapshots(t *testing.T) {
	t.Run("returns snapshots as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_snapshots": func(_ context.Context, _ string, _ ...truenas.ListOptions) ([]truenas.Snapshot, error) {
					return []truenas.Snapshot{{ID: "tank/data@snap1"}}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_snapshots", nil)
		assertResultJSON(t, res)
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_snapshots": func(_ context.Context, _ string, _ ...truenas.ListOptions) ([]truenas.Snapshot, error) {
					return nil, errors.New("API error")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_snapshots", nil)
		assertError(t, res, "API error")
	})
}
