package truenas

import (
	"context"
	"fmt"
)

// DatasetChecksumChoices returns the checksum algorithms available for a new or updated dataset.
func (c *Client) DatasetChecksumChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "pool.dataset.checksum_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing dataset checksum choices: %w", err)
	}
	return choices, nil
}

// DatasetCompressionChoices returns the compression algorithms available for a new or updated dataset.
func (c *Client) DatasetCompressionChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "pool.dataset.compression_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing dataset compression choices: %w", err)
	}
	return choices, nil
}

// DatasetEncryptionAlgorithmChoices returns the encryption algorithms available for a new or updated dataset.
func (c *Client) DatasetEncryptionAlgorithmChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "pool.dataset.encryption_algorithm_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing dataset encryption algorithm choices: %w", err)
	}
	return choices, nil
}

// DatasetRecordsizeChoices returns the record size choices available for a new or updated dataset.
func (c *Client) DatasetRecordsizeChoices(ctx context.Context) ([]string, error) {
	var choices []string
	if err := c.call(ctx, "pool.dataset.recordsize_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing dataset recordsize choices: %w", err)
	}
	return choices, nil
}

// DDTPruneParams controls how much of the dedup table to prune. Exactly one of
// Days or Percentage should be set.
type DDTPruneParams struct {
	Days       int `json:"days,omitempty"`
	Percentage int `json:"percentage,omitempty"`
}

// DDTPrune prunes the ZFS dedup table (DDT) for a pool.
func (c *Client) DDTPrune(ctx context.Context, poolID int, p DDTPruneParams) error {
	if err := c.call(ctx, "pool.ddt_prune", []any{poolID, p}, nil); err != nil {
		return fmt.Errorf("pruning dedup table for pool %d: %w", poolID, err)
	}
	return nil
}

// PoolFilesystemChoices returns filesystem paths available across the given pools
// (or all pools if poolIDs is empty), suitable as parent paths for a new dataset or zvol.
func (c *Client) PoolFilesystemChoices(ctx context.Context, poolIDs []int) ([]string, error) {
	var choices []string
	if err := c.call(ctx, "pool.filesystem_choices", []any{poolIDs}, &choices); err != nil {
		return nil, fmt.Errorf("listing pool filesystem choices: %w", err)
	}
	return choices, nil
}
