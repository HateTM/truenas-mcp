package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerAlertManagementTools registers alert.list_categories/list_policies/dismiss/restore,
// alertservice.*, and alertclasses.* MCP tools.
func registerAlertManagementTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "alert_list_categories",
		Description: "List the alert categories available on the server.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		categories, err := client.AlertListCategories(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("alert_list_categories: %w", err))
		}
		return jsonResult(categories)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "alert_list_policies",
		Description: "List the dismissal policies available for alerts (e.g. IMMEDIATELY, NEVER).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		policies, err := client.AlertListPolicies(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("alert_list_policies: %w", err))
		}
		return jsonResult(policies)
	})

	type alertUUIDInput struct {
		UUID string `json:"uuid" jsonschema:"Alert UUID (from list_alerts)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "alert_dismiss",
		Description: "Dismiss an alert by its UUID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p alertUUIDInput) (*mcp.CallToolResult, any, error) {
		if p.UUID == "" {
			return errorResult(errors.New("alert_dismiss: uuid must not be empty"))
		}
		if err := client.DismissAlert(ctx, p.UUID); err != nil {
			return errorResult(fmt.Errorf("alert_dismiss: %w", err))
		}
		return jsonResult(map[string]any{"dismissed": true, "uuid": p.UUID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "alert_restore",
		Description: "Restore (un-dismiss) an alert by its UUID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p alertUUIDInput) (*mcp.CallToolResult, any, error) {
		if p.UUID == "" {
			return errorResult(errors.New("alert_restore: uuid must not be empty"))
		}
		if err := client.RestoreAlert(ctx, p.UUID); err != nil {
			return errorResult(fmt.Errorf("alert_restore: %w", err))
		}
		return jsonResult(map[string]any{"restored": true, "uuid": p.UUID})
	})

	type listAlertServicesInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of services to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of services to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertservice_list",
		Description: "List configured alert notification services (e.g. email, Slack).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listAlertServicesInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListAlertServices(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("alertservice_list: %w", err))
		}
		return jsonResult(result)
	})

	type getAlertServiceInput struct {
		ID int `json:"id" jsonschema:"Numeric alert service ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertservice_get",
		Description: "Get a single alert notification service by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getAlertServiceInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("alertservice_get: id must be a positive integer"))
		}
		result, err := client.GetAlertService(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("alertservice_get: %w", err))
		}
		return jsonResult(result)
	})

	type alertServiceInput struct {
		Name       string         `json:"name"                 jsonschema:"Service name"`
		Type       string         `json:"type"                 jsonschema:"Service type, e.g. Mail, Slack, PagerDuty"`
		Attributes map[string]any `json:"attributes,omitempty" jsonschema:"Type-specific configuration fields, e.g. webhook URL for Slack"`
		Enabled    bool           `json:"enabled,omitempty"    jsonschema:"Whether the service is active"`
		Level      string         `json:"level,omitempty"      jsonschema:"Minimum alert level that triggers this service, e.g. WARNING"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertservice_create",
		Description: "Create a new alert notification service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p alertServiceInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || p.Type == "" {
			return errorResult(errors.New("alertservice_create: name and type are required"))
		}
		result, err := client.CreateAlertService(ctx, &truenas.CreateAlertServiceParams{
			Name: p.Name, Type: p.Type, Attributes: p.Attributes, Enabled: p.Enabled, Level: p.Level,
		})
		if err != nil {
			return errorResult(fmt.Errorf("alertservice_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateAlertServiceInput struct {
		ID         int            `json:"id"                   jsonschema:"Numeric alert service ID"`
		Name       string         `json:"name"                 jsonschema:"Service name"`
		Type       string         `json:"type"                 jsonschema:"Service type"`
		Attributes map[string]any `json:"attributes,omitempty" jsonschema:"Type-specific configuration fields"`
		Enabled    bool           `json:"enabled,omitempty"    jsonschema:"Whether the service is active"`
		Level      string         `json:"level,omitempty"      jsonschema:"Minimum alert level that triggers this service"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertservice_update",
		Description: "Update an existing alert notification service.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateAlertServiceInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("alertservice_update: id must be a positive integer"))
		}
		if p.Name == "" || p.Type == "" {
			return errorResult(errors.New("alertservice_update: name and type are required"))
		}
		result, err := client.UpdateAlertService(ctx, p.ID, &truenas.CreateAlertServiceParams{
			Name: p.Name, Type: p.Type, Attributes: p.Attributes, Enabled: p.Enabled, Level: p.Level,
		})
		if err != nil {
			return errorResult(fmt.Errorf("alertservice_update: %w", err))
		}
		return jsonResult(result)
	})

	// alertservice_delete is destructive and lives in destructive_alert.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertservice_test",
		Description: "Send a test notification through an alert service configuration, without saving it.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p alertServiceInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || p.Type == "" {
			return errorResult(errors.New("alertservice_test: name and type are required"))
		}
		ok, err := client.TestAlertService(ctx, &truenas.CreateAlertServiceParams{
			Name: p.Name, Type: p.Type, Attributes: p.Attributes, Enabled: p.Enabled, Level: p.Level,
		})
		if err != nil {
			return errorResult(fmt.Errorf("alertservice_test: %w", err))
		}
		return jsonResult(map[string]bool{"success": ok})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertclasses_config",
		Description: "Get per-alert-class severity/policy overrides.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.AlertClassesConfigGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("alertclasses_config: %w", err))
		}
		return jsonResult(cfg)
	})

	type updateAlertClassesInput struct {
		Classes map[string]any `json:"classes,omitempty" jsonschema:"Per-class overrides keyed by alert class name, e.g. {\"PoolUsage\": {\"level\": \"CRITICAL\", \"policy\": \"IMMEDIATELY\"}}"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "alertclasses_update",
		Description: "Update per-alert-class severity/policy overrides.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateAlertClassesInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.UpdateAlertClasses(ctx, &truenas.UpdateAlertClassesParams{Classes: p.Classes})
		if err != nil {
			return errorResult(fmt.Errorf("alertclasses_update: %w", err))
		}
		return jsonResult(cfg)
	})
}
