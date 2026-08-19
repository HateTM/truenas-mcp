package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// listToolNames registers cfg's domains onto a fresh server and returns the
// set of registered tool names.
func listToolNames(t *testing.T, cfg Config) map[string]bool {
	t.Helper()

	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	RegisterAll(s, &mockTruenasClient{}, cfg)

	ct, st := mcp.NewInMemoryTransports()
	ctx := context.Background()
	if _, err := s.Connect(ctx, st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	c := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.1"}, nil)
	cs, err := c.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = cs.Close() }()

	names := map[string]bool{}
	for tool, err := range cs.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("Tools iter: %v", err)
		}
		names[tool.Name] = true
	}
	return names
}

func TestRegisterAll_EmptyDomainsRegistersEverything(t *testing.T) {
	all := listToolNames(t, Config{AllowDestructive: true})

	for _, want := range []string{"list_pools", "iscsi_auth_list", "nvmeof_host_list", "sharing_nfs_list", "user_list", "cronjob_list", "dns_query"} {
		if !all[want] {
			t.Errorf("expected tool %q to be registered with empty Domains, got %d tools total", want, len(all))
		}
	}
}

func TestRegisterAll_DomainFilter(t *testing.T) {
	got := listToolNames(t, Config{AllowDestructive: true, Domains: []string{DomainISCSI}})

	for _, want := range []string{"iscsi_auth_list", "iscsi_target_create", "iscsi_auth_delete"} {
		if !got[want] {
			t.Errorf("expected iscsi tool %q to be registered when Domains=[iscsi]", want)
		}
	}
	for _, unwanted := range []string{"list_pools", "nvmeof_host_list", "sharing_nfs_list", "user_list", "delete_pool"} {
		if got[unwanted] {
			t.Errorf("did not expect tool %q to be registered when Domains=[iscsi]", unwanted)
		}
	}
}

func TestRegisterAll_LegacyDestructiveFollowsItsDomains(t *testing.T) {
	// delete_pool lives in the legacy tools/destructive.go, which spans
	// pool/dataset/vm/app/snapshot/cloudsync. It should register when any of
	// those domains is selected, and stay excluded otherwise.
	withPool := listToolNames(t, Config{AllowDestructive: true, Domains: []string{DomainPool}})
	if !withPool["delete_pool"] {
		t.Error("expected delete_pool to be registered when Domains=[pool]")
	}

	withoutPool := listToolNames(t, Config{AllowDestructive: true, Domains: []string{DomainISCSI}})
	if withoutPool["delete_pool"] {
		t.Error("did not expect delete_pool to be registered when Domains=[iscsi]")
	}
}

func TestRegisterAll_DestructiveStillGatedByAllowDestructive(t *testing.T) {
	got := listToolNames(t, Config{AllowDestructive: false, Domains: []string{DomainISCSI}})
	if got["iscsi_auth_delete"] {
		t.Error("did not expect iscsi_auth_delete to be registered when AllowDestructive=false")
	}
	if !got["iscsi_auth_list"] {
		t.Error("expected iscsi_auth_list to still be registered when AllowDestructive=false")
	}
}

func TestRegisterAll_MultipleDomains(t *testing.T) {
	got := listToolNames(t, Config{Domains: []string{DomainISCSI, DomainUser}})

	for _, want := range []string{"iscsi_auth_list", "user_list"} {
		if !got[want] {
			t.Errorf("expected tool %q with Domains=[iscsi,user]", want)
		}
	}
	if got["nvmeof_host_list"] {
		t.Error("did not expect nvmeof tool with Domains=[iscsi,user]")
	}
}
