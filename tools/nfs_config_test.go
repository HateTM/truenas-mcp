package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestNFSConfigTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"nfs_bindip_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"0.0.0.0": "0.0.0.0"}, nil
			},
			"nfs_client_count": func(_ context.Context) (int, error) { return 2, nil },
			"nfs_config": func(_ context.Context) (*truenas.NFSConfig, error) {
				return &truenas.NFSConfig{Servers: 4, V4: true}, nil
			},
			"nfs_get_nfs3_clients": func(_ context.Context) ([]truenas.NFSClient, error) {
				return []truenas.NFSClient{{Client: "10.0.0.5"}}, nil
			},
			"nfs_get_nfs4_clients": func(_ context.Context) ([]truenas.NFSClient, error) {
				return []truenas.NFSClient{{Client: "10.0.0.6"}}, nil
			},
			"nfs_update": func(_ context.Context, p *truenas.UpdateNFSConfigParams) (*truenas.NFSConfig, error) {
				return &truenas.NFSConfig{Servers: p.Servers, V4: p.V4}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	for _, tool := range []string{
		"nfs_bindip_choices",
		"nfs_client_count",
		"nfs_config",
		"nfs_get_nfs3_clients",
		"nfs_get_nfs4_clients",
	} {
		t.Run(tool, func(t *testing.T) {
			res := callTool(t, cs, tool, nil)
			assertResultJSON(t, res)
		})
	}

	t.Run("nfs_update", func(t *testing.T) {
		res := callTool(t, cs, "nfs_update", map[string]any{"servers": 4, "v4": true})
		assertResultJSON(t, res)
	})
}
