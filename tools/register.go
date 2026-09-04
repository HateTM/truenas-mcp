// Package tools registers all MCP tools onto the server.
package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Config holds optional feature flags for tool registration.
type Config struct {
	// AllowDestructive enables destructive tools (delete_dataset, delete_vm, etc.).
	AllowDestructive bool
	// Domains restricts registration to only these TrueNAS API domains (see
	// AllDomains for valid values), reducing the tool schema payload sent to MCP
	// clients. A nil or empty slice registers every domain — the default,
	// matching pre-filter behavior.
	Domains []string
}

// RegisterAll wires TrueNAS MCP tools onto the provided server, restricted to
// cfg.Domains (or all domains, if empty) and gated by cfg.AllowDestructive.
func RegisterAll(s *mcp.Server, client truenasClient, cfg Config) {
	domains := newDomainSet(cfg.Domains)

	if domains.enabled(DomainSystem) {
		registerSystemTools(s, client)
		registerNetworkTools(s, client)
		registerFilesystemTools(s, client)
		registerJobTools(s, client)
		registerMiscTools(s, client)
	}
	if domains.enabled(DomainPool) {
		registerPoolTools(s, client)
		RegisterPoolManagementTools(s, client)
		registerPoolDatasetChoicesTools(s, client)
	}
	if domains.enabled(DomainDataset) {
		registerDatasetTools(s, client)
	}
	if domains.enabled(DomainSnapshot) {
		registerSnapshotTools(s, client)
	}
	if domains.enabled(DomainVM) {
		registerVMTools(s, client)
	}
	if domains.enabled(DomainApp) {
		registerAppTools(s, client)
		registerAppManagementTools(s, client)
	}
	if domains.enabled(DomainISCSI) {
		registerISCSIAuthTools(s, client)
		registerISCSIExtentTools(s, client)
		registerISCSIPortalTools(s, client)
		registerISCSITargetTools(s, client)
	}
	if domains.enabled(DomainNVMeOF) {
		registerNVMeOFHostTools(s, client)
		registerNVMeOFSubsysTools(s, client)
		registerNVMeOFNamespaceTools(s, client)
		registerNVMeOFPortTools(s, client)
	}
	if domains.enabled(DomainSharing) {
		registerSharingTools(s, client)
		registerNFSConfigTools(s, client)
	}
	if domains.enabled(DomainCloudSync) {
		registerCloudSyncTools(s, client)
	}
	if domains.enabled(DomainCloudBackup) {
		registerCloudBackupTools(s, client)
	}
	if domains.enabled(DomainAuth) {
		registerAuthTools(s, client)
	}
	if domains.enabled(DomainAlert) {
		registerAlertTools(s, client)
		registerAlertManagementTools(s, client)
	}
	if domains.enabled(DomainCertificate) {
		registerCertificateTools(s, client)
	}
	if domains.enabled(DomainCore) {
		registerCoreTools(s, client)
	}
	if domains.enabled(DomainDirectoryServices) {
		registerDirectoryServicesTools(s, client)
	}
	if domains.enabled(DomainDisk) {
		registerDiskTools(s, client)
	}
	if domains.enabled(DomainReplication) {
		registerReplicationTools(s, client)
		registerReplicationEndpointTools(s, client)
	}
	if domains.enabled(DomainUser) {
		registerUserTools(s, client)
	}
	if domains.enabled(DomainCronJob) {
		registerCronJobTools(s, client)
	}

	if !cfg.AllowDestructive {
		return
	}
	// registerDestructiveTools predates the per-domain destructive_*.go files
	// and spans several domains in one function (see tools/destructive.go);
	// register it if any of those domains is enabled.
	if domains.enabledAny(DomainPool, DomainDataset, DomainVM, DomainApp, DomainSnapshot, DomainCloudSync) {
		registerDestructiveTools(s, client)
	}
	if domains.enabled(DomainISCSI) {
		registerISCSIDestructiveTools(s, client)
	}
	if domains.enabled(DomainNVMeOF) {
		registerNVMeOFDestructiveTools(s, client)
	}
	if domains.enabled(DomainSharing) {
		registerSharingDestructiveTools(s, client)
	}
	if domains.enabled(DomainApp) {
		registerAppManagementDestructiveTools(s, client)
	}
	if domains.enabled(DomainAlert) {
		registerAlertManagementDestructiveTools(s, client)
	}
	if domains.enabled(DomainCloudBackup) {
		registerCloudBackupDestructiveTools(s, client)
	}
	if domains.enabled(DomainCertificate) {
		registerCertificateDestructiveTools(s, client)
	}
	if domains.enabled(DomainDisk) {
		registerDiskDestructiveTools(s, client)
	}
	if domains.enabled(DomainReplication) {
		registerReplicationDestructiveTools(s, client)
	}
	if domains.enabled(DomainUser) {
		registerUserDestructiveTools(s, client)
	}
	if domains.enabled(DomainCronJob) {
		registerCronJobDestructiveTools(s, client)
	}
}
