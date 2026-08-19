package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestListAlerts(t *testing.T) {
	t.Run("returns alerts as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_alerts": func(_ context.Context) ([]truenas.Alert, error) {
					return []truenas.Alert{{Level: "WARNING", Formatted: "Pool NaS is low on space"}}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_alerts", nil)
		assertResultJSON(t, res)
	})

	t.Run("returns error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_alerts": func(_ context.Context) ([]truenas.Alert, error) {
					return nil, errors.New("API error")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_alerts", nil)
		assertError(t, res, "API error")
	})
}
