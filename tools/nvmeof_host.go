package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerNVMeOFHostTools registers nvmet.global.* and nvmet.host.* MCP tools.
func registerNVMeOFHostTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_global_config",
		Description: "Get the global NVMe-oF service configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.NVMetGlobalConfigGet(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_global_config: %w", err))
		}
		return jsonResult(cfg)
	})

	type updateNVMeOFGlobalInput struct {
		ANA  bool `json:"ana,omitempty"  jsonschema:"Enable Asymmetric Namespace Access"`
		RDMA bool `json:"rdma,omitempty" jsonschema:"Enable RDMA transport support"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_global_update",
		Description: "Update the global NVMe-oF service configuration.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNVMeOFGlobalInput) (*mcp.CallToolResult, any, error) {
		cfg, err := client.UpdateNVMetGlobal(ctx, &truenas.UpdateNVMetGlobalParams{ANA: p.ANA, RDMA: p.RDMA})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_global_update: %w", err))
		}
		return jsonResult(cfg)
	})

	type listNVMeOFHostsInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of hosts to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of hosts to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_list",
		Description: "List configured NVMe-oF hosts.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listNVMeOFHostsInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListNVMetHosts(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_list: %w", err))
		}
		return jsonResult(result)
	})

	type getNVMeOFHostInput struct {
		ID int `json:"id" jsonschema:"Numeric NVMe-oF host ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_get",
		Description: "Get a single NVMe-oF host by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getNVMeOFHostInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_host_get: id must be a positive integer"))
		}
		result, err := client.GetNVMetHost(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_get: %w", err))
		}
		return jsonResult(result)
	})

	type createNVMeOFHostInput struct {
		HostNQN       string `json:"hostnqn"                  jsonschema:"Host NVMe Qualified Name, e.g. nqn.2014-08.org.nvmexpress:uuid:..."`
		DHCHAPKey     string `json:"dhchap_key,omitempty"     jsonschema:"DH-HMAC-CHAP key (from nvmeof_host_generate_key) for authenticating this host"`
		DHCHAPDHGroup string `json:"dhchap_dhgroup,omitempty" jsonschema:"DH group for DH-HMAC-CHAP, from nvmeof_host_dhchap_dhgroup_choices"`
		DHCHAPHash    string `json:"dhchap_hash,omitempty"    jsonschema:"Hash algorithm for DH-HMAC-CHAP, from nvmeof_host_dhchap_hash_choices"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_create",
		Description: "Create a new NVMe-oF host entry, to be bound to subsystems via nvmeof_host_subsys_create.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p createNVMeOFHostInput) (*mcp.CallToolResult, any, error) {
		if p.HostNQN == "" {
			return errorResult(errors.New("nvmeof_host_create: hostnqn must not be empty"))
		}
		result, err := client.CreateNVMetHost(ctx, &truenas.CreateNVMetHostParams{
			HostNQN: p.HostNQN, DHCHAPKey: p.DHCHAPKey, DHCHAPDHGroup: p.DHCHAPDHGroup, DHCHAPHash: p.DHCHAPHash,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateNVMeOFHostInput struct {
		ID            int    `json:"id"                       jsonschema:"Numeric NVMe-oF host ID"`
		HostNQN       string `json:"hostnqn"                  jsonschema:"Host NVMe Qualified Name"`
		DHCHAPKey     string `json:"dhchap_key,omitempty"     jsonschema:"DH-HMAC-CHAP key"`
		DHCHAPDHGroup string `json:"dhchap_dhgroup,omitempty" jsonschema:"DH group for DH-HMAC-CHAP"`
		DHCHAPHash    string `json:"dhchap_hash,omitempty"    jsonschema:"Hash algorithm for DH-HMAC-CHAP"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_update",
		Description: "Update an existing NVMe-oF host.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateNVMeOFHostInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("nvmeof_host_update: id must be a positive integer"))
		}
		if p.HostNQN == "" {
			return errorResult(errors.New("nvmeof_host_update: hostnqn must not be empty"))
		}
		result, err := client.UpdateNVMetHost(ctx, p.ID, &truenas.CreateNVMetHostParams{
			HostNQN: p.HostNQN, DHCHAPKey: p.DHCHAPKey, DHCHAPDHGroup: p.DHCHAPDHGroup, DHCHAPHash: p.DHCHAPHash,
		})
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_update: %w", err))
		}
		return jsonResult(result)
	})

	// nvmeof_host_delete is destructive and lives in destructive.go, gated
	// behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_dhchap_dhgroup_choices",
		Description: "List the Diffie-Hellman groups available for NVMe-oF host DH-HMAC-CHAP authentication.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.NVMetHostDHCHAPDHGroupChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_dhchap_dhgroup_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_dhchap_hash_choices",
		Description: "List the hash algorithms available for NVMe-oF host DH-HMAC-CHAP authentication.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.NVMetHostDHCHAPHashChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_dhchap_hash_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "nvmeof_host_generate_key",
		Description: "Generate a new random DH-HMAC-CHAP key suitable for use as an NVMe-oF host's dhchap_key.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		key, err := client.NVMetHostGenerateKey(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("nvmeof_host_generate_key: %w", err))
		}
		return jsonResult(map[string]string{"key": key})
	})
}
