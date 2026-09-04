package truenas

import (
	"context"
	"errors"
	"fmt"
)

// CronJob represents a scheduled cron job.
type CronJob struct {
	ID          int            `json:"id"`
	Command     string         `json:"command"`
	Schedule    map[string]any `json:"schedule,omitempty"`
	User        string         `json:"user,omitempty"`
	Enabled     bool           `json:"enabled,omitempty"`
	Description string         `json:"description,omitempty"`
}

// CreateCronJobParams holds fields for creating or updating a cron job.
type CreateCronJobParams struct {
	Command     string         `json:"command"`
	Schedule    map[string]any `json:"schedule,omitempty"`
	User        string         `json:"user"`
	Enabled     bool           `json:"enabled,omitempty"`
	Description string         `json:"description,omitempty"`
}

// ListCronJobs lists scheduled cron jobs.
func (c *Client) ListCronJobs(ctx context.Context, opts ...ListOptions) ([]CronJob, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []CronJob
	if err := c.call(ctx, "cronjob.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing cron jobs: %w", err)
	}
	return result, nil
}

// GetCronJob returns a single cron job by ID.
func (c *Client) GetCronJob(ctx context.Context, id int) (*CronJob, error) {
	var result CronJob
	if err := c.call(ctx, "cronjob.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting cron job %d: %w", id, err)
	}
	return &result, nil
}

// CreateCronJob creates a new cron job.
func (c *Client) CreateCronJob(ctx context.Context, p *CreateCronJobParams) (*CronJob, error) {
	if p == nil {
		return nil, errors.New("create cron job: params required")
	}
	var result CronJob
	if err := c.call(ctx, "cronjob.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating cron job: %w", err)
	}
	return &result, nil
}

// UpdateCronJob updates an existing cron job.
func (c *Client) UpdateCronJob(ctx context.Context, id int, p *CreateCronJobParams) (*CronJob, error) {
	if p == nil {
		return nil, errors.New("update cron job: params required")
	}
	var result CronJob
	if err := c.call(ctx, "cronjob.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating cron job %d: %w", id, err)
	}
	return &result, nil
}

// DeleteCronJob deletes a cron job.
func (c *Client) DeleteCronJob(ctx context.Context, id int) error {
	if err := c.call(ctx, "cronjob.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting cron job %d: %w", id, err)
	}
	return nil
}

// RunCronJob triggers a cron job to run immediately, returning its async job ID.
func (c *Client) RunCronJob(ctx context.Context, id int) (int, error) {
	var jobID int
	if err := c.call(ctx, "cronjob.run", []any{id}, &jobID); err != nil {
		return 0, fmt.Errorf("running cron job %d: %w", id, err)
	}
	return jobID, nil
}
