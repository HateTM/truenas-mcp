package tools

import (
	"context"
	"errors"
	"testing"
)

func TestPoolDatasetChecksumChoices(t *testing.T) {
	t.Run("returns choices as JSON", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"dataset_checksum_choices": func(_ context.Context) (map[string]string, error) {
					return map[string]string{"SHA256": "SHA256"}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "pool_dataset_checksum_choices", nil)
		assertResultJSON(t, res)
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"dataset_checksum_choices": func(_ context.Context) (map[string]string, error) {
					return nil, errors.New("API error")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "pool_dataset_checksum_choices", nil)
		assertError(t, res, "API error")
	})
}

func TestPoolDatasetCompressionChoices(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"dataset_compression_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"LZ4": "LZ4"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "pool_dataset_compression_choices", nil)
	assertResultJSON(t, res)
}

func TestPoolDatasetEncryptionAlgorithmChoices(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"dataset_encryption_algorithm_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"AES-256-GCM": "AES-256-GCM"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "pool_dataset_encryption_algorithm_choices", nil)
	assertResultJSON(t, res)
}

func TestPoolDatasetRecordsizeChoices(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"dataset_recordsize_choices": func(_ context.Context) ([]string, error) {
				return []string{"128K", "1M"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "pool_dataset_recordsize_choices", nil)
	assertResultJSON(t, res)
}

func TestPoolDDTPrune(t *testing.T) {
	t.Run("prunes by days", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"ddt_prune": func(_ context.Context, poolID int, _ any) error {
					if poolID != 1 {
						t.Errorf("poolID = %d, want 1", poolID)
					}
					return nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "pool_ddt_prune", map[string]any{"pool_id": 1, "days": 30})
		assertResultJSON(t, res)
	})

	t.Run("requires days or percentage", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "pool_ddt_prune", map[string]any{"pool_id": 1})
		assertError(t, res, "one of days or percentage is required")
	})

	t.Run("invalid pool_id returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "pool_ddt_prune", map[string]any{"pool_id": 0, "days": 30})
		assertError(t, res, "pool_id must be a positive integer")
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"ddt_prune": func(_ context.Context, _ int, _ any) error {
					return errors.New("pool busy")
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "pool_ddt_prune", map[string]any{"pool_id": 1, "percentage": 10})
		assertError(t, res, "pool busy")
	})
}

func TestPoolFilesystemChoices(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"pool_filesystem_choices": func(_ context.Context, _ []int) ([]string, error) {
				return []string{"Storage", "Storage/backups"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "pool_filesystem_choices", map[string]any{"pool_ids": []int{1}})
	assertResultJSON(t, res)
}
