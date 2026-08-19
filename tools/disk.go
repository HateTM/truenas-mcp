package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerDiskTools registers disk.* MCP tools onto the server (except disk.wipe,
// which is destructive and lives in destructive_disk.go).
func registerDiskTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listDisksInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of disks to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of disks to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_query",
		Description: "List physical disks installed in the system.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listDisksInput) (*mcp.CallToolResult, any, error) {
		result, err := client.QueryDisks(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("disk_query: %w", err))
		}
		return jsonResult(result)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_details",
		Description: "Get detailed information about disks not currently in use by any pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		details, err := client.DiskDetails(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("disk_details: %w", err))
		}
		return jsonResult(details)
	})

	type diskNameInput struct {
		Name string `json:"name" jsonschema:"Disk device name, e.g. sda"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_get_used",
		Description: "Get the number of bytes used on a disk.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskNameInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" {
			return errorResult(errors.New("disk_get_used: name must not be empty"))
		}
		used, err := client.DiskGetUsed(ctx, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("disk_get_used: %w", err))
		}
		return jsonResult(map[string]int64{"used": used})
	})

	type diskNamesInput struct {
		Names []string `json:"names,omitempty" jsonschema:"Disk device names to restrict results to; omit for all disks"`
	}

	type temperatureAggInput struct {
		Names []string `json:"names,omitempty" jsonschema:"Disk device names to restrict results to; omit for all disks"`
		Days  int      `json:"days,omitempty"  jsonschema:"Number of trailing days to aggregate over; defaults to server default"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_temperature_agg",
		Description: "Get aggregated (min/max/avg) temperature stats for disks over a trailing period.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p temperatureAggInput) (*mcp.CallToolResult, any, error) {
		agg, err := client.DiskTemperatureAggGet(ctx, p.Names, p.Days)
		if err != nil {
			return errorResult(fmt.Errorf("disk_temperature_agg: %w", err))
		}
		return jsonResult(agg)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_temperature_alerts",
		Description: "List disks currently in a temperature alert state.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskNamesInput) (*mcp.CallToolResult, any, error) {
		alerts, err := client.DiskTemperatureAlerts(ctx, p.Names)
		if err != nil {
			return errorResult(fmt.Errorf("disk_temperature_alerts: %w", err))
		}
		return jsonResult(alerts)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_temperatures",
		Description: "Get the current temperature in Celsius for each disk.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskNamesInput) (*mcp.CallToolResult, any, error) {
		temps, err := client.DiskTemperatures(ctx, p.Names)
		if err != nil {
			return errorResult(fmt.Errorf("disk_temperatures: %w", err))
		}
		return jsonResult(temps)
	})

	type updateDiskInput struct {
		Identifier    string `json:"identifier"              jsonschema:"Disk identifier, e.g. {serial}sN12345"`
		Description   string `json:"description,omitempty"   jsonschema:"Free-text description"`
		ToggleSMART   bool   `json:"togglesmart,omitempty"   jsonschema:"Enable S.M.A.R.T. monitoring for this disk"`
		AdvPowerMgmt  string `json:"advpowermgmt,omitempty"  jsonschema:"Advanced power management level, e.g. DISABLED, 1-254"`
		AcousticLevel string `json:"acousticlevel,omitempty" jsonschema:"Acoustic management level, e.g. DISABLED, MINIMUM, MEDIUM, MAXIMUM"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "disk_update",
		Description: "Update a disk's settings (description, S.M.A.R.T., power management, acoustic level).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateDiskInput) (*mcp.CallToolResult, any, error) {
		if p.Identifier == "" {
			return errorResult(errors.New("disk_update: identifier must not be empty"))
		}
		result, err := client.UpdateDisk(ctx, p.Identifier, &truenas.UpdateDiskParams{
			Description: p.Description, ToggleSMART: p.ToggleSMART, AdvPowerMgmt: p.AdvPowerMgmt, AcousticLevel: p.AcousticLevel,
		})
		if err != nil {
			return errorResult(fmt.Errorf("disk_update: %w", err))
		}
		return jsonResult(result)
	})

	// disk_wipe is destructive and lives in destructive_disk.go, gated behind
	// Config.AllowDestructive.
}
