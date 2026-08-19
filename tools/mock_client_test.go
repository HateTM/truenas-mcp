package tools

import (
	"context"
	"reflect"
	"sync"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
)

type mockTruenasClient struct {
	mu          sync.RWMutex
	DispatchMap map[string]any
}

// call dispatches to the *DispatchMap* function registered for method, if any.
// The registered function's signature mirrors the corresponding truenasClient
// method exactly (context.Context, then the method's other arguments in
// order, returning (value, error) or just error). If no function is
// registered, call returns nil with out left at its zero value — this lets
// tests register only the methods they care about.
func (m *mockTruenasClient) call(ctx context.Context, method string, params, out any) error {
	m.mu.RLock()
	fn, ok := m.DispatchMap[method]
	m.mu.RUnlock()
	if !ok {
		return nil
	}

	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()

	var extra []any
	if params != nil {
		if list, ok := params.([]any); ok {
			extra = list
		} else {
			extra = []any{params}
		}
	}

	args := make([]reflect.Value, 0, len(extra)+1)
	args = append(args, reflect.ValueOf(ctx))
	for i, a := range extra {
		paramIdx := i + 1
		var pt reflect.Type
		switch {
		case fnType.IsVariadic() && paramIdx >= fnType.NumIn()-1:
			pt = fnType.In(fnType.NumIn() - 1)
		case paramIdx < fnType.NumIn():
			pt = fnType.In(paramIdx)
		}
		if a == nil && pt != nil {
			args = append(args, reflect.Zero(pt))
			continue
		}
		args = append(args, reflect.ValueOf(a))
	}

	var results []reflect.Value
	if fnType.IsVariadic() {
		results = fnVal.CallSlice(args)
	} else {
		results = fnVal.Call(args)
	}

	errVal := results[len(results)-1]
	if !errVal.IsNil() {
		return errVal.Interface().(error)
	}

	if out == nil || len(results) < 2 {
		return nil
	}
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Pointer || outVal.IsNil() {
		return nil
	}
	resVal := results[0]
	elemType := outVal.Elem().Type()
	switch {
	case resVal.Type() == elemType:
		outVal.Elem().Set(resVal)
	case resVal.Kind() == reflect.Pointer && resVal.Type().Elem() == elemType:
		if !resVal.IsNil() {
			outVal.Elem().Set(resVal.Elem())
		}
	case resVal.Type().AssignableTo(elemType):
		outVal.Elem().Set(resVal)
	}
	return nil
}

