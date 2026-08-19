package truenas

import (
	"context"
	"fmt"
	"strings"
)

// DirEntry represents a single file or directory returned by ListDirectory.
type DirEntry struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	RealPath string `json:"realpath,omitempty"`
	// Type is "FILE" or "DIRECTORY".
	Type string `json:"type"`
	Size int64  `json:"size"`
	Mode int    `json:"mode"`
}

// ListDirectory returns the contents of the given absolute path on the TrueNAS
// host filesystem (e.g. "/mnt/Storage/pbs").
func (c *Client) ListDirectory(ctx context.Context, path string) ([]DirEntry, error) {
	if path == "" {
		return nil, fmt.Errorf("listing directory: path must not be empty")
	}
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("listing directory: path must be absolute (start with /), got %q", path)
	}
	var entries []DirEntry
	if err := c.call(ctx, "filesystem.listdir", []any{path}, &entries); err != nil {
		return nil, fmt.Errorf("listing directory %q: %w", path, err)
	}
	return entries, nil
}

// Mkdir creates a directory.
func (c *Client) Mkdir(ctx context.Context, path, mode string) error {
	if err := c.call(ctx, "filesystem.mkdir", []any{path, mode}, nil); err != nil {
		return fmt.Errorf("creating directory %q: %w", path, err)
	}
	return nil
}

// WriteFile writes data to a file.
func (c *Client) WriteFile(ctx context.Context, path string, content []byte, appendFlag bool) error {
	// TrueNAS API for writing files might require base64 encoding or raw data depending on the version.
	// Assuming raw data for now.
	if err := c.call(ctx, "filesystem.file_write", []any{path, content, appendFlag}, nil); err != nil {
		return fmt.Errorf("writing file %q: %w", path, err)
	}
	return nil
}

// ReadFile reads the contents of a file.
func (c *Client) ReadFile(ctx context.Context, path string) ([]byte, error) {
	var content []byte
	if err := c.call(ctx, "filesystem.file_read", []any{path}, &content); err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}
	return content, nil
}

// Stat returns metadata for a file or directory.
func (c *Client) Stat(ctx context.Context, path string) (*FileStat, error) {
	var stat FileStat
	if err := c.call(ctx, "filesystem.stat", []any{path}, &stat); err != nil {
		return nil, fmt.Errorf("stating %q: %w", path, err)
	}
	return &stat, nil
}

// StatFS returns capacity information for a filesystem.
func (c *Client) StatFS(ctx context.Context, path string) (*FSStat, error) {
	var stat FSStat
	if err := c.call(ctx, "filesystem.statfs", []any{path}, &stat); err != nil {
		return nil, fmt.Errorf("stating filesystem %q: %w", path, err)
	}
	return &stat, nil
}

// FileStat represents metadata for a file or directory.
type FileStat struct {
	Type         string `json:"type"`
	Size         int64  `json:"size"`
	Mode         int    `json:"mode"`
	IsMountpoint bool   `json:"is_mountpoint"`
}

// FSStat represents capacity information for a filesystem.
type FSStat struct {
	TotalBytes     int64 `json:"total_bytes"`
	AvailableBytes int64 `json:"available_bytes"`
	FreeBytes      int64 `json:"free_bytes"`
}
