package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestISCSIExtentList(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_extent_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ISCSIExtent, error) {
				return []truenas.ISCSIExtent{{ID: 1, Name: "extent0", Type: "DISK"}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_extent_list", nil)
	assertResultJSON(t, res)
}

func TestISCSIExtentGet(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_extent_get": func(_ context.Context, id int) (*truenas.ISCSIExtent, error) {
				return &truenas.ISCSIExtent{ID: id, Name: "extent0"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_extent_get", map[string]any{"id": 1})
	assertResultJSON(t, res)
}

func TestISCSIExtentCreate(t *testing.T) {
	t.Run("creates DISK extent", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_extent_create": func(_ context.Context, p *truenas.CreateISCSIExtentParams) (*truenas.ISCSIExtent, error) {
					return &truenas.ISCSIExtent{ID: 1, Name: p.Name, Type: p.Type}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_extent_create", map[string]any{"name": "extent0", "type": "DISK", "disk": "zvol/Storage/vm-disk"})
		assertResultJSON(t, res)
	})

	t.Run("DISK type requires disk", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_extent_create", map[string]any{"name": "extent0", "type": "DISK"})
		assertError(t, res, "disk is required when type=DISK")
	})

	t.Run("FILE type requires path and filesize", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_extent_create", map[string]any{"name": "extent0", "type": "FILE"})
		assertError(t, res, "path and filesize are required when type=FILE")
	})
}

func TestISCSIExtentUpdate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_extent_update": func(_ context.Context, id int, _ *truenas.CreateISCSIExtentParams) (*truenas.ISCSIExtent, error) {
				return &truenas.ISCSIExtent{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_extent_update", map[string]any{"id": 1, "name": "extent0", "type": "DISK", "disk": "zvol/Storage/vm-disk"})
	assertResultJSON(t, res)
}

func TestISCSIExtentDiskChoices(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_extent_disk_choices": func(_ context.Context) (map[string]string, error) {
				return map[string]string{"zvol/Storage/vm-disk": "zvol/Storage/vm-disk"}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_extent_disk_choices", nil)
	assertResultJSON(t, res)
}

func TestISCSIInitiatorList(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_initiator_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ISCSIInitiator, error) {
				return []truenas.ISCSIInitiator{{ID: 1}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_initiator_list", nil)
	assertResultJSON(t, res)
}

func TestISCSIInitiatorGet(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_initiator_get": func(_ context.Context, id int) (*truenas.ISCSIInitiator, error) {
				return &truenas.ISCSIInitiator{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_initiator_get", map[string]any{"id": 1})
	assertResultJSON(t, res)
}

func TestISCSIInitiatorCreate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_initiator_create": func(_ context.Context, p *truenas.CreateISCSIInitiatorParams) (*truenas.ISCSIInitiator, error) {
				return &truenas.ISCSIInitiator{ID: 1, Initiators: p.Initiators}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_initiator_create", map[string]any{"initiators": []string{"iqn.initiator"}})
	assertResultJSON(t, res)
}

func TestISCSIInitiatorUpdate(t *testing.T) {
	t.Run("updates initiator group", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_initiator_update": func(_ context.Context, id int, _ *truenas.CreateISCSIInitiatorParams) (*truenas.ISCSIInitiator, error) {
					return &truenas.ISCSIInitiator{ID: id}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_initiator_update", map[string]any{"id": 1, "comment": "updated"})
		assertResultJSON(t, res)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_initiator_update", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})
}
