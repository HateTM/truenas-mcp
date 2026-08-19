// Package tools — auth.go intentionally omits tools for auth.login(), auth.login_ex(),
// auth.login_ex_continue(), and auth.login_with_token(): see the doc comment on
// internal/truenas/auth.go for the rationale (identity-spoofing risk, no legitimate use
// case for this server switching identity mid-session).
package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerAuthTools registers auth.* MCP tools onto the server.
func registerAuthTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_generate_onetime_password",
		Description: "Generate a one-time password for the current user.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		otp, err := client.GenerateOnetimePassword(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_generate_onetime_password: %w", err))
		}
		return jsonResult(map[string]string{"otp": otp})
	})

	type generateTokenInput struct {
		TTL int `json:"ttl,omitempty" jsonschema:"Token lifetime in seconds; 0 uses the server default"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_generate_token",
		Description: "Generate an API token for the current session.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p generateTokenInput) (*mcp.CallToolResult, any, error) {
		token, err := client.GenerateToken(ctx, p.TTL)
		if err != nil {
			return errorResult(fmt.Errorf("auth_generate_token: %w", err))
		}
		return jsonResult(map[string]string{"token": token})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_logout",
		Description: "Terminate the current session.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		ok, err := client.Logout(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_logout: %w", err))
		}
		return jsonResult(map[string]bool{"success": ok})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_me",
		Description: "Get the identity TrueNAS associates with the current session.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		identity, err := client.Me(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_me: %w", err))
		}
		return jsonResult(identity)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_mechanism_choices",
		Description: "List the login mechanisms supported by the server.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.AuthMechanismChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_mechanism_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_sessions",
		Description: "List all currently authenticated sessions on the server.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		sessions, err := client.AuthSessionsList(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_sessions: %w", err))
		}
		return jsonResult(sessions)
	})

	type setAuthAttributeInput struct {
		Key   string `json:"key"   jsonschema:"Attribute name"`
		Value any    `json:"value" jsonschema:"Attribute value (any JSON type)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_set_attribute",
		Description: "Set a key/value pair in the current user's session attributes.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p setAuthAttributeInput) (*mcp.CallToolResult, any, error) {
		if p.Key == "" {
			return errorResult(errors.New("auth_set_attribute: key must not be empty"))
		}
		ok, err := client.SetAuthAttribute(ctx, p.Key, p.Value)
		if err != nil {
			return errorResult(fmt.Errorf("auth_set_attribute: %w", err))
		}
		return jsonResult(map[string]bool{"success": ok})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_terminate_other_sessions",
		Description: "Terminate every session other than the current one.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		ok, err := client.TerminateOtherSessions(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_terminate_other_sessions: %w", err))
		}
		return jsonResult(map[string]bool{"success": ok})
	})

	type terminateSessionInput struct {
		ID string `json:"id" jsonschema:"Session ID (from auth_sessions)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_terminate_session",
		Description: "Terminate a specific session by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p terminateSessionInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("auth_terminate_session: id must not be empty"))
		}
		ok, err := client.TerminateSession(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("auth_terminate_session: %w", err))
		}
		return jsonResult(map[string]bool{"success": ok})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_twofactor",
		Description: "Check whether two-factor authentication is required for the current user.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		enabled, err := client.TwoFactorEnabled(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_twofactor: %w", err))
		}
		return jsonResult(map[string]bool{"enabled": enabled})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_twofactor_config",
		Description: "Get the two-factor authentication configuration for the current user.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.TwoFactorConfigGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("auth_twofactor_config: %w", err))
		}
		return jsonResult(cfg)
	})

	type updateTwoFactorInput struct {
		Enabled   bool `json:"enabled,omitempty"    jsonschema:"Require two-factor authentication for this user"`
		Interval  int  `json:"interval,omitempty"   jsonschema:"OTP validity window in seconds"`
		OTPDigits int  `json:"otp_digits,omitempty" jsonschema:"Number of digits in generated OTP codes"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "auth_twofactor_update",
		Description: "Update the two-factor authentication configuration for the current user.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateTwoFactorInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.UpdateTwoFactor(ctx, &truenas.UpdateTwoFactorParams{
			Enabled: p.Enabled, Interval: p.Interval, OTPDigits: p.OTPDigits,
		})
		if err != nil {
			return errorResult(fmt.Errorf("auth_twofactor_update: %w", err))
		}
		return jsonResult(cfg)
	})
}
