package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerJobTools registers generic TrueNAS job-control MCP tools onto the server.
func registerJobTools(s *mcp.Server, client truenasClient) {
	type getJobInput struct {
		ID int `json:"id" jsonschema:"ID of the job to look up"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_job",
		Description: "Get the current state of a TrueNAS job by ID. Useful for checking on a long-running operation after an MCP tool call times out client-side but the job continues running server-side.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getJobInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("get_job: id must be a positive integer"))
		}
		job, err := client.GetJob(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("get_job: %w", err))
		}
		return jsonResult(job)
	})

	type abortJobInput struct {
		ID int `json:"id" jsonschema:"ID of the job to abort"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "abort_job",
		Description: "Abort a running TrueNAS job by ID, regardless of which API domain created it. Use this when a domain-specific abort tool (e.g. abort_cloudsync_task) doesn't stop all jobs associated with a task.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p abortJobInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("abort_job: id must be a positive integer"))
		}
		if err := client.AbortJob(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("abort_job: %w", err))
		}
		return jsonResult(map[string]any{"aborted": true, "id": p.ID})
	})
}
