package truenas

import (
	"context"
	"encoding/json"
	"testing"
)

func TestListAlerts_success(t *testing.T) {
	t.Parallel()

	alerts := []Alert{
		{UUID: "1", Klass: "PoolUSBDisk", Level: "WARNING", Formatted: "Pool NaS is low on space", Dismissed: false},
	}

	srv := wsTestServer(t, map[string]methodHandler{
		"alert.list": func(_ json.RawMessage) (any, *rpcError) {
			return alerts, nil
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.ListAlerts(context.Background())
	if err != nil {
		t.Fatalf("ListAlerts: %v", err)
	}
	if len(got) != 1 || got[0].Level != "WARNING" {
		t.Errorf("got %+v, want one WARNING alert", got)
	}
}

func TestListAlerts_error(t *testing.T) {
	t.Parallel()

	srv := wsTestServer(t, map[string]methodHandler{
		"alert.list": func(_ json.RawMessage) (any, *rpcError) {
			return nil, &rpcError{Code: -32603, Message: "internal error"}
		},
	})
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if _, err := c.ListAlerts(context.Background()); err == nil {
		t.Error("expected error, got nil")
	}
}
