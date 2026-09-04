package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestISCSIPortalList(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_portal_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ISCSIPortal, error) {
				return []truenas.ISCSIPortal{{ID: 1}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_portal_list", nil)
	assertResultJSON(t, res)
}

func TestISCSIPortalGet(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_portal_get": func(_ context.Context, id int) (*truenas.ISCSIPortal, error) {
				return &truenas.ISCSIPortal{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_portal_get", map[string]any{"id": 1})
	assertResultJSON(t, res)
}

func TestISCSIPortalCreate(t *testing.T) {
	t.Run("creates portal", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_portal_create": func(_ context.Context, p *truenas.CreateISCSIPortalParams) (*truenas.ISCSIPortal, error) {
					if len(p.Listen) != 1 || p.Listen[0].IP != "0.0.0.0" {
						t.Errorf("unexpected listen: %+v", p.Listen)
					}
					return &truenas.ISCSIPortal{ID: 1, Listen: p.Listen}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_portal_create", map[string]any{
			"listen": []map[string]any{{"ip": "0.0.0.0", "port": 3260}},
		})
		assertResultJSON(t, res)
	})

	t.Run("requires at least one listen entry", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_portal_create", map[string]any{"listen": []map[string]any{}})
		assertError(t, res, "listen must have at least one entry")
	})
}

func TestISCSIPortalUpdate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_portal_update": func(_ context.Context, id int, _ *truenas.CreateISCSIPortalParams) (*truenas.ISCSIPortal, error) {
				return &truenas.ISCSIPortal{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_portal_update", map[string]any{
		"id": 1, "listen": []map[string]any{{"ip": "0.0.0.0"}},
	})
	assertResultJSON(t, res)
}

func TestISCSIPortalListenIPChoices(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_portal_listen_ip_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"0.0.0.0": "0.0.0.0"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_portal_listen_ip_choices", nil)
	assertResultJSON(t, res)
}
