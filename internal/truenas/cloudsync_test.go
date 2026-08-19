package truenas

import (
	"context"
	"encoding/json"
	"testing"
)

func TestListCloudProviders_success(t *testing.T) {
	t.Parallel()

	providers := []CloudProvider{{Name: "S3"}, {Name: "YANDEX"}}

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.providers": func(_ json.RawMessage) (any, *rpcError) {
			return providers, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.ListCloudProviders(context.Background())
	if err != nil {
		t.Fatalf("ListCloudProviders: %v", err)
	}
	if len(got) != 2 || got[0].Name != "S3" {
		t.Errorf("got %+v, want providers starting with S3", got)
	}
}

func TestListCloudCredentials_success(t *testing.T) {
	t.Parallel()

	creds := []CloudCredential{{ID: 1, Name: "minio", Provider: map[string]any{"type": "S3"}}}

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.credentials.query": func(_ json.RawMessage) (any, *rpcError) {
			return creds, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.ListCloudCredentials(context.Background(), ListOptions{})
	if err != nil {
		t.Fatalf("ListCloudCredentials: %v", err)
	}
	if len(got) != 1 || got[0].Name != "minio" {
		t.Errorf("got %+v, want one credential named minio", got)
	}
}

func TestCreateCloudCredential_success(t *testing.T) {
	t.Parallel()

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.credentials.create": func(params json.RawMessage) (any, *rpcError) {
			var args []map[string]any
			if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
				return nil, &rpcError{Code: -32600, Message: "wrong params"}
			}
			config, ok := args[0]["config"].(map[string]any)
			if !ok || config["type"] != "S3" || config["access_key_id"] != "kb_admin" {
				return nil, &rpcError{Code: -32600, Message: "wrong config payload"}
			}
			return CloudCredential{ID: 5, Name: args[0]["name"].(string), Provider: map[string]any{"type": "S3"}}, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.CreateCloudCredential(context.Background(), &CreateCloudCredentialParams{
		Name:     "minio",
		Provider: "S3",
		Attributes: map[string]any{
			"access_key_id": "kb_admin",
		},
	})
	if err != nil {
		t.Fatalf("CreateCloudCredential: %v", err)
	}
	if got.ID != 5 || got.Name != "minio" {
		t.Errorf("got %+v, want ID=5 Name=minio", got)
	}
}

