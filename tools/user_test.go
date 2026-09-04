package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestUserCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"user_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.User, error) {
				return []truenas.User{{ID: 1, Username: "backup-svc"}}, nil
			},
			"user_get": func(_ context.Context, id int) (*truenas.User, error) {
				return &truenas.User{ID: id}, nil
			},
			"user_create": func(_ context.Context, p *truenas.CreateUserParams) (*truenas.User, error) {
				return &truenas.User{ID: 1, Username: p.Username}, nil
			},
			"user_update": func(_ context.Context, id int, _ *truenas.CreateUserParams) (*truenas.User, error) {
				return &truenas.User{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "user_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "user_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "user_create", map[string]any{"username": "backup-svc", "full_name": "Backup Service"})
		assertResultJSON(t, res)
	})
	t.Run("create requires username and full_name", func(t *testing.T) {
		res := callTool(t, cs, "user_create", map[string]any{"username": "", "full_name": ""})
		assertError(t, res, "username and full_name are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "user_update", map[string]any{"id": 1, "username": "backup-svc", "full_name": "Backup Service"})
		assertResultJSON(t, res)
	})
}

func TestUserChoiceTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"user_schemas": func(_ context.Context) (map[string]any, error) {
				return map[string]any{"user_create": map[string]any{}}, nil
			},
			"user_home_directory_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"/mnt/Storage/home": "/mnt/Storage/home"}, nil
			},
			"user_shell_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"/usr/bin/bash": "bash"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	for _, tool := range []string{"user_schemas", "user_home_directory_choices", "user_shell_choices"} {
		t.Run(tool, func(t *testing.T) {
			res := callTool(t, cs, tool, nil)
			assertResultJSON(t, res)
		})
	}
}
