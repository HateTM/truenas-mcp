package truenas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// QuotaInfo represents quota information for a dataset.
type QuotaInfo struct {
	Quota     int64 `json:"quota"`
	Used      int64 `json:"used"`
	Guid      string `json:"guid"`
	Dataset   string `json:"dataset"`
}

// QuotaStatus represents the quota status of a dataset.
type QuotaStatus struct {
	Quota     int64 `json:"quota"`
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
}

// DiskInfo represents information about a disk in a pool.
type DiskInfo struct {
	ID        int     `json:"id"`
	Path      string  `json:"path"`
	Device    string  `json:"device"`
	Size      int64   `json:"size"`
	Status    string  `json:"status"`
	Role      string  `json:"role"`
	Faulted   bool    `json:"faulted"`
	Offline   bool    `json:"offline"`
}

// Dataset represents a ZFS dataset or zvol.
type Dataset struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	GUID        string     `json:"guid"`
	Type        string     `json:"type"`
	Properties  map[string]DatasetValue `json:"properties"`
	Mountpoint  string     `json:"mountpoint"`
	Quota       DatasetValue `json:"quota"`
	Used        DatasetValue `json:"used"`
}

// DatasetValue is a TrueNAS property wrapper returned for most dataset fields.
type DatasetValue struct {
	Value    string `json:"value"`
	RawValue string `json:"rawvalue"`
	Parsed   any    `json:"parsed"`
}

func (v DatasetValue) MarshalJSON() ([]byte, error) {
	if v.Value == "" {
		return []byte("null"), nil
	}
	return []byte("\"" + v.Value + "\""), nil
}

func (v *DatasetValue) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v.Value = s
	v.RawValue = s
	return nil
}

// CreateDatasetParams holds parameters for creating a new ZFS dataset.
type CreateDatasetParams struct {
	Name       string `json:"name"`
	Type       string `json:"type,omitempty"`
	Compression string `json:"compression,omitempty"`
	Comments   string `json:"comments,omitempty"`
	Quota      int64  `json:"quota,omitempty"`
	Volsize    int64  `json:"volsize,omitempty"`
}

// UpdateDatasetParams holds parameters for updating a dataset.
type UpdateDatasetParams struct {
	Comments string `json:"comments,omitempty"`
	Quota    int64  `json:"quota,omitempty"`
}

// ListOptions holds optional list parameters.
type ListOptions struct {
	Limit   int
	Offset  int
}

// DeleteDataset permanently removes the ZFS dataset.
func (c *Client) DeleteDataset(ctx context.Context, id string) error {
	return c.call(ctx, "/storage/dataset/"+url.QueryEscape(id)+"/delete", nil, nil)
}

func (c *Client) ListDatasets(ctx context.Context, pool string, opts ...ListOptions) ([]Dataset, error) {
	if len(opts) > 1 {
		return nil, errors.New("listing datasets: at most one ListOptions value may be provided")
	}
	var o ListOptions
	if len(opts) == 1 {
		o = opts[0]
	}
	if err := validateListOptions(o); err != nil {
		return nil, fmt.Errorf("listing datasets: %w", err)
	}
	var filter [][]string
	if pool != "" {
		filter = [][]string{{"pool", "=", pool}}
	}
	var datasets []Dataset
	if err := c.call(ctx, "pool.dataset.query", buildQueryParams(filter, o), &datasets); err != nil {
		return nil, fmt.Errorf("listing datasets: %w", err)
	}
	return datasets, nil
}

func (c *Client) GetDataset(ctx context.Context, id string) (*Dataset, error) {
	var dataset Dataset
	if err := c.call(ctx, "pool.dataset.get_instance", []any{id}, &dataset); err != nil {
		return nil, fmt.Errorf("getting dataset %q: %w", id, err)
	}
	return &dataset, nil
}

func (c *Client) CreateDataset(ctx context.Context, params *CreateDatasetParams) (*Dataset, error) {
	if params == nil {
		return nil, errors.New("create dataset: params required")
	}
	if params.Name == "" {
		return nil, errors.New("create dataset: name must not be empty")
	}
	if params.Type == "VOLUME" && params.Volsize <= 0 {
		return nil, errors.New("create dataset: volsize is required for VOLUME type")
	}
	var dataset Dataset
	if err := c.call(ctx, "pool.dataset.create", []any{params}, &dataset); err != nil {
		return nil, fmt.Errorf("creating dataset %q: %w", params.Name, err)
	}
	return &dataset, nil
}

