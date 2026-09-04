package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerFilesystemTools registers filesystem-related MCP tools onto the server.
func registerFilesystemTools(s *mcp.Server, client truenasClient) {
	type listDirectoryInput struct {
		// Path is the absolute path on the TrueNAS host to list, e.g. "/mnt/Storage/pbs".
		Path string `json:"path" jsonschema:"Absolute path on the TrueNAS host filesystem, e.g. /mnt/Storage/pbs"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_directory",
		Description: "List the contents of a directory on the TrueNAS host filesystem. Returns name, type (FILE or DIRECTORY), size, and path for each entry.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listDirectoryInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("list_directory: path must not be empty"))
		}
		entries, err := client.ListDirectory(ctx, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("list_directory: %w", err))
		}
		return jsonResult(entries)
	})

	type makeDirectoryInput struct {
		Path string `json:"path"           jsonschema:"Absolute path of the directory to create, e.g. /mnt/Storage/backups"`
		Mode string `json:"mode,omitempty" jsonschema:"Octal permission string to set (e.g. 755); defaults to TrueNAS's default (0o755) when omitted"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "make_directory",
		Description: "Create a directory at the given absolute path on the TrueNAS host filesystem.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p makeDirectoryInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("make_directory: path must not be empty"))
		}
		if err := client.Mkdir(ctx, p.Path, p.Mode); err != nil {
			return errorResult(fmt.Errorf("make_directory: %w", err))
		}
		return jsonResult(map[string]any{"created": true, "path": p.Path})
	})

	type readFileInput struct {
		Path string `json:"path" jsonschema:"Absolute path of the file to read, e.g. /mnt/Storage/rclone.conf"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "read_file",
		Description: "Read the full contents of a file on the TrueNAS host filesystem as UTF-8 text. Not suitable for binary files.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p readFileInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("read_file: path must not be empty"))
		}
		content, err := client.ReadFile(ctx, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("read_file: %w", err))
		}
		return jsonResult(map[string]any{"path": p.Path, "content": string(content)})
	})

	type writeFileInput struct {
		Path    string `json:"path"             jsonschema:"Absolute path of the file to write, e.g. /mnt/Storage/rclone.conf"`
		Content string `json:"content"          jsonschema:"UTF-8 text content to write"`
		Append  bool   `json:"append,omitempty" jsonschema:"When true, append content to any existing file instead of overwriting it"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "write_file",
		Description: "Write text content to a file on the TrueNAS host filesystem. Overwrites the file by default; set append=true to append instead.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p writeFileInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("write_file: path must not be empty"))
		}
		if err := client.WriteFile(ctx, p.Path, []byte(p.Content), p.Append); err != nil {
			return errorResult(fmt.Errorf("write_file: %w", err))
		}
		return jsonResult(map[string]any{"written": true, "path": p.Path, "bytes": len(p.Content)})
	})

	type statFileInput struct {
		Path string `json:"path" jsonschema:"Absolute path to stat, e.g. /mnt/Storage/rclone.conf"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "stat_file",
		Description: "Get metadata (existence, type, size, ownership, mtime) for a file or directory on the TrueNAS host filesystem, without reading its contents. Cheaper than list_directory or read_file when you only need to check whether a path exists or how big it is.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p statFileInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("stat_file: path must not be empty"))
		}
		stat, err := client.Stat(ctx, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("stat_file: %w", err))
		}
		return jsonResult(stat)
	})

	type statFilesystemInput struct {
		Path string `json:"path" jsonschema:"Absolute path on the filesystem to check capacity for, e.g. /mnt/NaS"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "stat_filesystem",
		Description: "Get capacity information (total/free/available bytes and inodes) for the filesystem mounted at or containing the given absolute path. Use to check free space before or during operations that consume disk space.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p statFilesystemInput) (*mcp.CallToolResult, any, error) {
		if p.Path == "" {
			return errorResult(errors.New("stat_filesystem: path must not be empty"))
		}
		stat, err := client.StatFS(ctx, p.Path)
		if err != nil {
			return errorResult(fmt.Errorf("stat_filesystem: %w", err))
		}
		return jsonResult(stat)
	})
}
