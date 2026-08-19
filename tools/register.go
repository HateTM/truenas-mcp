// Package tools registers all MCP tools onto the server.
package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Config holds optional feature flags for tool registration.
type Config struct {
	// AllowDestructive enables destructive tools (delete_dataset, delete_vm, etc.).
	AllowDestructive bool
}

// RegisterAll wires all TrueNAS MCP tools onto the provided server.
func RegisterAll(s *mcp.Server, client truenasClient, cfg Config) {
	registerSystemTools(s, client)
	registerNetworkTools(s, client)
	registerFilesystemTools(s, client)
	registerPoolTools(s, client)
	registerDatasetTools(s, client)
	registerVMTools(s, client)
	registerAppTools(s, client)
	registerSnapshotTools(s, client)
	registerCloudSyncTools(s, client)
	registerJobTools(s, client)
	registerAlertTools(s, client)
	// Register pool management tools
	RegisterPoolManagementTools(s, client)
	registerPoolDatasetChoicesTools(s, client)
	registerISCSIAuthTools(s, client)
	registerISCSIExtentTools(s, client)
	registerISCSIPortalTools(s, client)
	registerISCSITargetTools(s, client)
	registerNVMeOFHostTools(s, client)
	registerNVMeOFSubsysTools(s, client)
	registerNVMeOFNamespaceTools(s, client)
	registerNVMeOFPortTools(s, client)
	registerSharingTools(s, client)
	registerNFSConfigTools(s, client)
	registerAppManagementTools(s, client)
	registerAuthTools(s, client)
	registerAlertManagementTools(s, client)
	registerCloudBackupTools(s, client)
	if cfg.AllowDestructive {
		registerDestructiveTools(s, client)
		registerISCSIDestructiveTools(s, client)
		registerNVMeOFDestructiveTools(s, client)
		registerSharingDestructiveTools(s, client)
		registerAppManagementDestructiveTools(s, client)
		registerAlertManagementDestructiveTools(s, client)
		registerCloudBackupDestructiveTools(s, client)
	}
}
