package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestAlertMiscTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"alert_list_categories": func(_ context.Context) ([]truenas.AlertCategory, error) {
				return []truenas.AlertCategory{{ID: "ZFS", Title: "Storage"}}, nil
			},
			"alert_list_policies": func(_ context.Context) ([]string, error) {
				return []string{"IMMEDIATELY", "NEVER"}, nil
			},
			"alert_dismiss": func(_ context.Context, uuid string) error { return nil },
			"alert_restore": func(_ context.Context, uuid string) error { return nil },
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list_categories", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "alert_list_categories", nil))
	})
	t.Run("list_policies", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "alert_list_policies", nil))
	})
	t.Run("dismiss", func(t *testing.T) {
		res := callTool(t, cs, "alert_dismiss", map[string]any{"uuid": "abc-123"})
		assertResultJSON(t, res)
	})
	t.Run("dismiss requires uuid", func(t *testing.T) {
		res := callTool(t, cs, "alert_dismiss", map[string]any{"uuid": ""})
		assertError(t, res, "uuid must not be empty")
	})
	t.Run("restore", func(t *testing.T) {
		res := callTool(t, cs, "alert_restore", map[string]any{"uuid": "abc-123"})
		assertResultJSON(t, res)
	})
}

func TestAlertServiceCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"alertservice_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.AlertService, error) {
				return []truenas.AlertService{{ID: 1, Name: "ops-slack", Type: "Slack"}}, nil
			},
			"alertservice_get": func(_ context.Context, id int) (*truenas.AlertService, error) {
				return &truenas.AlertService{ID: id}, nil
			},
			"alertservice_create": func(_ context.Context, p *truenas.CreateAlertServiceParams) (*truenas.AlertService, error) {
				return &truenas.AlertService{ID: 1, Name: p.Name, Type: p.Type}, nil
			},
			"alertservice_update": func(_ context.Context, id int, _ *truenas.CreateAlertServiceParams) (*truenas.AlertService, error) {
				return &truenas.AlertService{ID: id}, nil
			},
			"alertservice_test": func(_ context.Context, p *truenas.CreateAlertServiceParams) (bool, error) {
				return true, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "alertservice_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "alertservice_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "alertservice_create", map[string]any{"name": "ops-slack", "type": "Slack"})
		assertResultJSON(t, res)
	})
	t.Run("create requires name and type", func(t *testing.T) {
		res := callTool(t, cs, "alertservice_create", map[string]any{"name": "", "type": ""})
		assertError(t, res, "name and type are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "alertservice_update", map[string]any{"id": 1, "name": "ops-slack", "type": "Slack"})
		assertResultJSON(t, res)
	})
	t.Run("test", func(t *testing.T) {
		res := callTool(t, cs, "alertservice_test", map[string]any{"name": "ops-slack", "type": "Slack"})
		assertResultJSON(t, res)
	})
}

func TestAlertClassesTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"alertclasses_config": func(_ context.Context) (*truenas.AlertClassesConfig, error) {
				return &truenas.AlertClassesConfig{Classes: map[string]any{"PoolUsage": map[string]any{"level": "CRITICAL"}}}, nil
			},
			"alertclasses_update": func(_ context.Context, p *truenas.UpdateAlertClassesParams) (*truenas.AlertClassesConfig, error) {
				return &truenas.AlertClassesConfig{Classes: p.Classes}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("config", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "alertclasses_config", nil))
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "alertclasses_update", map[string]any{
			"classes": map[string]any{"PoolUsage": map[string]any{"level": "CRITICAL"}},
		})
		assertResultJSON(t, res)
	})
}
