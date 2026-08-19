package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestDirectoryServicesTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"directoryservices_cache_refresh": func(_ context.Context) error { return nil },
			"directoryservices_certificate_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"1": "wildcard"}, nil
			},
			"directoryservices_config": func(_ context.Context) (*truenas.DirectoryServicesConfig, error) {
				return &truenas.DirectoryServicesConfig{Service: "ACTIVEDIRECTORY", Enable: true}, nil
			},
			"directoryservices_leave": func(_ context.Context, _ map[string]any) error { return nil },
			"directoryservices_status": func(_ context.Context) (*truenas.DirectoryServicesStatus, error) {
				return &truenas.DirectoryServicesStatus{Type: "ACTIVEDIRECTORY", Status: "HEALTHY"}, nil
			},
			"directoryservices_sync_keytab": func(_ context.Context) error { return nil },
			"directoryservices_update": func(_ context.Context, p *truenas.UpdateDirectoryServicesParams) (*truenas.DirectoryServicesConfig, error) {
				return &truenas.DirectoryServicesConfig{Service: p.Service, Enable: p.Enable}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("cache_refresh", func(t *testing.T) {
		res := callTool(t, cs, "directoryservices_cache_refresh", nil)
		assertResultJSON(t, res)
	})
	t.Run("certificate_choices", func(t *testing.T) {
		res := callTool(t, cs, "directoryservices_certificate_choices", nil)
		assertResultJSON(t, res)
	})
	t.Run("config", func(t *testing.T) {
		res := callTool(t, cs, "directoryservices_config", nil)
		assertResultJSON(t, res)
	})
	t.Run("leave", func(t *testing.T) {
		res := callTool(t, cs, "directoryservices_leave", map[string]any{
			"params": map[string]any{"username": "admin", "password": "supersecret1"},
		})
		assertResultJSON(t, res)
	})
	t.Run("status", func(t *testing.T) {
		res := callTool(t, cs, "directoryservices_status", nil)
		assertResultJSON(t, res)
	})
	t.Run("sync_keytab", func(t *testing.T) {
		res := callTool(t, cs, "directoryservices_sync_keytab", nil)
		assertResultJSON(t, res)
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "directoryservices_update", map[string]any{"service": "ACTIVEDIRECTORY", "enable": true})
		assertResultJSON(t, res)
	})
}
