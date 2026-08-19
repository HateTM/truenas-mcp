package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestNVMeOFPortCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_port_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.NVMetPort, error) {
				return []truenas.NVMetPort{{ID: 1, AddrTrtype: "TCP", AddrTraddr: "0.0.0.0"}}, nil
			},
			"nvmeof_port_get": func(_ context.Context, id int) (*truenas.NVMetPort, error) {
				return &truenas.NVMetPort{ID: id}, nil
			},
			"nvmeof_port_create": func(_ context.Context, p *truenas.CreateNVMetPortParams) (*truenas.NVMetPort, error) {
				return &truenas.NVMetPort{ID: 1, AddrTrtype: p.AddrTrtype, AddrTraddr: p.AddrTraddr}, nil
			},
			"nvmeof_port_update": func(_ context.Context, id int, _ *truenas.CreateNVMetPortParams) (*truenas.NVMetPort, error) {
				return &truenas.NVMetPort{ID: id}, nil
			},
			"nvmeof_port_transport_address_choices": func(_ context.Context, _ string) (map[string]string, error) {
				return map[string]string{"0.0.0.0": "0.0.0.0"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_port_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_port_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_port_create", map[string]any{"addr_trtype": "TCP", "addr_traddr": "0.0.0.0"})
		assertResultJSON(t, res)
	})
	t.Run("create requires addr_trtype and addr_traddr", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_port_create", map[string]any{"addr_trtype": "", "addr_traddr": ""})
		assertError(t, res, "addr_trtype and addr_traddr are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_port_update", map[string]any{"id": 1, "addr_trtype": "TCP", "addr_traddr": "0.0.0.0"})
		assertResultJSON(t, res)
	})
	t.Run("transport_address_choices", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_port_transport_address_choices", map[string]any{"trtype": "TCP"})
		assertResultJSON(t, res)
	})
}
