package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCoreTools registers core.* connectivity, introspection, and job-utility MCP tools.
func registerCoreTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_ping",
		Description: "Check connectivity to the TrueNAS server.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		pong, err := client.Ping(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("core_ping: %w", err))
		}
		return jsonResult(map[string]string{"pong": pong})
	})

	type pingRemoteInput struct {
		Params map[string]any `json:"params" jsonschema:"Remote ping parameters, e.g. {\"type\": \"ICMP\", \"hosts\": [\"1.1.1.1\"]}"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_ping_remote",
		Description: "Ping a remote host via the TrueNAS server, relaying the request server-side.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p pingRemoteInput) (*mcp.CallToolResult, any, error) {
		pong, err := client.PingRemote(ctx, p.Params)
		if err != nil {
			return errorResult(fmt.Errorf("core_ping_remote: %w", err))
		}
		return jsonResult(map[string]string{"pong": pong})
	})

	type getMethodsInput struct {
		App string `json:"app,omitempty" jsonschema:"Restrict results to methods belonging to this app/namespace; omit for all methods"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_get_methods",
		Description: "List the RPC methods available on the server, for introspection.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getMethodsInput) (*mcp.CallToolResult, any, error) {
		methods, err := client.GetMethods(ctx, p.App)
		if err != nil {
			return errorResult(fmt.Errorf("core_get_methods: %w", err))
		}
		return jsonResult(methods)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_get_services",
		Description: "List the middleware services running on the server.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		services, err := client.GetServices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("core_get_services: %w", err))
		}
		return jsonResult(services)
	})

	type jobIDInput struct {
		ID int `json:"id" jsonschema:"Numeric job ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_job_download_logs",
		Description: "Get a one-time URL to download the log file for a job.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p jobIDInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("core_job_download_logs: id must be a positive integer"))
		}
		url, err := client.JobDownloadLogs(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("core_job_download_logs: %w", err))
		}
		return jsonResult(map[string]string{"url": url})
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_job_wait",
		Description: "Block until a job reaches a terminal state (SUCCESS, FAILED, or ABORTED) and return its final state. Prefer get_job for a single non-blocking check.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p jobIDInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("core_job_wait: id must be a positive integer"))
		}
		job, err := client.JobWait(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("core_job_wait: %w", err))
		}
		return jsonResult(job)
	})

	type arpInput struct {
		Interface string `json:"interface,omitempty" jsonschema:"Restrict results to this network interface; omit for all interfaces"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_arp",
		Description: "Get the ARP table, optionally filtered to a single network interface.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p arpInput) (*mcp.CallToolResult, any, error) {
		table, err := client.ARP(ctx, p.Interface)
		if err != nil {
			return errorResult(fmt.Errorf("core_arp: %w", err))
		}
		return jsonResult(table)
	})

	type bulkInput struct {
		Method string  `json:"method" jsonschema:"RPC method name to invoke once per params entry, e.g. pool.dataset.update"`
		Params [][]any `json:"params" jsonschema:"One parameter list per call, e.g. [[\"tank/a\", {\"comments\": \"x\"}], [\"tank/b\", {\"comments\": \"y\"}]]"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_bulk",
		Description: "Invoke an RPC method once per entry in params, in a single request. Returns each call's raw result.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p bulkInput) (*mcp.CallToolResult, any, error) {
		if p.Method == "" {
			return errorResult(errors.New("core_bulk: method must not be empty"))
		}
		results, err := client.Bulk(ctx, p.Method, p.Params)
		if err != nil {
			return errorResult(fmt.Errorf("core_bulk: %w", err))
		}
		return jsonResult(results)
	})

	type downloadInput struct {
		Method   string `json:"method"   jsonschema:"RPC method whose output should be streamed as a downloadable file"`
		Args     []any  `json:"args,omitempty" jsonschema:"Arguments to pass to method"`
		Filename string `json:"filename" jsonschema:"Name for the downloaded file"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_download",
		Description: "Start a job that produces downloadable file content from an RPC method's output. Returns the job ID and a one-time download URL.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p downloadInput) (*mcp.CallToolResult, any, error) {
		if p.Method == "" {
			return errorResult(errors.New("core_download: method must not be empty"))
		}
		if p.Filename == "" {
			return errorResult(errors.New("core_download: filename must not be empty"))
		}
		result, err := client.Download(ctx, p.Method, p.Args, p.Filename)
		if err != nil {
			return errorResult(fmt.Errorf("core_download: %w", err))
		}
		return jsonResult(result)
	})

	type resizeShellInput struct {
		ID   string `json:"id"   jsonschema:"Shell session ID"`
		Cols int    `json:"cols" jsonschema:"Terminal width in columns"`
		Rows int    `json:"rows" jsonschema:"Terminal height in rows"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_resize_shell",
		Description: "Resize an interactive shell session's terminal dimensions.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p resizeShellInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("core_resize_shell: id must not be empty"))
		}
		if p.Cols <= 0 || p.Rows <= 0 {
			return errorResult(errors.New("core_resize_shell: cols and rows must be positive integers"))
		}
		if err := client.ResizeShell(ctx, p.ID, p.Cols, p.Rows); err != nil {
			return errorResult(fmt.Errorf("core_resize_shell: %w", err))
		}
		return jsonResult(map[string]any{"resized": true, "id": p.ID})
	})

	type subscribeInput struct {
		Event string `json:"event" jsonschema:"Event name to subscribe to, e.g. pool.query"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_subscribe",
		Description: "Subscribe to a server-side event stream. Returns a subscription ID for use with core_unsubscribe.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p subscribeInput) (*mcp.CallToolResult, any, error) {
		if p.Event == "" {
			return errorResult(errors.New("core_subscribe: event must not be empty"))
		}
		subID, err := client.Subscribe(ctx, p.Event)
		if err != nil {
			return errorResult(fmt.Errorf("core_subscribe: %w", err))
		}
		return jsonResult(map[string]string{"subscription_id": subID})
	})

	type unsubscribeInput struct {
		SubscriptionID string `json:"subscription_id" jsonschema:"Subscription ID returned by core_subscribe"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "core_unsubscribe",
		Description: "Cancel a subscription previously created with core_subscribe.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p unsubscribeInput) (*mcp.CallToolResult, any, error) {
		if p.SubscriptionID == "" {
			return errorResult(errors.New("core_unsubscribe: subscription_id must not be empty"))
		}
		if err := client.Unsubscribe(ctx, p.SubscriptionID); err != nil {
			return errorResult(fmt.Errorf("core_unsubscribe: %w", err))
		}
		return jsonResult(map[string]any{"unsubscribed": true, "subscription_id": p.SubscriptionID})
	})
}
