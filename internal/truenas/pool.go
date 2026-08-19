package truenas

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Pool represents a ZFS pool.
type Pool struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	GUID           string  `json:"guid"`
	Properties     map[string]DatasetValue `json:"properties"`
	Allocated      DatasetValue `json:"allocated"`
	Free           DatasetValue `json:"free"`
	Status         string `json:"status"`
	ExpansionState string `json:"expansion_state"`
}

// PoolStatus represents the status of a pool.
type PoolStatus struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	Health        string  `json:"health"`
	Expanding     bool    `json:"expanding"`
	Repairing     bool    `json:"repairing"`
	Scanning      bool    `json:"scanning"`
	Verifying     bool    `json:"verifying"`
	Relabelling   bool    `json:"relabelling"`
	DetachDisk    bool    `json:"detach_disk"`
	Importing     bool    `json:"importing"`
	Offline       bool    `json:"offline"`
}

// CreatePoolParams holds parameters for creating a new ZFS pool.
type CreatePoolParams struct {
	Name        string `json:"name"`
	Layout      string `json:"layout,omitempty"`
	Disks       []string `json:"disks,omitempty"`
	Compression string `json:"compression,omitempty"`
}

// UpdatePoolParams holds parameters for updating a pool.
type UpdatePoolParams struct {
	Name        string `json:"name,omitempty"`
	Layout      string `json:"layout,omitempty"`
	Compression string `json:"compression,omitempty"`
}

// ListOptions holds optional list parameters.

// ExpandPoolParams holds parameters for expanding a pool.

