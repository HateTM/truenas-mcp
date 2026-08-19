package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestCertificateCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"certificate_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.Certificate, error) {
				return []truenas.Certificate{{ID: 1, Name: "wildcard"}}, nil
			},
			"certificate_get": func(_ context.Context, id int) (*truenas.Certificate, error) {
				return &truenas.Certificate{ID: id}, nil
			},
			"certificate_create": func(_ context.Context, p *truenas.CreateCertificateParams) (*truenas.Certificate, error) {
				return &truenas.Certificate{ID: 1, Name: p.Name}, nil
			},
			"certificate_update": func(_ context.Context, id int, _ *truenas.CreateCertificateParams) (*truenas.Certificate, error) {
				return &truenas.Certificate{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "certificate_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "certificate_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "certificate_create", map[string]any{
			"name": "wildcard", "create_type": "CERTIFICATE_CREATE_INTERNAL", "common": "example.com",
		})
		assertResultJSON(t, res)
	})
	t.Run("create requires name and create_type", func(t *testing.T) {
		res := callTool(t, cs, "certificate_create", map[string]any{"name": "", "create_type": ""})
		assertError(t, res, "name and create_type are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "certificate_update", map[string]any{"id": 1, "name": "wildcard-renamed"})
		assertResultJSON(t, res)
	})
}

func TestCertificateChoiceTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"certificate_acme_server_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"https://acme-v02.api.letsencrypt.org/directory": "Let's Encrypt"}, nil
			},
			"certificate_country_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"US": "United States"}, nil
			},
			"certificate_ec_curve_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"SECP384R1": "SECP384R1"}, nil
			},
			"certificate_extended_key_usage_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"SERVER_AUTH": "Server Authentication"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	for _, tool := range []string{
		"certificate_acme_server_choices",
		"certificate_country_choices",
		"certificate_ec_curve_choices",
		"certificate_extended_key_usage_choices",
	} {
		t.Run(tool, func(t *testing.T) {
			res := callTool(t, cs, tool, nil)
			assertResultJSON(t, res)
		})
	}
}
