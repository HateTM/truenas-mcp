package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestCronJobCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"cronjob_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.CronJob, error) {
				return []truenas.CronJob{{ID: 1, Command: "df -h"}}, nil
			},
			"cronjob_get": func(_ context.Context, id int) (*truenas.CronJob, error) {
				return &truenas.CronJob{ID: id}, nil
			},
			"cronjob_create": func(_ context.Context, p *truenas.CreateCronJobParams) (*truenas.CronJob, error) {
				return &truenas.CronJob{ID: 1, Command: p.Command, User: p.User}, nil
			},
			"cronjob_update": func(_ context.Context, id int, _ *truenas.CreateCronJobParams) (*truenas.CronJob, error) {
				return &truenas.CronJob{ID: id}, nil
			},
			"cronjob_run": func(_ context.Context, id int) (int, error) {
				return 88, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "cronjob_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "cronjob_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "cronjob_create", map[string]any{"command": "df -h", "user": "root"})
		assertResultJSON(t, res)
	})
	t.Run("create requires command and user", func(t *testing.T) {
		res := callTool(t, cs, "cronjob_create", map[string]any{"command": "", "user": ""})
		assertError(t, res, "command and user are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "cronjob_update", map[string]any{"id": 1, "command": "df -h", "user": "root"})
		assertResultJSON(t, res)
	})
	t.Run("run", func(t *testing.T) {
		res := callTool(t, cs, "cronjob_run", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})
	t.Run("run invalid id", func(t *testing.T) {
		res := callTool(t, cs, "cronjob_run", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})
}
