package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerUserTools registers user.* MCP tools onto the server, except user.delete and
// user.update_password, which are destructive/sensitive and live in destructive_user.go.
func registerUserTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listUsersInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of users to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of users to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_list",
		Description: "List local TrueNAS user accounts.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listUsersInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListUsers(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("user_list: %w", err))
		}
		return jsonResult(result)
	})

	type getUserInput struct {
		ID int `json:"id" jsonschema:"Numeric user ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_get",
		Description: "Get a single local user account by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getUserInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("user_get: id must be a positive integer"))
		}
		result, err := client.GetUser(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("user_get: %w", err))
		}
		return jsonResult(result)
	})

	type createUserInput struct {
		Username      string `json:"username"                jsonschema:"Login username"`
		FullName      string `json:"full_name"                jsonschema:"Full display name"`
		GroupID       int    `json:"group,omitempty"          jsonschema:"Primary group ID; a new private group is created if omitted"`
		Password      string `json:"password,omitempty"       jsonschema:"Login password; omit to create a password-disabled account"`
		HomeDirectory string `json:"home,omitempty"           jsonschema:"Home directory path, from user_home_directory_choices"`
		Shell         string `json:"shell,omitempty"          jsonschema:"Login shell path, from user_shell_choices"`
		Email         string `json:"email,omitempty"          jsonschema:"Email address"`
		Locked        bool   `json:"locked,omitempty"         jsonschema:"Whether the account is locked (login disabled)"`
		SMB           bool   `json:"smb,omitempty"            jsonschema:"Whether this account can authenticate to SMB shares"`
		SSHPubKey     string `json:"sshpubkey,omitempty"      jsonschema:"SSH public key for key-based authentication"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_create",
		Description: "Create a new local user account.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createUserInput) (*mcp.CallToolResult, any, error) {
		if p.Username == "" || p.FullName == "" {
			return errorResult(errors.New("user_create: username and full_name are required"))
		}
		result, err := client.CreateUser(ctx, &truenas.CreateUserParams{
			Username: p.Username, FullName: p.FullName, GroupID: p.GroupID, Password: p.Password,
			HomeDirectory: p.HomeDirectory, Shell: p.Shell, Email: p.Email, Locked: p.Locked,
			SMB: p.SMB, SSHPubKey: p.SSHPubKey,
		})
		if err != nil {
			return errorResult(fmt.Errorf("user_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateUserInput struct {
		ID            int    `json:"id"                       jsonschema:"Numeric user ID"`
		Username      string `json:"username"                jsonschema:"Login username"`
		FullName      string `json:"full_name"                jsonschema:"Full display name"`
		GroupID       int    `json:"group,omitempty"          jsonschema:"Primary group ID"`
		HomeDirectory string `json:"home,omitempty"           jsonschema:"Home directory path"`
		Shell         string `json:"shell,omitempty"          jsonschema:"Login shell path"`
		Email         string `json:"email,omitempty"          jsonschema:"Email address"`
		Locked        bool   `json:"locked,omitempty"         jsonschema:"Whether the account is locked"`
		SMB           bool   `json:"smb,omitempty"            jsonschema:"Whether this account can authenticate to SMB shares"`
		SSHPubKey     string `json:"sshpubkey,omitempty"      jsonschema:"SSH public key"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_update",
		Description: "Update an existing local user account. Use user_update_password to change the password.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateUserInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("user_update: id must be a positive integer"))
		}
		if p.Username == "" || p.FullName == "" {
			return errorResult(errors.New("user_update: username and full_name are required"))
		}
		result, err := client.UpdateUser(ctx, p.ID, &truenas.CreateUserParams{
			Username: p.Username, FullName: p.FullName, GroupID: p.GroupID, HomeDirectory: p.HomeDirectory,
			Shell: p.Shell, Email: p.Email, Locked: p.Locked, SMB: p.SMB, SSHPubKey: p.SSHPubKey,
		})
		if err != nil {
			return errorResult(fmt.Errorf("user_update: %w", err))
		}
		return jsonResult(result)
	})

	// user_delete and user_update_password are destructive/sensitive and live in
	// destructive_user.go, gated behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_schemas",
		Description: "Get the JSON schemas describing valid user account configurations.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		schemas, err := client.UserSchemas(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("user_schemas: %w", err))
		}
		return jsonResult(schemas)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_home_directory_choices",
		Description: "List the paths available as a new user's home directory.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.UserHomeDirectoryChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("user_home_directory_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "user_shell_choices",
		Description: "List the shells available for a user account.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.UserShellChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("user_shell_choices: %w", err))
		}
		return jsonResult(choices)
	})
}
