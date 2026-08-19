package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCloudSyncTools registers Cloud Sync (rclone-backed) MCP tools onto the server.
func registerCloudSyncTools(s *mcp.Server, client truenasClient) {
	type listCloudProvidersInput struct{}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_cloud_providers",
		Description: "List the Cloud Sync providers supported by this TrueNAS host (S3, backup/cloud services, etc.), for use as the `provider` value when creating a Cloud Sync credential.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ listCloudProvidersInput) (*mcp.CallToolResult, any, error) {
		providers, err := client.ListCloudProviders(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("list_cloud_providers: %w", err))
		}
		return jsonResult(providers)
	})

	type listCloudCredentialsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of credentials to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of credentials to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_cloud_credentials",
		Description: "List all stored Cloud Sync credentials.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listCloudCredentialsInput) (*mcp.CallToolResult, any, error) {
		creds, err := client.ListCloudCredentials(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("list_cloud_credentials: %w", err))
		}
		return jsonResult(creds)
	})

	type createCloudCredentialInput struct {
		Name       string         `json:"name"       jsonschema:"Human-readable label for the credential"`
		Provider   string         `json:"provider"   jsonschema:"Provider identifier as returned by list_cloud_providers, e.g. S3"`
		Attributes map[string]any `json:"attributes,omitempty" jsonschema:"Provider-specific auth fields (e.g. access_key_id, secret_access_key, endpoint for S3/MinIO; token for OAuth-based providers)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_cloud_credential",
		Description: "Create a Cloud Sync credential (authentication for a remote such as S3-compatible storage or a cloud drive). Use list_cloud_providers first to get valid provider identifiers.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createCloudCredentialInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" {
			return errorResult(errors.New("create_cloud_credential: name must not be empty"))
		}
		if p.Provider == "" {
			return errorResult(errors.New("create_cloud_credential: provider must not be empty"))
		}
		cred, err := client.CreateCloudCredential(ctx, &truenas.CreateCloudCredentialParams{
			Name:       p.Name,
			Provider:   p.Provider,
			Attributes: p.Attributes,
		})
		if err != nil {
			return errorResult(fmt.Errorf("create_cloud_credential: %w", err))
		}
		return jsonResult(cred)
	})

	type listCloudSyncTasksInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of tasks to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of tasks to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_cloudsync_tasks",
		Description: "List all configured Cloud Sync Tasks.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listCloudSyncTasksInput) (*mcp.CallToolResult, any, error) {
		tasks, err := client.ListCloudSyncTasks(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("list_cloudsync_tasks: %w", err))
		}
		return jsonResult(tasks)
	})

	type createCloudSyncTaskInput struct {
		Description  string         `json:"description,omitempty" jsonschema:"Name of the task to display in the TrueNAS UI"`
		Path         string         `json:"path"                  jsonschema:"Local path on the TrueNAS host, beginning with /mnt or /dev/zvol"`
		Credentials  int            `json:"credentials"            jsonschema:"ID of a Cloud Sync credential (from create_cloud_credential or list_cloud_credentials)"`
		Direction    string         `json:"direction"              jsonschema:"PUSH (local to remote) or PULL (remote to local)"`
		TransferMode string         `json:"transfer_mode"          jsonschema:"SYNC, COPY, or MOVE"`
		Attributes   map[string]any `json:"attributes,omitempty"   jsonschema:"Provider-specific task fields, e.g. bucket, folder"`
		Extra        map[string]any `json:"extra,omitempty"        jsonschema:"Additional top-level fields not otherwise modeled (e.g. schedule)"`
		Enabled      bool           `json:"enabled,omitempty"      jsonschema:"Whether the task becomes active on its own schedule immediately. Defaults to false (disabled) — TrueNAS would otherwise default new tasks to enabled with an hourly schedule. Set true only when the task should run unattended immediately."`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_cloudsync_task",
		Description: "Create a Cloud Sync Task that syncs a local TrueNAS path with a remote via a previously created credential. Created disabled (enabled=false) by default — pass enabled=true only if the task should run unattended on its schedule immediately. Does not run the task itself — use run_cloudsync_task for that.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createCloudSyncTaskInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("create_cloudsync_task: path must not be empty"))
		}
		if p.Credentials <= 0 {
			return errorResult(errors.New("create_cloudsync_task: credentials must be a positive integer"))
		}
		if p.Direction == "" {
			return errorResult(errors.New("create_cloudsync_task: direction must not be empty"))
		}
		if p.TransferMode == "" {
			return errorResult(errors.New("create_cloudsync_task: transfer_mode must not be empty"))
		}
		task, err := client.CreateCloudSyncTask(ctx, &truenas.CreateCloudSyncTaskParams{
			Description:  p.Description,
			Path:         p.Path,
			Credentials:  p.Credentials,
			Direction:    p.Direction,
			TransferMode: p.TransferMode,
			Attributes:   p.Attributes,
			Extra:        p.Extra,
			Enabled:      p.Enabled,
		})
		if err != nil {
			return errorResult(fmt.Errorf("create_cloudsync_task: %w", err))
		}
		return jsonResult(task)
	})

	type runCloudSyncTaskInput struct {
		ID     int  `json:"id"                jsonschema:"ID of the Cloud Sync Task to run"`
		DryRun bool `json:"dry_run,omitempty" jsonschema:"Perform a dry run without making actual changes"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "run_cloudsync_task",
		Description: "Run a Cloud Sync Task, syncing local data to/from the remote. Returns the async job ID immediately (non-blocking) — this can be a long-running transfer.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p runCloudSyncTaskInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("run_cloudsync_task: id must be a positive integer"))
		}
		jobID, err := client.RunCloudSyncTask(ctx, p.ID, p.DryRun)
		if err != nil {
			return errorResult(fmt.Errorf("run_cloudsync_task: %w", err))
		}
		return jsonResult(map[string]int{"job_id": jobID})
	})

	type abortCloudSyncTaskInput struct {
		ID int `json:"id" jsonschema:"ID of the Cloud Sync Task whose running job should be aborted"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "abort_cloudsync_task",
		Description: "Abort a currently running Cloud Sync Task job.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p abortCloudSyncTaskInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("abort_cloudsync_task: id must be a positive integer"))
		}
		if err := client.AbortCloudSyncTask(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("abort_cloudsync_task: %w", err))
		}
		return jsonResult(map[string]any{"aborted": true, "id": p.ID})
	})
}
