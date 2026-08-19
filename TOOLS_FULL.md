# Полный список инструментов TrueNAS MCP сервера

## Реализованные инструменты (32+ MCP Tools)

### System (1 инструмент)
| MCP Tool | Описание |
|----------|----------|
| `get_system_info` | Получение информации о системе: версия, hostname, CPU, память, uptime, load averages |

### Network (3 инструмента)
| MCP Tool | Описание |
|----------|----------|
| `list_interfaces` | Список всех сетевых интерфейсов (bridges, physical ports, bonds, VLANs) |
| `get_interface` | Детальная информация об интерфейсе по ID |
| `update_interface` | Обновление description сетевого интерфейса |

### Filesystem (6 инструментов)
| MCP Tool | Описание |
|----------|----------|
| `list_directory` | Перечисление содержимого директории на TrueNAS |
| `make_directory` | Создание новой директории по абсолютному пути |
| `read_file` | Чтение файла как UTF-8 текста (не подходит для бинарных файлов) |
| `write_file` | Запись содержимого в файл; append=true для добавления вместо перезаписи |
| `stat_file` | Метаданные файла: существование, тип (FILE/DIRECTORY), размер, uid/gid |
| `stat_filesystem` | Информация о файловой системе: total/free/available bytes и inode count |

### Pool/ZFS Storage (4 инструмента)
| MCP Tool | Описание |
|----------|----------|
| `list_pools` | Список всех ZFS пулов со статусом, размером, health; поддержка pagination |
| `get_pool` | Детальная информация о пуле по numeric ID |
| `create_pool` | Создание нового ZFS пула (name, disks array, layout: STRIPE/RAIDZ1/MIRROR/etc.) → job object |
| `update_pool` | Обновление name существующего пула по ID |

### Dataset/ZFS (4 инструмента)
| MCP Tool | Описание |
|----------|----------|
| `list_datasets` | Список ZFS datasets и zvols, опционально filtered by pool name |
| `get_dataset` | Детальная информация о dataset по full path ID (e.g., "Storage/backups") |
| `create_dataset` | Создание нового ZFS filesystem или zvol; FILESYSTEM по умолчанию, VOLUME требует volsize |
| `update_dataset` | Обновление comments или quota существующего dataset |

### Virtual Machines (10 инструментов)
| MCP Tool | Описание |
|----------|----------|
| `list_vms` | Список всех configured VMs со state (RUNNING/STOPPED), CPU count, memory in MiB |
| `get_vm` | Детальная информация о VM по numeric ID |
| `start_vm` | Запуск VM → возвращает `{job_id}` для async tracking |
| `stop_vm` | Остановка VM; force=true для forceful termination; возвращает `{job_id}` |
| `restart_vm` | Перезапуск VM → возвращает `{job_id}` |
| `create_vm` | Создание новой VM (name и memory обязательны) |
| `update_vm` | Обновление конфигурации существующей VM по ID |
| `list_vm_devices` | Список всех hardware devices, attached to VM (disks, CDROMs, NICs, displays) |
| `add_vm_device` | Добавление device к VM (types: DISK, CDROM, NIC, DISPLAY, RAW) |
| `delete_vm` | Удаление VM по ID (только при TRUENAS_ALLOW_DESTRUCTIVE=true, требует confirmed=true) |

### Applications/Docker Images (12 инструментов)
| MCP Tool | Описание |
|----------|----------|
| `list_apps` | Список всех installed apps на TrueNAS SCALE |
| `get_app` | Детальная информация об app по name |
| `start_app` | Запуск app → возвращает `{job_id}` |
| `stop_app` | Остановка app → возвращает `{job_id}` |
| `restart_app` | Перезапуск app → возвращает `{job_id}` |
| `list_images` | Список Docker images на TrueNAS SCALE |
| `install_app` | Установка app из catalog → возвращает `{job_id}` |
| `install_custom_app` | Установка custom Docker Compose app → возвращает `{job_id}` |
| `upgrade_app` | Обновление app до указанной версии или latest → возвращает `{job_id}` |
| `upgrade_summary` | Информация об available upgrade и changelog |
| `rollback_app` | Откат app к предыдущей версии → возвращает `{job_id}` |
| `delete_app` | Удаление app по name (только при TRUENAS_ALLOW_DESTRUCTIVE=true, требует confirmed=true) |

### Cloud Sync (9 инструментов)
| MCP Tool | Описание |
|----------|----------|
| `list_cloud_providers` | Список supported Cloud Sync providers (S3, backup/cloud services, etc.) |
| `list_cloud_credentials` | Список всех stored Cloud Sync credentials; поддержка pagination |
| `create_cloud_credential` | Создание нового Cloud Sync credential |
| `delete_cloud_credential` | Удаление credential по ID (только при TRUENAS_ALLOW_DESTRUCTIVE=true, требует confirmed=true) |
| `list_cloudsync_tasks` | Список всех configured Cloud Sync Tasks; поддержка pagination |
| `create_cloudsync_task` | Создание нового Cloud Sync Task → job object |
| `run_cloudsync_task` | Запуск задачи синхронизации (dry_run опция) → возвращает `{job_id}` |
| `abort_cloudsync_task` | Остановка задачи синхронизации → возвращает `{job_id}` |
| `delete_cloudsync_task` | Удаление задачи по ID (только при TRUENAS_ALLOW_DESTRUCTIVE=true, требует confirmed=true) |

