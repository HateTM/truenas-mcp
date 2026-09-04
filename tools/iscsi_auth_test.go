package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestISCSIAuthList(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_auth_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ISCSIAuth, error) {
				return []truenas.ISCSIAuth{{ID: 1, Tag: 1, User: "chapuser"}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_auth_list", nil)
	assertResultJSON(t, res)
}

func TestISCSIAuthGet(t *testing.T) {
	t.Run("returns credential as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_auth_get": func(_ context.Context, id int) (*truenas.ISCSIAuth, error) {
					return &truenas.ISCSIAuth{ID: id, User: "chapuser"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_auth_get", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_auth_get", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_auth_get": func(_ context.Context, _ int) (*truenas.ISCSIAuth, error) {
					return nil, errors.New("not found")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_auth_get", map[string]any{"id": 1})
		assertError(t, res, "not found")
	})
}

func TestISCSIAuthCreate(t *testing.T) {
	t.Run("creates credential", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_auth_create": func(_ context.Context, p *truenas.CreateISCSIAuthParams) (*truenas.ISCSIAuth, error) {
					if p.User != "chapuser" || p.Secret != "supersecret1" {
						t.Errorf("unexpected params: %+v", p)
					}
					return &truenas.ISCSIAuth{ID: 1, Tag: p.Tag, User: p.User}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_auth_create", map[string]any{"tag": 1, "user": "chapuser", "secret": "supersecret1"})
		assertResultJSON(t, res)
	})

	t.Run("requires user and secret", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_auth_create", map[string]any{"tag": 1, "user": "", "secret": ""})
		assertError(t, res, "user and secret are required")
	})
}

func TestISCSIAuthUpdate(t *testing.T) {
	t.Run("updates credential", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_auth_update": func(_ context.Context, id int, _ *truenas.CreateISCSIAuthParams) (*truenas.ISCSIAuth, error) {
					return &truenas.ISCSIAuth{ID: id}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_auth_update", map[string]any{"id": 1, "tag": 1, "user": "chapuser", "secret": "supersecret1"})
		assertResultJSON(t, res)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_auth_update", map[string]any{"id": 0, "tag": 1, "user": "u", "secret": "s"})
		assertError(t, res, "id must be a positive integer")
	})
}

func TestISCSIGlobalTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_global_alua_enabled": func(_ context.Context) (bool, error) { return true, nil },
			"iscsi_global_client_count": func(_ context.Context) (int, error) { return 3, nil },
			"iscsi_global_config": func(_ context.Context) (*truenas.ISCSIGlobalConfig, error) {
				return &truenas.ISCSIGlobalConfig{Basename: "iqn.2005-10.org.freenas.ctl"}, nil
			},
			"iscsi_global_iser_enabled": func(_ context.Context) (bool, error) { return false, nil },
			"iscsi_global_sessions": func(_ context.Context) ([]truenas.ISCSISession, error) {
				return []truenas.ISCSISession{{ID: 1, Initiator: "iqn.initiator", Target: "iqn.target"}}, nil
			},
			"iscsi_global_update": func(_ context.Context, _ *truenas.UpdateISCSIGlobalParams) (*truenas.ISCSIGlobalConfig, error) {
				return &truenas.ISCSIGlobalConfig{Basename: "iqn.2005-10.org.freenas.ctl"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	for _, tool := range []string{
		"iscsi_global_alua_enabled",
		"iscsi_global_client_count",
		"iscsi_global_config",
		"iscsi_global_iser_enabled",
		"iscsi_global_sessions",
	} {
		t.Run(tool, func(t *testing.T) {
			res := callTool(t, cs, tool, nil)
			assertResultJSON(t, res)
		})
	}

	t.Run("iscsi_global_update", func(t *testing.T) {
		res := callTool(t, cs, "iscsi_global_update", map[string]any{"basename": "iqn.2005-10.org.freenas.ctl"})
		assertResultJSON(t, res)
	})
}
