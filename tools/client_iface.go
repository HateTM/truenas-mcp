package tools

import (
	"context"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

// truenasClient is the subset of *truenas.Client used by the tools layer.
// It exists to allow mock implementations in tests without requiring a live TrueNAS server.
type truenasClient interface {
	// System
	GetSystemInfo(ctx context.Context) (*truenas.SystemInfo, error)

	// Pools
	ListPools(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.Pool, error)
	GetPool(ctx context.Context, id int) (*truenas.Pool, error)
	CreatePool(ctx context.Context, params *truenas.CreatePoolParams) (*truenas.Pool, error)
	UpdatePool(ctx context.Context, id int, params *truenas.UpdatePoolParams) (*truenas.Pool, error)
	DeletePool(ctx context.Context, id int) error
	AttachPool(ctx context.Context, id int) (*truenas.Pool, error)
	DetachPool(ctx context.Context, id int) error
	ExpandPool(ctx context.Context, params *truenas.ExpandPoolParams) (*truenas.Pool, error)
	ExportPool(ctx context.Context, id int, path string) (int64, error)
	ListDisks(ctx context.Context, id int) ([]truenas.DiskInfo, error)
	ExportPoolHTML(ctx context.Context, id int) (string, error)
	FindDatasetsForImport(ctx context.Context, id int, opts truenas.ListOptions) ([]truenas.Dataset, error)
	ImportPool(ctx context.Context, id int, path, name string) error
	IsPoolUpgraded(ctx context.Context, id int) (bool, error)
	OfflinePool(ctx context.Context, id int) error
	OnlinePool(ctx context.Context, id int) error
	RelabelPool(ctx context.Context, id int, key string, force bool) error
	ScanPool(ctx context.Context, id int, force bool) error
	SplitPool(ctx context.Context, id int, path, name string) error
	GetPoolStatus(ctx context.Context, id int) (*truenas.PoolStatus, error)
	SyncPool(ctx context.Context, id int) error
	VerifyPool(ctx context.Context, id int) error
	DatasetChecksumChoices(ctx context.Context) (map[string]string, error)
	DatasetCompressionChoices(ctx context.Context) (map[string]string, error)
	DatasetEncryptionAlgorithmChoices(ctx context.Context) (map[string]string, error)
	DatasetRecordsizeChoices(ctx context.Context) ([]string, error)
	DDTPrune(ctx context.Context, poolID int, p truenas.DDTPruneParams) error
	PoolFilesystemChoices(ctx context.Context, poolIDs []int) ([]string, error)

	// Datasets
	ListDatasets(ctx context.Context, pool string, opts ...truenas.ListOptions) ([]truenas.Dataset, error)
	GetDataset(ctx context.Context, id string) (*truenas.Dataset, error)
	CreateDataset(ctx context.Context, params *truenas.CreateDatasetParams) (*truenas.Dataset, error)
	UpdateDataset(ctx context.Context, id string, params *truenas.UpdateDatasetParams) (*truenas.Dataset, error)
	DeleteDataset(ctx context.Context, id string) error
	DeleteDatasetByPath(ctx context.Context, path string, recursive bool) error
	ExportDatasetKey(ctx context.Context, poolID int, datasetName, path string) (string, error)
	ExportDatasetKeys(ctx context.Context, poolID int, path string) (int64, error)
	ExportDatasetKeysForReplication(ctx context.Context, poolID int, remote, path string) (int64, error)
	LockDatasetKey(ctx context.Context, poolID int, datasetName string) error
	UnlockDatasetKey(ctx context.Context, poolID int, datasetName, key string) error
	PromoteDataset(ctx context.Context, poolID int, datasetName string) error
	DDTPrefetch(ctx context.Context, poolID int, enable bool) error
	GetDatasetInstance(ctx context.Context, datasetName string) (*truenas.Dataset, error)
	GetDatasetQuota(ctx context.Context, datasetName string) (*truenas.QuotaInfo, error)
	SetDatasetQuota(ctx context.Context, datasetName string, quota int64) error
	RenameDataset(ctx context.Context, poolID int, oldName, newName string) (*truenas.Dataset, error)
	ImportDataset(ctx context.Context, poolID int, datasetName, path string) error
	CheckQuota(ctx context.Context, poolID int, datasetName string) (*truenas.QuotaStatus, error)
	SyncDataset(ctx context.Context, poolID int, datasetName string) error
	OfflineDataset(ctx context.Context, poolID int, datasetName string) error
	OnlineDataset(ctx context.Context, poolID int, datasetName string) error
	GetDiskInfo(ctx context.Context, poolID int) (*truenas.DiskInfo, error)
	OfflineDisk(ctx context.Context, poolID int, diskID int) error
	OnlineDisk(ctx context.Context, poolID int, diskID int) error
	RemoveDisk(ctx context.Context, poolID int, diskID int) error
	ReplaceDisk(ctx context.Context, poolID int, diskID int, path string) error
	ImportPoolHTML(ctx context.Context, path string) error

	// Snapshots
	ListSnapshots(ctx context.Context, dataset string, opts ...truenas.ListOptions) ([]truenas.Snapshot, error)
	GetSnapshot(ctx context.Context, id string) (*truenas.Snapshot, error)
	CreateSnapshot(ctx context.Context, params truenas.CreateSnapshotParams) (*truenas.Snapshot, error)
	RollbackSnapshot(ctx context.Context, id string, params truenas.RollbackSnapshotParams) error
	DeleteSnapshot(ctx context.Context, id string) error

	// Virtual Machines
	ListVMs(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.VM, error)
	GetVM(ctx context.Context, id int) (*truenas.VM, error)
	StartVM(ctx context.Context, id int) (int, error)
	StopVM(ctx context.Context, id int, force bool) (int, error)
	RestartVM(ctx context.Context, id int) (int, error)
	CreateVM(ctx context.Context, params *truenas.CreateVMParams) (*truenas.VM, error)
	UpdateVM(ctx context.Context, id int, params *truenas.UpdateVMParams) (*truenas.VM, error)
	DeleteVM(ctx context.Context, id int) error
	ListVMDevices(ctx context.Context, vmID int) ([]truenas.VMDevice, error)
	AddVMDevice(ctx context.Context, params *truenas.AddVMDeviceParams) (*truenas.VMDevice, error)
	DeleteVMDevice(ctx context.Context, deviceID int) error

	// Apps
	ListApps(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.App, error)
	GetApp(ctx context.Context, name string) (*truenas.App, error)
	StartApp(ctx context.Context, name string) (int, error)
	StopApp(ctx context.Context, name string) (int, error)
	RestartApp(ctx context.Context, name string) (int, error)
	ListImages(ctx context.Context) ([]truenas.Image, error)
	CreateApp(ctx context.Context, params *truenas.CreateAppParams) (int, error)
	DeleteApp(ctx context.Context, name string) error
	UpgradeApp(ctx context.Context, name, version string) (int, error)
	GetUpgradeSummary(ctx context.Context, name string) (*truenas.AppUpgradeSummary, error)
	RollbackApp(ctx context.Context, name, version string) (int, error)

	// Network
	ListInterfaces(ctx context.Context) ([]truenas.Interface, error)
	GetInterface(ctx context.Context, id string) (*truenas.Interface, error)
	UpdateInterface(ctx context.Context, id string, params *truenas.UpdateInterfaceParams) (*truenas.Interface, error)

	// Filesystem
	ListDirectory(ctx context.Context, path string) ([]truenas.DirEntry, error)
	Mkdir(ctx context.Context, path, mode string) error
	WriteFile(ctx context.Context, path string, content []byte, appendFlag bool) error
	ReadFile(ctx context.Context, path string) ([]byte, error)
	Stat(ctx context.Context, path string) (*truenas.FileStat, error)
	StatFS(ctx context.Context, path string) (*truenas.FSStat, error)

	// Cloud Sync
	ListCloudProviders(ctx context.Context) ([]truenas.CloudProvider, error)
	ListCloudCredentials(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CloudCredential, error)
	CreateCloudCredential(ctx context.Context, params *truenas.CreateCloudCredentialParams) (*truenas.CloudCredential, error)
	DeleteCloudCredential(ctx context.Context, id int) error
	ListCloudSyncTasks(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CloudSyncTask, error)
	CreateCloudSyncTask(ctx context.Context, params *truenas.CreateCloudSyncTaskParams) (*truenas.CloudSyncTask, error)
	RunCloudSyncTask(ctx context.Context, id int, dryRun bool) (int, error)
	AbortCloudSyncTask(ctx context.Context, id int) error
	DeleteCloudSyncTask(ctx context.Context, id int) error

	// Jobs
	GetJob(ctx context.Context, id int) (*truenas.Job, error)
	AbortJob(ctx context.Context, id int) error

	// Alerts
	ListAlerts(ctx context.Context) ([]truenas.Alert, error)

	// iSCSI
	ListISCSIAuth(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIAuth, error)
	GetISCSIAuth(ctx context.Context, id int) (*truenas.ISCSIAuth, error)
	CreateISCSIAuth(ctx context.Context, p *truenas.CreateISCSIAuthParams) (*truenas.ISCSIAuth, error)
	UpdateISCSIAuth(ctx context.Context, id int, p *truenas.CreateISCSIAuthParams) (*truenas.ISCSIAuth, error)
	DeleteISCSIAuth(ctx context.Context, id int) error
	ISCSIGlobalALUAEnabled(ctx context.Context) (bool, error)
	ISCSIGlobalClientCount(ctx context.Context) (int, error)
	ISCSIGlobalConfigGet(ctx context.Context) (*truenas.ISCSIGlobalConfig, error)
	ISCSIGlobalISEREnabled(ctx context.Context) (bool, error)
	ISCSIGlobalSessions(ctx context.Context) ([]truenas.ISCSISession, error)
	UpdateISCSIGlobal(ctx context.Context, p *truenas.UpdateISCSIGlobalParams) (*truenas.ISCSIGlobalConfig, error)
	ListISCSIExtents(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIExtent, error)
	GetISCSIExtent(ctx context.Context, id int) (*truenas.ISCSIExtent, error)
	CreateISCSIExtent(ctx context.Context, p *truenas.CreateISCSIExtentParams) (*truenas.ISCSIExtent, error)
	UpdateISCSIExtent(ctx context.Context, id int, p *truenas.CreateISCSIExtentParams) (*truenas.ISCSIExtent, error)
	DeleteISCSIExtent(ctx context.Context, id int) error
	ISCSIExtentDiskChoices(ctx context.Context) (map[string]string, error)
	ListISCSIInitiators(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIInitiator, error)
	GetISCSIInitiator(ctx context.Context, id int) (*truenas.ISCSIInitiator, error)
	CreateISCSIInitiator(ctx context.Context, p *truenas.CreateISCSIInitiatorParams) (*truenas.ISCSIInitiator, error)
	UpdateISCSIInitiator(ctx context.Context, id int, p *truenas.CreateISCSIInitiatorParams) (*truenas.ISCSIInitiator, error)
	DeleteISCSIInitiator(ctx context.Context, id int) error
	ListISCSIPortals(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIPortal, error)
	GetISCSIPortal(ctx context.Context, id int) (*truenas.ISCSIPortal, error)
	CreateISCSIPortal(ctx context.Context, p *truenas.CreateISCSIPortalParams) (*truenas.ISCSIPortal, error)
	UpdateISCSIPortal(ctx context.Context, id int, p *truenas.CreateISCSIPortalParams) (*truenas.ISCSIPortal, error)
	DeleteISCSIPortal(ctx context.Context, id int) error
	ISCSIPortalListenIPChoices(ctx context.Context) (map[string]string, error)
	ListISCSITargets(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSITarget, error)
	GetISCSITarget(ctx context.Context, id int) (*truenas.ISCSITarget, error)
	CreateISCSITarget(ctx context.Context, p *truenas.CreateISCSITargetParams) (*truenas.ISCSITarget, error)
	UpdateISCSITarget(ctx context.Context, id int, p *truenas.CreateISCSITargetParams) (*truenas.ISCSITarget, error)
	DeleteISCSITarget(ctx context.Context, id int) error
	ValidateISCSITargetName(ctx context.Context, name string) (bool, error)
	ListISCSITargetExtents(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSITargetExtent, error)
	GetISCSITargetExtent(ctx context.Context, id int) (*truenas.ISCSITargetExtent, error)
	CreateISCSITargetExtent(ctx context.Context, p *truenas.CreateISCSITargetExtentParams) (*truenas.ISCSITargetExtent, error)
	UpdateISCSITargetExtent(ctx context.Context, id int, p *truenas.CreateISCSITargetExtentParams) (*truenas.ISCSITargetExtent, error)
	DeleteISCSITargetExtent(ctx context.Context, id int) error

	// NVMe-oF
	NVMetGlobalConfigGet(ctx context.Context) (*truenas.NVMetGlobalConfig, error)
	UpdateNVMetGlobal(ctx context.Context, p *truenas.UpdateNVMetGlobalParams) (*truenas.NVMetGlobalConfig, error)
	ListNVMetHosts(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetHost, error)
	GetNVMetHost(ctx context.Context, id int) (*truenas.NVMetHost, error)
	CreateNVMetHost(ctx context.Context, p *truenas.CreateNVMetHostParams) (*truenas.NVMetHost, error)
	UpdateNVMetHost(ctx context.Context, id int, p *truenas.CreateNVMetHostParams) (*truenas.NVMetHost, error)
	DeleteNVMetHost(ctx context.Context, id int) error
	NVMetHostDHCHAPDHGroupChoices(ctx context.Context) (map[string]string, error)
	NVMetHostDHCHAPHashChoices(ctx context.Context) (map[string]string, error)
	NVMetHostGenerateKey(ctx context.Context) (string, error)
	ListNVMetHostSubsys(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetHostSubsys, error)
	GetNVMetHostSubsys(ctx context.Context, id int) (*truenas.NVMetHostSubsys, error)
	CreateNVMetHostSubsys(ctx context.Context, p *truenas.CreateNVMetHostSubsysParams) (*truenas.NVMetHostSubsys, error)
	UpdateNVMetHostSubsys(ctx context.Context, id int, p *truenas.CreateNVMetHostSubsysParams) (*truenas.NVMetHostSubsys, error)
	DeleteNVMetHostSubsys(ctx context.Context, id int) error
	ListNVMetSubsys(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetSubsys, error)
	GetNVMetSubsys(ctx context.Context, id int) (*truenas.NVMetSubsys, error)
	CreateNVMetSubsys(ctx context.Context, p *truenas.CreateNVMetSubsysParams) (*truenas.NVMetSubsys, error)
	UpdateNVMetSubsys(ctx context.Context, id int, p *truenas.CreateNVMetSubsysParams) (*truenas.NVMetSubsys, error)
	DeleteNVMetSubsys(ctx context.Context, id int) error
	ListNVMetNamespaces(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetNamespace, error)
	GetNVMetNamespace(ctx context.Context, id int) (*truenas.NVMetNamespace, error)
	CreateNVMetNamespace(ctx context.Context, p *truenas.CreateNVMetNamespaceParams) (*truenas.NVMetNamespace, error)
	UpdateNVMetNamespace(ctx context.Context, id int, p *truenas.CreateNVMetNamespaceParams) (*truenas.NVMetNamespace, error)
	DeleteNVMetNamespace(ctx context.Context, id int) error
	ListNVMetPortSubsys(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetPortSubsys, error)
	GetNVMetPortSubsys(ctx context.Context, id int) (*truenas.NVMetPortSubsys, error)
	CreateNVMetPortSubsys(ctx context.Context, p *truenas.CreateNVMetPortSubsysParams) (*truenas.NVMetPortSubsys, error)
	UpdateNVMetPortSubsys(ctx context.Context, id int, p *truenas.CreateNVMetPortSubsysParams) (*truenas.NVMetPortSubsys, error)
	DeleteNVMetPortSubsys(ctx context.Context, id int) error
	ListNVMetPorts(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetPort, error)
	GetNVMetPort(ctx context.Context, id int) (*truenas.NVMetPort, error)
	CreateNVMetPort(ctx context.Context, p *truenas.CreateNVMetPortParams) (*truenas.NVMetPort, error)
	UpdateNVMetPort(ctx context.Context, id int, p *truenas.CreateNVMetPortParams) (*truenas.NVMetPort, error)
	DeleteNVMetPort(ctx context.Context, id int) error
	NVMetPortTransportAddressChoices(ctx context.Context, trtype string) (map[string]string, error)

	// Sharing (NFS/SMB/WebDAV)
	ListNFSShares(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NFSShare, error)
	GetNFSShare(ctx context.Context, id int) (*truenas.NFSShare, error)
	CreateNFSShare(ctx context.Context, p *truenas.CreateNFSShareParams) (*truenas.NFSShare, error)
	UpdateNFSShare(ctx context.Context, id int, p *truenas.CreateNFSShareParams) (*truenas.NFSShare, error)
	DeleteNFSShare(ctx context.Context, id int) error
	ListSMBShares(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.SMBShare, error)
	GetSMBShare(ctx context.Context, id int) (*truenas.SMBShare, error)
	CreateSMBShare(ctx context.Context, p *truenas.CreateSMBShareParams) (*truenas.SMBShare, error)
	UpdateSMBShare(ctx context.Context, id int, p *truenas.CreateSMBShareParams) (*truenas.SMBShare, error)
	DeleteSMBShare(ctx context.Context, id int) error
	ListWebDAVShares(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.WebDAVShare, error)
	GetWebDAVShare(ctx context.Context, id int) (*truenas.WebDAVShare, error)
	CreateWebDAVShare(ctx context.Context, p *truenas.CreateWebDAVShareParams) (*truenas.WebDAVShare, error)
	UpdateWebDAVShare(ctx context.Context, id int, p *truenas.CreateWebDAVShareParams) (*truenas.WebDAVShare, error)
	DeleteWebDAVShare(ctx context.Context, id int) error

	// NFS global service configuration
	NFSBindIPChoices(ctx context.Context) (map[string]string, error)
	NFSClientCount(ctx context.Context) (int, error)
	NFSConfigGet(ctx context.Context) (*truenas.NFSConfig, error)
	GetNFS3Clients(ctx context.Context) ([]truenas.NFSClient, error)
	GetNFS4Clients(ctx context.Context) ([]truenas.NFSClient, error)
	UpdateNFSConfig(ctx context.Context, p *truenas.UpdateNFSConfigParams) (*truenas.NFSConfig, error)

	// App management (registries, images, catalog metadata)
	ListAppRegistries(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.AppRegistry, error)
	GetAppRegistry(ctx context.Context, id int) (*truenas.AppRegistry, error)
	CreateAppRegistry(ctx context.Context, p *truenas.CreateAppRegistryParams) (*truenas.AppRegistry, error)
	UpdateAppRegistry(ctx context.Context, id int, p *truenas.CreateAppRegistryParams) (*truenas.AppRegistry, error)
	DeleteAppRegistry(ctx context.Context, id int) error
	PullAppImage(ctx context.Context, p *truenas.PullAppImageParams) (int, error)
	AppImageDockerHubRateLimitGet(ctx context.Context) (*truenas.DockerHubRateLimit, error)
	GetAppImage(ctx context.Context, id string) (*truenas.Image, error)
	AppCategories(ctx context.Context) ([]string, error)
	AppAvailableSpaceGet(ctx context.Context) (*truenas.AppAvailableSpace, error)
	AppConfigGet(ctx context.Context) (*truenas.AppGlobalConfig, error)
	AppContainerIDs(ctx context.Context, appName string) ([]string, error)
	ConvertAppToCustom(ctx context.Context, appName string) (*truenas.App, error)
	AppOutdatedDockerImages(ctx context.Context, appName string) ([]string, error)
	PullAppImages(ctx context.Context, appName string) (int, error)

	// Authentication
	GenerateOnetimePassword(ctx context.Context) (string, error)
	GenerateToken(ctx context.Context, ttl int) (string, error)
	Logout(ctx context.Context) (bool, error)
	Me(ctx context.Context) (*truenas.UserIdentity, error)
	AuthMechanismChoices(ctx context.Context) ([]string, error)
	AuthSessionsList(ctx context.Context) ([]truenas.AuthSession, error)
	SetAuthAttribute(ctx context.Context, key string, value any) (bool, error)
	TerminateOtherSessions(ctx context.Context) (bool, error)
	TerminateSession(ctx context.Context, id string) (bool, error)
	TwoFactorEnabled(ctx context.Context) (bool, error)
	TwoFactorConfigGet(ctx context.Context) (*truenas.TwoFactorConfig, error)
	UpdateTwoFactor(ctx context.Context, p *truenas.UpdateTwoFactorParams) (*truenas.TwoFactorConfig, error)

	// Alert management (categories, policies, services, classes)
	AlertListCategories(ctx context.Context) ([]truenas.AlertCategory, error)
	AlertListPolicies(ctx context.Context) ([]string, error)
	DismissAlert(ctx context.Context, uuid string) error
	RestoreAlert(ctx context.Context, uuid string) error
	ListAlertServices(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.AlertService, error)
	GetAlertService(ctx context.Context, id int) (*truenas.AlertService, error)
	CreateAlertService(ctx context.Context, p *truenas.CreateAlertServiceParams) (*truenas.AlertService, error)
	UpdateAlertService(ctx context.Context, id int, p *truenas.CreateAlertServiceParams) (*truenas.AlertService, error)
	DeleteAlertService(ctx context.Context, id int) error
	TestAlertService(ctx context.Context, p *truenas.CreateAlertServiceParams) (bool, error)
	AlertClassesConfigGet(ctx context.Context) (*truenas.AlertClassesConfig, error)
	UpdateAlertClasses(ctx context.Context, p *truenas.UpdateAlertClassesParams) (*truenas.AlertClassesConfig, error)

	// Cloud Backup (restic)
	ListCloudBackups(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CloudBackupTask, error)
	GetCloudBackup(ctx context.Context, id int) (*truenas.CloudBackupTask, error)
	CreateCloudBackup(ctx context.Context, p *truenas.CreateCloudBackupParams) (*truenas.CloudBackupTask, error)
	UpdateCloudBackup(ctx context.Context, id int, p *truenas.CreateCloudBackupParams) (*truenas.CloudBackupTask, error)
	DeleteCloudBackup(ctx context.Context, id int) error
	AbortCloudBackup(ctx context.Context, id int) error
	SyncCloudBackup(ctx context.Context, id int) (int, error)
	RestoreCloudBackup(ctx context.Context, id int, p *truenas.CloudBackupRestoreParams) (int, error)
	DeleteCloudBackupSnapshot(ctx context.Context, id int, snapshotID string) error
	ListCloudBackupSnapshots(ctx context.Context, id int) ([]truenas.CloudBackupSnapshot, error)
	ListCloudBackupSnapshotDirectory(ctx context.Context, id int, snapshotID, path string) ([]truenas.DirEntry, error)
	CloudBackupTransferSettingChoices(ctx context.Context) ([]string, error)

	// Certificate management
	ListCertificates(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.Certificate, error)
	GetCertificate(ctx context.Context, id int) (*truenas.Certificate, error)
	CreateCertificate(ctx context.Context, p *truenas.CreateCertificateParams) (*truenas.Certificate, error)
	UpdateCertificate(ctx context.Context, id int, p *truenas.CreateCertificateParams) (*truenas.Certificate, error)
	DeleteCertificate(ctx context.Context, id int) error
	CertificateACMEServerChoices(ctx context.Context) (map[string]string, error)
	CertificateCountryChoices(ctx context.Context) (map[string]string, error)
	CertificateECCurveChoices(ctx context.Context) (map[string]string, error)
	CertificateExtendedKeyUsageChoices(ctx context.Context) (map[string]string, error)

	// Core system (connectivity, introspection, job utilities)
	Ping(ctx context.Context) (string, error)
	PingRemote(ctx context.Context, params map[string]any) (string, error)
	GetMethods(ctx context.Context, app string) (map[string]any, error)
	GetServices(ctx context.Context) ([]map[string]any, error)
	JobDownloadLogs(ctx context.Context, jobID int) (string, error)
	JobWait(ctx context.Context, jobID int) (*truenas.Job, error)
	ARP(ctx context.Context, iface string) (map[string]string, error)
	Bulk(ctx context.Context, method string, paramsList [][]any) ([]any, error)
	Download(ctx context.Context, method string, args []any, filename string) (*truenas.CoreDownloadResult, error)
	ResizeShell(ctx context.Context, id string, cols, rows int) error
	Subscribe(ctx context.Context, event string) (string, error)
	Unsubscribe(ctx context.Context, subscriptionID string) error

	// Directory services
	DirectoryServicesCacheRefresh(ctx context.Context) error
	DirectoryServicesCertificateChoices(ctx context.Context) (map[string]string, error)
	DirectoryServicesConfigGet(ctx context.Context) (*truenas.DirectoryServicesConfig, error)
	DirectoryServicesLeave(ctx context.Context, params map[string]any) error
	DirectoryServicesStatusGet(ctx context.Context) (*truenas.DirectoryServicesStatus, error)
	DirectoryServicesSyncKeytab(ctx context.Context) error
	UpdateDirectoryServices(ctx context.Context, p *truenas.UpdateDirectoryServicesParams) (*truenas.DirectoryServicesConfig, error)

	// Disk management
	QueryDisks(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.Disk, error)
	DiskDetails(ctx context.Context) (map[string]any, error)
	DiskGetUsed(ctx context.Context, name string) (int64, error)
	DiskTemperatureAggGet(ctx context.Context, names []string, days int) (map[string]truenas.DiskTemperatureAgg, error)
	DiskTemperatureAlerts(ctx context.Context, names []string) ([]string, error)
	DiskTemperatures(ctx context.Context, names []string) (map[string]int, error)
	UpdateDisk(ctx context.Context, identifier string, p *truenas.UpdateDiskParams) (*truenas.Disk, error)
	WipeDisk(ctx context.Context, identifier string, p *truenas.WipeDiskParams) (int, error)

	// Replication tasks and endpoints
	ListReplicationTasks(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ReplicationTask, error)
	GetReplicationTask(ctx context.Context, id int) (*truenas.ReplicationTask, error)
	CreateReplicationTask(ctx context.Context, p *truenas.CreateReplicationParams) (*truenas.ReplicationTask, error)
	UpdateReplicationTask(ctx context.Context, id int, p *truenas.CreateReplicationParams) (*truenas.ReplicationTask, error)
	DeleteReplicationTask(ctx context.Context, id int) error
	TestReplicationTask(ctx context.Context, id int) (map[string]any, error)
	ReplicationSchemas(ctx context.Context) (map[string]any, error)
	ReplicationGlobalConfigGet(ctx context.Context) (*truenas.ReplicationGlobalConfig, error)
	UpdateReplicationGlobal(ctx context.Context, p *truenas.UpdateReplicationGlobalParams) (*truenas.ReplicationGlobalConfig, error)
	ListReplicationEndpoints(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ReplicationEndpoint, error)
	GetReplicationEndpoint(ctx context.Context, id int) (*truenas.ReplicationEndpoint, error)
	CreateReplicationEndpoint(ctx context.Context, p *truenas.CreateReplicationEndpointParams) (*truenas.ReplicationEndpoint, error)
	UpdateReplicationEndpoint(ctx context.Context, id int, p *truenas.CreateReplicationEndpointParams) (*truenas.ReplicationEndpoint, error)
	DeleteReplicationEndpoint(ctx context.Context, id int) error

	// User management
	ListUsers(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.User, error)
	GetUser(ctx context.Context, id int) (*truenas.User, error)
	CreateUser(ctx context.Context, p *truenas.CreateUserParams) (*truenas.User, error)
	UpdateUser(ctx context.Context, id int, p *truenas.CreateUserParams) (*truenas.User, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUserPassword(ctx context.Context, id int, password string) error
	UserSchemas(ctx context.Context) (map[string]any, error)
	UserHomeDirectoryChoices(ctx context.Context) (map[string]string, error)
	UserShellChoices(ctx context.Context) (map[string]string, error)

	// Cron jobs
	ListCronJobs(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CronJob, error)
	GetCronJob(ctx context.Context, id int) (*truenas.CronJob, error)
	CreateCronJob(ctx context.Context, p *truenas.CreateCronJobParams) (*truenas.CronJob, error)
	UpdateCronJob(ctx context.Context, id int, p *truenas.CreateCronJobParams) (*truenas.CronJob, error)
	DeleteCronJob(ctx context.Context, id int) error
	RunCronJob(ctx context.Context, id int) (int, error)

	// Device info and DNS
	DeviceGetInfo(ctx context.Context, deviceType string) (map[string]any, error)
	DNSQuery(ctx context.Context) (*truenas.DNSConfig, error)
}
