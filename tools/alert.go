package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerAlertTools registers TrueNAS system-alert MCP tools onto the server.
func registerAlertTools(s *mcp.Server, client truenasClient) {
	type listAlertsInput struct{}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_alerts",
		Description: "List current TrueNAS system alerts (both active and dismissed), e.g. pool capacity warnings, app updates, service failures. Check this proactively before/during long-running operations that consume disk space or system resources.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ listAlertsInput) (*mcp.CallToolResult, any, error) {
		alerts, err := client.ListAlerts(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("list_alerts: %w", err))
		}
		return jsonResult(alerts)
	})
}
