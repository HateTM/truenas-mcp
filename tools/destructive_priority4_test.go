package tools

import (
	"testing"
)

func TestPriority4DeleteTools(t *testing.T) {
	tools := []string{
		"user_delete",
		"cronjob_delete",
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

func TestUserUpdatePassword(t *testing.T) {
	t.Run("requires confirmed", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "user_update_password", map[string]any{"id": 1, "password": "newsecret1", "confirmed": false})
		assertError(t, res, "confirmed must be true")
	})

	t.Run("invalid id", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "user_update_password", map[string]any{"id": 0, "password": "newsecret1", "confirmed": true})
		assertError(t, res, "id must be a positive integer")
	})

	t.Run("requires password", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "user_update_password", map[string]any{"id": 1, "password": "", "confirmed": true})
		assertError(t, res, "password must not be empty")
	})

	t.Run("succeeds", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "user_update_password", map[string]any{"id": 1, "password": "newsecret1", "confirmed": true})
		assertResultJSON(t, res)
	})
}
