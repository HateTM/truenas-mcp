package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestNVMeOFNamespaceCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_namespace_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.NVMetNamespace, error) {
				return []truenas.NVMetNamespace{{ID: 1, Subsys: 1}}, nil
			},
			"nvmeof_namespace_get": func(_ context.Context, id int) (*truenas.NVMetNamespace, error) {
				return &truenas.NVMetNamespace{ID: id}, nil
			},
			"nvmeof_namespace_create": func(_ context.Context, p *truenas.CreateNVMetNamespaceParams) (*truenas.NVMetNamespace, error) {
				return &truenas.NVMetNamespace{ID: 1, Subsys: p.Subsys, DeviceType: p.DeviceType}, nil
			},
			"nvmeof_namespace_update": func(_ context.Context, id int, _ *truenas.CreateNVMetNamespaceParams) (*truenas.NVMetNamespace, error) {
				return &truenas.NVMetNamespace{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_namespace_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_namespace_get", map[string]any{"id": 1}))
	})
	t.Run("create ZVOL", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_namespace_create", map[string]any{
			"subsys": 1, "device_type": "ZVOL", "device_path": "zvol/Storage/vm-disk",
		})
		assertResultJSON(t, res)
	})
	t.Run("create FILE requires filesize", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_namespace_create", map[string]any{
			"subsys": 1, "device_type": "FILE", "device_path": "/mnt/Storage/ns0",
		})
		assertError(t, res, "filesize is required when device_type=FILE")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_namespace_update", map[string]any{
			"id": 1, "subsys": 1, "device_type": "ZVOL", "device_path": "zvol/Storage/vm-disk",
		})
		assertResultJSON(t, res)
	})
}

func TestNVMeOFPortSubsysCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_port_subsys_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.NVMetPortSubsys, error) {
				return []truenas.NVMetPortSubsys{{ID: 1, Port: 1, Subsys: 1}}, nil
			},
			"nvmeof_port_subsys_get": func(_ context.Context, id int) (*truenas.NVMetPortSubsys, error) {
				return &truenas.NVMetPortSubsys{ID: id}, nil
			},
			"nvmeof_port_subsys_create": func(_ context.Context, p *truenas.CreateNVMetPortSubsysParams) (*truenas.NVMetPortSubsys, error) {
				return &truenas.NVMetPortSubsys{ID: 1, Port: p.Port, Subsys: p.Subsys}, nil
			},
			"nvmeof_port_subsys_update": func(_ context.Context, id int, _ *truenas.CreateNVMetPortSubsysParams) (*truenas.NVMetPortSubsys, error) {
				return &truenas.NVMetPortSubsys{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_port_subsys_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_port_subsys_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_port_subsys_create", map[string]any{"port": 1, "subsys": 1}))
	})
	t.Run("update", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_port_subsys_update", map[string]any{"id": 1, "port": 1, "subsys": 1}))
	})
}
