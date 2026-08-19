package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerPoolDatasetChoicesTools registers the remaining pool/dataset choices
// and maintenance MCP tools not covered by tools/pool_management.go.
func registerPoolDatasetChoicesTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_checksum_choices",
		Description: "List the checksum algorithms available for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.DatasetChecksumChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_checksum_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_compression_choices",
		Description: "List the compression algorithms available for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.DatasetCompressionChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_compression_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_encryption_algorithm_choices",
		Description: "List the encryption algorithms available for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.DatasetEncryptionAlgorithmChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_encryption_algorithm_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_recordsize_choices",
		Description: "List the record size choices available for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.DatasetRecordsizeChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_recordsize_choices: %w", err))
		}
		return jsonResult(choices)
	})

	type ddtPruneInput struct {
		PoolID     int `json:"pool_id"              jsonschema:"Numeric pool ID"`
		Days       int `json:"days,omitempty"       jsonschema:"Prune dedup table entries older than this many days"`
		Percentage int `json:"percentage,omitempty" jsonschema:"Prune this percentage of the dedup table"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_ddt_prune",
		Description: "Prune the ZFS dedup table (DDT) for a pool, by age in days or by percentage of entries.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p ddtPruneInput) (*mcp.CallToolResult, any, error) {
		if p.PoolID <= 0 {
			return errorResult(errors.New("pool_ddt_prune: pool_id must be a positive integer"))
		}
		if p.Days <= 0 && p.Percentage <= 0 {
			return errorResult(errors.New("pool_ddt_prune: one of days or percentage is required"))
		}
		if err := client.DDTPrune(ctx, p.PoolID, truenas.DDTPruneParams{Days: p.Days, Percentage: p.Percentage}); err != nil {
			return errorResult(fmt.Errorf("pool_ddt_prune: %w", err))
		}
		return jsonResult(map[string]any{"status": "pruned", "pool_id": p.PoolID})
	})

	type poolFilesystemChoicesInput struct {
		PoolIDs []int `json:"pool_ids,omitempty" jsonschema:"Restrict results to these pool IDs; omit for all pools"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_filesystem_choices",
		Description: "List filesystem paths available across pools, suitable as parent paths for a new dataset or zvol.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p poolFilesystemChoicesInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.PoolFilesystemChoices(ctx, p.PoolIDs)
		if err != nil {
			return errorResult(fmt.Errorf("pool_filesystem_choices: %w", err))
		}
		return jsonResult(choices)
	})
}
