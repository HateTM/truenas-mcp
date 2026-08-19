package tools

// Domain names for Config.Domains. Each corresponds to one or more
// register*Tools functions grouped by TrueNAS API subsystem — see RegisterAll.
const (
	DomainSystem            = "system"
	DomainPool              = "pool"
	DomainDataset           = "dataset"
	DomainSnapshot          = "snapshot"
	DomainVM                = "vm"
	DomainApp               = "app"
	DomainISCSI             = "iscsi"
	DomainNVMeOF            = "nvmeof"
	DomainSharing           = "sharing"
	DomainCloudSync         = "cloudsync"
	DomainCloudBackup       = "cloudbackup"
	DomainAuth              = "auth"
	DomainAlert             = "alert"
	DomainCertificate       = "certificate"
	DomainCore              = "core"
	DomainDirectoryServices = "directoryservices"
	DomainDisk              = "disk"
	DomainReplication       = "replication"
	DomainUser              = "user"
	DomainCronJob           = "cronjob"
)

// AllDomains lists every valid Config.Domains value, in registration order.
var AllDomains = []string{
	DomainSystem, DomainPool, DomainDataset, DomainSnapshot, DomainVM, DomainApp,
	DomainISCSI, DomainNVMeOF, DomainSharing, DomainCloudSync, DomainCloudBackup,
	DomainAuth, DomainAlert, DomainCertificate, DomainCore, DomainDirectoryServices,
	DomainDisk, DomainReplication, DomainUser, DomainCronJob,
}

// domainSet controls which domains RegisterAll registers. A nil domainSet means
// "all domains enabled" (Config.Domains was empty), matching the pre-filter
// default of registering every tool.
type domainSet map[string]bool

// newDomainSet builds a domainSet from Config.Domains. An empty slice yields a
// nil domainSet (all domains enabled).
func newDomainSet(domains []string) domainSet {
	if len(domains) == 0 {
		return nil
	}
	s := make(domainSet, len(domains))
	for _, d := range domains {
		s[d] = true
	}
	return s
}

// enabled reports whether domain should be registered.
func (s domainSet) enabled(domain string) bool {
	return s == nil || s[domain]
}

// enabledAny reports whether any of the given domains should be registered.
// Used for the legacy destructive.go tools, which predate per-domain
// destructive files and span several domains (pool, dataset, vm, app,
// snapshot, cloudsync) in one function.
func (s domainSet) enabledAny(domains ...string) bool {
	if s == nil {
		return true
	}
	for _, d := range domains {
		if s[d] {
			return true
		}
	}
	return false
}
