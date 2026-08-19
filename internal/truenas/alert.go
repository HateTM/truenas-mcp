package truenas

import (
	"context"
	"fmt"
)

// Alert represents a TrueNAS system alert (e.g. pool capacity, app updates,
// deprecated API usage). Field names and shapes verified against a live
// alert.list response on TrueNAS SCALE 25.10.
type Alert struct {
	UUID      string `json:"uuid"`
	Source    string `json:"source"`
	Klass     string `json:"klass"`
	Level     string `json:"level"` // e.g. INFO, WARNING, CRITICAL
	Formatted string `json:"formatted"`
	Dismissed bool   `json:"dismissed"`
	OneShot   bool   `json:"one_shot"`
}

// ListAlerts returns all current TrueNAS alerts, both active and dismissed.
func (c *Client) ListAlerts(ctx context.Context) ([]Alert, error) {
	var alerts []Alert
	if err := c.call(ctx, "alert.list", []any{}, &alerts); err != nil {
		return nil, fmt.Errorf("listing alerts: %w", err)
	}
	return alerts, nil
}
