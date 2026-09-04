package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestListApps(t *testing.T) {
	t.Run("returns apps as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_apps": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.App, error) {
					return []truenas.App{{Name: "jellyfin"}}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_apps", nil)
		assertResultJSON(t, res)
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_apps": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.App, error) {
					return nil, errors.New("API error")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_apps", nil)
		assertError(t, res, "API error")
	})
}

func TestGetApp(t *testing.T) {
	t.Run("returns app as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"get_app": func(_ context.Context, name string) (*truenas.App, error) {
					return &truenas.App{Name: name}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "get_app", map[string]any{"name": "jellyfin"})
		assertResultJSON(t, res)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "get_app", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"get_app": func(_ context.Context, _ string) (*truenas.App, error) {
					return nil, errors.New("app not found")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "get_app", map[string]any{"name": "noexist"})
		assertError(t, res, "app not found")
	})
}

func TestStartApp(t *testing.T) {
	t.Run("returns success as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"start_app": func(_ context.Context, _ string) error {
					return nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "start_app", map[string]any{"name": "jellyfin"})
		assertResultJSON(t, res)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "start_app", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
}

func TestStopApp(t *testing.T) {
	t.Run("returns success as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"stop_app": func(_ context.Context, _ string) error {
					return nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "stop_app", map[string]any{"name": "jellyfin"})
		assertResultJSON(t, res)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "stop_app", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
}

func TestRestartApp(t *testing.T) {
	t.Run("returns app entry as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"restart_app": func(_ context.Context, name string) (*truenas.App, error) {
					return &truenas.App{Name: name, State: "DEPLOYING"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "restart_app", map[string]any{"name": "jellyfin"})
		assertResultJSON(t, res)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "restart_app", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
}

func TestListImages(t *testing.T) {
	t.Run("returns images as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_images": func(_ context.Context) ([]truenas.Image, error) {
					return []truenas.Image{{ID: "jellyfin/jellyfin:latest"}}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_images", nil)
		assertResultJSON(t, res)
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"list_images": func(_ context.Context) ([]truenas.Image, error) {
					return nil, errors.New("docker not available")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "list_images", nil)
		assertError(t, res, "docker not available")
	})
}

func TestInstallApp(t *testing.T) {
	t.Run("returns app entry as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"create_app": func(_ context.Context, params *truenas.CreateAppParams) (*truenas.App, error) {
					return &truenas.App{Name: params.AppName, State: "DEPLOYING"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "install_app", map[string]any{
			"app_name":    "my-jellyfin",
			"catalog_app": "jellyfin",
		})
		assertResultJSON(t, res)
	})

	t.Run("empty app_name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "install_app", map[string]any{"app_name": "", "catalog_app": "jellyfin"})
		assertError(t, res, "app_name must not be empty")
	})

	t.Run("empty catalog_app returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "install_app", map[string]any{"app_name": "my-app", "catalog_app": ""})
		assertError(t, res, "catalog_app must not be empty")
	})
}

func TestInstallCustomApp(t *testing.T) {
	t.Run("returns app entry as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"create_app": func(_ context.Context, params *truenas.CreateAppParams) (*truenas.App, error) {
					return &truenas.App{Name: params.AppName, State: "DEPLOYING"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "install_custom_app", map[string]any{
			"app_name":                     "my-app",
			"custom_compose_config_string": "services:\n  web:\n    image: nginx\n",
		})
		assertResultJSON(t, res)
	})

	t.Run("empty app_name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "install_custom_app", map[string]any{
			"app_name":                     "",
			"custom_compose_config_string": "services: {}",
		})
		assertError(t, res, "app_name must not be empty")
	})
}

func TestUpgradeApp(t *testing.T) {
	t.Run("returns app entry as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"upgrade_app": func(_ context.Context, name, _ string) (*truenas.App, error) {
					return &truenas.App{Name: name, State: "DEPLOYING"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "upgrade_app", map[string]any{"name": "jellyfin"})
		assertResultJSON(t, res)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "upgrade_app", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
}

func TestUpgradeSummary(t *testing.T) {
	t.Run("returns summary as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"get_upgrade_summary": func(_ context.Context, name string) (*truenas.AppUpgradeSummary, error) {
					return &truenas.AppUpgradeSummary{LatestVersion: "1.2.0"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "upgrade_summary", map[string]any{"name": "jellyfin"})
		assertResultJSON(t, res)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "upgrade_summary", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
}

func TestRollbackApp(t *testing.T) {
	t.Run("returns app entry as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"rollback_app": func(_ context.Context, name, _ string) (*truenas.App, error) {
					return &truenas.App{Name: name, State: "DEPLOYING"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "rollback_app", map[string]any{"name": "jellyfin", "version": "1.1.0"})
		assertResultJSON(t, res)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "rollback_app", map[string]any{"name": "", "version": "1.1.0"})
		assertError(t, res, "name must not be empty")
	})

	t.Run("empty version returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "rollback_app", map[string]any{"name": "jellyfin", "version": ""})
		assertError(t, res, "version must not be empty")
	})
}
