package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestISCSITargetList(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_target_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ISCSITarget, error) {
				return []truenas.ISCSITarget{{ID: 1, Name: "target0"}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_target_list", nil)
	assertResultJSON(t, res)
}

func TestISCSITargetGet(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_target_get": func(_ context.Context, id int) (*truenas.ISCSITarget, error) {
				return &truenas.ISCSITarget{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_target_get", map[string]any{"id": 1})
	assertResultJSON(t, res)
}

func TestISCSITargetCreate(t *testing.T) {
	t.Run("creates target", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_target_create": func(_ context.Context, p *truenas.CreateISCSITargetParams) (*truenas.ISCSITarget, error) {
					return &truenas.ISCSITarget{ID: 1, Name: p.Name, Groups: p.Groups}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_target_create", map[string]any{
			"name":   "target0",
			"groups": []map[string]any{{"portal": 1}},
		})
		assertResultJSON(t, res)
	})

	t.Run("requires name", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_target_create", map[string]any{"name": ""})
		assertError(t, res, "name must not be empty")
	})
}

func TestISCSITargetUpdate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_target_update": func(_ context.Context, id int, _ *truenas.CreateISCSITargetParams) (*truenas.ISCSITarget, error) {
				return &truenas.ISCSITarget{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_target_update", map[string]any{"id": 1, "name": "target0"})
	assertResultJSON(t, res)
}

func TestISCSITargetValidateName(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_target_validate_name": func(_ context.Context, name string) (bool, error) {
				return name == "available-name", nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_target_validate_name", map[string]any{"name": "available-name"})
	assertResultJSON(t, res)
}

func TestISCSITargetExtentList(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_targetextent_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.ISCSITargetExtent, error) {
				return []truenas.ISCSITargetExtent{{ID: 1, Target: 1, Extent: 1}}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_targetextent_list", nil)
	assertResultJSON(t, res)
}

func TestISCSITargetExtentGet(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_targetextent_get": func(_ context.Context, id int) (*truenas.ISCSITargetExtent, error) {
				return &truenas.ISCSITargetExtent{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_targetextent_get", map[string]any{"id": 1})
	assertResultJSON(t, res)
}

func TestISCSITargetExtentCreate(t *testing.T) {
	t.Run("creates mapping", func(t *testing.T) {
		mock := &mockTruenasClient{
			DispatchMap: map[string]any{
				"iscsi_targetextent_create": func(_ context.Context, p *truenas.CreateISCSITargetExtentParams) (*truenas.ISCSITargetExtent, error) {
					return &truenas.ISCSITargetExtent{ID: 1, Target: p.Target, Extent: p.Extent}, nil
				},
			},
		}
		cs, cleanup := connectTestServer(t, mock)
		defer cleanup()

		res := callTool(t, cs, "iscsi_targetextent_create", map[string]any{"target": 1, "extent": 1})
		assertResultJSON(t, res)
	})

	t.Run("requires target and extent", func(t *testing.T) {
		cs, cleanup := connectTestServer(t, &mockTruenasClient{})
		defer cleanup()

		res := callTool(t, cs, "iscsi_targetextent_create", map[string]any{"target": 0, "extent": 0})
		assertError(t, res, "target and extent are required")
	})
}

func TestISCSITargetExtentUpdate(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"iscsi_targetextent_update": func(_ context.Context, id int, _ *truenas.CreateISCSITargetExtentParams) (*truenas.ISCSITargetExtent, error) {
				return &truenas.ISCSITargetExtent{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	res := callTool(t, cs, "iscsi_targetextent_update", map[string]any{"id": 1, "target": 1, "extent": 1})
	assertResultJSON(t, res)
}
