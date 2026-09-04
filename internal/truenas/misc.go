// Package truenas — misc.go implements the two smallest remaining tracked categories,
// device.get_info() and dns.query(), sharing one file since each is a single method.
package truenas

import (
	"context"
	"fmt"
)

// DNSConfig holds the currently configured DNS resolvers.
type DNSConfig struct {
	Nameservers []string `json:"nameservers,omitempty"`
}

// DeviceGetInfo returns information about devices of the given type
// (e.g. "SERIAL", "DISK", "GPU").
func (c *Client) DeviceGetInfo(ctx context.Context, deviceType string) (map[string]any, error) {
	var info map[string]any
	if err := c.call(ctx, "device.get_info", []any{deviceType}, &info); err != nil {
		return nil, fmt.Errorf("getting device info for %q: %w", deviceType, err)
	}
	return info, nil
}

// DNSQuery returns the currently configured DNS resolvers.
func (c *Client) DNSQuery(ctx context.Context) (*DNSConfig, error) {
	var cfg DNSConfig
	if err := c.call(ctx, "dns.query", nil, &cfg); err != nil {
		return nil, fmt.Errorf("querying dns config: %w", err)
	}
	return &cfg, nil
}