func (c *Client) UpdateDataset(ctx context.Context, id string, params *UpdateDatasetParams) (*Dataset, error) {
	if params == nil {
		return nil, errors.New("update dataset: params required")
	}
	u := "/storage/dataset/" + url.QueryEscape(id) + "/update"
	if params.Comments != "" {
		u += "?comments=" + url.QueryEscape(params.Comments)
	}
	if params.Quota > 0 {
		u += "&quota=" + strconv.FormatInt(params.Quota, 10)
	}
	var dataset Dataset
	err := c.call(ctx, u, nil, &dataset)
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

// DeleteDatasetByPath permanently removes the ZFS dataset by path.
func (c *Client) DeleteDatasetByPath(ctx context.Context, path string, recursive bool) error {
	return c.call(ctx, "/storage/dataset/"+url.QueryEscape(path)+"/delete", nil, nil)
}

// ExportDatasetKey exports a dataset's encryption key.
func (c *Client) ExportDatasetKey(ctx context.Context, poolID int, datasetName, path string) (string, error) {
	var result struct{ KeyID string }
	err := c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/export_key", nil, &result)
	if err != nil {
		return "", err
	}
	return result.KeyID, nil
}

// ExportDatasetKeys exports all encryption keys from datasets in a pool.
func (c *Client) ExportDatasetKeys(ctx context.Context, poolID int, path string) (int64, error) {
	var result struct{ Count int64 }
	err := c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/export_keys", nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

// ExportDatasetKeysForReplication exports encryption keys for replication.
func (c *Client) ExportDatasetKeysForReplication(ctx context.Context, poolID int, remote, path string) (int64, error) {
	u := "/storage/dataset/" + strconv.Itoa(poolID) + "/export_keys_for_replication"
	if remote != "" {
		u += "?remote=" + url.QueryEscape(remote)
	}
	if path != "" {
		u += "&path=" + url.QueryEscape(path)
	}
	var result struct{ Count int64 }
	err := c.call(ctx, u, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

// LockDatasetKey locks a dataset's encryption key.
func (c *Client) LockDatasetKey(ctx context.Context, poolID int, datasetName string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/lock", nil, nil)
}

// UnlockDatasetKey unlocks a dataset's encryption key.
func (c *Client) UnlockDatasetKey(ctx context.Context, poolID int, datasetName, key string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/unlock", map[string]string{"key": key}, nil)
}

// PromoteDataset promotes a cloned dataset to standalone.
func (c *Client) PromoteDataset(ctx context.Context, poolID int, datasetName string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/promote", nil, nil)
}

// DDTPrefetch configures DDT prefetch for a pool.
func (c *Client) DDTPrefetch(ctx context.Context, poolID int, enable bool) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/ddt_prefetch", map[string]bool{"enable": enable}, nil)
}

// GetDatasetInstance gets instance information for a dataset.
func (c *Client) GetDatasetInstance(ctx context.Context, datasetName string) (*Dataset, error) {
	var instance Dataset
	err := c.call(ctx, "/storage/dataset/"+url.QueryEscape(datasetName)+"/instance", nil, &instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// GetDatasetQuota gets quota information for a dataset.
func (c *Client) GetDatasetQuota(ctx context.Context, datasetName string) (*QuotaInfo, error) {
	var quota QuotaInfo
	err := c.call(ctx, "/storage/dataset/"+url.QueryEscape(datasetName)+"/quota", nil, &quota)
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// SetDatasetQuota sets a quota for a dataset.
func (c *Client) SetDatasetQuota(ctx context.Context, datasetName string, quota int64) error {
	return c.call(ctx, "/storage/dataset/"+url.QueryEscape(datasetName)+"/set_quota", map[string]int64{"quota": quota}, nil)
}

// RenameDataset renames a ZFS dataset.
func (c *Client) RenameDataset(ctx context.Context, poolID int, oldName, newName string) (*Dataset, error) {
	var result Dataset
	err := c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(oldName)+"/rename", map[string]string{"new_name": newName}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ImportDataset imports a dataset from backup.
func (c *Client) ImportDataset(ctx context.Context, poolID int, datasetName, path string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/import", map[string]string{"path": path}, nil)
}

// CheckQuota checks quota status for a dataset.
func (c *Client) CheckQuota(ctx context.Context, poolID int, datasetName string) (*QuotaStatus, error) {
	var status QuotaStatus
	err := c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/quota_check", nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// SyncDataset syncs a dataset to ensure consistency.
func (c *Client) SyncDataset(ctx context.Context, poolID int, datasetName string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/sync", nil, nil)
}

// OfflineDataset marks a dataset as offline.
func (c *Client) OfflineDataset(ctx context.Context, poolID int, datasetName string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/offline", nil, nil)
}

// OnlineDataset brings a dataset back online.
func (c *Client) OnlineDataset(ctx context.Context, poolID int, datasetName string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/"+url.QueryEscape(datasetName)+"/online", nil, nil)
}

// GetDiskInfo gets disk information for a pool.
func (c *Client) GetDiskInfo(ctx context.Context, poolID int) (*DiskInfo, error) {
	var info DiskInfo
	err := c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/disk_info", nil, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// OfflineDisk marks a disk as offline.
func (c *Client) OfflineDisk(ctx context.Context, poolID int, diskID int) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/disk_offline", map[string]int{"disk_id": diskID}, nil)
}

// OnlineDisk brings a disk back online.
func (c *Client) OnlineDisk(ctx context.Context, poolID int, diskID int) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/disk_online", map[string]int{"disk_id": diskID}, nil)
}

// RemoveDisk removes a disk from a pool.
func (c *Client) RemoveDisk(ctx context.Context, poolID int, diskID int) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/disk_remove", map[string]int{"disk_id": diskID}, nil)
}

// ReplaceDisk replaces a disk in a pool.
func (c *Client) ReplaceDisk(ctx context.Context, poolID int, diskID int, path string) error {
	return c.call(ctx, "/storage/dataset/"+strconv.Itoa(poolID)+"/disk_replace", map[string]any{"disk_id": diskID, "path": path}, nil)
}

// ImportPoolHTML imports a pool from HTML configuration.
func (c *Client) ImportPoolHTML(ctx context.Context, path string) error {
	return c.call(ctx, "/storage/pool/html_import", map[string]string{"path": path}, nil)
}


// QueryDataset queries a dataset's properties.
func (c *Client) QueryDataset(ctx context.Context, datasetName string) (*Dataset, error) {
	return c.GetDataset(ctx, datasetName)
}