func (m *mockTruenasClient) GetSystemInfo(ctx context.Context) (*truenas.SystemInfo, error) {
	var res truenas.SystemInfo
	if err := m.call(ctx, "get_system_info", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ListPools(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.Pool, error) {
	var res []truenas.Pool
	if err := m.call(ctx, "list_pools", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetPool(ctx context.Context, id int) (*truenas.Pool, error) {
	var res truenas.Pool
	if err := m.call(ctx, "get_pool", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreatePool(ctx context.Context, params *truenas.CreatePoolParams) (*truenas.Pool, error) {
	var res truenas.Pool
	if err := m.call(ctx, "create_pool", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdatePool(ctx context.Context, id int, params *truenas.UpdatePoolParams) (*truenas.Pool, error) {
	var res truenas.Pool
	if err := m.call(ctx, "update_pool", []any{id, params}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeletePool(ctx context.Context, id int) error {
	return m.call(ctx, "delete_pool", id, nil)
}

func (m *mockTruenasClient) AttachPool(ctx context.Context, id int) (*truenas.Pool, error) {
	var res truenas.Pool
	if err := m.call(ctx, "attach_pool", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DetachPool(ctx context.Context, id int) error {
	return m.call(ctx, "detach_pool", id, nil)
}

func (m *mockTruenasClient) ExpandPool(ctx context.Context, params *truenas.ExpandPoolParams) (*truenas.Pool, error) {
	var res truenas.Pool
	if err := m.call(ctx, "expand_pool", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ExportPool(ctx context.Context, id int, path string) (int64, error) {
	var res int64
	if err := m.call(ctx, "export_pool", []any{id, path}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListDisks(ctx context.Context, id int) ([]truenas.DiskInfo, error) {
	var res []truenas.DiskInfo
	if err := m.call(ctx, "list_disks", id, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ExportPoolHTML(ctx context.Context, id int) (string, error) {
	var res string
	if err := m.call(ctx, "export_pool_html", id, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) FindDatasetsForImport(ctx context.Context, id int, opts truenas.ListOptions) ([]truenas.Dataset, error) {
	var res []truenas.Dataset
	if err := m.call(ctx, "find_datasets_for_import", []any{id, opts}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ImportPool(ctx context.Context, id int, path, name string) error {
	return m.call(ctx, "import_pool", []any{id, path, name}, nil)
}

func (m *mockTruenasClient) IsPoolUpgraded(ctx context.Context, id int) (bool, error) {
	var res bool
	if err := m.call(ctx, "is_pool_upgraded", id, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) OfflinePool(ctx context.Context, id int) error {
	return m.call(ctx, "offline_pool", id, nil)
}

func (m *mockTruenasClient) OnlinePool(ctx context.Context, id int) error {
	return m.call(ctx, "online_pool", id, nil)
}

func (m *mockTruenasClient) RelabelPool(ctx context.Context, id int, key string, force bool) error {
	return m.call(ctx, "relabel_pool", []any{id, key, force}, nil)
}

func (m *mockTruenasClient) ScanPool(ctx context.Context, id int, force bool) error {
	return m.call(ctx, "scan_pool", []any{id, force}, nil)
}

func (m *mockTruenasClient) SplitPool(ctx context.Context, id int, path, name string) error {
	return m.call(ctx, "split_pool", []any{id, path, name}, nil)
}

func (m *mockTruenasClient) GetPoolStatus(ctx context.Context, id int) (*truenas.PoolStatus, error) {
	var res truenas.PoolStatus
	if err := m.call(ctx, "get_pool_status", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) SyncPool(ctx context.Context, id int) error {
	return m.call(ctx, "sync_pool", id, nil)
}

func (m *mockTruenasClient) VerifyPool(ctx context.Context, id int) error {
	return m.call(ctx, "verify_pool", id, nil)
}

func (m *mockTruenasClient) ListDatasets(ctx context.Context, pool string, opts ...truenas.ListOptions) ([]truenas.Dataset, error) {
	var res []truenas.Dataset
	if err := m.call(ctx, "list_datasets", []any{pool, opts}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetDataset(ctx context.Context, id string) (*truenas.Dataset, error) {
	var res truenas.Dataset
	if err := m.call(ctx, "get_dataset", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateDataset(ctx context.Context, params *truenas.CreateDatasetParams) (*truenas.Dataset, error) {
	var res truenas.Dataset
	if err := m.call(ctx, "create_dataset", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateDataset(ctx context.Context, id string, params *truenas.UpdateDatasetParams) (*truenas.Dataset, error) {
	var res truenas.Dataset
	if err := m.call(ctx, "update_dataset", []any{id, params}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteDataset(ctx context.Context, id string) error {
	return m.call(ctx, "delete_dataset", id, nil)
}

func (m *mockTruenasClient) DeleteDatasetByPath(ctx context.Context, path string, recursive bool) error {
	return m.call(ctx, "delete_dataset_by_path", []any{path, recursive}, nil)
}

func (m *mockTruenasClient) ExportDatasetKey(ctx context.Context, poolID int, datasetName, path string) (string, error) {
	var res string
	if err := m.call(ctx, "export_dataset_key", []any{poolID, datasetName, path}, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) ExportDatasetKeys(ctx context.Context, poolID int, path string) (int64, error) {
	var res int64
	if err := m.call(ctx, "export_dataset_keys", []any{poolID, path}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) ExportDatasetKeysForReplication(ctx context.Context, poolID int, remote, path string) (int64, error) {
	var res int64
	if err := m.call(ctx, "export_dataset_keys_for_replication", []any{poolID, remote, path}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) LockDatasetKey(ctx context.Context, poolID int, datasetName string) error {
	return m.call(ctx, "lock_dataset_key", []any{poolID, datasetName}, nil)
}

func (m *mockTruenasClient) UnlockDatasetKey(ctx context.Context, poolID int, datasetName, key string) error {
	return m.call(ctx, "unlock_dataset_key", []any{poolID, datasetName, key}, nil)
}

func (m *mockTruenasClient) PromoteDataset(ctx context.Context, poolID int, datasetName string) error {
	return m.call(ctx, "promote_dataset", []any{poolID, datasetName}, nil)
}

func (m *mockTruenasClient) DDTPrefetch(ctx context.Context, poolID int, enable bool) error {
	return m.call(ctx, "ddt_prefetch", []any{poolID, enable}, nil)
}

func (m *mockTruenasClient) GetDatasetInstance(ctx context.Context, datasetName string) (*truenas.Dataset, error) {
	var res truenas.Dataset
	if err := m.call(ctx, "get_dataset_instance", datasetName, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) GetDatasetQuota(ctx context.Context, datasetName string) (*truenas.QuotaInfo, error) {
	var res truenas.QuotaInfo
	if err := m.call(ctx, "get_dataset_quota", datasetName, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) SetDatasetQuota(ctx context.Context, datasetName string, quota int64) error {
	return m.call(ctx, "set_dataset_quota", []any{datasetName, quota}, nil)
}

func (m *mockTruenasClient) RenameDataset(ctx context.Context, poolID int, oldName, newName string) (*truenas.Dataset, error) {
	var res truenas.Dataset
	if err := m.call(ctx, "rename_dataset", []any{poolID, oldName, newName}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ImportDataset(ctx context.Context, poolID int, datasetName, path string) error {
	return m.call(ctx, "import_dataset", []any{poolID, datasetName, path}, nil)
}

func (m *mockTruenasClient) CheckQuota(ctx context.Context, poolID int, datasetName string) (*truenas.QuotaStatus, error) {
	var res truenas.QuotaStatus
	if err := m.call(ctx, "check_quota", []any{poolID, datasetName}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) SyncDataset(ctx context.Context, poolID int, datasetName string) error {
	return m.call(ctx, "sync_dataset", []any{poolID, datasetName}, nil)
}

func (m *mockTruenasClient) OfflineDataset(ctx context.Context, poolID int, datasetName string) error {
	return m.call(ctx, "offline_dataset", []any{poolID, datasetName}, nil)
}

func (m *mockTruenasClient) OnlineDataset(ctx context.Context, poolID int, datasetName string) error {
	return m.call(ctx, "online_dataset", []any{poolID, datasetName}, nil)
}

func (m *mockTruenasClient) GetDiskInfo(ctx context.Context, poolID int) (*truenas.DiskInfo, error) {
	var res truenas.DiskInfo
	if err := m.call(ctx, "get_disk_info", poolID, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) OfflineDisk(ctx context.Context, poolID int, diskID int) error {
	return m.call(ctx, "offline_disk", []any{poolID, diskID}, nil)
}

func (m *mockTruenasClient) OnlineDisk(ctx context.Context, poolID int, diskID int) error {
	return m.call(ctx, "online_disk", []any{poolID, diskID}, nil)
}

func (m *mockTruenasClient) RemoveDisk(ctx context.Context, poolID int, diskID int) error {
	return m.call(ctx, "remove_disk", []any{poolID, diskID}, nil)
}

func (m *mockTruenasClient) ReplaceDisk(ctx context.Context, poolID int, diskID int, path string) error {
	return m.call(ctx, "replace_disk", []any{poolID, diskID, path}, nil)
}

func (m *mockTruenasClient) ImportPoolHTML(ctx context.Context, path string) error {
	return m.call(ctx, "import_pool_html", path, nil)
}

func (m *mockTruenasClient) ListSnapshots(ctx context.Context, dataset string, opts ...truenas.ListOptions) ([]truenas.Snapshot, error) {
	var res []truenas.Snapshot
	if err := m.call(ctx, "list_snapshots", []any{dataset, opts}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetSnapshot(ctx context.Context, id string) (*truenas.Snapshot, error) {
	var res truenas.Snapshot
	if err := m.call(ctx, "get_snapshot", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateSnapshot(ctx context.Context, params truenas.CreateSnapshotParams) (*truenas.Snapshot, error) {
	var res truenas.Snapshot
	if err := m.call(ctx, "create_snapshot", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) RollbackSnapshot(ctx context.Context, id string, params truenas.RollbackSnapshotParams) error {
	return m.call(ctx, "rollback_snapshot", []any{id, params}, nil)
}

func (m *mockTruenasClient) DeleteSnapshot(ctx context.Context, id string) error {
	return m.call(ctx, "delete_snapshot", id, nil)
}

func (m *mockTruenasClient) ListVMs(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.VM, error) {
	var res []truenas.VM
	if err := m.call(ctx, "list_vms", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetVM(ctx context.Context, id int) (*truenas.VM, error) {
	var res truenas.VM
	if err := m.call(ctx, "get_vm", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) StartVM(ctx context.Context, id int) (int, error) {
	var res int
	if err := m.call(ctx, "start_vm", id, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) StopVM(ctx context.Context, id int, force bool) (int, error) {
	var res int
	if err := m.call(ctx, "stop_vm", []any{id, force}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) RestartVM(ctx context.Context, id int) (int, error) {
	var res int
	if err := m.call(ctx, "restart_vm", id, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) CreateVM(ctx context.Context, params *truenas.CreateVMParams) (*truenas.VM, error) {
	var res truenas.VM
	if err := m.call(ctx, "create_vm", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateVM(ctx context.Context, id int, params *truenas.UpdateVMParams) (*truenas.VM, error) {
	var res truenas.VM
	if err := m.call(ctx, "update_vm", []any{id, params}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteVM(ctx context.Context, id int) error {
	return m.call(ctx, "delete_vm", id, nil)
}

func (m *mockTruenasClient) ListVMDevices(ctx context.Context, vmID int) ([]truenas.VMDevice, error) {
	var res []truenas.VMDevice
	if err := m.call(ctx, "list_vm_devices", vmID, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) AddVMDevice(ctx context.Context, params *truenas.AddVMDeviceParams) (*truenas.VMDevice, error) {
	var res truenas.VMDevice
	if err := m.call(ctx, "add_vm_device", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteVMDevice(ctx context.Context, deviceID int) error {
	return m.call(ctx, "delete_vm_device", deviceID, nil)
}

func (m *mockTruenasClient) ListApps(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.App, error) {
	var res []truenas.App
	if err := m.call(ctx, "list_apps", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetApp(ctx context.Context, name string) (*truenas.App, error) {
	var res truenas.App
	if err := m.call(ctx, "get_app", name, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) StartApp(ctx context.Context, name string) (int, error) {
	var res int
	if err := m.call(ctx, "start_app", name, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) StopApp(ctx context.Context, name string) (int, error) {
	var res int
	if err := m.call(ctx, "stop_app", name, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) RestartApp(ctx context.Context, name string) (int, error) {
	var res int
	if err := m.call(ctx, "restart_app", name, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListImages(ctx context.Context) ([]truenas.Image, error) {
	var res []truenas.Image
	if err := m.call(ctx, "list_images", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) CreateApp(ctx context.Context, params *truenas.CreateAppParams) (int, error) {
	var res int
	if err := m.call(ctx, "create_app", params, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) DeleteApp(ctx context.Context, name string) error {
	return m.call(ctx, "delete_app", name, nil)
}

func (m *mockTruenasClient) UpgradeApp(ctx context.Context, name, version string) (int, error) {
	var res int
	if err := m.call(ctx, "upgrade_app", []any{name, version}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetUpgradeSummary(ctx context.Context, name string) (*truenas.AppUpgradeSummary, error) {
	var res truenas.AppUpgradeSummary
	if err := m.call(ctx, "get_upgrade_summary", name, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) RollbackApp(ctx context.Context, name, version string) (int, error) {
	var res int
	if err := m.call(ctx, "rollback_app", []any{name, version}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListInterfaces(ctx context.Context) ([]truenas.Interface, error) {
	var res []truenas.Interface
	if err := m.call(ctx, "list_interfaces", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetInterface(ctx context.Context, id string) (*truenas.Interface, error) {
	var res truenas.Interface
	if err := m.call(ctx, "get_interface", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateInterface(ctx context.Context, id string, params *truenas.UpdateInterfaceParams) (*truenas.Interface, error) {
	var res truenas.Interface
	if err := m.call(ctx, "update_interface", []any{id, params}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ListDirectory(ctx context.Context, path string) ([]truenas.DirEntry, error) {
	var res []truenas.DirEntry
	if err := m.call(ctx, "list_directory", path, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) Mkdir(ctx context.Context, path, mode string) error {
	return m.call(ctx, "mkdir", []any{path, mode}, nil)
}

func (m *mockTruenasClient) WriteFile(ctx context.Context, path string, content []byte, appendFlag bool) error {
	return m.call(ctx, "write_file", []any{path, content, appendFlag}, nil)
}

func (m *mockTruenasClient) ReadFile(ctx context.Context, path string) ([]byte, error) {
	var res []byte
	if err := m.call(ctx, "read_file", path, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) Stat(ctx context.Context, path string) (*truenas.FileStat, error) {
	var res truenas.FileStat
	if err := m.call(ctx, "stat", path, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) StatFS(ctx context.Context, path string) (*truenas.FSStat, error) {
	var res truenas.FSStat
	if err := m.call(ctx, "stat_fs", path, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ListCloudProviders(ctx context.Context) ([]truenas.CloudProvider, error) {
	var res []truenas.CloudProvider
	if err := m.call(ctx, "list_cloud_providers", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListCloudCredentials(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CloudCredential, error) {
	var res []truenas.CloudCredential
	if err := m.call(ctx, "list_cloud_credentials", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) CreateCloudCredential(ctx context.Context, params *truenas.CreateCloudCredentialParams) (*truenas.CloudCredential, error) {
	var res truenas.CloudCredential
	if err := m.call(ctx, "create_cloud_credential", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteCloudCredential(ctx context.Context, id int) error {
	return m.call(ctx, "delete_cloud_credential", id, nil)
}

func (m *mockTruenasClient) ListCloudSyncTasks(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CloudSyncTask, error) {
	var res []truenas.CloudSyncTask
	if err := m.call(ctx, "list_cloudsync_tasks", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) CreateCloudSyncTask(ctx context.Context, params *truenas.CreateCloudSyncTaskParams) (*truenas.CloudSyncTask, error) {
	var res truenas.CloudSyncTask
	if err := m.call(ctx, "create_cloudsync_task", params, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) RunCloudSyncTask(ctx context.Context, id int, dryRun bool) (int, error) {
	var res int
	if err := m.call(ctx, "run_cloudsync_task", []any{id, dryRun}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) AbortCloudSyncTask(ctx context.Context, id int) error {
	return m.call(ctx, "abort_cloudsync_task", id, nil)
}

func (m *mockTruenasClient) DeleteCloudSyncTask(ctx context.Context, id int) error {
	return m.call(ctx, "delete_cloudsync_task", id, nil)
}

func (m *mockTruenasClient) GetJob(ctx context.Context, id int) (*truenas.Job, error) {
	var res truenas.Job
	if err := m.call(ctx, "get_job", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) AbortJob(ctx context.Context, id int) error {
	return m.call(ctx, "abort_job", id, nil)
}

func (m *mockTruenasClient) ListAlerts(ctx context.Context) ([]truenas.Alert, error) {
	var res []truenas.Alert
	if err := m.call(ctx, "list_alerts", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) QueryDataset(ctx context.Context, datasetName string) (*truenas.Dataset, error) {
	var res truenas.Dataset
	if err := m.call(ctx, "get_dataset", datasetName, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DatasetChecksumChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "dataset_checksum_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DatasetCompressionChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "dataset_compression_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DatasetEncryptionAlgorithmChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "dataset_encryption_algorithm_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DatasetRecordsizeChoices(ctx context.Context) ([]string, error) {
	var res []string
	if err := m.call(ctx, "dataset_recordsize_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DDTPrune(ctx context.Context, poolID int, p truenas.DDTPruneParams) error {
	return m.call(ctx, "ddt_prune", []any{poolID, p}, nil)
}

func (m *mockTruenasClient) PoolFilesystemChoices(ctx context.Context, poolIDs []int) ([]string, error) {
	var res []string
	if err := m.call(ctx, "pool_filesystem_choices", poolIDs, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// --- iSCSI ---

func (m *mockTruenasClient) ListISCSIAuth(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIAuth, error) {
	var res []truenas.ISCSIAuth
	if err := m.call(ctx, "iscsi_auth_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetISCSIAuth(ctx context.Context, id int) (*truenas.ISCSIAuth, error) {
	var res truenas.ISCSIAuth
	if err := m.call(ctx, "iscsi_auth_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateISCSIAuth(ctx context.Context, p *truenas.CreateISCSIAuthParams) (*truenas.ISCSIAuth, error) {
	var res truenas.ISCSIAuth
	if err := m.call(ctx, "iscsi_auth_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateISCSIAuth(ctx context.Context, id int, p *truenas.CreateISCSIAuthParams) (*truenas.ISCSIAuth, error) {
	var res truenas.ISCSIAuth
	if err := m.call(ctx, "iscsi_auth_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteISCSIAuth(ctx context.Context, id int) error {
	return m.call(ctx, "iscsi_auth_delete", id, nil)
}

func (m *mockTruenasClient) ISCSIGlobalALUAEnabled(ctx context.Context) (bool, error) {
	var res bool
	if err := m.call(ctx, "iscsi_global_alua_enabled", nil, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) ISCSIGlobalClientCount(ctx context.Context) (int, error) {
	var res int
	if err := m.call(ctx, "iscsi_global_client_count", nil, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) ISCSIGlobalConfigGet(ctx context.Context) (*truenas.ISCSIGlobalConfig, error) {
	var res truenas.ISCSIGlobalConfig
	if err := m.call(ctx, "iscsi_global_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ISCSIGlobalISEREnabled(ctx context.Context) (bool, error) {
	var res bool
	if err := m.call(ctx, "iscsi_global_iser_enabled", nil, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) ISCSIGlobalSessions(ctx context.Context) ([]truenas.ISCSISession, error) {
	var res []truenas.ISCSISession
	if err := m.call(ctx, "iscsi_global_sessions", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) UpdateISCSIGlobal(ctx context.Context, p *truenas.UpdateISCSIGlobalParams) (*truenas.ISCSIGlobalConfig, error) {
	var res truenas.ISCSIGlobalConfig
	if err := m.call(ctx, "iscsi_global_update", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ListISCSIExtents(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIExtent, error) {
	var res []truenas.ISCSIExtent
	if err := m.call(ctx, "iscsi_extent_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetISCSIExtent(ctx context.Context, id int) (*truenas.ISCSIExtent, error) {
	var res truenas.ISCSIExtent
	if err := m.call(ctx, "iscsi_extent_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateISCSIExtent(ctx context.Context, p *truenas.CreateISCSIExtentParams) (*truenas.ISCSIExtent, error) {
	var res truenas.ISCSIExtent
	if err := m.call(ctx, "iscsi_extent_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateISCSIExtent(ctx context.Context, id int, p *truenas.CreateISCSIExtentParams) (*truenas.ISCSIExtent, error) {
	var res truenas.ISCSIExtent
	if err := m.call(ctx, "iscsi_extent_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteISCSIExtent(ctx context.Context, id int) error {
	return m.call(ctx, "iscsi_extent_delete", id, nil)
}

func (m *mockTruenasClient) ISCSIExtentDiskChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "iscsi_extent_disk_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListISCSIInitiators(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIInitiator, error) {
	var res []truenas.ISCSIInitiator
	if err := m.call(ctx, "iscsi_initiator_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetISCSIInitiator(ctx context.Context, id int) (*truenas.ISCSIInitiator, error) {
	var res truenas.ISCSIInitiator
	if err := m.call(ctx, "iscsi_initiator_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateISCSIInitiator(ctx context.Context, p *truenas.CreateISCSIInitiatorParams) (*truenas.ISCSIInitiator, error) {
	var res truenas.ISCSIInitiator
	if err := m.call(ctx, "iscsi_initiator_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateISCSIInitiator(ctx context.Context, id int, p *truenas.CreateISCSIInitiatorParams) (*truenas.ISCSIInitiator, error) {
	var res truenas.ISCSIInitiator
	if err := m.call(ctx, "iscsi_initiator_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteISCSIInitiator(ctx context.Context, id int) error {
	return m.call(ctx, "iscsi_initiator_delete", id, nil)
}

func (m *mockTruenasClient) ListISCSIPortals(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSIPortal, error) {
	var res []truenas.ISCSIPortal
	if err := m.call(ctx, "iscsi_portal_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetISCSIPortal(ctx context.Context, id int) (*truenas.ISCSIPortal, error) {
	var res truenas.ISCSIPortal
	if err := m.call(ctx, "iscsi_portal_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateISCSIPortal(ctx context.Context, p *truenas.CreateISCSIPortalParams) (*truenas.ISCSIPortal, error) {
	var res truenas.ISCSIPortal
	if err := m.call(ctx, "iscsi_portal_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateISCSIPortal(ctx context.Context, id int, p *truenas.CreateISCSIPortalParams) (*truenas.ISCSIPortal, error) {
	var res truenas.ISCSIPortal
	if err := m.call(ctx, "iscsi_portal_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteISCSIPortal(ctx context.Context, id int) error {
	return m.call(ctx, "iscsi_portal_delete", id, nil)
}

func (m *mockTruenasClient) ISCSIPortalListenIPChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "iscsi_portal_listen_ip_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListISCSITargets(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSITarget, error) {
	var res []truenas.ISCSITarget
	if err := m.call(ctx, "iscsi_target_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetISCSITarget(ctx context.Context, id int) (*truenas.ISCSITarget, error) {
	var res truenas.ISCSITarget
	if err := m.call(ctx, "iscsi_target_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateISCSITarget(ctx context.Context, p *truenas.CreateISCSITargetParams) (*truenas.ISCSITarget, error) {
	var res truenas.ISCSITarget
	if err := m.call(ctx, "iscsi_target_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateISCSITarget(ctx context.Context, id int, p *truenas.CreateISCSITargetParams) (*truenas.ISCSITarget, error) {
	var res truenas.ISCSITarget
	if err := m.call(ctx, "iscsi_target_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteISCSITarget(ctx context.Context, id int) error {
	return m.call(ctx, "iscsi_target_delete", id, nil)
}

func (m *mockTruenasClient) ValidateISCSITargetName(ctx context.Context, name string) (bool, error) {
	var res bool
	if err := m.call(ctx, "iscsi_target_validate_name", name, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListISCSITargetExtents(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ISCSITargetExtent, error) {
	var res []truenas.ISCSITargetExtent
	if err := m.call(ctx, "iscsi_targetextent_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetISCSITargetExtent(ctx context.Context, id int) (*truenas.ISCSITargetExtent, error) {
	var res truenas.ISCSITargetExtent
	if err := m.call(ctx, "iscsi_targetextent_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateISCSITargetExtent(ctx context.Context, p *truenas.CreateISCSITargetExtentParams) (*truenas.ISCSITargetExtent, error) {
	var res truenas.ISCSITargetExtent
	if err := m.call(ctx, "iscsi_targetextent_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateISCSITargetExtent(ctx context.Context, id int, p *truenas.CreateISCSITargetExtentParams) (*truenas.ISCSITargetExtent, error) {
	var res truenas.ISCSITargetExtent
	if err := m.call(ctx, "iscsi_targetextent_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteISCSITargetExtent(ctx context.Context, id int) error {
	return m.call(ctx, "iscsi_targetextent_delete", id, nil)
}

// --- NVMe-oF ---

func (m *mockTruenasClient) NVMetGlobalConfigGet(ctx context.Context) (*truenas.NVMetGlobalConfig, error) {
	var res truenas.NVMetGlobalConfig
	if err := m.call(ctx, "nvmeof_global_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNVMetGlobal(ctx context.Context, p *truenas.UpdateNVMetGlobalParams) (*truenas.NVMetGlobalConfig, error) {
	var res truenas.NVMetGlobalConfig
	if err := m.call(ctx, "nvmeof_global_update", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ListNVMetHosts(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetHost, error) {
	var res []truenas.NVMetHost
	if err := m.call(ctx, "nvmeof_host_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNVMetHost(ctx context.Context, id int) (*truenas.NVMetHost, error) {
	var res truenas.NVMetHost
	if err := m.call(ctx, "nvmeof_host_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateNVMetHost(ctx context.Context, p *truenas.CreateNVMetHostParams) (*truenas.NVMetHost, error) {
	var res truenas.NVMetHost
	if err := m.call(ctx, "nvmeof_host_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNVMetHost(ctx context.Context, id int, p *truenas.CreateNVMetHostParams) (*truenas.NVMetHost, error) {
	var res truenas.NVMetHost
	if err := m.call(ctx, "nvmeof_host_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteNVMetHost(ctx context.Context, id int) error {
	return m.call(ctx, "nvmeof_host_delete", id, nil)
}

func (m *mockTruenasClient) NVMetHostDHCHAPDHGroupChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "nvmeof_host_dhchap_dhgroup_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) NVMetHostDHCHAPHashChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "nvmeof_host_dhchap_hash_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) NVMetHostGenerateKey(ctx context.Context) (string, error) {
	var res string
	if err := m.call(ctx, "nvmeof_host_generate_key", nil, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) ListNVMetHostSubsys(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetHostSubsys, error) {
	var res []truenas.NVMetHostSubsys
	if err := m.call(ctx, "nvmeof_host_subsys_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNVMetHostSubsys(ctx context.Context, id int) (*truenas.NVMetHostSubsys, error) {
	var res truenas.NVMetHostSubsys
	if err := m.call(ctx, "nvmeof_host_subsys_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateNVMetHostSubsys(ctx context.Context, p *truenas.CreateNVMetHostSubsysParams) (*truenas.NVMetHostSubsys, error) {
	var res truenas.NVMetHostSubsys
	if err := m.call(ctx, "nvmeof_host_subsys_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNVMetHostSubsys(ctx context.Context, id int, p *truenas.CreateNVMetHostSubsysParams) (*truenas.NVMetHostSubsys, error) {
	var res truenas.NVMetHostSubsys
	if err := m.call(ctx, "nvmeof_host_subsys_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteNVMetHostSubsys(ctx context.Context, id int) error {
	return m.call(ctx, "nvmeof_host_subsys_delete", id, nil)
}

func (m *mockTruenasClient) ListNVMetSubsys(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetSubsys, error) {
	var res []truenas.NVMetSubsys
	if err := m.call(ctx, "nvmeof_subsys_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNVMetSubsys(ctx context.Context, id int) (*truenas.NVMetSubsys, error) {
	var res truenas.NVMetSubsys
	if err := m.call(ctx, "nvmeof_subsys_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateNVMetSubsys(ctx context.Context, p *truenas.CreateNVMetSubsysParams) (*truenas.NVMetSubsys, error) {
	var res truenas.NVMetSubsys
	if err := m.call(ctx, "nvmeof_subsys_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNVMetSubsys(ctx context.Context, id int, p *truenas.CreateNVMetSubsysParams) (*truenas.NVMetSubsys, error) {
	var res truenas.NVMetSubsys
	if err := m.call(ctx, "nvmeof_subsys_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteNVMetSubsys(ctx context.Context, id int) error {
	return m.call(ctx, "nvmeof_subsys_delete", id, nil)
}

func (m *mockTruenasClient) ListNVMetNamespaces(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetNamespace, error) {
	var res []truenas.NVMetNamespace
	if err := m.call(ctx, "nvmeof_namespace_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNVMetNamespace(ctx context.Context, id int) (*truenas.NVMetNamespace, error) {
	var res truenas.NVMetNamespace
	if err := m.call(ctx, "nvmeof_namespace_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateNVMetNamespace(ctx context.Context, p *truenas.CreateNVMetNamespaceParams) (*truenas.NVMetNamespace, error) {
	var res truenas.NVMetNamespace
	if err := m.call(ctx, "nvmeof_namespace_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNVMetNamespace(ctx context.Context, id int, p *truenas.CreateNVMetNamespaceParams) (*truenas.NVMetNamespace, error) {
	var res truenas.NVMetNamespace
	if err := m.call(ctx, "nvmeof_namespace_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteNVMetNamespace(ctx context.Context, id int) error {
	return m.call(ctx, "nvmeof_namespace_delete", id, nil)
}

func (m *mockTruenasClient) ListNVMetPortSubsys(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetPortSubsys, error) {
	var res []truenas.NVMetPortSubsys
	if err := m.call(ctx, "nvmeof_port_subsys_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNVMetPortSubsys(ctx context.Context, id int) (*truenas.NVMetPortSubsys, error) {
	var res truenas.NVMetPortSubsys
	if err := m.call(ctx, "nvmeof_port_subsys_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateNVMetPortSubsys(ctx context.Context, p *truenas.CreateNVMetPortSubsysParams) (*truenas.NVMetPortSubsys, error) {
	var res truenas.NVMetPortSubsys
	if err := m.call(ctx, "nvmeof_port_subsys_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNVMetPortSubsys(ctx context.Context, id int, p *truenas.CreateNVMetPortSubsysParams) (*truenas.NVMetPortSubsys, error) {
	var res truenas.NVMetPortSubsys
	if err := m.call(ctx, "nvmeof_port_subsys_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteNVMetPortSubsys(ctx context.Context, id int) error {
	return m.call(ctx, "nvmeof_port_subsys_delete", id, nil)
}

func (m *mockTruenasClient) ListNVMetPorts(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NVMetPort, error) {
	var res []truenas.NVMetPort
	if err := m.call(ctx, "nvmeof_port_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNVMetPort(ctx context.Context, id int) (*truenas.NVMetPort, error) {
	var res truenas.NVMetPort
	if err := m.call(ctx, "nvmeof_port_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateNVMetPort(ctx context.Context, p *truenas.CreateNVMetPortParams) (*truenas.NVMetPort, error) {
	var res truenas.NVMetPort
	if err := m.call(ctx, "nvmeof_port_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNVMetPort(ctx context.Context, id int, p *truenas.CreateNVMetPortParams) (*truenas.NVMetPort, error) {
	var res truenas.NVMetPort
	if err := m.call(ctx, "nvmeof_port_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteNVMetPort(ctx context.Context, id int) error {
	return m.call(ctx, "nvmeof_port_delete", id, nil)
}

func (m *mockTruenasClient) NVMetPortTransportAddressChoices(ctx context.Context, trtype string) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "nvmeof_port_transport_address_choices", trtype, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// --- Sharing (NFS/SMB/WebDAV) ---

func (m *mockTruenasClient) ListNFSShares(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.NFSShare, error) {
	var res []truenas.NFSShare
	if err := m.call(ctx, "sharing_nfs_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNFSShare(ctx context.Context, id int) (*truenas.NFSShare, error) {
	var res truenas.NFSShare
	if err := m.call(ctx, "sharing_nfs_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateNFSShare(ctx context.Context, p *truenas.CreateNFSShareParams) (*truenas.NFSShare, error) {
	var res truenas.NFSShare
	if err := m.call(ctx, "sharing_nfs_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateNFSShare(ctx context.Context, id int, p *truenas.CreateNFSShareParams) (*truenas.NFSShare, error) {
	var res truenas.NFSShare
	if err := m.call(ctx, "sharing_nfs_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteNFSShare(ctx context.Context, id int) error {
	return m.call(ctx, "sharing_nfs_delete", id, nil)
}

func (m *mockTruenasClient) ListSMBShares(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.SMBShare, error) {
	var res []truenas.SMBShare
	if err := m.call(ctx, "sharing_smb_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetSMBShare(ctx context.Context, id int) (*truenas.SMBShare, error) {
	var res truenas.SMBShare
	if err := m.call(ctx, "sharing_smb_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateSMBShare(ctx context.Context, p *truenas.CreateSMBShareParams) (*truenas.SMBShare, error) {
	var res truenas.SMBShare
	if err := m.call(ctx, "sharing_smb_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateSMBShare(ctx context.Context, id int, p *truenas.CreateSMBShareParams) (*truenas.SMBShare, error) {
	var res truenas.SMBShare
	if err := m.call(ctx, "sharing_smb_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteSMBShare(ctx context.Context, id int) error {
	return m.call(ctx, "sharing_smb_delete", id, nil)
}

func (m *mockTruenasClient) ListWebDAVShares(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.WebDAVShare, error) {
	var res []truenas.WebDAVShare
	if err := m.call(ctx, "sharing_webdav_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetWebDAVShare(ctx context.Context, id int) (*truenas.WebDAVShare, error) {
	var res truenas.WebDAVShare
	if err := m.call(ctx, "sharing_webdav_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateWebDAVShare(ctx context.Context, p *truenas.CreateWebDAVShareParams) (*truenas.WebDAVShare, error) {
	var res truenas.WebDAVShare
	if err := m.call(ctx, "sharing_webdav_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateWebDAVShare(ctx context.Context, id int, p *truenas.CreateWebDAVShareParams) (*truenas.WebDAVShare, error) {
	var res truenas.WebDAVShare
	if err := m.call(ctx, "sharing_webdav_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteWebDAVShare(ctx context.Context, id int) error {
	return m.call(ctx, "sharing_webdav_delete", id, nil)
}

// --- NFS global service configuration ---

func (m *mockTruenasClient) NFSBindIPChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "nfs_bindip_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) NFSClientCount(ctx context.Context) (int, error) {
	var res int
	if err := m.call(ctx, "nfs_client_count", nil, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) NFSConfigGet(ctx context.Context) (*truenas.NFSConfig, error) {
	var res truenas.NFSConfig
	if err := m.call(ctx, "nfs_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) GetNFS3Clients(ctx context.Context) ([]truenas.NFSClient, error) {
	var res []truenas.NFSClient
	if err := m.call(ctx, "nfs_get_nfs3_clients", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetNFS4Clients(ctx context.Context) ([]truenas.NFSClient, error) {
	var res []truenas.NFSClient
	if err := m.call(ctx, "nfs_get_nfs4_clients", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) UpdateNFSConfig(ctx context.Context, p *truenas.UpdateNFSConfigParams) (*truenas.NFSConfig, error) {
	var res truenas.NFSConfig
	if err := m.call(ctx, "nfs_update", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// --- App management (registries, images, catalog metadata) ---

func (m *mockTruenasClient) ListAppRegistries(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.AppRegistry, error) {
	var res []truenas.AppRegistry
	if err := m.call(ctx, "app_registry_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetAppRegistry(ctx context.Context, id int) (*truenas.AppRegistry, error) {
	var res truenas.AppRegistry
	if err := m.call(ctx, "app_registry_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateAppRegistry(ctx context.Context, p *truenas.CreateAppRegistryParams) (*truenas.AppRegistry, error) {
	var res truenas.AppRegistry
	if err := m.call(ctx, "app_registry_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateAppRegistry(ctx context.Context, id int, p *truenas.CreateAppRegistryParams) (*truenas.AppRegistry, error) {
	var res truenas.AppRegistry
	if err := m.call(ctx, "app_registry_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteAppRegistry(ctx context.Context, id int) error {
	return m.call(ctx, "app_registry_delete", id, nil)
}

func (m *mockTruenasClient) PullAppImage(ctx context.Context, p *truenas.PullAppImageParams) (int, error) {
	var res int
	if err := m.call(ctx, "app_image_pull", p, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) AppImageDockerHubRateLimitGet(ctx context.Context) (*truenas.DockerHubRateLimit, error) {
	var res truenas.DockerHubRateLimit
	if err := m.call(ctx, "app_image_dockerhub_rate_limit", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) GetAppImage(ctx context.Context, id string) (*truenas.Image, error) {
	var res truenas.Image
	if err := m.call(ctx, "app_image_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) AppCategories(ctx context.Context) ([]string, error) {
	var res []string
	if err := m.call(ctx, "app_categories", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) AppAvailableSpaceGet(ctx context.Context) (*truenas.AppAvailableSpace, error) {
	var res truenas.AppAvailableSpace
	if err := m.call(ctx, "app_available_space", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) AppConfigGet(ctx context.Context) (*truenas.AppGlobalConfig, error) {
	var res truenas.AppGlobalConfig
	if err := m.call(ctx, "app_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) AppContainerIDs(ctx context.Context, appName string) ([]string, error) {
	var res []string
	if err := m.call(ctx, "app_container_ids", appName, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ConvertAppToCustom(ctx context.Context, appName string) (*truenas.App, error) {
	var res truenas.App
	if err := m.call(ctx, "app_convert_to_custom", appName, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) AppOutdatedDockerImages(ctx context.Context, appName string) ([]string, error) {
	var res []string
	if err := m.call(ctx, "app_outdated_docker_images", appName, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) PullAppImages(ctx context.Context, appName string) (int, error) {
	var res int
	if err := m.call(ctx, "app_pull_images", appName, &res); err != nil {
		return 0, err
	}
	return res, nil
}

// --- Authentication ---

func (m *mockTruenasClient) GenerateOnetimePassword(ctx context.Context) (string, error) {
	var res string
	if err := m.call(ctx, "auth_generate_onetime_password", nil, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) GenerateToken(ctx context.Context, ttl int) (string, error) {
	var res string
	if err := m.call(ctx, "auth_generate_token", ttl, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) Logout(ctx context.Context) (bool, error) {
	var res bool
	if err := m.call(ctx, "auth_logout", nil, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) Me(ctx context.Context) (*truenas.UserIdentity, error) {
	var res truenas.UserIdentity
	if err := m.call(ctx, "auth_me", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) AuthMechanismChoices(ctx context.Context) ([]string, error) {
	var res []string
	if err := m.call(ctx, "auth_mechanism_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) AuthSessionsList(ctx context.Context) ([]truenas.AuthSession, error) {
	var res []truenas.AuthSession
	if err := m.call(ctx, "auth_sessions", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) SetAuthAttribute(ctx context.Context, key string, value any) (bool, error) {
	var res bool
	if err := m.call(ctx, "auth_set_attribute", []any{key, value}, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) TerminateOtherSessions(ctx context.Context) (bool, error) {
	var res bool
	if err := m.call(ctx, "auth_terminate_other_sessions", nil, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) TerminateSession(ctx context.Context, id string) (bool, error) {
	var res bool
	if err := m.call(ctx, "auth_terminate_session", id, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) TwoFactorEnabled(ctx context.Context) (bool, error) {
	var res bool
	if err := m.call(ctx, "auth_twofactor", nil, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) TwoFactorConfigGet(ctx context.Context) (*truenas.TwoFactorConfig, error) {
	var res truenas.TwoFactorConfig
	if err := m.call(ctx, "auth_twofactor_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateTwoFactor(ctx context.Context, p *truenas.UpdateTwoFactorParams) (*truenas.TwoFactorConfig, error) {
	var res truenas.TwoFactorConfig
	if err := m.call(ctx, "auth_twofactor_update", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// --- Alert management (categories, policies, services, classes) ---

func (m *mockTruenasClient) AlertListCategories(ctx context.Context) ([]truenas.AlertCategory, error) {
	var res []truenas.AlertCategory
	if err := m.call(ctx, "alert_list_categories", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) AlertListPolicies(ctx context.Context) ([]string, error) {
	var res []string
	if err := m.call(ctx, "alert_list_policies", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DismissAlert(ctx context.Context, uuid string) error {
	return m.call(ctx, "alert_dismiss", uuid, nil)
}

func (m *mockTruenasClient) RestoreAlert(ctx context.Context, uuid string) error {
	return m.call(ctx, "alert_restore", uuid, nil)
}

func (m *mockTruenasClient) ListAlertServices(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.AlertService, error) {
	var res []truenas.AlertService
	if err := m.call(ctx, "alertservice_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetAlertService(ctx context.Context, id int) (*truenas.AlertService, error) {
	var res truenas.AlertService
	if err := m.call(ctx, "alertservice_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateAlertService(ctx context.Context, p *truenas.CreateAlertServiceParams) (*truenas.AlertService, error) {
	var res truenas.AlertService
	if err := m.call(ctx, "alertservice_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateAlertService(ctx context.Context, id int, p *truenas.CreateAlertServiceParams) (*truenas.AlertService, error) {
	var res truenas.AlertService
	if err := m.call(ctx, "alertservice_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteAlertService(ctx context.Context, id int) error {
	return m.call(ctx, "alertservice_delete", id, nil)
}

func (m *mockTruenasClient) TestAlertService(ctx context.Context, p *truenas.CreateAlertServiceParams) (bool, error) {
	var res bool
	if err := m.call(ctx, "alertservice_test", p, &res); err != nil {
		return false, err
	}
	return res, nil
}

func (m *mockTruenasClient) AlertClassesConfigGet(ctx context.Context) (*truenas.AlertClassesConfig, error) {
	var res truenas.AlertClassesConfig
	if err := m.call(ctx, "alertclasses_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateAlertClasses(ctx context.Context, p *truenas.UpdateAlertClassesParams) (*truenas.AlertClassesConfig, error) {
	var res truenas.AlertClassesConfig
	if err := m.call(ctx, "alertclasses_update", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// --- Cloud Backup (restic) ---

func (m *mockTruenasClient) ListCloudBackups(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CloudBackupTask, error) {
	var res []truenas.CloudBackupTask
	if err := m.call(ctx, "cloud_backup_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetCloudBackup(ctx context.Context, id int) (*truenas.CloudBackupTask, error) {
	var res truenas.CloudBackupTask
	if err := m.call(ctx, "cloud_backup_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateCloudBackup(ctx context.Context, p *truenas.CreateCloudBackupParams) (*truenas.CloudBackupTask, error) {
	var res truenas.CloudBackupTask
	if err := m.call(ctx, "cloud_backup_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateCloudBackup(ctx context.Context, id int, p *truenas.CreateCloudBackupParams) (*truenas.CloudBackupTask, error) {
	var res truenas.CloudBackupTask
	if err := m.call(ctx, "cloud_backup_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteCloudBackup(ctx context.Context, id int) error {
	return m.call(ctx, "cloud_backup_delete", id, nil)
}

func (m *mockTruenasClient) AbortCloudBackup(ctx context.Context, id int) error {
	return m.call(ctx, "cloud_backup_abort", id, nil)
}

func (m *mockTruenasClient) SyncCloudBackup(ctx context.Context, id int) (int, error) {
	var res int
	if err := m.call(ctx, "cloud_backup_sync", id, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) RestoreCloudBackup(ctx context.Context, id int, p *truenas.CloudBackupRestoreParams) (int, error) {
	var res int
	if err := m.call(ctx, "cloud_backup_restore", []any{id, p}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) DeleteCloudBackupSnapshot(ctx context.Context, id int, snapshotID string) error {
	return m.call(ctx, "cloud_backup_delete_snapshot", []any{id, snapshotID}, nil)
}

func (m *mockTruenasClient) ListCloudBackupSnapshots(ctx context.Context, id int) ([]truenas.CloudBackupSnapshot, error) {
	var res []truenas.CloudBackupSnapshot
	if err := m.call(ctx, "cloud_backup_list_snapshots", id, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ListCloudBackupSnapshotDirectory(ctx context.Context, id int, snapshotID, path string) ([]truenas.DirEntry, error) {
	var res []truenas.DirEntry
	if err := m.call(ctx, "cloud_backup_list_snapshot_directory", []any{id, snapshotID, path}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) CloudBackupTransferSettingChoices(ctx context.Context) ([]string, error) {
	var res []string
	if err := m.call(ctx, "cloud_backup_transfer_setting_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// --- Certificate management ---

func (m *mockTruenasClient) ListCertificates(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.Certificate, error) {
	var res []truenas.Certificate
	if err := m.call(ctx, "certificate_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetCertificate(ctx context.Context, id int) (*truenas.Certificate, error) {
	var res truenas.Certificate
	if err := m.call(ctx, "certificate_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateCertificate(ctx context.Context, p *truenas.CreateCertificateParams) (*truenas.Certificate, error) {
	var res truenas.Certificate
	if err := m.call(ctx, "certificate_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateCertificate(ctx context.Context, id int, p *truenas.CreateCertificateParams) (*truenas.Certificate, error) {
	var res truenas.Certificate
	if err := m.call(ctx, "certificate_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteCertificate(ctx context.Context, id int) error {
	return m.call(ctx, "certificate_delete", id, nil)
}

func (m *mockTruenasClient) CertificateACMEServerChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "certificate_acme_server_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) CertificateCountryChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "certificate_country_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) CertificateECCurveChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "certificate_ec_curve_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) CertificateExtendedKeyUsageChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "certificate_extended_key_usage_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// --- Core system (connectivity, introspection, job utilities) ---

func (m *mockTruenasClient) Ping(ctx context.Context) (string, error) {
	var res string
	if err := m.call(ctx, "core_ping", nil, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) PingRemote(ctx context.Context, params map[string]any) (string, error) {
	var res string
	if err := m.call(ctx, "core_ping_remote", params, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) GetMethods(ctx context.Context, app string) (map[string]any, error) {
	var res map[string]any
	if err := m.call(ctx, "core_get_methods", app, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetServices(ctx context.Context) ([]map[string]any, error) {
	var res []map[string]any
	if err := m.call(ctx, "core_get_services", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) JobDownloadLogs(ctx context.Context, jobID int) (string, error) {
	var res string
	if err := m.call(ctx, "core_job_download_logs", jobID, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) JobWait(ctx context.Context, jobID int) (*truenas.Job, error) {
	var res truenas.Job
	if err := m.call(ctx, "core_job_wait", jobID, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ARP(ctx context.Context, iface string) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "core_arp", iface, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) Bulk(ctx context.Context, method string, paramsList [][]any) ([]any, error) {
	var res []any
	if err := m.call(ctx, "core_bulk", []any{method, paramsList}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) Download(ctx context.Context, method string, args []any, filename string) (*truenas.CoreDownloadResult, error) {
	var res truenas.CoreDownloadResult
	if err := m.call(ctx, "core_download", []any{method, args, filename}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ResizeShell(ctx context.Context, id string, cols, rows int) error {
	return m.call(ctx, "core_resize_shell", []any{id, cols, rows}, nil)
}

func (m *mockTruenasClient) Subscribe(ctx context.Context, event string) (string, error) {
	var res string
	if err := m.call(ctx, "core_subscribe", event, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (m *mockTruenasClient) Unsubscribe(ctx context.Context, subscriptionID string) error {
	return m.call(ctx, "core_unsubscribe", subscriptionID, nil)
}

// --- Directory services ---

func (m *mockTruenasClient) DirectoryServicesCacheRefresh(ctx context.Context) error {
	return m.call(ctx, "directoryservices_cache_refresh", nil, nil)
}

func (m *mockTruenasClient) DirectoryServicesCertificateChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "directoryservices_certificate_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DirectoryServicesConfigGet(ctx context.Context) (*truenas.DirectoryServicesConfig, error) {
	var res truenas.DirectoryServicesConfig
	if err := m.call(ctx, "directoryservices_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DirectoryServicesLeave(ctx context.Context, params map[string]any) error {
	return m.call(ctx, "directoryservices_leave", params, nil)
}

func (m *mockTruenasClient) DirectoryServicesStatusGet(ctx context.Context) (*truenas.DirectoryServicesStatus, error) {
	var res truenas.DirectoryServicesStatus
	if err := m.call(ctx, "directoryservices_status", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DirectoryServicesSyncKeytab(ctx context.Context) error {
	return m.call(ctx, "directoryservices_sync_keytab", nil, nil)
}

func (m *mockTruenasClient) UpdateDirectoryServices(ctx context.Context, p *truenas.UpdateDirectoryServicesParams) (*truenas.DirectoryServicesConfig, error) {
	var res truenas.DirectoryServicesConfig
	if err := m.call(ctx, "directoryservices_update", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// --- Disk management ---

func (m *mockTruenasClient) QueryDisks(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.Disk, error) {
	var res []truenas.Disk
	if err := m.call(ctx, "disk_query", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DiskDetails(ctx context.Context) (map[string]any, error) {
	var res map[string]any
	if err := m.call(ctx, "disk_details", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DiskGetUsed(ctx context.Context, name string) (int64, error) {
	var res int64
	if err := m.call(ctx, "disk_get_used", name, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func (m *mockTruenasClient) DiskTemperatureAggGet(ctx context.Context, names []string, days int) (map[string]truenas.DiskTemperatureAgg, error) {
	var res map[string]truenas.DiskTemperatureAgg
	if err := m.call(ctx, "disk_temperature_agg", []any{names, days}, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DiskTemperatureAlerts(ctx context.Context, names []string) ([]string, error) {
	var res []string
	if err := m.call(ctx, "disk_temperature_alerts", names, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DiskTemperatures(ctx context.Context, names []string) (map[string]int, error) {
	var res map[string]int
	if err := m.call(ctx, "disk_temperatures", names, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) UpdateDisk(ctx context.Context, identifier string, p *truenas.UpdateDiskParams) (*truenas.Disk, error) {
	var res truenas.Disk
	if err := m.call(ctx, "disk_update", []any{identifier, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) WipeDisk(ctx context.Context, identifier string, p *truenas.WipeDiskParams) (int, error) {
	var res int
	if err := m.call(ctx, "disk_wipe", []any{identifier, p}, &res); err != nil {
		return 0, err
	}
	return res, nil
}

// --- Replication tasks and endpoints ---

func (m *mockTruenasClient) ListReplicationTasks(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ReplicationTask, error) {
	var res []truenas.ReplicationTask
	if err := m.call(ctx, "replication_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetReplicationTask(ctx context.Context, id int) (*truenas.ReplicationTask, error) {
	var res truenas.ReplicationTask
	if err := m.call(ctx, "replication_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateReplicationTask(ctx context.Context, p *truenas.CreateReplicationParams) (*truenas.ReplicationTask, error) {
	var res truenas.ReplicationTask
	if err := m.call(ctx, "replication_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateReplicationTask(ctx context.Context, id int, p *truenas.CreateReplicationParams) (*truenas.ReplicationTask, error) {
	var res truenas.ReplicationTask
	if err := m.call(ctx, "replication_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteReplicationTask(ctx context.Context, id int) error {
	return m.call(ctx, "replication_delete", id, nil)
}

func (m *mockTruenasClient) TestReplicationTask(ctx context.Context, id int) (map[string]any, error) {
	var res map[string]any
	if err := m.call(ctx, "replication_test", id, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ReplicationSchemas(ctx context.Context) (map[string]any, error) {
	var res map[string]any
	if err := m.call(ctx, "replication_schemas", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) ReplicationGlobalConfigGet(ctx context.Context) (*truenas.ReplicationGlobalConfig, error) {
	var res truenas.ReplicationGlobalConfig
	if err := m.call(ctx, "replication_global_config", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateReplicationGlobal(ctx context.Context, p *truenas.UpdateReplicationGlobalParams) (*truenas.ReplicationGlobalConfig, error) {
	var res truenas.ReplicationGlobalConfig
	if err := m.call(ctx, "replication_global_update", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) ListReplicationEndpoints(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.ReplicationEndpoint, error) {
	var res []truenas.ReplicationEndpoint
	if err := m.call(ctx, "replication_endpoint_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetReplicationEndpoint(ctx context.Context, id int) (*truenas.ReplicationEndpoint, error) {
	var res truenas.ReplicationEndpoint
	if err := m.call(ctx, "replication_endpoint_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateReplicationEndpoint(ctx context.Context, p *truenas.CreateReplicationEndpointParams) (*truenas.ReplicationEndpoint, error) {
	var res truenas.ReplicationEndpoint
	if err := m.call(ctx, "replication_endpoint_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateReplicationEndpoint(ctx context.Context, id int, p *truenas.CreateReplicationEndpointParams) (*truenas.ReplicationEndpoint, error) {
	var res truenas.ReplicationEndpoint
	if err := m.call(ctx, "replication_endpoint_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteReplicationEndpoint(ctx context.Context, id int) error {
	return m.call(ctx, "replication_endpoint_delete", id, nil)
}

// --- User management ---

func (m *mockTruenasClient) ListUsers(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.User, error) {
	var res []truenas.User
	if err := m.call(ctx, "user_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetUser(ctx context.Context, id int) (*truenas.User, error) {
	var res truenas.User
	if err := m.call(ctx, "user_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateUser(ctx context.Context, p *truenas.CreateUserParams) (*truenas.User, error) {
	var res truenas.User
	if err := m.call(ctx, "user_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateUser(ctx context.Context, id int, p *truenas.CreateUserParams) (*truenas.User, error) {
	var res truenas.User
	if err := m.call(ctx, "user_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteUser(ctx context.Context, id int) error {
	return m.call(ctx, "user_delete", id, nil)
}

func (m *mockTruenasClient) UpdateUserPassword(ctx context.Context, id int, password string) error {
	return m.call(ctx, "user_update_password", []any{id, password}, nil)
}

func (m *mockTruenasClient) UserSchemas(ctx context.Context) (map[string]any, error) {
	var res map[string]any
	if err := m.call(ctx, "user_schemas", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) UserHomeDirectoryChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "user_home_directory_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) UserShellChoices(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	if err := m.call(ctx, "user_shell_choices", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// --- Cron jobs ---

func (m *mockTruenasClient) ListCronJobs(ctx context.Context, opts ...truenas.ListOptions) ([]truenas.CronJob, error) {
	var res []truenas.CronJob
	if err := m.call(ctx, "cronjob_list", opts, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) GetCronJob(ctx context.Context, id int) (*truenas.CronJob, error) {
	var res truenas.CronJob
	if err := m.call(ctx, "cronjob_get", id, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) CreateCronJob(ctx context.Context, p *truenas.CreateCronJobParams) (*truenas.CronJob, error) {
	var res truenas.CronJob
	if err := m.call(ctx, "cronjob_create", p, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) UpdateCronJob(ctx context.Context, id int, p *truenas.CreateCronJobParams) (*truenas.CronJob, error) {
	var res truenas.CronJob
	if err := m.call(ctx, "cronjob_update", []any{id, p}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *mockTruenasClient) DeleteCronJob(ctx context.Context, id int) error {
	return m.call(ctx, "cronjob_delete", id, nil)
}

func (m *mockTruenasClient) RunCronJob(ctx context.Context, id int) (int, error) {
	var res int
	if err := m.call(ctx, "cronjob_run", id, &res); err != nil {
		return 0, err
	}
	return res, nil
}

// --- Device info and DNS ---

func (m *mockTruenasClient) DeviceGetInfo(ctx context.Context, deviceType string) (map[string]any, error) {
	var res map[string]any
	if err := m.call(ctx, "device_get_info", deviceType, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mockTruenasClient) DNSQuery(ctx context.Context) (*truenas.DNSConfig, error) {
	var res truenas.DNSConfig
	if err := m.call(ctx, "dns_query", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
