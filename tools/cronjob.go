package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCronJobTools registers cronjob.* MCP tools onto the server.
func registerCronJobTools(s *mcp.Server, client truenasClient) {
	type listCronJobsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of cron jobs to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of cron jobs to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cronjob_list",
		Description: "List scheduled cron jobs.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listCronJobsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListCronJobs(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("cronjob_list: %w", err))
		}
		return jsonResult(result)
	})

	type getCronJobInput struct {
		ID int `json:"id" jsonschema:"Numeric cron job ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cronjob_get",
		Description: "Get a single cron job by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getCronJobInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cronjob_get: id must be a positive integer"))
		}
		result, err := client.GetCronJob(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("cronjob_get: %w", err))
		}
		return jsonResult(result)
	})

	type cronJobInput struct {
		Command     string         `json:"command"               jsonschema:"Shell command to run"`
		Schedule    map[string]any `json:"schedule,omitempty"    jsonschema:"Cron schedule fields, e.g. {\"minute\": \"0\", \"hour\": \"3\"}"`
		User        string         `json:"user"                  jsonschema:"System user to run the command as, e.g. root"`
		Enabled     bool           `json:"enabled,omitempty"     jsonschema:"Whether the job runs on its schedule"`
		Description string         `json:"description,omitempty" jsonschema:"Free-text description"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cronjob_create",
		Description: "Create a new scheduled cron job.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p cronJobInput) (*mcp.CallToolResult, any, error) {
		if p.Command == "" || p.User == "" {
			return errorResult(errors.New("cronjob_create: command and user are required"))
		}
		result, err := client.CreateCronJob(ctx, &truenas.CreateCronJobParams{
			Command: p.Command, Schedule: p.Schedule, User: p.User, Enabled: p.Enabled, Description: p.Description,
		})
		if err != nil {
			return errorResult(fmt.Errorf("cronjob_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateCronJobInput struct {
		ID          int            `json:"id"                     jsonschema:"Numeric cron job ID"`
		Command     string         `json:"command"                jsonschema:"Shell command to run"`
		Schedule    map[string]any `json:"schedule,omitempty"     jsonschema:"Cron schedule fields"`
		User        string         `json:"user"                   jsonschema:"System user to run the command as"`
		Enabled     bool           `json:"enabled,omitempty"      jsonschema:"Whether the job runs on its schedule"`
		Description string         `json:"description,omitempty"  jsonschema:"Free-text description"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "cronjob_update",
		Description: "Update an existing scheduled cron job.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateCronJobInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cronjob_update: id must be a positive integer"))
		}
		if p.Command == "" || p.User == "" {
			return errorResult(errors.New("cronjob_update: command and user are required"))
		}
		result, err := client.UpdateCronJob(ctx, p.ID, &truenas.CreateCronJobParams{
			Command: p.Command, Schedule: p.Schedule, User: p.User, Enabled: p.Enabled, Description: p.Description,
		})
		if err != nil {
			return errorResult(fmt.Errorf("cronjob_update: %w", err))
		}
		return jsonResult(result)
	})

	// cronjob_delete is destructive and lives in destructive_cronjob.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "cronjob_run",
		Description: "Trigger a cron job to run immediately, outside its schedule. Returns the async job ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getCronJobInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("cronjob_run: id must be a positive integer"))
		}
		jobID, err := client.RunCronJob(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("cronjob_run: %w", err))
		}
		return jsonResult(map[string]int{"job_id": jobID})
	})
}
