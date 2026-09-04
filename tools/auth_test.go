package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestAuthTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"auth_generate_onetime_password": func(_ context.Context) (string, error) { return "123456", nil },
			"auth_generate_token": func(_ context.Context, ttl int) (string, error) {
				return "tok-abc", nil
			},
			"auth_logout": func(_ context.Context) (bool, error) { return true, nil },
			"auth_me": func(_ context.Context) (*truenas.UserIdentity, error) {
				return &truenas.UserIdentity{Username: "admin"}, nil
			},
			"auth_mechanism_choices": func(_ context.Context) ([]string, error) {
				return []string{"PASSWORD_PLAIN", "API_KEY_PLAIN"}, nil
			},
			"auth_sessions": func(_ context.Context) ([]truenas.AuthSession, error) {
				return []truenas.AuthSession{{ID: "sess-1", Current: true}}, nil
			},
			"auth_set_attribute": func(_ context.Context, key string, value any) (bool, error) {
				if key != "theme" {
					t.Errorf("key = %q, want theme", key)
				}
				return true, nil
			},
			"auth_terminate_other_sessions": func(_ context.Context) (bool, error) { return true, nil },
			"auth_terminate_session": func(_ context.Context, id string) (bool, error) {
				return true, nil
			},
			"auth_twofactor": func(_ context.Context) (bool, error) { return true, nil },
			"auth_twofactor_config": func(_ context.Context) (*truenas.TwoFactorConfig, error) {
				return &truenas.TwoFactorConfig{Enabled: true, OTPDigits: 6}, nil
			},
			"auth_twofactor_update": func(_ context.Context, p *truenas.UpdateTwoFactorParams) (*truenas.TwoFactorConfig, error) {
				return &truenas.TwoFactorConfig{Enabled: p.Enabled}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	for _, tool := range []string{
		"auth_generate_onetime_password",
		"auth_logout",
		"auth_me",
		"auth_mechanism_choices",
		"auth_sessions",
		"auth_terminate_other_sessions",
		"auth_twofactor",
		"auth_twofactor_config",
	} {
		t.Run(tool, func(t *testing.T) {
			res := callTool(t, cs, tool, nil)
			assertResultJSON(t, res)
		})
	}

	t.Run("auth_generate_token", func(t *testing.T) {
		res := callTool(t, cs, "auth_generate_token", map[string]any{"ttl": 3600})
		assertResultJSON(t, res)
	})

	t.Run("auth_set_attribute", func(t *testing.T) {
		res := callTool(t, cs, "auth_set_attribute", map[string]any{"key": "theme", "value": "dark"})
		assertResultJSON(t, res)
	})
	t.Run("auth_set_attribute requires key", func(t *testing.T) {
		res := callTool(t, cs, "auth_set_attribute", map[string]any{"key": "", "value": "dark"})
		assertError(t, res, "key must not be empty")
	})

	t.Run("auth_terminate_session", func(t *testing.T) {
		res := callTool(t, cs, "auth_terminate_session", map[string]any{"id": "sess-1"})
		assertResultJSON(t, res)
	})
	t.Run("auth_terminate_session requires id", func(t *testing.T) {
		res := callTool(t, cs, "auth_terminate_session", map[string]any{"id": ""})
		assertError(t, res, "id must not be empty")
	})

	t.Run("auth_twofactor_update", func(t *testing.T) {
		res := callTool(t, cs, "auth_twofactor_update", map[string]any{"enabled": true, "otp_digits": 6})
		assertResultJSON(t, res)
	})
}
