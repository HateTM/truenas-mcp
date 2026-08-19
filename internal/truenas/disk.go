package truenas

import (
	"context"
	"errors"
	"fmt"
)

// Disk represents a physical disk installed in the system.
type Disk struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Serial     string `json:"serial,omitempty"`
	Size       int64  `json:"size,omitempty"`
	Model      string `json:"model,omitempty"`
	Pool       string `json:"pool,omitempty"`
	Type       string `json:"type,omitempty"` // HDD, SSD, NVME
}

// UpdateDiskParams holds fields for updating a disk's settings.
type UpdateDiskParams struct {
	Description   string `json:"description,omitempty"`
	ToggleSMART   bool   `json:"togglesmart,omitempty"`
	AdvPowerMgmt  string `json:"advpowermgmt,omitempty"`
	AcousticLevel string `json:"acousticlevel,omitempty"`
}

// DiskTemperatureAgg holds aggregated temperature stats for a disk over a time period.
type DiskTemperatureAgg struct {
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`
	Avg float64 `json:"avg,omitempty"`
}

// WipeDiskParams holds fields for wiping a disk.
type WipeDiskParams struct {
	// Mode is QUICK, FULL, or FULL_WITH_ZEROS.
	Mode string `json:"mode"`
}

// QueryDisks lists physical disks installed in the system.
func (c *Client) QueryDisks(ctx context.Context, opts ...ListOptions) ([]Disk, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []Disk
	if err := c.call(ctx, "disk.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing disks: %w", err)
	}
	return result, nil
}

// DiskDetails returns detailed information about disks not currently in use by any pool.
func (c *Client) DiskDetails(ctx context.Context) (map[string]any, error) {
	var details map[string]any
	if err := c.call(ctx, "disk.details", nil, &details); err != nil {
		return nil, fmt.Errorf("getting disk details: %w", err)
	}
	return details, nil
}

// DiskGetUsed returns the number of bytes used on a disk.
func (c *Client) DiskGetUsed(ctx context.Context, name string) (int64, error) {
	if name == "" {
		return 0, errors.New("get disk used space: name must not be empty")
	}
	var used int64
	if err := c.call(ctx, "disk.get_used", []any{name}, &used); err != nil {
		return 0, fmt.Errorf("getting used space for disk %q: %w", name, err)
	}
	return used, nil
}

// DiskTemperatureAggGet returns aggregated (min/max/avg) temperature stats for the given
// disks over the trailing period of days days. An empty names slice covers all disks.
func (c *Client) DiskTemperatureAggGet(ctx context.Context, names []string, days int) (map[string]DiskTemperatureAgg, error) {
	var agg map[string]DiskTemperatureAgg
	if err := c.call(ctx, "disk.temperature_agg", []any{names, map[string]any{"days": days}}, &agg); err != nil {
		return nil, fmt.Errorf("getting disk temperature aggregates: %w", err)
	}
	return agg, nil
}

// DiskTemperatureAlerts returns the names of disks currently in a temperature alert state.
func (c *Client) DiskTemperatureAlerts(ctx context.Context, names []string) ([]string, error) {
	var alerts []string
	if err := c.call(ctx, "disk.temperature_alerts", []any{names}, &alerts); err != nil {
		return nil, fmt.Errorf("getting disk temperature alerts: %w", err)
	}
	return alerts, nil
}

// DiskTemperatures returns the current temperature in Celsius for each of the given disks.
// An empty names slice covers all disks.
func (c *Client) DiskTemperatures(ctx context.Context, names []string) (map[string]int, error) {
	var temps map[string]int
	if err := c.call(ctx, "disk.temperatures", []any{names}, &temps); err != nil {
		return nil, fmt.Errorf("getting disk temperatures: %w", err)
	}
	return temps, nil
}

// UpdateDisk updates a disk's settings (description, SMART, power management, acoustic level).
func (c *Client) UpdateDisk(ctx context.Context, identifier string, p *UpdateDiskParams) (*Disk, error) {
	if identifier == "" {
		return nil, errors.New("update disk: identifier must not be empty")
	}
	if p == nil {
		return nil, errors.New("update disk: params required")
	}
	var result Disk
	if err := c.call(ctx, "disk.update", []any{identifier, p}, &result); err != nil {
		return nil, fmt.Errorf("updating disk %q: %w", identifier, err)
	}
	return &result, nil
}

// WipeDisk wipes a disk and returns the async job ID.
func (c *Client) WipeDisk(ctx context.Context, identifier string, p *WipeDiskParams) (int, error) {
	if identifier == "" {
		return 0, errors.New("wipe disk: identifier must not be empty")
	}
	if p == nil || p.Mode == "" {
		return 0, errors.New("wipe disk: mode is required")
	}
	var jobID int
	if err := c.call(ctx, "disk.wipe", []any{identifier, p.Mode}, &jobID); err != nil {
		return 0, fmt.Errorf("wiping disk %q: %w", identifier, err)
	}
	return jobID, nil
}
