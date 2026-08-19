package truenas

import (
	"context"
	"encoding/json"
	"testing"
)

func TestGetJob_success(t *testing.T) {
	t.Parallel()

	srv := wsTestServer(t, map[string]methodHandler{
		"core.get_jobs": func(_ json.RawMessage) (any, *rpcError) {
			return []Job{{ID: 42, Method: "cloudsync.sync", State: "RUNNING"}}, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetJob(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetJob: %v", err)
	}
	if got.ID != 42 || got.State != "RUNNING" {
		t.Errorf("got %+v, want ID=42 State=RUNNING", got)
	}
}

func TestGetJob_validation(t *testing.T) {
	t.Parallel()

	c, err := NewClient("http://localhost:19999", "test-api-key", false) //nolint:gosec // G101: fake placeholder
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if _, err := c.GetJob(context.Background(), 0); err == nil {
		t.Error("expected error for id=0, got nil")
	}
}

func TestAbortJob_success(t *testing.T) {
	t.Parallel()

	var gotID float64
	srv := wsTestServer(t, map[string]methodHandler{
		"core.job_abort": func(params json.RawMessage) (any, *rpcError) {
			var args []float64
			if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
				return nil, &rpcError{Code: -32600, Message: "wrong params"}
			}
			gotID = args[0]
			return nil, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if err := c.AbortJob(context.Background(), 10818); err != nil {
		t.Fatalf("AbortJob: %v", err)
	}
	if gotID != 10818 {
		t.Errorf("job id sent = %v, want 10818", gotID)
	}
}

func TestAbortJob_validation(t *testing.T) {
	t.Parallel()

	c, err := NewClient("http://localhost:19999", "test-api-key", false) //nolint:gosec // G101: fake placeholder
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if err := c.AbortJob(context.Background(), -1); err == nil {
		t.Error("expected error for id=-1, got nil")
	}
}
