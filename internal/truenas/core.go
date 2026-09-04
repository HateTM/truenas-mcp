package truenas

import (
	"context"
	"errors"
	"fmt"
)

// CoreDownloadResult holds the async job and one-time download URL returned by core.download.
type CoreDownloadResult struct {
	JobID int    `json:"job_id"`
	URL   string `json:"url"`
}

// Ping checks connectivity to the TrueNAS server.
func (c *Client) Ping(ctx context.Context) (string, error) {
	var pong string
	if err := c.call(ctx, "core.ping", []any{}, &pong); err != nil {
		return "", fmt.Errorf("pinging server: %w", err)
	}
	return pong, nil
}

// PingRemote pings a remote host via the TrueNAS server, relaying the request server-side.
// params is passed through verbatim to core.ping_remote, e.g. {"type": "ICMP", "hosts": [...]}.
func (c *Client) PingRemote(ctx context.Context, params map[string]any) (string, error) {
	var pong string
	if err := c.call(ctx, "core.ping_remote", []any{params}, &pong); err != nil {
		return "", fmt.Errorf("pinging remote host: %w", err)
	}
	return pong, nil
}

// GetMethods returns the RPC methods available on the server, optionally filtered to those
// belonging to app (an empty app returns all methods).
func (c *Client) GetMethods(ctx context.Context, app string) (map[string]any, error) {
	var params []any
	if app != "" {
		params = []any{app}
	}
	var methods map[string]any
	if err := c.call(ctx, "core.get_methods", params, &methods); err != nil {
		return nil, fmt.Errorf("listing rpc methods: %w", err)
	}
	return methods, nil
}

// GetServices returns the middleware services running on the server.
func (c *Client) GetServices(ctx context.Context) ([]map[string]any, error) {
	var services []map[string]any
	if err := c.call(ctx, "core.get_services", nil, &services); err != nil {
		return nil, fmt.Errorf("listing services: %w", err)
	}
	return services, nil
}

// JobDownloadLogs returns a one-time URL to download the log file for a job.
func (c *Client) JobDownloadLogs(ctx context.Context, jobID int) (string, error) {
	var url string
	if err := c.call(ctx, "core.job_download_logs", []any{jobID}, &url); err != nil {
		return "", fmt.Errorf("getting log download url for job %d: %w", jobID, err)
	}
	return url, nil
}

// JobWait blocks until job jobID reaches a terminal state, reusing the existing PollJob
// long-poll helper rather than issuing a separate server-side blocking RPC.
func (c *Client) JobWait(ctx context.Context, jobID int) (*Job, error) {
	return c.PollJob(ctx, jobID)
}

// ARP returns the ARP table, optionally filtered to a single network interface.
func (c *Client) ARP(ctx context.Context, iface string) (map[string]string, error) {
	var params []any
	if iface != "" {
		params = []any{map[string]any{"interface": iface}}
	}
	var table map[string]string
	if err := c.call(ctx, "core.arp", params, &table); err != nil {
		return nil, fmt.Errorf("getting arp table: %w", err)
	}
	return table, nil
}

// Bulk invokes method once per entry in paramsList and returns each call's raw result.
func (c *Client) Bulk(ctx context.Context, method string, paramsList [][]any) ([]any, error) {
	if method == "" {
		return nil, errors.New("bulk call: method must not be empty")
	}
	var results []any
	if err := c.call(ctx, "core.bulk", []any{method, paramsList}, &results); err != nil {
		return nil, fmt.Errorf("bulk-calling %q: %w", method, err)
	}
	return results, nil
}

// Download starts a job that produces file content for retrieval, returning the job ID
// and a one-time download URL. method and args describe the underlying RPC call whose
// output should be streamed as a file; filename names the downloaded file.
func (c *Client) Download(ctx context.Context, method string, args []any, filename string) (*CoreDownloadResult, error) {
	if method == "" {
		return nil, errors.New("download: method must not be empty")
	}
	var result CoreDownloadResult
	if err := c.call(ctx, "core.download", []any{method, args, filename}, &result); err != nil {
		return nil, fmt.Errorf("starting download for %q: %w", method, err)
	}
	return &result, nil
}

// ResizeShell resizes an interactive shell session identified by id to cols x rows.
func (c *Client) ResizeShell(ctx context.Context, id string, cols, rows int) error {
	if id == "" {
		return errors.New("resize shell: id must not be empty")
	}
	if err := c.call(ctx, "core.resize_shell", []any{id, cols, rows}, nil); err != nil {
		return fmt.Errorf("resizing shell %q: %w", id, err)
	}
	return nil
}

// Subscribe subscribes to a server-side event stream and returns the subscription ID.
func (c *Client) Subscribe(ctx context.Context, event string) (string, error) {
	if event == "" {
		return "", errors.New("subscribe: event must not be empty")
	}
	var subID string
	if err := c.call(ctx, "core.subscribe", []any{event}, &subID); err != nil {
		return "", fmt.Errorf("subscribing to %q: %w", event, err)
	}
	return subID, nil
}

// Unsubscribe cancels a subscription previously created with Subscribe.
func (c *Client) Unsubscribe(ctx context.Context, subscriptionID string) error {
	if subscriptionID == "" {
		return errors.New("unsubscribe: subscription_id must not be empty")
	}
	if err := c.call(ctx, "core.unsubscribe", []any{subscriptionID}, nil); err != nil {
		return fmt.Errorf("unsubscribing %q: %w", subscriptionID, err)
	}
	return nil
}
