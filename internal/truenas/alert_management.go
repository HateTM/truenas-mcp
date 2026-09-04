package truenas

import (
	"context"
	"errors"
	"fmt"
)

// AlertCategory groups related alert classes for display/filtering purposes.
type AlertCategory struct {
	ID      string   `json:"id"`
	Title   string   `json:"title,omitempty"`
	Classes []string `json:"classes,omitempty"`
}

// AlertService represents a configured alert notification target (e.g. email, Slack).
type AlertService struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes,omitempty"`
	Enabled    bool           `json:"enabled,omitempty"`
	Level      string         `json:"level,omitempty"`
}

// CreateAlertServiceParams holds fields for creating, updating, or test-sending an alert service.
type CreateAlertServiceParams struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes,omitempty"`
	Enabled    bool           `json:"enabled,omitempty"`
	Level      string         `json:"level,omitempty"`
}

// AlertClassesConfig holds per-alert-class severity/policy overrides.
type AlertClassesConfig struct {
	Classes map[string]any `json:"classes,omitempty"`
}

// UpdateAlertClassesParams holds fields for updating per-alert-class overrides.
type UpdateAlertClassesParams struct {
	Classes map[string]any `json:"classes,omitempty"`
}

// AlertListCategories returns the alert categories available on the server.
func (c *Client) AlertListCategories(ctx context.Context) ([]AlertCategory, error) {
	var categories []AlertCategory
	if err := c.call(ctx, "alert.list_categories", nil, &categories); err != nil {
		return nil, fmt.Errorf("listing alert categories: %w", err)
	}
	return categories, nil
}

// AlertListPolicies returns the dismissal policies available for alerts (e.g. IMMEDIATELY, NEVER).
func (c *Client) AlertListPolicies(ctx context.Context) ([]string, error) {
	var policies []string
	if err := c.call(ctx, "alert.list_policies", nil, &policies); err != nil {
		return nil, fmt.Errorf("listing alert policies: %w", err)
	}
	return policies, nil
}

// DismissAlert dismisses an alert by its UUID.
func (c *Client) DismissAlert(ctx context.Context, uuid string) error {
	if uuid == "" {
		return errors.New("dismiss alert: uuid must not be empty")
	}
	if err := c.call(ctx, "alert.dismiss", []any{uuid}, nil); err != nil {
		return fmt.Errorf("dismissing alert %q: %w", uuid, err)
	}
	return nil
}

// RestoreAlert restores (un-dismisses) an alert by its UUID.
func (c *Client) RestoreAlert(ctx context.Context, uuid string) error {
	if uuid == "" {
		return errors.New("restore alert: uuid must not be empty")
	}
	if err := c.call(ctx, "alert.restore", []any{uuid}, nil); err != nil {
		return fmt.Errorf("restoring alert %q: %w", uuid, err)
	}
	return nil
}

// ListAlertServices lists configured alert notification services.
func (c *Client) ListAlertServices(ctx context.Context, opts ...ListOptions) ([]AlertService, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []AlertService
	if err := c.call(ctx, "alertservice.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing alert services: %w", err)
	}
	return result, nil
}

// GetAlertService returns a single alert notification service by ID.
func (c *Client) GetAlertService(ctx context.Context, id int) (*AlertService, error) {
	var result AlertService
	if err := c.call(ctx, "alertservice.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting alert service %d: %w", id, err)
	}
	return &result, nil
}

// CreateAlertService creates a new alert notification service.
func (c *Client) CreateAlertService(ctx context.Context, p *CreateAlertServiceParams) (*AlertService, error) {
	if p == nil {
		return nil, errors.New("create alert service: params required")
	}
	var result AlertService
	if err := c.call(ctx, "alertservice.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating alert service %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateAlertService updates an existing alert notification service.
func (c *Client) UpdateAlertService(ctx context.Context, id int, p *CreateAlertServiceParams) (*AlertService, error) {
	if p == nil {
		return nil, errors.New("update alert service: params required")
	}
	var result AlertService
	if err := c.call(ctx, "alertservice.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating alert service %d: %w", id, err)
	}
	return &result, nil
}

// DeleteAlertService deletes an alert notification service.
func (c *Client) DeleteAlertService(ctx context.Context, id int) error {
	if err := c.call(ctx, "alertservice.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting alert service %d: %w", id, err)
	}
	return nil
}

// TestAlertService sends a test notification through an (unsaved) alert service configuration.
func (c *Client) TestAlertService(ctx context.Context, p *CreateAlertServiceParams) (bool, error) {
	if p == nil {
		return false, errors.New("test alert service: params required")
	}
	var ok bool
	if err := c.call(ctx, "alertservice.test", []any{p}, &ok); err != nil {
		return false, fmt.Errorf("testing alert service %q: %w", p.Name, err)
	}
	return ok, nil
}

// AlertClassesConfigGet returns per-alert-class severity/policy overrides.
func (c *Client) AlertClassesConfigGet(ctx context.Context) (*AlertClassesConfig, error) {
	var cfg AlertClassesConfig
	if err := c.call(ctx, "alertclasses.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting alert classes config: %w", err)
	}
	return &cfg, nil
}

// UpdateAlertClasses updates per-alert-class severity/policy overrides.
func (c *Client) UpdateAlertClasses(ctx context.Context, p *UpdateAlertClassesParams) (*AlertClassesConfig, error) {
	if p == nil {
		return nil, errors.New("update alert classes: params required")
	}
	var cfg AlertClassesConfig
	if err := c.call(ctx, "alertclasses.update", []any{p}, &cfg); err != nil {
		return nil, fmt.Errorf("updating alert classes config: %w", err)
	}
	return &cfg, nil
}
