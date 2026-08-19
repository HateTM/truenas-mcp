package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestMiscTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"device_get_info": func(_ context.Context, deviceType string) (map[string]any, error) {
				if deviceType != "DISK" {
					t.Errorf("deviceType = %q, want DISK", deviceType)
				}
				return map[string]any{"sda": map[string]any{"size": 1000}}, nil
			},
			"dns_query": func(_ context.Context) (*truenas.DNSConfig, error) {
				return &truenas.DNSConfig{Nameservers: []string{"1.1.1.1", "8.8.8.8"}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("device_get_info", func(t *testing.T) {
		res := callTool(t, cs, "device_get_info", map[string]any{"type": "DISK"})
		assertResultJSON(t, res)
	})
	t.Run("device_get_info requires type", func(t *testing.T) {
		res := callTool(t, cs, "device_get_info", map[string]any{"type": ""})
		assertError(t, res, "type must not be empty")
	})
	t.Run("dns_query", func(t *testing.T) {
		res := callTool(t, cs, "dns_query", nil)
		assertResultJSON(t, res)
	})
}