func TestCreateCloudCredential_validation(t *testing.T) {
	t.Parallel()

	c, err := NewClient("http://localhost:19999", "test-api-key", false) //nolint:gosec // G101: fake placeholder
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tests := []struct {
		name   string
		params *CreateCloudCredentialParams
	}{
		{"nil params", nil},
		{"empty name", &CreateCloudCredentialParams{Provider: "S3"}},
		{"empty provider", &CreateCloudCredentialParams{Name: "minio"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := c.CreateCloudCredential(context.Background(), tc.params); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestDeleteCloudCredential_validation(t *testing.T) {
	t.Parallel()

	c, err := NewClient("http://localhost:19999", "test-api-key", false) //nolint:gosec // G101: fake placeholder
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if err := c.DeleteCloudCredential(context.Background(), 0); err == nil {
		t.Error("expected error for id=0, got nil")
	}
}

func TestListCloudSyncTasks_success(t *testing.T) {
	t.Parallel()

	tasks := []CloudSyncTask{{ID: 1, Path: "/mnt/NaS/yandex_disk", Direction: "PUSH"}}

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.query": func(_ json.RawMessage) (any, *rpcError) {
			return tasks, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.ListCloudSyncTasks(context.Background(), ListOptions{})
	if err != nil {
		t.Fatalf("ListCloudSyncTasks: %v", err)
	}
	if len(got) != 1 || got[0].Path != "/mnt/NaS/yandex_disk" {
		t.Errorf("got %+v, want one task for /mnt/NaS/yandex_disk", got)
	}
}

func TestCreateCloudSyncTask_success(t *testing.T) {
	t.Parallel()

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.create": func(params json.RawMessage) (any, *rpcError) {
			var args []map[string]any
			if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
				return nil, &rpcError{Code: -32600, Message: "wrong params"}
			}
			if args[0]["path"] != "/mnt/NaS/yandex_disk" || args[0]["credentials"] != float64(5) {
				return nil, &rpcError{Code: -32600, Message: "wrong task payload"}
			}
			return CloudSyncTask{ID: 9, Path: "/mnt/NaS/yandex_disk", Direction: "PUSH"}, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.CreateCloudSyncTask(context.Background(), &CreateCloudSyncTaskParams{
		Path:         "/mnt/NaS/yandex_disk",
		Credentials:  5,
		Direction:    "PUSH",
		TransferMode: "COPY",
		Attributes:   map[string]any{"bucket": "site-archive", "folder": "raw"},
	})
	if err != nil {
		t.Fatalf("CreateCloudSyncTask: %v", err)
	}
	if got.ID != 9 {
		t.Errorf("ID = %d, want 9", got.ID)
	}
}

func TestCreateCloudSyncTask_enabledDefaultsFalse(t *testing.T) {
	t.Parallel()

	var gotEnabled any
	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.create": func(params json.RawMessage) (any, *rpcError) {
			var args []map[string]any
			if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
				return nil, &rpcError{Code: -32600, Message: "wrong params"}
			}
			enabled, present := args[0]["enabled"]
			if !present {
				return nil, &rpcError{Code: -32600, Message: "enabled field missing from request"}
			}
			gotEnabled = enabled
			return CloudSyncTask{ID: 1, Path: args[0]["path"].(string)}, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	// Enabled deliberately left at its zero value (false) to mirror a caller
	// that doesn't explicitly opt in — TrueNAS's own default (enabled=true,
	// hourly schedule) must never take over silently.
	if _, err := c.CreateCloudSyncTask(context.Background(), &CreateCloudSyncTaskParams{
		Path: "/mnt/NaS/yandex_disk", Credentials: 5, Direction: "PUSH", TransferMode: "MOVE",
	}); err != nil {
		t.Fatalf("CreateCloudSyncTask: %v", err)
	}
	if gotEnabled != false {
		t.Errorf("enabled sent = %v, want false", gotEnabled)
	}
}

func TestCreateCloudSyncTask_enabledTruePassesThrough(t *testing.T) {
	t.Parallel()

	var gotEnabled any
	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.create": func(params json.RawMessage) (any, *rpcError) {
			var args []map[string]any
			if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
				return nil, &rpcError{Code: -32600, Message: "wrong params"}
			}
			gotEnabled = args[0]["enabled"]
			return CloudSyncTask{ID: 1, Path: args[0]["path"].(string)}, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if _, err := c.CreateCloudSyncTask(context.Background(), &CreateCloudSyncTaskParams{
		Path: "/mnt/NaS/yandex_disk", Credentials: 5, Direction: "PUSH", TransferMode: "MOVE", Enabled: true,
	}); err != nil {
		t.Fatalf("CreateCloudSyncTask: %v", err)
	}
	if gotEnabled != true {
		t.Errorf("enabled sent = %v, want true", gotEnabled)
	}
}

func TestCreateCloudSyncTask_validation(t *testing.T) {
	t.Parallel()

	c, err := NewClient("http://localhost:19999", "test-api-key", false) //nolint:gosec // G101: fake placeholder
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tests := []struct {
		name   string
		params *CreateCloudSyncTaskParams
	}{
		{"nil params", nil},
		{"empty path", &CreateCloudSyncTaskParams{Credentials: 1, Direction: "PUSH", TransferMode: "COPY"}},
		{"missing credentials", &CreateCloudSyncTaskParams{Path: "/mnt/x", Direction: "PUSH", TransferMode: "COPY"}},
		{"missing direction", &CreateCloudSyncTaskParams{Path: "/mnt/x", Credentials: 1, TransferMode: "COPY"}},
		{"missing transfer_mode", &CreateCloudSyncTaskParams{Path: "/mnt/x", Credentials: 1, Direction: "PUSH"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := c.CreateCloudSyncTask(context.Background(), tc.params); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestRunCloudSyncTask_success(t *testing.T) {
	t.Parallel()

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.sync": func(params json.RawMessage) (any, *rpcError) {
			var args []any
			if err := json.Unmarshal(params, &args); err != nil || len(args) != 2 {
				return nil, &rpcError{Code: -32600, Message: "wrong params"}
			}
			return 33, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	jobID, err := c.RunCloudSyncTask(context.Background(), 9, false)
	if err != nil {
		t.Fatalf("RunCloudSyncTask: %v", err)
	}
	if jobID != 33 {
		t.Errorf("jobID = %d, want 33", jobID)
	}
}

func TestRunCloudSyncTask_validation(t *testing.T) {
	t.Parallel()

	c, err := NewClient("http://localhost:19999", "test-api-key", false) //nolint:gosec // G101: fake placeholder
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if _, err := c.RunCloudSyncTask(context.Background(), 0, false); err == nil {
		t.Error("expected error for id=0, got nil")
	}
}

func TestAbortCloudSyncTask_success(t *testing.T) {
	t.Parallel()

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.abort": func(_ json.RawMessage) (any, *rpcError) {
			return nil, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if err := c.AbortCloudSyncTask(context.Background(), 9); err != nil {
		t.Fatalf("AbortCloudSyncTask: %v", err)
	}
}

func TestDeleteCloudSyncTask_success(t *testing.T) {
	t.Parallel()

	srv := wsTestServer(t, map[string]methodHandler{
		"cloudsync.delete": func(_ json.RawMessage) (any, *rpcError) {
			return nil, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if err := c.DeleteCloudSyncTask(context.Background(), 9); err != nil {
		t.Fatalf("DeleteCloudSyncTask: %v", err)
	}
}
