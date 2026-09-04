package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestDiskQueryTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"disk_query": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.Disk, error) {
				return []truenas.Disk{{Identifier: "{serial}sN12345", Name: "sda"}}, nil
			},
			"disk_details": func(_ context.Context) (map[string]any, error) {
				return map[string]any{"sda": map[string]any{"size": 1000}}, nil
			},
			"disk_get_used": func(_ context.Context, name string) (int64, error) {
				return 12345, nil
			},
			"disk_temperature_agg": func(_ context.Context, names []string, days int) (map[string]truenas.DiskTemperatureAgg, error) {
				return map[string]truenas.DiskTemperatureAgg{"sda": {Min: 30, Max: 45, Avg: 37}}, nil
			},
			"disk_temperature_alerts": func(_ context.Context, names []string) ([]string, error) {
				return []string{}, nil
			},
			"disk_temperatures": func(_ context.Context, names []string) (map[string]int, error) {
				return map[string]int{"sda": 37}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("query", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "disk_query", nil))
	})
	t.Run("details", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "disk_details", nil))
	})
	t.Run("get_used", func(t *testing.T) {
		res := callTool(t, cs, "disk_get_used", map[string]any{"name": "sda"})
		assertResultJSON(t, res)
	})
	t.Run("get_used requires name", func(t *testing.T) {
		res := callTool(t, cs, "disk_get_used", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
	t.Run("temperature_agg", func(t *testing.T) {
		res := callTool(t, cs, "disk_temperature_agg", map[string]any{"names": []string{"sda"}, "days": 7})
		assertResultJSON(t, res)
	})
	t.Run("temperature_alerts", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "disk_temperature_alerts", nil))
	})
	t.Run("temperatures", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "disk_temperatures", nil))
	})
}

func TestDiskUpdate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"disk_update": func(_ context.Context, identifier string, _ *truenas.UpdateDiskParams) (*truenas.Disk, error) {
				return &truenas.Disk{Identifier: identifier}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "disk_update", map[string]any{"identifier": "{serial}sN12345", "description": "backup disk"})
		assertResultJSON(t, res)
	})
	t.Run("update requires identifier", func(t *testing.T) {
		res := callTool(t, cs, "disk_update", map[string]any{"identifier": ""})
		assertError(t, res, "identifier must not be empty")
	})
}

func TestDiskWipe(t *testing.T) {
	t.Run("requires confirmed", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "disk_wipe", map[string]any{"identifier": "{serial}sN12345", "mode": "QUICK", "confirmed": false})
		assertError(t, res, "confirmed must be true")
	})

	t.Run("requires identifier and mode", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "disk_wipe", map[string]any{"identifier": "", "mode": "", "confirmed": true})
		assertError(t, res, "identifier must not be empty")
	})

	t.Run("succeeds", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"disk_wipe": func(_ context.Context, identifier string, p *truenas.WipeDiskParams) (int, error) {
					if p.Mode != "QUICK" {
						t.Errorf("mode = %q, want QUICK", p.Mode)
					}
					return 55, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "disk_wipe", map[string]any{"identifier": "{serial}sN12345", "mode": "QUICK", "confirmed": true})
		assertResultJSON(t, res)
	})
}
