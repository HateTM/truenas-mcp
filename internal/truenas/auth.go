// Package truenas — auth.go intentionally does NOT wrap auth.login(), auth.login_ex(),
// auth.login_ex_continue(), or auth.login_with_token(): all four let a caller authenticate
// as an arbitrary identity, and this MCP server already authenticates once via API key at
// Connect() (see auth.login_with_api_key usage in client.go). There is no legitimate use
// case for the server switching identity mid-session, so exposing them as callable tools
// would only add an identity-spoofing risk. This mirrors the existing decision not to
// expose auth.login_with_api_key as a tool.
package truenas

import (
	"context"
	"errors"
	"fmt"
)

// UserIdentity represents the identity TrueNAS associates with the current session (auth.me).
type UserIdentity struct {
	Username          string   `json:"username"`
	AccountAttributes []string `json:"account_attributes,omitempty"`
}

// AuthSession represents an authenticated session on the TrueNAS host.
type AuthSession struct {
	ID          string `json:"id"`
	Internal    bool   `json:"internal,omitempty"`
	Origin      string `json:"origin,omitempty"`
	Credentials string `json:"credentials,omitempty"`
	Current     bool   `json:"current,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// TwoFactorConfig represents the two-factor authentication configuration for the current user.
type TwoFactorConfig struct {
	Enabled   bool `json:"enabled,omitempty"`
	Interval  int  `json:"interval,omitempty"`
	OTPDigits int  `json:"otp_digits,omitempty"`
}

// UpdateTwoFactorParams holds fields for updating two-factor authentication configuration.
type UpdateTwoFactorParams struct {
	Enabled   bool `json:"enabled,omitempty"`
	Interval  int  `json:"interval,omitempty"`
	OTPDigits int  `json:"otp_digits,omitempty"`
}

// GenerateOnetimePassword generates a one-time password for the current user.
func (c *Client) GenerateOnetimePassword(ctx context.Context) (string, error) {
	var otp string
	if err := c.call(ctx, "auth.generate_onetime_password", nil, &otp); err != nil {
		return "", fmt.Errorf("generating onetime password: %w", err)
	}
	return otp, nil
}

// GenerateToken generates an API token for the current session, valid for ttl seconds
// (0 uses the server default).
func (c *Client) GenerateToken(ctx context.Context, ttl int) (string, error) {
	var token string
	if err := c.call(ctx, "auth.generate_token", []any{map[string]any{"ttl": ttl}}, &token); err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}
	return token, nil
}

// Logout terminates the current session.
func (c *Client) Logout(ctx context.Context) (bool, error) {
	var ok bool
	if err := c.call(ctx, "auth.logout", nil, &ok); err != nil {
		return false, fmt.Errorf("logging out: %w", err)
	}
	return ok, nil
}

// Me returns the identity TrueNAS associates with the current session.
func (c *Client) Me(ctx context.Context) (*UserIdentity, error) {
	var identity UserIdentity
	if err := c.call(ctx, "auth.me", nil, &identity); err != nil {
		return nil, fmt.Errorf("getting current identity: %w", err)
	}
	return &identity, nil
}

// AuthMechanismChoices returns the login mechanisms supported by the server.
func (c *Client) AuthMechanismChoices(ctx context.Context) ([]string, error) {
	var choices []string
	if err := c.call(ctx, "auth.mechanism_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing auth mechanism choices: %w", err)
	}
	return choices, nil
}

// AuthSessionsList returns all currently authenticated sessions.
func (c *Client) AuthSessionsList(ctx context.Context) ([]AuthSession, error) {
	var sessions []AuthSession
	if err := c.call(ctx, "auth.sessions", nil, &sessions); err != nil {
		return nil, fmt.Errorf("listing auth sessions: %w", err)
	}
	return sessions, nil
}

// SetAuthAttribute sets a key/value pair in the current user's session attributes.
func (c *Client) SetAuthAttribute(ctx context.Context, key string, value any) (bool, error) {
	if key == "" {
		return false, errors.New("set auth attribute: key must not be empty")
	}
	var ok bool
	if err := c.call(ctx, "auth.set_attribute", []any{key, value}, &ok); err != nil {
		return false, fmt.Errorf("setting auth attribute %q: %w", key, err)
	}
	return ok, nil
}

// TerminateOtherSessions terminates every session other than the current one.
func (c *Client) TerminateOtherSessions(ctx context.Context) (bool, error) {
	var ok bool
	if err := c.call(ctx, "auth.terminate_other_sessions", nil, &ok); err != nil {
		return false, fmt.Errorf("terminating other sessions: %w", err)
	}
	return ok, nil
}

// TerminateSession terminates a specific session by ID.
func (c *Client) TerminateSession(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("terminate session: id must not be empty")
	}
	var ok bool
	if err := c.call(ctx, "auth.terminate_session", []any{id}, &ok); err != nil {
		return false, fmt.Errorf("terminating session %q: %w", id, err)
	}
	return ok, nil
}

// TwoFactorEnabled reports whether two-factor authentication is required for the current user.
func (c *Client) TwoFactorEnabled(ctx context.Context) (bool, error) {
	var enabled bool
	if err := c.call(ctx, "auth.twofactor", nil, &enabled); err != nil {
		return false, fmt.Errorf("checking two-factor status: %w", err)
	}
	return enabled, nil
}

// TwoFactorConfigGet returns the two-factor authentication configuration for the current user.
func (c *Client) TwoFactorConfigGet(ctx context.Context) (*TwoFactorConfig, error) {
	var cfg TwoFactorConfig
	if err := c.call(ctx, "auth.twofactor.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting two-factor config: %w", err)
	}
	return &cfg, nil
}

// UpdateTwoFactor updates the two-factor authentication configuration for the current user.
func (c *Client) UpdateTwoFactor(ctx context.Context, p *UpdateTwoFactorParams) (*TwoFactorConfig, error) {
	if p == nil {
		return nil, errors.New("update two-factor config: params required")
	}
	var cfg TwoFactorConfig
	if err := c.call(ctx, "auth.twofactor.update", []any{p}, &cfg); err != nil {
		return nil, fmt.Errorf("updating two-factor config: %w", err)
	}
	return &cfg, nil
}