### Jobs (2 инструмента)
| MCP Tool | Описание |
|----------|----------|
| `get_job` | Получение текущего state job по ID (для проверки long-running operations) |
| `abort_job` | Остановка job по ID (для jobs, не связанных с облачными операциями) |

### Alerts (1 инструмент)
| MCP Tool | Описание |
|----------|----------|
| `list_alerts` | Список текущих system alerts (active и dismissed), e.g., pool capacity warnings, app updates, service failures |

### Snapshots (5 инструментов)
| MCP Tool | Описание |
|----------|----------|
| `list_snapshots` | Список ZFS snapshots; опционально filtered by dataset path |
| `get_snapshot` | Детальная информация о snapshot по full ID (e.g., "Storage/backups@before-upgrade") |
| `create_snapshot` | Создание нового ZFS snapshot по dataset path и snapshot name |
| `rollback_snapshot` | Rollback dataset к previous ZFS snapshot (только при TRUENAS_ALLOW_DESTRUCTIVE=true, требует confirmed=true) |
| `delete_snapshot` | Удаление snapshot по ID (только при TRUENAS_ALLOW_DESTRUCTIVE=true, требует confirmed=true) |

### Destructive Tools (9 инструментов, опциональны)
Эти инструменты доступны только при флага `TRUENAS_ALLOW_DESTRUCTIVE=true`:

| MCP Tool | Описание |
|----------|----------|
| `delete_vm` | Permanently delete VM by ID (VM must be stopped first) |
| `delete_app` | Permanently delete app by name (app must be stopped first) |
| `delete_snapshot` | Permanently delete ZFS snapshot |
| `delete_dataset` | Permanently delete ZFS dataset or zvol |
| `delete_pool` | Permanently delete ZFS pool |
| `delete_vm_device` | Remove hardware device from VM (changes take effect on next boot) |
| `rollback_snapshot` | Roll dataset back to previous ZFS snapshot (ALL data after snapshot destroyed) |
| `delete_cloud_credential` | Permanently delete Cloud Sync credential (fails if any task references it) |
| `delete_cloudsync_task` | Permanently delete Cloud Sync Task |

---

## Архитектура проекта

### Структура
```
/home/hatetm/truenas-mcp/
├── tools/
│   ├── register.go          # Регистрация всех инструментов
│   ├── alert.go             # Alert инструменты
│   ├── app.go               # App инструменты
│   ├── cloudsync.go         # Cloud Sync инструменты
│   ├── dataset.go           # Dataset инструменты
│   ├── filesystem.go        # Filesystem инструменты
│   ├── jobs.go              # Jobs инструменты
│   ├── network.go           # Network инструменты
│   ├── pool.go              # Pool инструменты
│   ├── snapshot.go          # Snapshot инструменты
│   ├── system.go            # System инструменты
│   ├── vm.go                # VM инструменты
│   └── destructive.go       # Destructive инструменты (opt-in)
├── internal/truenas/
│   ├── client.go            # WebSocket JSON-RPC 2.0 клиент
│   ├── alert.go             # Alert API методы
│   ├── app.go               # App API методы
│   ├── cloudsync.go         # Cloud Sync API методы
│   ├── dataset.go           # Dataset API методы
│   ├── filesystem.go        # Filesystem API методы
│   ├── jobs.go              # Jobs API методы
│   ├── network.go           # Network API методы
│   ├── pool.go              # Pool API методы
│   ├── query.go             # Query методы
│   ├── snapshot.go          # Snapshot API методы
│   ├── system.go            # System API методы
│   ├── types.go             # Общие типы
│   ├── vm.go                # VM API методы
│   ├── alert_test.go        # Alert тесты
│   ├── app_test.go          # App тесты
│   ├── client_test.go       # Client тесты
│   ├── cloudsync_test.go    # Cloud Sync тесты
│   ├── dataset_test.go      # Dataset тесты
│   ├── filesystem_test.go   # Filesystem тесты
│   ├── jobs_test.go         # Jobs тесты
│   ├── network_test.go      # Network тесты
│   ├── pool_test.go         # Pool тесты
│   ├── query_test.go        # Query тесты
│   ├── snapshot_test.go     # Snapshot тесты
│   ├── system_test.go       # System тесты
│   └── vm_test.go           # VM тесты
└── README.md                # Документация
```

