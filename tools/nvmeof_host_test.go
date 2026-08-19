package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestNVMeOFGlobalConfig(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_global_config": func(_ context.Context) (*truenas.NVMetGlobalConfig, error) {
				return &truenas.NVMetGlobalConfig{ANA: true}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "nvmeof_global_config", nil)
	assertResultJSON(t, res)
}

func TestNVMeOFGlobalUpdate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_global_update": func(_ context.Context, p *truenas.UpdateNVMetGlobalParams) (*truenas.NVMetGlobalConfig, error) {
				return &truenas.NVMetGlobalConfig{ANA: p.ANA}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "nvmeof_global_update", map[string]any{"ana": true})
	assertResultJSON(t, res)
}

func TestNVMeOFHostList(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_host_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.NVMetHost, error) {
				return []truenas.NVMetHost{{ID: 1, HostNQN: "nqn.host"}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "nvmeof_host_list", nil)
	assertResultJSON(t, res)
}

func TestNVMeOFHostGet(t *testing.T) {
	t.Run("returns host as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"nvmeof_host_get": func(_ context.Context, id int) (*truenas.NVMetHost, error) {
					return &truenas.NVMetHost{ID: id}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "nvmeof_host_get", map[string]any{"id": 1})
		assertResultJSON(t, res)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "nvmeof_host_get", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})
}

func TestNVMeOFHostCreate(t *testing.T) {
	t.Run("creates host", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"nvmeof_host_create": func(_ context.Context, p *truenas.CreateNVMetHostParams) (*truenas.NVMetHost, error) {
					return &truenas.NVMetHost{ID: 1, HostNQN: p.HostNQN}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "nvmeof_host_create", map[string]any{"hostnqn": "nqn.2014-08.org.nvmexpress:uuid:1"})
		assertResultJSON(t, res)
	})

	t.Run("requires hostnqn", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "nvmeof_host_create", map[string]any{"hostnqn": ""})
		assertError(t, res, "hostnqn must not be empty")
	})
}

func TestNVMeOFHostUpdate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_host_update": func(_ context.Context, id int, _ *truenas.CreateNVMetHostParams) (*truenas.NVMetHost, error) {
				return &truenas.NVMetHost{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "nvmeof_host_update", map[string]any{"id": 1, "hostnqn": "nqn.host"})
	assertResultJSON(t, res)
}

func TestNVMeOFHostMiscTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nvmeof_host_dhchap_dhgroup_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"2048-BIT": "2048-BIT"}, nil
			},
			"nvmeof_host_dhchap_hash_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"SHA-256": "SHA-256"}, nil
			},
			"nvmeof_host_generate_key": func(_ context.Context) (string, error) {
				return "DHHC-1:00:...", nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	for _, tool := range []string{
		"nvmeof_host_dhchap_dhgroup_choices",
		"nvmeof_host_dhchap_hash_choices",
		"nvmeof_host_generate_key",
	} {
		t.Run(tool, func(t *testing.T) {
			res := callTool(t, cs, tool, nil)
			assertResultJSON(t, res)
		})
	}
}
