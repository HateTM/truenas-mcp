package truenas

import (
	"context"
	"errors"
	"fmt"
)

// AppRegistry represents a private container registry credential used by the App catalog.
type AppRegistry struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URI      string `json:"uri"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// CreateAppRegistryParams holds fields for creating or updating a private container registry.
type CreateAppRegistryParams struct {
	Name     string `json:"name"`
	URI      string `json:"uri"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// PullAppImageParams holds fields for pulling a container image.
type PullAppImageParams struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag,omitempty"`
	RegistryID int    `json:"registry_id,omitempty"`
}

// DockerHubRateLimit reports the App catalog's current Docker Hub pull rate limit usage.
type DockerHubRateLimit struct {
	TotalLimit     int `json:"total_limit,omitempty"`
	RemainingLimit int `json:"remaining_limit,omitempty"`
}

// AppAvailableSpace reports free space available to the App catalog's storage pool.
type AppAvailableSpace struct {
	FreeBytes int64 `json:"free_bytes,omitempty"`
}

// AppGlobalConfig represents the App catalog's global configuration.
type AppGlobalConfig struct {
	Pool string `json:"pool,omitempty"`
}

// ListAppRegistries lists configured private container registries.
func (c *Client) ListAppRegistries(ctx context.Context, opts ...ListOptions) ([]AppRegistry, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []AppRegistry
	if err := c.call(ctx, "app.registry.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing app registries: %w", err)
	}
	return result, nil
}

// GetAppRegistry returns a single private container registry by ID.
func (c *Client) GetAppRegistry(ctx context.Context, id int) (*AppRegistry, error) {
	var result AppRegistry
	if err := c.call(ctx, "app.registry.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting app registry %d: %w", id, err)
	}
	return &result, nil
}

// CreateAppRegistry creates a new private container registry credential.
func (c *Client) CreateAppRegistry(ctx context.Context, p *CreateAppRegistryParams) (*AppRegistry, error) {
	if p == nil {
		return nil, errors.New("create app registry: params required")
	}
	var result AppRegistry
	if err := c.call(ctx, "app.registry.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating app registry %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateAppRegistry updates an existing private container registry credential.
func (c *Client) UpdateAppRegistry(ctx context.Context, id int, p *CreateAppRegistryParams) (*AppRegistry, error) {
	if p == nil {
		return nil, errors.New("update app registry: params required")
	}
	var result AppRegistry
	if err := c.call(ctx, "app.registry.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating app registry %d: %w", id, err)
	}
	return &result, nil
}

// DeleteAppRegistry deletes a private container registry credential.
func (c *Client) DeleteAppRegistry(ctx context.Context, id int) error {
	if err := c.call(ctx, "app.registry.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting app registry %d: %w", id, err)
	}
	return nil
}

// PullAppImage pulls a container image and returns the async job ID.
func (c *Client) PullAppImage(ctx context.Context, p *PullAppImageParams) (int, error) {
	if p == nil {
		return 0, errors.New("pull app image: params required")
	}
	var jobID int
	if err := c.call(ctx, "app.image.pull", []any{p}, &jobID); err != nil {
		return 0, fmt.Errorf("pulling app image %q: %w", p.Repository, err)
	}
	return jobID, nil
}

// AppImageDockerHubRateLimitGet returns the App catalog's current Docker Hub pull rate limit usage.
func (c *Client) AppImageDockerHubRateLimitGet(ctx context.Context) (*DockerHubRateLimit, error) {
	var limit DockerHubRateLimit
	if err := c.call(ctx, "app.image.dockerhub_rate_limit", nil, &limit); err != nil {
		return nil, fmt.Errorf("getting app image dockerhub rate limit: %w", err)
	}
	return &limit, nil
}

// GetAppImage returns a single pulled container image by ID.
func (c *Client) GetAppImage(ctx context.Context, id string) (*Image, error) {
	var result Image
	if err := c.call(ctx, "app.image.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting app image %q: %w", id, err)
	}
	return &result, nil
}

// AppCategories returns the categories apps in the catalog can belong to.
func (c *Client) AppCategories(ctx context.Context) ([]string, error) {
	var categories []string
	if err := c.call(ctx, "app.categories", nil, &categories); err != nil {
		return nil, fmt.Errorf("listing app categories: %w", err)
	}
	return categories, nil
}

// AppAvailableSpaceGet returns free space available to the App catalog's storage pool.
func (c *Client) AppAvailableSpaceGet(ctx context.Context) (*AppAvailableSpace, error) {
	var space AppAvailableSpace
	if err := c.call(ctx, "app.available_space", nil, &space); err != nil {
		return nil, fmt.Errorf("getting app available space: %w", err)
	}
	return &space, nil
}

// AppConfigGet returns the App catalog's global configuration.
func (c *Client) AppConfigGet(ctx context.Context) (*AppGlobalConfig, error) {
	var cfg AppGlobalConfig
	if err := c.call(ctx, "app.config", nil, &cfg); err != nil {
		return nil, fmt.Errorf("getting app config: %w", err)
	}
	return &cfg, nil
}

// AppContainerIDs returns the container IDs backing a running app.
func (c *Client) AppContainerIDs(ctx context.Context, appName string) ([]string, error) {
	if err := validateAppName(appName, "app_name"); err != nil {
		return nil, fmt.Errorf("listing app container ids: %w", err)
	}
	var ids []string
	if err := c.call(ctx, "app.container_ids", []any{appName}, &ids); err != nil {
		return nil, fmt.Errorf("listing container ids for app %q: %w", appName, err)
	}
	return ids, nil
}

// ConvertAppToCustom converts a catalog-managed app to a custom (unmanaged) app.
func (c *Client) ConvertAppToCustom(ctx context.Context, appName string) (*App, error) {
	if err := validateAppName(appName, "app_name"); err != nil {
		return nil, fmt.Errorf("converting app to custom: %w", err)
	}
	var result App
	if err := c.call(ctx, "app.convert_to_custom", []any{appName}, &result); err != nil {
		return nil, fmt.Errorf("converting app %q to custom: %w", appName, err)
	}
	return &result, nil
}

// AppOutdatedDockerImages returns the container images used by an app that have a newer version available.
func (c *Client) AppOutdatedDockerImages(ctx context.Context, appName string) ([]string, error) {
	if err := validateAppName(appName, "app_name"); err != nil {
		return nil, fmt.Errorf("listing outdated app images: %w", err)
	}
	var images []string
	if err := c.call(ctx, "app.outdated_docker_images", []any{appName}, &images); err != nil {
		return nil, fmt.Errorf("listing outdated images for app %q: %w", appName, err)
	}
	return images, nil
}

// PullAppImages pulls the latest images for an app's containers and returns the async job ID.
func (c *Client) PullAppImages(ctx context.Context, appName string) (int, error) {
	if err := validateAppName(appName, "app_name"); err != nil {
		return 0, fmt.Errorf("pulling app images: %w", err)
	}
	var jobID int
	if err := c.call(ctx, "app.pull_images", []any{appName}, &jobID); err != nil {
		return 0, fmt.Errorf("pulling images for app %q: %w", appName, err)
	}
	return jobID, nil
}
