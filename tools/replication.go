package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerReplicationTools registers replication.* (task and global-config) MCP tools.
// replication.endpoint.* tools live in replication_endpoint.go.
func registerReplicationTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listReplicationTasksInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of tasks to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of tasks to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_list",
		Description: "List configured ZFS replication tasks.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listReplicationTasksInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListReplicationTasks(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("replication_list: %w", err))
		}
		return jsonResult(result)
	})

	type getReplicationTaskInput struct {
		ID int `json:"id" jsonschema:"Numeric replication task ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_get",
		Description: "Get a single replication task by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getReplicationTaskInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("replication_get: id must be a positive integer"))
		}
		result, err := client.GetReplicationTask(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("replication_get: %w", err))
		}
		return jsonResult(result)
	})

	type replicationTaskInput struct {
		Name           string         `json:"name"                     jsonschema:"Task name"`
		SourceDatasets []string       `json:"source_datasets"          jsonschema:"Local datasets to replicate, e.g. [\"Storage/backups\"]"`
		TargetDataset  string         `json:"target_dataset"           jsonschema:"Destination dataset path on the target system"`
		Direction      string         `json:"direction"                jsonschema:"PUSH (local to remote/local) or PULL (remote to local)"`
		SSHCredentials int            `json:"ssh_credentials,omitempty" jsonschema:"Replication endpoint ID (from replication_endpoint_create) for the remote system; omit for local-to-local replication"`
		Schedule       map[string]any `json:"schedule,omitempty"       jsonschema:"Cron schedule fields, e.g. {\"minute\": \"0\", \"hour\": \"3\"}"`
		Enabled        bool           `json:"enabled,omitempty"        jsonschema:"Whether the task becomes active on its own schedule immediately. Defaults to false — set true only if it should run unattended."`
		ReadOnly       string         `json:"readonly,omitempty"       jsonschema:"SET, REQUIRE, or IGNORE — how to handle the readonly property on the target dataset"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_create",
		Description: "Create a new ZFS replication task. Created disabled (enabled=false) by default. Does not run the task itself.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p replicationTaskInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" {
			return errorResult(errors.New("replication_create: name must not be empty"))
		}
		if len(p.SourceDatasets) == 0 {
			return errorResult(errors.New("replication_create: source_datasets must have at least one entry"))
		}
		if p.TargetDataset == "" || p.Direction == "" {
			return errorResult(errors.New("replication_create: target_dataset and direction are required"))
		}
		result, err := client.CreateReplicationTask(ctx, &truenas.CreateReplicationParams{
			Name: p.Name, SourceDatasets: p.SourceDatasets, TargetDataset: p.TargetDataset, Direction: p.Direction,
			SSHCredentials: p.SSHCredentials, Schedule: p.Schedule, Enabled: p.Enabled, ReadOnly: p.ReadOnly,
		})
		if err != nil {
			return errorResult(fmt.Errorf("replication_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateReplicationTaskInput struct {
		ID             int            `json:"id"                       jsonschema:"Numeric replication task ID"`
		Name           string         `json:"name"                     jsonschema:"Task name"`
		SourceDatasets []string       `json:"source_datasets"          jsonschema:"Local datasets to replicate"`
		TargetDataset  string         `json:"target_dataset"           jsonschema:"Destination dataset path"`
		Direction      string         `json:"direction"                jsonschema:"PUSH or PULL"`
		SSHCredentials int            `json:"ssh_credentials,omitempty" jsonschema:"Replication endpoint ID for the remote system"`
		Schedule       map[string]any `json:"schedule,omitempty"       jsonschema:"Cron schedule fields"`
		Enabled        bool           `json:"enabled,omitempty"        jsonschema:"Whether the task is active"`
		ReadOnly       string         `json:"readonly,omitempty"       jsonschema:"SET, REQUIRE, or IGNORE"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_update",
		Description: "Update an existing ZFS replication task.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateReplicationTaskInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("replication_update: id must be a positive integer"))
		}
		if p.Name == "" || len(p.SourceDatasets) == 0 || p.TargetDataset == "" || p.Direction == "" {
			return errorResult(errors.New("replication_update: name, source_datasets, target_dataset, and direction are required"))
		}
		result, err := client.UpdateReplicationTask(ctx, p.ID, &truenas.CreateReplicationParams{
			Name: p.Name, SourceDatasets: p.SourceDatasets, TargetDataset: p.TargetDataset, Direction: p.Direction,
			SSHCredentials: p.SSHCredentials, Schedule: p.Schedule, Enabled: p.Enabled, ReadOnly: p.ReadOnly,
		})
		if err != nil {
			return errorResult(fmt.Errorf("replication_update: %w", err))
		}
		return jsonResult(result)
	})

	// replication_delete is destructive and lives in destructive_replication.go,
	// gated behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_test",
		Description: "Run a dry-run validation of a replication task's configuration without transferring data.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getReplicationTaskInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("replication_test: id must be a positive integer"))
		}
		result, err := client.TestReplicationTask(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("replication_test: %w", err))
		}
		return jsonResult(result)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_schemas",
		Description: "Get the JSON schemas describing valid replication task configurations.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		schemas, err := client.ReplicationSchemas(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("replication_schemas: %w", err))
		}
		return jsonResult(schemas)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_global_config",
		Description: "Get global replication settings.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.ReplicationGlobalConfigGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("replication_global_config: %w", err))
		}
		return jsonResult(cfg)
	})

	type updateReplicationGlobalInput struct {
		MaxParallelReplicationTasks int `json:"max_parallel_replication_tasks,omitempty" jsonschema:"Maximum number of replication tasks that may run concurrently"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "replication_global_update",
		Description: "Update global replication settings.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateReplicationGlobalInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.UpdateReplicationGlobal(ctx, &truenas.UpdateReplicationGlobalParams{
			MaxParallelReplicationTasks: p.MaxParallelReplicationTasks,
		})
		if err != nil {
			return errorResult(fmt.Errorf("replication_global_update: %w", err))
		}
		return jsonResult(cfg)
	})
}