### Технологический стек
- **Language**: Go
- **Protocol**: WebSocket JSON-RPC 2.0
- **API**: TrueNAS SCALE API v2.0
- **Client**: Custom WebSocket client с reconnect logic
- **Connection**: wss://<host>/api/current

### Безопасность
- **Authentication**: API Key через параметры WebSocket
- **Destructive Operations**: Доступны только при `TRUENAS_ALLOW_DESTRUCTIVE=true`
- **TLS**: Опционально отключение certificate verification (insecure=true)

---

## API Methods TrueNAS SCALE v25.10

### Реализованные методы (сопоставление с TrueNAS API)

| TrueNAS API Namespace | MCP Tool(s) | Status |
|----------------------|-------------|--------|
| system.get_system_info | get_system_info | ✅ |
| network.list_interfaces | list_interfaces | ✅ |
| network.get_interface | get_interface | ✅ |
| network.update_interface | update_interface | ✅ |
| filesystem.list_directory | list_directory | ✅ |
| filesystem.mkdir | make_directory | ✅ |
| filesystem.read_file | read_file | ✅ |
| filesystem.write_file | write_file | ✅ |
| filesystem.stat_file | stat_file | ✅ |
| filesystem.stat_filesystem | stat_filesystem | ✅ |
| pool.list_pools | list_pools | ✅ |
| pool.get_pool | get_pool | ✅ |
| pool.create_pool | create_pool | ✅ |
| pool.update_pool | update_pool | ✅ |
| dataset.list_datasets | list_datasets | ✅ |
| dataset.get_dataset | get_dataset | ✅ |
| dataset.create_dataset | create_dataset | ✅ |
| dataset.update_dataset | update_dataset | ✅ |
| vm.list_vms | list_vms | ✅ |
| vm.get_vm | get_vm | ✅ |
| vm.start | start_vm | ✅ |
| vm.stop | stop_vm | ✅ |
| vm.restart | restart_vm | ✅ |
| vm.create | create_vm | ✅ |
| vm.update | update_vm | ✅ |
| vm.device.list | list_vm_devices | ✅ |
| vm.device.add | add_vm_device | ✅ |
| vm.device.delete | delete_vm | ✅ |
| app.list_apps | list_apps | ✅ |
| app.get_app | get_app | ✅ |
| app.start | start_app | ✅ |
| app.stop | stop_app | ✅ |
| app.restart | restart_app | ✅ |
| app.list_images | list_images | ✅ |
| app.install | install_app | ✅ |
| app.install_custom | install_custom_app | ✅ |
| app.upgrade | upgrade_app | ✅ |
| app.upgrade_summary | upgrade_summary | ✅ |
| app.rollback | rollback_app | ✅ |
| app.delete | delete_app | ✅ |
| cloudsync.list_providers | list_cloud_providers | ✅ |
| cloudsync.list_credentials | list_cloud_credentials | ✅ |
| cloudsync.create_credential | create_cloud_credential | ✅ |
| cloudsync.delete_credential | delete_cloud_credential | ✅ |
| cloudsync.list_tasks | list_cloudsync_tasks | ✅ |
| cloudsync.create_task | create_cloudsync_task | ✅ |
| cloudsync.run_task | run_cloudsync_task | ✅ |
| cloudsync.abort_task | abort_cloudsync_task | ✅ |
| cloudsync.delete_task | delete_cloudsync_task | ✅ |
| jobs.get_job | get_job | ✅ |
| jobs.abort_job | abort_job | ✅ |
| alert.list_alerts | list_alerts | ✅ |
| snapshot.list_snapshots | list_snapshots | ✅ |
| snapshot.get_snapshot | get_snapshot | ✅ |
| snapshot.create_snapshot | create_snapshot | ✅ |
| snapshot.rollback | rollback_snapshot | ✅ |
| snapshot.delete | delete_snapshot | ✅ |

---

## Usage Examples

### Создание клиента
```go
client, err := truenas.NewClient("https://truenas.local", "your-api-key", false)
if err != nil {
    log.Fatal(err)
}

// Подключение к WebSocket
if err := client.Connect(ctx); err != nil {
    log.Fatal(err)
}
```

### Получение информации о системе
```go
info, err := client.GetSystemInfo(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Version: %s\n", info.Version)
```

### Создание VM
```go
vm, err := client.CreateVM(ctx, truenas.VMConfig{
    Name: "my-vm",
    Memory: 4096,
    Description: "Test VM",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("VM created: %s\n", vm.Name)
```

### Запуск приложения
```go
jobID, err := client.StartApp(ctx, "jellyfin")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("App starting, job_id: %d\n", jobID)
```

---

## Notes

- **All async operations** return `{job_id}` for tracking
- **Destructive tools** require `TRUENAS_ALLOW_DESTRUCTIVE=true` environment variable
- **WebSocket connection** must be established before using any API methods
- **Context cancellation** properly propagates to WebSocket operations
- **Reconnection logic** handles automatic reconnection on connection failure
