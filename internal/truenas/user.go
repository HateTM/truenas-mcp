package truenas

import (
	"context"
	"errors"
	"fmt"
)

// User represents a local TrueNAS user account.
type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	FullName      string `json:"full_name,omitempty"`
	UID           int    `json:"uid,omitempty"`
	GroupID       int    `json:"group,omitempty"`
	HomeDirectory string `json:"home,omitempty"`
	Shell         string `json:"shell,omitempty"`
	Email         string `json:"email,omitempty"`
	Locked        bool   `json:"locked,omitempty"`
	SMB           bool   `json:"smb,omitempty"`
	SSHPubKey     string `json:"sshpubkey,omitempty"`
}

// CreateUserParams holds fields for creating or updating a local user.
type CreateUserParams struct {
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	GroupID       int    `json:"group,omitempty"`
	Password      string `json:"password,omitempty"`
	HomeDirectory string `json:"home,omitempty"`
	Shell         string `json:"shell,omitempty"`
	Email         string `json:"email,omitempty"`
	Locked        bool   `json:"locked,omitempty"`
	SMB           bool   `json:"smb,omitempty"`
	SSHPubKey     string `json:"sshpubkey,omitempty"`
}

// ListUsers lists local user accounts.
func (c *Client) ListUsers(ctx context.Context, opts ...ListOptions) ([]User, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []User
	if err := c.call(ctx, "user.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	return result, nil
}

// GetUser returns a single local user account by ID.
func (c *Client) GetUser(ctx context.Context, id int) (*User, error) {
	var result User
	if err := c.call(ctx, "user.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}
	return &result, nil
}

// CreateUser creates a new local user account.
func (c *Client) CreateUser(ctx context.Context, p *CreateUserParams) (*User, error) {
	if p == nil {
		return nil, errors.New("create user: params required")
	}
	var result User
	if err := c.call(ctx, "user.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating user %q: %w", p.Username, err)
	}
	return &result, nil
}

// UpdateUser updates an existing local user account's non-password fields.
func (c *Client) UpdateUser(ctx context.Context, id int, p *CreateUserParams) (*User, error) {
	if p == nil {
		return nil, errors.New("update user: params required")
	}
	var result User
	if err := c.call(ctx, "user.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating user %d: %w", id, err)
	}
	return &result, nil
}

// DeleteUser deletes a local user account.
func (c *Client) DeleteUser(ctx context.Context, id int) error {
	if err := c.call(ctx, "user.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting user %d: %w", id, err)
	}
	return nil
}

// UpdateUserPassword sets a new password for a local user account.
func (c *Client) UpdateUserPassword(ctx context.Context, id int, password string) error {
	if password == "" {
		return errors.New("update user password: password must not be empty")
	}
	if err := c.call(ctx, "user.update_password", []any{id, map[string]string{"new_password": password}}, nil); err != nil {
		return fmt.Errorf("updating password for user %d: %w", id, err)
	}
	return nil
}

// UserSchemas returns the JSON schemas describing valid user account configurations.
func (c *Client) UserSchemas(ctx context.Context) (map[string]any, error) {
	var schemas map[string]any
	if err := c.call(ctx, "user.user_schemas", nil, &schemas); err != nil {
		return nil, fmt.Errorf("getting user schemas: %w", err)
	}
	return schemas, nil
}

// UserHomeDirectoryChoices returns the paths available as a new user's home directory.
func (c *Client) UserHomeDirectoryChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "user.home_directory_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing user home directory choices: %w", err)
	}
	return choices, nil
}

// UserShellChoices returns the shells available for a user account.
func (c *Client) UserShellChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "user.shell_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing user shell choices: %w", err)
	}
	return choices, nil
}
