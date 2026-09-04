package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestNVMeOFHostSubsysCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_host_subsys_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.NVMetHostSubsys, error) {
				return []truenas.NVMetHostSubsys{{ID: 1, Host: 1, Subsys: 1}}, nil
			},
			"nvmeof_host_subsys_get": func(_ context.Context, id int) (*truenas.NVMetHostSubsys, error) {
				return &truenas.NVMetHostSubsys{ID: id}, nil
			},
			"nvmeof_host_subsys_create": func(_ context.Context, p *truenas.CreateNVMetHostSubsysParams) (*truenas.NVMetHostSubsys, error) {
				return &truenas.NVMetHostSubsys{ID: 1, Host: p.Host, Subsys: p.Subsys}, nil
			},
			"nvmeof_host_subsys_update": func(_ context.Context, id int, _ *truenas.CreateNVMetHostSubsysParams) (*truenas.NVMetHostSubsys, error) {
				return &truenas.NVMetHostSubsys{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_host_subsys_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_host_subsys_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_host_subsys_create", map[string]any{"host": 1, "subsys": 1}))
	})
	t.Run("create requires host and subsys", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_host_subsys_create", map[string]any{"host": 0, "subsys": 0})
		assertError(t, res, "host and subsys are required")
	})
	t.Run("update", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_host_subsys_update", map[string]any{"id": 1, "host": 1, "subsys": 1}))
	})
}

func TestNVMeOFSubsysCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_subsys_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.NVMetSubsys, error) {
				return []truenas.NVMetSubsys{{ID: 1, Name: "subsys0"}}, nil
			},
			"nvmeof_subsys_get": func(_ context.Context, id int) (*truenas.NVMetSubsys, error) {
				return &truenas.NVMetSubsys{ID: id}, nil
			},
			"nvmeof_subsys_create": func(_ context.Context, p *truenas.CreateNVMetSubsysParams) (*truenas.NVMetSubsys, error) {
				return &truenas.NVMetSubsys{ID: 1, Name: p.Name}, nil
			},
			"nvmeof_subsys_update": func(_ context.Context, id int, _ *truenas.CreateNVMetSubsysParams) (*truenas.NVMetSubsys, error) {
				return &truenas.NVMetSubsys{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_subsys_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_subsys_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_subsys_create", map[string]any{"name": "subsys0"}))
	})
	t.Run("create requires name", func(t *testing.T) {
		res := callTool(t, cs, "nvmeof_subsys_create", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
	t.Run("update", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "nvmeof_subsys_update", map[string]any{"id": 1, "name": "subsys0"}))
	})
}