// AttachPool attaches a disk to a pool.
func (c *Client) AttachPool(ctx context.Context, id int) (*Pool, error) {
	var pool Pool
	err := c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/attach", nil, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

// DetachPool detaches a disk from a pool.
func (c *Client) DetachPool(ctx context.Context, id int) error {
	return c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/detach", nil, nil)
}

// ExpandPool expands a pool to use all available space.
func (c *Client) ExpandPool(ctx context.Context, params *ExpandPoolParams) (*Pool, error) {
	if params == nil {
		return nil, errors.New("expand pool: params required")
	}
	var pool Pool
	err := c.call(ctx, "/storage/pool/"+strconv.Itoa(params.ID)+"/expand", nil, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

// ExportPool exports a pool's configuration.
func (c *Client) ExportPool(ctx context.Context, id int, path string) (int64, error) {
	u := "/storage/pool/" + strconv.Itoa(id) + "/export"
	if path != "" {
		u += "?path=" + url.QueryEscape(path)
	}
	var result struct{ JobID int64 }
	err := c.call(ctx, u, nil, &result)
	if err != nil {
		return 0, err
	}
	return result.JobID, nil
}

// ListDisks lists all disks in a pool.
func (c *Client) ListDisks(ctx context.Context, id int) ([]DiskInfo, error) {
	var disks []DiskInfo
	err := c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/disks", nil, &disks)
	if err != nil {
		return nil, err
	}
	return disks, nil
}

// ExportPoolHTML exports a pool's configuration as HTML.
func (c *Client) ExportPoolHTML(ctx context.Context, id int) (string, error) {
	var result struct{ HTML string }
	err := c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/html", nil, &result)
	if err != nil {
		return "", err
	}
	return result.HTML, nil
}

// FindDatasetsForImport finds datasets suitable for import into a pool.
func (c *Client) FindDatasetsForImport(ctx context.Context, id int, opts ListOptions) ([]Dataset, error) {
	var datasets []Dataset
	err := c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/import/find", nil, &datasets)
	if err != nil {
		return nil, err
	}
	return datasets, nil
}

// ImportPool imports a pool from configuration.
func (c *Client) ImportPool(ctx context.Context, id int, path, name string) error {
	u := "/storage/pool/" + strconv.Itoa(id) + "/import"
	if name != "" {
		u += "?name=" + url.QueryEscape(name)
	}
	return c.call(ctx, u, nil, nil)
}

// IsPoolUpgraded checks if a pool has been upgraded.
func (c *Client) IsPoolUpgraded(ctx context.Context, id int) (bool, error) {
	var result struct{ IsUpgraded bool }
	err := c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/is_upgraded", nil, &result)
	if err != nil {
		return false, err
	}
	return result.IsUpgraded, nil
}

// OfflinePool marks a pool as offline.
func (c *Client) OfflinePool(ctx context.Context, id int) error {
	return c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/offline", nil, nil)
}

// OnlinePool brings a pool back online.
func (c *Client) OnlinePool(ctx context.Context, id int) error {
	return c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/online", nil, nil)
}

// RelabelPool relabels a pool with a key.
func (c *Client) RelabelPool(ctx context.Context, id int, key string, force bool) error {
	u := "/storage/pool/" + strconv.Itoa(id) + "/relabel"
	if key != "" {
		u += "?key=" + url.QueryEscape(key)
	}
	return c.call(ctx, u, nil, nil)
}

// ScanPool scans a pool for errors.
func (c *Client) ScanPool(ctx context.Context, id int, force bool) error {
	return c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/scan", nil, nil)
}

// SplitPool splits a pool into multiple pools.
func (c *Client) SplitPool(ctx context.Context, id int, path, name string) error {
	u := "/storage/pool/" + strconv.Itoa(id) + "/split"
	if path != "" {
		u += "?path=" + url.QueryEscape(path)
	}
	if name != "" {
		u += "&name=" + url.QueryEscape(name)
	}
	return c.call(ctx, u, nil, nil)
}

// GetPoolStatus gets the status of a pool.
func (c *Client) GetPoolStatus(ctx context.Context, id int) (*PoolStatus, error) {
	var status PoolStatus
	err := c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/status", nil, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// SyncPool syncs a pool.
func (c *Client) SyncPool(ctx context.Context, id int) error {
	return c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/sync", nil, nil)
}

// VerifyPool verifies a pool.
func (c *Client) VerifyPool(ctx context.Context, id int) error {
	return c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/verify", nil, nil)
}

// ListPools lists all pools.
func (c *Client) ListPools(ctx context.Context, opts ...ListOptions) ([]Pool, error) {
	if len(opts) > 1 {
		return nil, errors.New("listing pools: at most one ListOptions value may be provided")
	}
	var o ListOptions
	if len(opts) == 1 {
		o = opts[0]
	}
	if err := validateListOptions(o); err != nil {
		return nil, fmt.Errorf("listing pools: %w", err)
	}
	var pools []Pool
	if err := c.call(ctx, "pool.query", buildQueryParams(nil, o), &pools); err != nil {
		return nil, fmt.Errorf("listing pools: %w", err)
	}
	return pools, nil
}

// GetPool gets a pool by ID.
func (c *Client) GetPool(ctx context.Context, id int) (*Pool, error) {
	var pool Pool
	if err := c.call(ctx, "pool.get_instance", []any{id}, &pool); err != nil {
		return nil, fmt.Errorf("getting pool %d: %w", id, err)
	}
	return &pool, nil
}

// CreatePool creates a new ZFS pool.
func (c *Client) CreatePool(ctx context.Context, params *CreatePoolParams) (*Pool, error) {
	if params == nil {
		return nil, errors.New("create pool: params required")
	}
	u := "/storage/pool/" + url.QueryEscape(params.Name)
	if params.Layout != "" {
		u += "?layout=" + url.QueryEscape(params.Layout)
	}
	if params.Compression != "" {
		u += "&compression=" + url.QueryEscape(params.Compression)
	}
	var pool Pool
	err := c.call(ctx, u, nil, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

// UpdatePool updates a pool's configuration.
func (c *Client) UpdatePool(ctx context.Context, id int, params *UpdatePoolParams) (*Pool, error) {
	if params == nil {
		return nil, errors.New("update pool: params required")
	}
	u := "/storage/pool/" + strconv.Itoa(id) + "/update"
	if params.Name != "" {
		u += "?name=" + url.QueryEscape(params.Name)
	}
	if params.Layout != "" {
		u += "&layout=" + url.QueryEscape(params.Layout)
	}
	if params.Compression != "" {
		u += "&compression=" + url.QueryEscape(params.Compression)
	}
	var pool Pool
	err := c.call(ctx, u, nil, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

// DeletePool deletes a pool.
func (c *Client) DeletePool(ctx context.Context, id int) error {
	return c.call(ctx, "/storage/pool/"+strconv.Itoa(id)+"/delete", nil, nil)
}

// ExpandPoolParams holds parameters for expanding a pool.


// ExpandPoolParams holds parameters for expanding a pool.
type ExpandPoolParams struct {
	ID int `json:"id"`
}
