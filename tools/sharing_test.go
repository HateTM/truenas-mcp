package tools

import (
	"context"
	"testing"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

func TestSharingNFSCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"sharing_nfs_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.NFSShare, error) {
				return []truenas.NFSShare{{ID: 1, Path: "/mnt/Storage/backups"}}, nil
			},
			"sharing_nfs_get": func(_ context.Context, id int) (*truenas.NFSShare, error) {
				return &truenas.NFSShare{ID: id}, nil
			},
			"sharing_nfs_create": func(_ context.Context, p *truenas.CreateNFSShareParams) (*truenas.NFSShare, error) {
				return &truenas.NFSShare{ID: 1, Path: p.Path}, nil
			},
			"sharing_nfs_update": func(_ context.Context, id int, _ *truenas.CreateNFSShareParams) (*truenas.NFSShare, error) {
				return &truenas.NFSShare{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "sharing_nfs_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "sharing_nfs_get", map[string]any{"id": 1}))
	})
	t.Run("get invalid id", func(t *testing.T) {
		res := callTool(t, cs, "sharing_nfs_get", map[string]any{"id": 0})
		assertError(t, res, "id must be a positive integer")
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "sharing_nfs_create", map[string]any{"path": "/mnt/Storage/backups"})
		assertResultJSON(t, res)
	})
	t.Run("create requires path", func(t *testing.T) {
		res := callTool(t, cs, "sharing_nfs_create", map[string]any{"path": ""})
		assertError(t, res, "path must not be empty")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "sharing_nfs_update", map[string]any{"id": 1, "path": "/mnt/Storage/backups"})
		assertResultJSON(t, res)
	})
}

func TestSharingSMBCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"sharing_smb_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.SMBShare, error) {
				return []truenas.SMBShare{{ID: 1, Path: "/mnt/Storage/media"}}, nil
			},
			"sharing_smb_get": func(_ context.Context, id int) (*truenas.SMBShare, error) {
				return &truenas.SMBShare{ID: id}, nil
			},
			"sharing_smb_create": func(_ context.Context, p *truenas.CreateSMBShareParams) (*truenas.SMBShare, error) {
				return &truenas.SMBShare{ID: 1, Path: p.Path}, nil
			},
			"sharing_smb_update": func(_ context.Context, id int, _ *truenas.CreateSMBShareParams) (*truenas.SMBShare, error) {
				return &truenas.SMBShare{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "sharing_smb_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "sharing_smb_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "sharing_smb_create", map[string]any{"path": "/mnt/Storage/media"})
		assertResultJSON(t, res)
	})
	t.Run("create requires path", func(t *testing.T) {
		res := callTool(t, cs, "sharing_smb_create", map[string]any{"path": ""})
		assertError(t, res, "path must not be empty")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "sharing_smb_update", map[string]any{"id": 1, "path": "/mnt/Storage/media"})
		assertResultJSON(t, res)
	})
}

func TestSharingWebDAVCRUD(t *testing.T) {
	mock := &mockTruenasClient{
		DispatchMap: map[string]any{
			"sharing_webdav_list": func(_ context.Context, _ ...truenas.ListOptions) ([]truenas.WebDAVShare, error) {
				return []truenas.WebDAVShare{{ID: 1, Name: "docs", Path: "/mnt/Storage/docs"}}, nil
			},
			"sharing_webdav_get": func(_ context.Context, id int) (*truenas.WebDAVShare, error) {
				return &truenas.WebDAVShare{ID: id}, nil
			},
			"sharing_webdav_create": func(_ context.Context, p *truenas.CreateWebDAVShareParams) (*truenas.WebDAVShare, error) {
				return &truenas.WebDAVShare{ID: 1, Name: p.Name, Path: p.Path}, nil
			},
			"sharing_webdav_update": func(_ context.Context, id int, _ *truenas.CreateWebDAVShareParams) (*truenas.WebDAVShare, error) {
				return &truenas.WebDAVShare{ID: id}, nil
			},
		},
	}
	cs, cleanup := connectTestServer(t, mock)
	defer cleanup()

	t.Run("list", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "sharing_webdav_list", nil))
	})
	t.Run("get", func(t *testing.T) {
		assertResultJSON(t, callTool(t, cs, "sharing_webdav_get", map[string]any{"id": 1}))
	})
	t.Run("create", func(t *testing.T) {
		res := callTool(t, cs, "sharing_webdav_create", map[string]any{"name": "docs", "path": "/mnt/Storage/docs"})
		assertResultJSON(t, res)
	})
	t.Run("create requires name and path", func(t *testing.T) {
		res := callTool(t, cs, "sharing_webdav_create", map[string]any{"name": "", "path": ""})
		assertError(t, res, "name and path are required")
	})
	t.Run("update", func(t *testing.T) {
		res := callTool(t, cs, "sharing_webdav_update", map[string]any{"id": 1, "name": "docs", "path": "/mnt/Storage/docs"})
		assertResultJSON(t, res)
	})
}
