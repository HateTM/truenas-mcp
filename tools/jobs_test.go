package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestGetJob(t *testing.T) {
	t.Run("returns job as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"get_job": func(_ context.Context, id int) (*truenas.Job, error) {
					return &truenas.Job{ID: id, State: "RUNNING"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "get_job", map[string]any{"id": 10904})
		assertResultJSON(t, res)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "get_job", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"get_job": func(_ context.Context, _ int) (*truenas.Job, error) {
					return nil, errors.New("job not found")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "get_job", map[string]any{"id": 1})
		assertError(t, res, "job not found")
	})
}

func TestAbortJob(t *testing.T) {
	t.Run("aborts job", func(t *testing.T) {
		var gotID int
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"abort_job": func(_ context.Context, id int) error {
					gotID = id
					return nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "abort_job", map[string]any{"id": 10818})
		assertResultJSON(t, res)
		if gotID != 10818 {
			t.Errorf("id passed to client = %d, want 10818", gotID)
		}
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "abort_job", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"abort_job": func(_ context.Context, _ int) error {
					return errors.New("job already finished")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "abort_job", map[string]any{"id": 1})
		assertError(t, res, "job already finished")
	})
}
