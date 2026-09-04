
package tools

import (
	"strconv"
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterPoolManagementTools registers all pool management MCP tools.
func RegisterPoolManagementTools(s *mcp.Server, client truenasClient) {
	// pool.attach() - Attach a detached pool
	type attachPoolInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_attach",
		Description: "Attach a detached ZFS pool to the system.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p attachPoolInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_attach: id must be a positive integer"))
		}
		pool, err := client.AttachPool(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_attach: %w", err))
		}
		return jsonResult(pool)
	})

	// pool.detach() - Detach pool from system
	type detachPoolInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_detach",
		Description: "Detach a ZFS pool from the system.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p detachPoolInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_detach: id must be a positive integer"))
		}
		err := client.DetachPool(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_detach: %w", err))
		}
		return jsonResult(map[string]string{"status": "detached"})
	})

	// pool.expand() - Expand pool capacity
	type expandPoolInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Path string `json:"path" jsonschema:"Path to new disk or zvol"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_expand",
		Description: "Expand a ZFS pool with additional capacity.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p expandPoolInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Path == "" {
			return errorResult(errors.New("pool_expand: id and path are required"))
		}
		pool, err := client.ExpandPool(ctx, &truenas.ExpandPoolParams{ID: p.ID})
		if err != nil {
			return errorResult(fmt.Errorf("pool_expand: %w", err))
		}
		return jsonResult(pool)
	})

	// pool.export() - Export pool data
	type exportPoolInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Path string `json:"path,omitempty" jsonschema:"Export destination path"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_export",
		Description: "Export a ZFS pool to an external location.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p exportPoolInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_export: id is required"))
		}
		jobID, err := client.ExportPool(ctx, p.ID, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("pool_export: %w", err))
		}
		return jsonResult(map[string]any{"job_id": jobID})
	})

	// pool.get_disks() - List all disks in pool with details
	type getDisksInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_get_disks",
		Description: "List all disks in a ZFS pool with detailed information.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getDisksInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_get_disks: id must be a positive integer"))
		}
		disks, err := client.ListDisks(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_get_disks: %w", err))
		}
		return jsonResult(disks)
	})

	// pool.html() - Get HTML representation of pool
	type htmlInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_html",
		Description: "Get an HTML representation of a ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p htmlInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_html: id must be a positive integer"))
		}
		html, err := client.ExportPoolHTML(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_html: %w", err))
		}
		return jsonResult(map[string]string{"html": html})
	})

	// pool.import_find() - Find datasets from import
	type importFindInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Name   string `json:"name,omitempty" jsonschema:"Dataset name filter"`
		Limit  int    `json:"limit,omitempty" jsonschema:"Maximum number of datasets to return"`
		Offset int    `json:"offset,omitempty" jsonschema:"Number of datasets to skip"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_import_find",
		Description: "Find datasets that can be imported into a pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p importFindInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_import_find: id must be a positive integer"))
		}
		datasets, err := client.FindDatasetsForImport(ctx, p.ID, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("pool_import_find: %w", err))
		}
		return jsonResult(datasets)
	})

	// pool.import_pool() - Import a pool from external storage
	type importPoolInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Path string `json:"path" jsonschema:"Path to import file or directory"`
		Name string `json:"name,omitempty" jsonschema:"New pool name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_import_pool",
		Description: "Import a ZFS pool from external storage.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p importPoolInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Path == "" {
			return errorResult(errors.New("pool_import_pool: id and path are required"))
		}
		err := client.ImportPool(ctx, p.ID, p.Path, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_import_pool: %w", err))
		}
		return jsonResult(map[string]string{"status": "imported"})
	})

	// pool.is_upgraded() - Check if pool is upgraded
	type isUpgradedInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_is_upgraded",
		Description: "Check if a ZFS pool has been upgraded.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p isUpgradedInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_is_upgraded: id must be a positive integer"))
		}
		isUpgraded, err := client.IsPoolUpgraded(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_is_upgraded: %w", err))
		}
		return jsonResult(map[string]bool{"is_upgraded": isUpgraded})
	})

	// pool.offline() - Mark pool as offline
	type offlineInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_offline",
		Description: "Mark a ZFS pool as offline.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p offlineInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_offline: id must be a positive integer"))
		}
		err := client.OfflinePool(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_offline: %w", err))
		}
		return jsonResult(map[string]string{"status": "offline"})
	})


	// pool.online() - Bring pool online
	type onlineInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_online",
		Description: "Bring a ZFS pool online.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p onlineInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_online: id must be a positive integer"))
		}
		err := client.OnlinePool(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_online: %w", err))
		}
		return jsonResult(map[string]string{"status": "online"})
	})

	// pool.relabel() - Relabel pool dataset files
	type relabelInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Key    string `json:"key,omitempty" jsonschema:"Relabel key"`
		Force  bool   `json:"force,omitempty" jsonschema:"Force relabel operation"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_relabel",
		Description: "Relabel dataset files with a new checksum key.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p relabelInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_relabel: id must be a positive integer"))
		}
		err := client.RelabelPool(ctx, p.ID, p.Key, p.Force)
		if err != nil {
			return errorResult(fmt.Errorf("pool_relabel: %w", err))
		}
		return jsonResult(map[string]string{"status": "relabelled"})
	})

	// pool.scan() - Scan pool for inconsistencies
	type scanInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Force  bool   `json:"force,omitempty" jsonschema:"Force scan operation"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_scan",
		Description: "Scan a ZFS pool for inconsistencies.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p scanInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_scan: id must be a positive integer"))
		}
		err := client.ScanPool(ctx, p.ID, p.Force)
		if err != nil {
			return errorResult(fmt.Errorf("pool_scan: %w", err))
		}
		return jsonResult(map[string]string{"status": "scan completed"})
	})

	// pool.split() - Split pool into new pools
	type splitInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Path   string `json:"path" jsonschema:"Path to split configuration"`
		Name   string `json:"name,omitempty" jsonschema:"New pool name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_split",
		Description: "Split a ZFS pool into new pools.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p splitInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Path == "" {
			return errorResult(errors.New("pool_split: id and path are required"))
		}
		err := client.SplitPool(ctx, p.ID, p.Path, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_split: %w", err))
		}
		return jsonResult(map[string]string{"status": "split"})
	})

	// pool.status() - Get pool status
	type statusInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_status",
		Description: "Get the current status of a ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p statusInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_status: id must be a positive integer"))
		}
		status, err := client.GetPoolStatus(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_status: %w", err))
		}
		return jsonResult(status)
	})

	// pool.sync() - Sync pool data
	type syncInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_sync",
		Description: "Synchronize a ZFS pool to ensure consistency.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p syncInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_sync: id must be a positive integer"))
		}
		err := client.SyncPool(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_sync: %w", err))
		}
		return jsonResult(map[string]string{"status": "synced"})
	})

	// pool.verify() - Verify pool integrity
	type verifyInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_verify",
		Description: "Verify a ZFS pool's integrity.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p verifyInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_verify: id must be a positive integer"))
		}
		err := client.VerifyPool(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_verify: %w", err))
		}
		return jsonResult(map[string]string{"status": "verified"})
	})

	// pool.dataset.create() - Create a new ZFS dataset
	type createDatasetInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Name   string `json:"name" jsonschema:"Full dataset path (e.g. Storage/backups)"`
		Type   string `json:"type,omitempty" jsonschema:"Dataset type: FILESYSTEM or VOLUME"`
		Compression string `json:"compression,omitempty" jsonschema:"Compression algorithm"`
		Comments string `json:"comments,omitempty" jsonschema:"Dataset description"`
		Quota   int64  `json:"quota,omitempty" jsonschema:"Maximum size in bytes (0 for no quota)"`
		Volsize int64  `json:"volsize,omitempty" jsonschema:"VZOL size in bytes (required for VOLUME type)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_create",
		Description: "Create a new ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createDatasetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_dataset_create: id and name are required"))
		}
	dataset, err := client.CreateDataset(ctx, &truenas.CreateDatasetParams{
			Name:      p.Name,
			Type:      p.Type,
			Compression: p.Compression,
			Comments:  p.Comments,
			Quota:     p.Quota,
			Volsize:   p.Volsize,
		})
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_create: %w", err))
		}
		return jsonResult(dataset)
	})

	// pool.dataset.delete() - Delete a dataset by path
	type deleteDatasetPathInput struct {
		Path     string `json:"path" jsonschema:"Full dataset path (e.g. Storage/backups/data)"`
		Recursive bool   `json:"recursive,omitempty" jsonschema:"Also delete child datasets"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_delete",
		Description: "Delete a ZFS dataset by its full path.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p deleteDatasetPathInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("pool_dataset_delete: path is required"))
		}
		err := client.DeleteDatasetByPath(ctx, p.Path, p.Recursive)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_delete: %w", err))
		}
		return jsonResult(map[string]string{"status": "dataset deleted"})
	})

	// pool.dataset.details() - Get detailed dataset information
	type detailsInput struct {
		ID string `json:"id" jsonschema:"Full dataset path (e.g. Storage/backups)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_details",
		Description: "Get detailed information about a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p detailsInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("pool_dataset_details: id is required"))
		}
		dataset, err := client.GetDataset(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_details: %w", err))
		}
		return jsonResult(dataset)
	})

	// pool.dataset.update() - Update dataset description or quota
	type updateDatasetInput struct {
		ID      string `json:"id" jsonschema:"Full dataset path"`
		Comments string `json:"comments,omitempty" jsonschema:"New dataset description"`
		Quota   int64  `json:"quota,omitempty" jsonschema:"New quota in bytes (0 for no quota)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_update",
		Description: "Update a ZFS dataset's description or quota.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateDatasetInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("pool_dataset_update: id is required"))
		}
	_, err := client.UpdateDataset(ctx, p.ID, &truenas.UpdateDatasetParams{
			Comments: p.Comments,
			Quota:    p.Quota,
		})
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_update: %w", err))
		}
		return jsonResult(map[string]string{"status": "updated"})
	})

	// pool.dataset.query() - Query dataset information
	type queryInput struct {
		ID string `json:"id" jsonschema:"Full dataset path"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_query",
		Description: "Query a ZFS dataset for information.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p queryInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("pool_dataset_query: id is required"))
		}
		query, err := client.GetDataset(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_query: %w", err))
		}
		return jsonResult(query)
	})

	// pool.dataset.get_instance() - Get dataset instance information
	type getInstanceInput struct {
		ID string `json:"id" jsonschema:"Full dataset path"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_get_instance",
		Description: "Get instance information for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getInstanceInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("pool_dataset_get_instance: id is required"))
		}
		instance, err := client.GetDatasetInstance(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_get_instance: %w", err))
		}
		return jsonResult(instance)
	})

	// pool.dataset.get_quota() - Get dataset quota information
	type getQuotaInput struct {
		ID string `json:"id" jsonschema:"Full dataset path"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_get_quota",
		Description: "Get quota information for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getQuotaInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" {
			return errorResult(errors.New("pool_dataset_get_quota: id is required"))
		}
		quota, err := client.GetDatasetQuota(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_get_quota: %w", err))
		}
		return jsonResult(quota)
	})

	// pool.dataset.set_quota() - Set dataset quota
	type setQuotaInput struct {
		ID    string `json:"id" jsonschema:"Full dataset path"`
		Quota int64  `json:"quota" jsonschema:"New quota in bytes (0 for no quota)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_set_quota",
		Description: "Set a quota for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p setQuotaInput) (*mcp.CallToolResult, any, error) {
		if p.ID == "" || p.Quota < 0 {
			return errorResult(errors.New("pool_dataset_set_quota: id and quota (>= 0) are required"))
		}
		err := client.SetDatasetQuota(ctx, p.ID, p.Quota)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_set_quota: %w", err))
		}
		return jsonResult(map[string]string{"status": "quota set"})
	})

	// pool.dataset.rename() - Rename dataset
	type renameDatasetInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Name   string `json:"name" jsonschema:"Current dataset name"`
		NewName string `json:"new_name" jsonschema:"New dataset name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_rename",
		Description: "Rename a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p renameDatasetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" || p.NewName == "" {
			return errorResult(errors.New("pool_dataset_rename: id, name, and new_name are required"))
		}
		dataset, err := client.RenameDataset(ctx, p.ID, p.Name, p.NewName)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_rename: %w", err))
		}
		return jsonResult(dataset)
	})

	// pool.dataset.export_key() - Export dataset encryption key
	type exportKeyInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
		Path string `json:"path" jsonschema:"Export path for key file"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_export_key",
		Description: "Export a dataset's encryption key to a file.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p exportKeyInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_dataset_export_key: id and name are required"))
		}
		keyID, err := client.ExportDatasetKey(ctx, p.ID, p.Name, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_export_key: %w", err))
		}
		return jsonResult(map[string]any{"key_id": keyID})
	})

	// pool.dataset.export_keys() - Export all dataset encryption keys
	type exportKeysInput struct {
		ID    int    `json:"id" jsonschema:"Numeric pool ID"`
		Path  string `json:"path,omitempty" jsonschema:"Export directory path"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_export_keys",
		Description: "Export all encryption keys from datasets in a pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p exportKeysInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_dataset_export_keys: id is required"))
		}
		count, err := client.ExportDatasetKeys(ctx, p.ID, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_export_keys: %w", err))
		}
		return jsonResult(map[string]any{"count": count})
	})

	// pool.dataset.export_keys_for_replication() - Export keys for replication
	type exportKeysForReplicationInput struct {
		ID    int    `json:"id" jsonschema:"Numeric pool ID"`
		Remote string `json:"remote" jsonschema:"Remote pool identifier"`
		Path  string `json:"path,omitempty" jsonschema:"Export directory path"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_export_keys_for_replication",
		Description: "Export encryption keys for replication setup.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p exportKeysForReplicationInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Remote == "" {
			return errorResult(errors.New("pool_dataset_export_keys_for_replication: id and remote are required"))
		}
		count, err := client.ExportDatasetKeysForReplication(ctx, p.ID, p.Remote, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_export_keys_for_replication: %w", err))
		}
		return jsonResult(map[string]any{"count": count})
	})

	// pool.dataset.lock() - Lock dataset encryption key
	type lockDatasetKeyInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_lock",
		Description: "Lock a dataset's encryption key.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p lockDatasetKeyInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_dataset_lock: id and name are required"))
		}
		err := client.LockDatasetKey(ctx, p.ID, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_lock: %w", err))
		}
		return jsonResult(map[string]string{"status": "key locked"})
	})

	// pool.dataset.unlock() - Unlock dataset encryption key
	type unlockDatasetKeyInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
		Key  string `json:"key" jsonschema:"Encryption key"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_unlock",
		Description: "Unlock a dataset's encryption key.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p unlockDatasetKeyInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" || p.Key == "" {
			return errorResult(errors.New("pool_dataset_unlock: id, name, and key are required"))
		}
		err := client.UnlockDatasetKey(ctx, p.ID, p.Name, p.Key)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_unlock: %w", err))
		}
		return jsonResult(map[string]string{"status": "key unlocked"})
	})

	// pool.dataset.promote() - Promote dataset from clone
	type promoteDatasetInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_dataset_promote",
		Description: "Promote a cloned dataset to standalone.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p promoteDatasetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_dataset_promote: id and name are required"))
		}
		err := client.PromoteDataset(ctx, p.ID, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_dataset_promote: %w", err))
		}
		return jsonResult(map[string]string{"status": "dataset promoted"})
	})

	// pool.ddt_prefetch() - Enable DDT prefetch
	type ddtPrefetchInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Enable bool   `json:"enable" jsonschema:"Enable or disable DDT prefetch"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_ddt_prefetch",
		Description: "Configure DDT prefetch for a ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p ddtPrefetchInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_ddt_prefetch: id is required"))
		}
		err := client.DDTPrefetch(ctx, p.ID, p.Enable)
		if err != nil {
			return errorResult(fmt.Errorf("pool_ddt_prefetch: %w", err))
		}
		return jsonResult(map[string]any{"status": "ddt prefetch configured", "enabled": p.Enable})
	})

	// pool.delete_dataset() - Delete a dataset by path
	type deleteDatasetByPathInput struct {
		Path     string `json:"path" jsonschema:"Full dataset path (e.g. Storage/backups/data)"`
		Recursive bool   `json:"recursive,omitempty" jsonschema:"Also delete child datasets"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_delete_dataset",
		Description: "Delete a ZFS dataset by its full path.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p deleteDatasetPathInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("pool_delete_dataset: path is required"))
		}
		err := client.DeleteDatasetByPath(ctx, p.Path, p.Recursive)
		if err != nil {
			return errorResult(fmt.Errorf("pool_delete_dataset: %w", err))
		}
		return jsonResult(map[string]string{"status": "dataset deleted"})
	})

	// pool.disk_info() - Get disk information
	type diskInfoInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_disk_info",
		Description: "Get detailed information about disks in a ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskInfoInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_disk_info: id must be a positive integer"))
		}
		info, err := client.GetDiskInfo(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_disk_info: %w", err))
		}
		return jsonResult(info)
	})

	// pool.disk_offline() - Mark disk as offline
	type diskOfflineInput struct {
		ID     int `json:"id" jsonschema:"Numeric pool ID"`
		DiskID int `json:"disk_id" jsonschema:"Disk identifier in pool"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_disk_offline",
		Description: "Mark a disk as offline.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskOfflineInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.DiskID <= 0 {
			return errorResult(errors.New("pool_disk_offline: id and disk_id must be positive integers"))
		}
		err := client.OfflineDisk(ctx, p.ID, p.DiskID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_disk_offline: %w", err))
		}
		return jsonResult(map[string]string{"status": "disk offline"})
	})

	// pool.disk_online() - Bring disk online
	type diskOnlineInput struct {
		ID     int `json:"id" jsonschema:"Numeric pool ID"`
		DiskID int `json:"disk_id" jsonschema:"Disk identifier in pool"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_disk_online",
		Description: "Bring a disk back online.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskOnlineInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.DiskID <= 0 {
			return errorResult(errors.New("pool_disk_online: id and disk_id must be positive integers"))
		}
		err := client.OnlineDisk(ctx, p.ID, p.DiskID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_disk_online: %w", err))
		}
		return jsonResult(map[string]string{"status": "disk online"})
	})

	// pool.disk_remove() - Remove disk from pool
	type diskRemoveInput struct {
		ID     int `json:"id" jsonschema:"Numeric pool ID"`
		DiskID int `json:"disk_id" jsonschema:"Disk identifier in pool"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_disk_remove",
		Description: "Remove a disk from a ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskRemoveInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.DiskID <= 0 {
			return errorResult(errors.New("pool_disk_remove: id and disk_id must be positive integers"))
		}
		err := client.RemoveDisk(ctx, p.ID, p.DiskID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_disk_remove: %w", err))
		}
		return jsonResult(map[string]string{"status": "disk removed"})
	})

	// pool.disk_replace() - Replace disk in pool
	type diskReplaceInput struct {
		ID      int    `json:"id" jsonschema:"Numeric pool ID"`
		DiskID  int    `json:"disk_id" jsonschema:"Disk identifier to replace"`
		Path    string `json:"path" jsonschema:"Path to replacement disk"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_disk_replace",
		Description: "Replace a disk in a ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p diskReplaceInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.DiskID <= 0 || p.Path == "" {
			return errorResult(errors.New("pool_disk_replace: id, disk_id, and path are required"))
		}
		err := client.ReplaceDisk(ctx, p.ID, p.DiskID, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("pool_disk_replace: %w", err))
		}
		return jsonResult(map[string]string{"status": "disk replaced"})
	})

	// pool.html_export() - Export pool HTML configuration
	type htmlExportInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_html_export",
		Description: "Export pool configuration as HTML.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p htmlExportInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_html_export: id must be a positive integer"))
		}
		html, err := client.ExportPoolHTML(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_html_export: %w", err))
		}
		return jsonResult(map[string]string{"html": html})
	})

	// pool.html_import() - Import from HTML configuration
	type htmlImportInput struct {
		Path string `json:"path" jsonschema:"Path to HTML file"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_html_import",
		Description: "Import pool configuration from HTML.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p htmlImportInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("pool_html_import: path is required"))
		}
		err := client.ImportPoolHTML(ctx, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("pool_html_import: %w", err))
		}
		return jsonResult(map[string]string{"status": "imported"})
	})

	// pool.import_dataset() - Import dataset from backup
	type importDatasetInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
		Path string `json:"path" jsonschema:"Path to backup file"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_import_dataset",
		Description: "Import a dataset from backup.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p importDatasetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" || p.Path == "" {
			return errorResult(errors.New("pool_import_dataset: id, name, and path are required"))
		}
		err := client.ImportDataset(ctx, p.ID, p.Name, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("pool_import_dataset: %w", err))
		}
		return jsonResult(map[string]string{"status": "dataset imported"})
	})

	// pool.list_datasets() - List all datasets in pool
	type listDatasetsInput struct {
		ID     int    `json:"id" jsonschema:"Numeric pool ID"`
		Limit  int    `json:"limit,omitempty" jsonschema:"Maximum number of datasets to return"`
		Offset int    `json:"offset,omitempty" jsonschema:"Number of datasets to skip"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_list_datasets",
		Description: "List all ZFS datasets in a pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listDatasetsInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_list_datasets: id must be a positive integer"))
		}
		datasets, err := client.ListDatasets(ctx, strconv.Itoa(p.ID), truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("pool_list_datasets: %w", err))
		}
		return jsonResult(datasets)
	})

	// pool.list_disks() - List all disks in pool
	type listDisksInput struct {
		ID int `json:"id" jsonschema:"Numeric pool ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_list_disks",
		Description: "List all disks in a ZFS pool.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listDisksInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("pool_list_disks: id must be a positive integer"))
		}
		disks, err := client.ListDisks(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("pool_list_disks: %w", err))
		}
		return jsonResult(disks)
	})

	// pool.offline_dataset() - Mark dataset as offline
	type offlineDatasetInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_offline_dataset",
		Description: "Mark a dataset as offline.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p offlineDatasetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_offline_dataset: id and name are required"))
		}
		err := client.OfflineDataset(ctx, p.ID, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_offline_dataset: %w", err))
		}
		return jsonResult(map[string]string{"status": "dataset offline"})
	})

	// pool.online_dataset() - Bring dataset online
	type onlineDatasetInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_online_dataset",
		Description: "Bring a dataset back online.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p onlineDatasetInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_online_dataset: id and name are required"))
		}
		err := client.OnlineDataset(ctx, p.ID, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_online_dataset: %w", err))
		}
		return jsonResult(map[string]string{"status": "dataset online"})
	})

	// pool.quota_check() - Check quota status
	type quotaCheckInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_quota_check",
		Description: "Check quota status for a ZFS dataset.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p quotaCheckInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_quota_check: id and name are required"))
		}
		status, err := client.CheckQuota(ctx, p.ID, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_quota_check: %w", err))
		}
		return jsonResult(status)
	})

	// pool.sync() - Sync dataset to ensure consistency
	type datasetSyncInput struct {
		ID   int    `json:"id" jsonschema:"Numeric pool ID"`
		Name string `json:"name" jsonschema:"Dataset name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "pool_sync",
		Description: "Sync a ZFS dataset to ensure consistency.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p datasetSyncInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 || p.Name == "" {
			return errorResult(errors.New("pool_sync: id and name are required"))
		}
		err := client.SyncDataset(ctx, p.ID, p.Name)
		if err != nil {
			return errorResult(fmt.Errorf("pool_sync: %w", err))
		}
		return jsonResult(map[string]string{"status": "dataset synced"})
	})

}
