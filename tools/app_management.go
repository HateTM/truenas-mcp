package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerAppManagementTools registers app.registry.*, app.image.*, and remaining app.* MCP tools.
func registerAppManagementTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listAppRegistriesInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of registries to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of registries to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_registry_list",
		Description: "List configured private container registries used by the App catalog.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listAppRegistriesInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListAppRegistries(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("app_registry_list: %w", err))
		}
		return jsonResult(result)
	})

	type getAppRegistryInput struct {
		ID int `json:"id" jsonschema:"Numeric registry ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_registry_get",
		Description: "Get a single private container registry by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getAppRegistryInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("app_registry_get: id must be a positive integer"))
		}
		result, err := client.GetAppRegistry(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("app_registry_get: %w", err))
		}
		return jsonResult(result)
	})

	type appRegistryInput struct {
		Name     string `json:"name"               jsonschema:"Registry name"`
		URI      string `json:"uri"                jsonschema:"Registry URI, e.g. https://ghcr.io"`
		Username string `json:"username,omitempty" jsonschema:"Registry login username"`
		Password string `json:"password,omitempty" jsonschema:"Registry login password or access token"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_registry_create",
		Description: "Create a new private container registry credential for the App catalog.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p appRegistryInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || p.URI == "" {
			return errorResult(errors.New("app_registry_create: name and uri are required"))
		}
		result, err := client.CreateAppRegistry(ctx, &truenas.CreateAppRegistryParams{
			Name: p.Name, URI: p.URI, Username: p.Username, Password: p.Password,
		})
		if err != nil {
			return errorResult(fmt.Errorf("app_registry_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateAppRegistryInput struct {
		ID       int    `json:"id"                 jsonschema:"Numeric registry ID"`
		Name     string `json:"name"               jsonschema:"Registry name"`
		URI      string `json:"uri"                jsonschema:"Registry URI"`
		Username string `json:"username,omitempty" jsonschema:"Registry login username"`
		Password string `json:"password,omitempty" jsonschema:"Registry login password or access token"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_registry_update",
		Description: "Update an existing private container registry credential.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateAppRegistryInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("app_registry_update: id must be a positive integer"))
		}
		if p.Name == "" || p.URI == "" {
			return errorResult(errors.New("app_registry_update: name and uri are required"))
		}
		result, err := client.UpdateAppRegistry(ctx, p.ID, &truenas.CreateAppRegistryParams{
			Name: p.Name, URI: p.URI, Username: p.Username, Password: p.Password,
		})
		if err != nil {
			return errorResult(fmt.Errorf("app_registry_update: %w", err))
		}
		return jsonResult(result)
	})

	// app_registry_delete is destructive and lives in destructive_app.go, gated
	// behind Config.AllowDestructive.

	type pullAppImageInput struct {
		Repository string `json:"repository"           jsonschema:"Image repository, e.g. library/nginx"`
		Tag        string `json:"tag,omitempty"        jsonschema:"Image tag; defaults to latest"`
		RegistryID int    `json:"registry_id,omitempty" jsonschema:"Private registry ID (from app_registry_create) to pull through; omit for public images"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_image_pull",
		Description: "Pull a container image, optionally through a configured private registry. Returns the async job ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p pullAppImageInput) (*mcp.CallToolResult, any, error) {
		if p.Repository == "" {
			return errorResult(errors.New("app_image_pull: repository must not be empty"))
		}
		jobID, err := client.PullAppImage(ctx, &truenas.PullAppImageParams{
			Repository: p.Repository, Tag: p.Tag, RegistryID: p.RegistryID,
		})
		if err != nil {
			return errorResult(fmt.Errorf("app_image_pull: %w", err))
		}
		return jsonResult(map[string]int{"job_id": jobID})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_image_dockerhub_rate_limit",
		Description: "Get the App catalog's current Docker Hub pull rate limit usage.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		limit, err := client.AppImageDockerHubRateLimitGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("app_image_dockerhub_rate_limit: %w", err))
		}
		return jsonResult(limit)
	})

	type getAppImageInput struct {
		ID string `json:"id" jsonschema:"Container image ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_image_get",
		Description: "Get details of a single pulled container image by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getAppImageInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("app_image_get: id must not be empty"))
		}
		result, err := client.GetAppImage(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("app_image_get: %w", err))
		}
		return jsonResult(result)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_categories",
		Description: "List the categories apps in the catalog can belong to.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		categories, err := client.AppCategories(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("app_categories: %w", err))
		}
		return jsonResult(categories)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_available_space",
		Description: "Get free space available to the App catalog's storage pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		space, err := client.AppAvailableSpaceGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("app_available_space: %w", err))
		}
		return jsonResult(space)
	})

	type appConfigInput struct {
		AppName string `json:"app_name" jsonschema:"App name (as shown in the TrueNAS UI)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_config",
		Description: "Get the user-specified configuration (values) of a single app by name.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p appConfigInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.AppConfigGet(ctx, p.AppName)
		if err != nil {
			return errorResult(fmt.Errorf("app_config: %w", err))
		}
		return jsonResult(cfg)
	})

	type appNameInput struct {
		AppName string `json:"app_name" jsonschema:"App name (as shown in the TrueNAS UI)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_container_ids",
		Description: "List the container IDs backing a running app.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p appNameInput) (*mcp.CallToolResult, any, error) {
		if p.AppName == "" {
			return errorResult(errors.New("app_container_ids: app_name must not be empty"))
		}
		ids, err := client.AppContainerIDs(ctx, p.AppName)
		if err != nil {
			return errorResult(fmt.Errorf("app_container_ids: %w", err))
		}
		return jsonResult(ids)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_convert_to_custom",
		Description: "Convert a catalog-managed app to a custom (unmanaged) app, detaching it from catalog upgrades.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p appNameInput) (*mcp.CallToolResult, any, error) {
		if p.AppName == "" {
			return errorResult(errors.New("app_convert_to_custom: app_name must not be empty"))
		}
		result, err := client.ConvertAppToCustom(ctx, p.AppName)
		if err != nil {
			return errorResult(fmt.Errorf("app_convert_to_custom: %w", err))
		}
		return jsonResult(result)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_outdated_docker_images",
		Description: "List the container images used by an app that have a newer version available.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p appNameInput) (*mcp.CallToolResult, any, error) {
		if p.AppName == "" {
			return errorResult(errors.New("app_outdated_docker_images: app_name must not be empty"))
		}
		images, err := client.AppOutdatedDockerImages(ctx, p.AppName)
		if err != nil {
			return errorResult(fmt.Errorf("app_outdated_docker_images: %w", err))
		}
		return jsonResult(images)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "app_pull_images",
		Description: "Pull the latest images for an app's containers. Returns the async job ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p appNameInput) (*mcp.CallToolResult, any, error) {
		if p.AppName == "" {
			return errorResult(errors.New("app_pull_images: app_name must not be empty"))
		}
		jobID, err := client.PullAppImages(ctx, p.AppName)
		if err != nil {
			return errorResult(fmt.Errorf("app_pull_images: %w", err))
		}
		return jsonResult(map[string]int{"job_id": jobID})
	})
}
