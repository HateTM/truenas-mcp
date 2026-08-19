package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestAppRegistryCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"app_registry_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.AppRegistry, error) {
				return []truenas.AppRegistry{{ID: 1, Name: "ghcr", URI: "https://ghcr.io"}}, nil
			},
			"app_registry_get": func(_ context.Context, id int) (*truenas.AppRegistry, error) {
				return &truenas.AppRegistry{ID: id}, nil
			},
			"app_registry_create": func(_ context.Context, p *truenas.CreateAppRegistryParams) (*truenas.AppRegistry, error) {
				return &truenas.AppRegistry{ID: 1, Name: p.Name, URI: p.URI}, nil
			},
			"app_registry_update": func(_ context.Context, id int, _ *truenas.CreateAppRegistryParams) (*truenas.AppRegistry, error) {
				return &truenas.AppRegistry{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "app_registry_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "app_registry_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "app_registry_create", map[string]any{"name": "ghcr", "uri": "https://ghcr.io"})
		assertResultJSON(t, res)
	})
	t.Run("create requires name and uri", func(t *testing.T) {
		res := callTool(t, cs, "app_registry_create", map[string]any{"name": "", "uri": ""})
		assertError(t, res, "name and uri are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "app_registry_update", map[string]any{"id": 1, "name": "ghcr", "uri": "https://ghcr.io"})
		assertResultJSON(t, res)
	})
}

func TestAppImageTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"app_image_pull": func(_ context.Context, p *truenas.PullAppImageParams) (int, error) {
				if p.Repository != "library/nginx" {
					t.Errorf("Repository = %q, want library/nginx", p.Repository)
				}
				return 42, nil
			},
			"app_image_dockerhub_rate_limit": func(_ context.Context) (*truenas.DockerHubRateLimit, error) {
				return &truenas.DockerHubRateLimit{TotalLimit: 100, RemainingLimit: 50}, nil
			},
			"app_image_get": func(_ context.Context, id string) (*truenas.Image, error) {
				return &truenas.Image{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("pull", func(t *testing.T) {
		res := callTool(t, cs, "app_image_pull", map[string]any{"repository": "library/nginx"})
		assertResultJSON(t, res)
	})
	t.Run("pull requires repository", func(t *testing.T) {
		res := callTool(t, cs, "app_image_pull", map[string]any{"repository": ""})
		assertError(t, res, "repository must not be empty")
	})
	t.Run("dockerhub_rate_limit", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "app_image_dockerhub_rate_limit", nil))
	})
	t.Run("get", func(t *testing.T) {
		res := callTool(t, cs, "app_image_get", map[string]any{"id": "sha256:abc"})
		assertResultJSON(t, res)
	})
	t.Run("get requires id", func(t *testing.T) {
		res := callTool(t, cs, "app_image_get", map[string]any{"id": ""})
		assertError(t, res, "id must not be empty")
	})
}

func TestAppMetadataTools(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"app_categories": func(_ context.Context) ([]string, error) {
				return []string{"media", "networking"}, nil
			},
			"app_available_space": func(_ context.Context) (*truenas.AppAvailableSpace, error) {
				return &truenas.AppAvailableSpace{FreeBytes: 1024}, nil
			},
			"app_config": func(_ context.Context) (*truenas.AppGlobalConfig, error) {
				return &truenas.AppGlobalConfig{Pool: "Storage"}, nil
			},
			"app_container_ids": func(_ context.Context, appName string) ([]string, error) {
				return []string{appName + "-container-1"}, nil
			},
			"app_convert_to_custom": func(_ context.Context, appName string) (*truenas.App, error) {
				return &truenas.App{Name: appName}, nil
			},
			"app_outdated_docker_images": func(_ context.Context, _ string) ([]string, error) {
				return []string{"library/nginx:1.24"}, nil
			},
			"app_pull_images": func(_ context.Context, _ string) (int, error) {
				return 7, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("categories", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "app_categories", nil))
	})
	t.Run("available_space", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "app_available_space", nil))
	})
	t.Run("config", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "app_config", nil))
	})
	t.Run("container_ids", func(t *testing.T) {
		res := callTool(t, cs, "app_container_ids", map[string]any{"app_name": "jellyfin"})
		assertResultJSON(t, res)
	})
	t.Run("container_ids requires app_name", func(t *testing.T) {
		res := callTool(t, cs, "app_container_ids", map[string]any{"app_name": ""})
		assertError(t, res, "app_name must not be empty")
	})
	t.Run("convert_to_custom", func(t *testing.T) {
		res := callTool(t, cs, "app_convert_to_custom", map[string]any{"app_name": "jellyfin"})
		assertResultJSON(t, res)
	})
	t.Run("outdated_docker_images", func(t *testing.T) {
		res := callTool(t, cs, "app_outdated_docker_images", map[string]any{"app_name": "jellyfin"})
		assertResultJSON(t, res)
	})
	t.Run("pull_images", func(t *testing.T) {
		res := callTool(t, cs, "app_pull_images", map[string]any{"app_name": "jellyfin"})
		assertResultJSON(t, res)
	})
}
