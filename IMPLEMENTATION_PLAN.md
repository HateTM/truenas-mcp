# План реализации всех недостающих инструментов TrueNAS MCP

## Статус: 138 реализовано, 122 остаются (из 260 методов, отслеживаемых этим чек-листом)

**Приоритет 1 (CRITICAL) закрыт целиком** — Pool Management (остаток), iSCSI, NVMe-oF,
Sharing NFS/SMB/WebDAV (102 метода) реализованы: клиентский слой (`internal/truenas/`),
mock-интерфейс (`tools/client_iface.go`), MCP-инструменты (`tools/*.go`) и тесты для всех.
Деструктивные операции (delete/*) вынесены в `tools/destructive_iscsi.go`,
`tools/destructive_nvmeof.go`, `tools/destructive_sharing.go` — каждый файл в пределах
правила «300 строк», гейт `Config.AllowDestructive` сохранён. Осталось: Приоритеты 2-4
(App Management, Authentication, Alert Management, Cloud Backup, Certificate Management,
Core System, Directory Services, Disk Management, Replication, User Management, Cron Jobs,
Device Management, DNS).

Пересчитано по факту: для каждого метода ниже проверено, есть ли в `internal/truenas/*.go` клиентский метод, который
семантически выполняет именно эту TrueNAS RPC-операцию (по строке, передаваемой в `c.call(...)`, либо — для части
`pool.*`/`storage/pool`-кода — по REST-style пути, который используется как временная замена настоящего имени метода;
это существующая особенность клиентского слоя, не то, что нужно чинить в рамках этого документа).

В коде всего 146 вызовов `mcp.AddTool` в `tools/*.go`, но не все из них соответствуют пунктам исходного чек-листа —
часть строит дополнительную функциональность сверх него (например `pool_status`, `pool_scan`, `create_vm`,
`install_app`, `create_cloud_credential` и т.д.), которая в этом документе не учитывается, так как изначально не была
в списке. Прежний заголовок «32+ реализовано, 400+ остаются» был устаревшим ещё до подсчёта готовности: сумма самих
пунктов чек-листа (по всем 17 категориям) — 260, а не 400+; это поправлено ниже вместе со статусом.

### Обнаруженные попутно проблемы (не относятся к статусу методов, но важны для будущих тикетов)

- **`tools/pool_api.go` — мёртвый код.** Функция `RegisterTrueNASManagementTools` (38 вызовов `mcp.AddTool`,
  ~700 строк) нигде не вызывается — ни из `tools.RegisterAll`, ни из `cmd/truenas-mcp/main.go`, ни из тестов. Ни один
  из её инструментов не подключается к реальному серверу. Вся эта функциональность (`pool_create`, `pool_list`,
  `app_list`, `vm_list`, `cloudsync_list` и т.д. под другими именами) уже продублирована в `pool.go`, `dataset.go`,
  `snapshot.go`, `cloudsync.go`, `app.go`, `vm.go`, `network.go`, `filesystem.go`, `jobs.go`, `alert.go` — поэтому на
  подсчёт готовности ниже это не влияет, но файл стоит либо удалить, либо явно исключить из грепа при следующей
  сверке, чтобы он не вводил в заблуждение.
- **Двойная регистрация `RegisterPoolManagementTools`.** `tools.RegisterAll` (вызывается из `main.go`) уже вызывает
  `RegisterPoolManagementTools` изнутри `register.go`, но `cmd/truenas-mcp/main.go` затем вызывает её ещё раз
  напрямую. Это, вероятно, приводит к конфликту имён инструментов (`pool_attach`, `pool_expand` и т. д.
  регистрируются дважды) при старте сервера — стоит завести отдельный баг-тикет вне этого документа.
- **`tools/pool_management.go` уже 987 строк** при жёстком правиле «не более 300 строк на файл» — это существующее
  нарушение, а не прогноз на будущее. Рефакторинг существующих 48 инструментов в этом документе не планируется (это
  отдельная задача), но новые методы (см. Приоритет 1, п.1) размещаются в отдельном новом файле, а не дописываются
  в этот.

---

## Приоритет 1 (CRITICAL)

### 1. Pool Management (6 методов остаётся из 34)

Готово (28/34) в `tools/pool_management.go` (+ `tools/pool.go`, `tools/dataset.go`): `pool.attach()`,
`pool.detach()`, `pool.create()`, `pool.expand()`, `pool.export()`, `pool.get_disks()`, `pool.get_instance()`,
`pool.html()`, `pool.import_find()`, `pool.import_pool()`, `pool.is_upgraded()`, `pool.offline()`,
`pool.dataset.create()`, `pool.dataset.delete()`, `pool.dataset.details()`, `pool.dataset.update()`,
`pool.dataset.query()`, `pool.dataset.get_instance()`, `pool.dataset.get_quota()`, `pool.dataset.set_quota()`,
`pool.dataset.rename()`, `pool.dataset.export_key()`, `pool.dataset.export_keys()`,
`pool.dataset.export_keys_for_replication()`, `pool.dataset.lock()`, `pool.dataset.unlock()`,
`pool.dataset.promote()`, `pool.ddt_prefetch()`.

Осталось:
```go
// tools/pool_dataset_choices.go (новый файл — не дописывать в переполненный pool_management.go)
- pool.dataset.checksum_choices()
- pool.dataset.compression_choices()
- pool.dataset.encryption_algorithm_choices()
- pool.dataset.recordsize_choices()
- pool.ddt_prune()
- pool.filesystem_choices()
```

### 2. iSCSI (39 методов, не реализовано)
```go
// tools/iscsi_auth.go — iscsi.auth.* + iscsi.global.* (11 методов)
- iscsi.auth.create()
- iscsi.auth.delete()
- iscsi.auth.get_instance()
- iscsi.auth.query()
- iscsi.auth.update()
- iscsi.global.alua_enabled()
- iscsi.global.client_count()
- iscsi.global.config()
- iscsi.global.iser_enabled()
- iscsi.global.sessions()
- iscsi.global.update()
```
```go
// tools/iscsi_extent.go — iscsi.extent.* + iscsi.initiator.* (11 методов)
- iscsi.extent.create()
- iscsi.extent.delete()
- iscsi.extent.disk_choices()
- iscsi.extent.get_instance()
- iscsi.extent.query()
- iscsi.extent.update()
- iscsi.initiator.create()
- iscsi.initiator.delete()
- iscsi.initiator.get_instance()
- iscsi.initiator.query()
- iscsi.initiator.update()
```
```go
// tools/iscsi_portal.go — iscsi.portal.* (6 методов)
- iscsi.portal.create()
- iscsi.portal.delete()
- iscsi.portal.get_instance()
- iscsi.portal.listen_ip_choices()
- iscsi.portal.query()
- iscsi.portal.update()
```
```go
// tools/iscsi_target.go — iscsi.target.* + iscsi.targetextent.* (11 методов)
- iscsi.target.create()
- iscsi.target.delete()
- iscsi.target.get_instance()
- iscsi.target.query()
- iscsi.target.update()
- iscsi.target.validate_name()
- iscsi.targetextent.create()
- iscsi.targetextent.delete()
- iscsi.targetextent.get_instance()
- iscsi.targetextent.query()
- iscsi.targetextent.update()
```

### 3. NVMe-oF (36 методов, не реализовано)
```go
// tools/nvmeof_host.go — nvmet.global.* + nvmet.host.* (10 методов)
- nvmet.global.config()
- nvmet.global.update()
- nvmet.host.create()
- nvmet.host.delete()
- nvmet.host.dhchap_dhgroup_choices()
- nvmet.host.dhchap_hash_choices()
- nvmet.host.generate_key()
- nvmet.host.get_instance()
- nvmet.host.query()
- nvmet.host.update()
```
```go
// tools/nvmeof_subsys.go — nvmet.host_subsys.* + nvmet.subsys.* (10 методов)
- nvmet.host_subsys.create()
- nvmet.host_subsys.delete()
- nvmet.host_subsys.get_instance()
- nvmet.host_subsys.query()
- nvmet.host_subsys.update()
- nvmet.subsys.create()
- nvmet.subsys.delete()
- nvmet.subsys.get_instance()
- nvmet.subsys.query()
- nvmet.subsys.update()
```
```go
// tools/nvmeof_namespace.go — nvmet.namespace.* + nvmet.port_subsys.* (10 методов)
- nvmet.namespace.create()
- nvmet.namespace.delete()
- nvmet.namespace.get_instance()
- nvmet.namespace.query()
- nvmet.namespace.update()
- nvmet.port_subsys.create()
- nvmet.port_subsys.delete()
- nvmet.port_subsys.get_instance()
- nvmet.port_subsys.query()
- nvmet.port_subsys.update()
```
```go
// tools/nvmeof_port.go — nvmet.port.* (6 методов)
- nvmet.port.create()
- nvmet.port.delete()
- nvmet.port.get_instance()
- nvmet.port.query()
- nvmet.port.transport_address_choices()
- nvmet.port.update()
```

### 4. Sharing NFS/SMB/WebDAV (21 метод, не реализовано)
```go
// tools/sharing.go — sharing.nfs.* + sharing.smb.* + sharing.webdav.* (15 методов)
- sharing.nfs.create()
- sharing.nfs.delete()
- sharing.nfs.get_instance()
- sharing.nfs.query()
- sharing.nfs.update()
- sharing.smb.create()
- sharing.smb.delete()
- sharing.smb.get_instance()
- sharing.smb.query()
- sharing.smb.update()
- sharing.webdav.create()
- sharing.webdav.delete()
- sharing.webdav.get_instance()
- sharing.webdav.query()
- sharing.webdav.update()
```
```go
// tools/nfs_config.go — глобальные настройки NFS-сервиса, отдельный namespace от sharing.nfs.* (6 методов)
- nfs.bindip_choices()
- nfs.client_count()
- nfs.config()
- nfs.get_nfs3_clients()
- nfs.get_nfs4_clients()
- nfs.update()
```

---

## Приоритет 2 (HIGH)

### 5. App Management (15 методов остаётся из 18)

Готово (3/18) в `tools/app.go` (`list_apps` → `app.query()`, `list_images` → `app.image.query()`,
`restart_app` → `app.redeploy()`). Сверх чек-листа там же уже реализованы `app.get_instance()`, `app.start()`,
`app.stop()`, `app.create()`, `app.delete()`, `app.upgrade()`, `app.upgrade_summary()`, `app.rollback()` — эти
RPC-методы не входили в исходный список категории, поэтому не учтены в счётчике.

Осталось:
```go
// tools/app_management.go (15 методов)
- app.registry.create()
- app.registry.delete()
- app.registry.get_instance()
- app.registry.query()
- app.registry.update()
- app.image.pull()
- app.image.dockerhub_rate_limit()
- app.image.get_instance()
- app.categories()
- app.available_space()
- app.config()
- app.container_ids()
- app.convert_to_custom()
- app.outdated_docker_images()
- app.pull_images()
```

### 6. Authentication (16 методов остаётся из 17)

`auth.login_with_api_key()` фактически уже выполняется — `internal/truenas/client.go` вызывает его при каждом
подключении (`Connect()`), — но только для аутентификации самой MCP-сессии, а не как отдельный вызываемый
инструмент. Осознанно не выставлять произвольный `auth_login` MCP-инструментом (риск: даёт модели логиниться
чужим API-ключом) — этот пункт закрыт архитектурно, отдельная реализация не нужна.

```go
// tools/auth.go (16 методов)
- auth.generate_onetime_password()
- auth.generate_token()
- auth.login()
- auth.login_ex()
- auth.login_ex_continue()
- auth.login_with_token()
- auth.logout()
- auth.me()
- auth.mechanism_choices()
- auth.sessions()
- auth.set_attribute()
- auth.terminate_other_sessions()
- auth.terminate_session()
- auth.twofactor()
- auth.twofactor.config()
- auth.twofactor.update()
```

### 7. Alert Management (12 методов остаётся из 13)

Готово (1/13): `alert.list()` — `tools/alert.go` (`list_alerts`).

```go
// tools/alert_management.go (12 методов)
- alert.list_categories()
- alert.list_policies()
- alert.dismiss()
- alert.restore()
- alertservice.create()
- alertservice.delete()
- alertservice.get_instance()
- alertservice.query()
- alertservice.test()
- alertservice.update()
- alertclasses.config()
- alertclasses.update()
```

### 8. Cloud Backup (12 методов, не реализовано)

**Важно:** `tools/cloudsync.go` / `internal/truenas/cloudsync.go` реализуют домен `cloudsync.*` (rclone-based Cloud
Sync Tasks) — это отдельный API-namespace от `cloud_backup.*` (restic-based Cloud Backup) ниже. Они не пересекаются;
несмотря на то, что появление `cloudsync.go` могло выглядеть как закрытие этой категории, она по-прежнему полностью
не реализована.

```go
// tools/cloud_backup.go (12 методов)
- cloud_backup.abort()
- cloud_backup.create()
- cloud_backup.delete()
- cloud_backup.delete_snapshot()
- cloud_backup.get_instance()
- cloud_backup.list_snapshot_directory()
- cloud_backup.list_snapshots()
- cloud_backup.query()
- cloud_backup.restore()
- cloud_backup.sync()
- cloud_backup.transfer_setting_choices()
- cloud_backup.update()
```

---

## Приоритет 3 (MEDIUM)

### 9. Certificate Management (9 методов, не реализовано)
```go
// tools/certificate.go (9 методов)
- certificate.create()
- certificate.delete()
- certificate.get_instance()
- certificate.query()
- certificate.update()
- certificate.acme_server_choices()
- certificate.country_choices()
- certificate.ec_curve_choices()
- certificate.extended_key_usage_choices()
```

### 10. Core System (12 методов остаётся из 15)

Готово (3/15) в `tools/jobs.go` (`get_job()` → `core.get_jobs()`, `abort_job()` → `core.job_abort()`) и
`internal/truenas/client.go` (`core.set_options()` — вызывается один раз при подключении для включения
`legacy_jobs: false`; не выставлен отдельным инструментом по тем же причинам, что и `auth.login_with_api_key()`
выше — это конфигурация сессии, а не операция, которую должен вызывать пользователь/модель).

```go
// tools/core.go (12 методов)
- core.ping()
- core.ping_remote()
- core.get_methods()
- core.get_services()
- core.job_download_logs()
- core.job_wait()
- core.arp()
- core.bulk()
- core.download()
- core.resize_shell()
- core.subscribe()
- core.unsubscribe()
```

Примечание: в `internal/truenas/jobs.go` есть неиспользуемый `PollJob()` (клиентский long-poll поверх
`core.get_jobs()`), но он нигде не вызывается ни одним MCP-инструментом — ближайший кандидат на реализацию
`core.job_wait()`, но сейчас это мёртвый код, а не рабочий инструмент.

### 11. Directory Services (7 методов, не реализовано)
```go
// tools/directory_services.go (7 методов)
- directoryservices.cache_refresh()
- directoryservices.certificate_choices()
- directoryservices.config()
- directoryservices.leave()
- directoryservices.status()
- directoryservices.sync_keytab()
- directoryservices.update()
```

### 12. Disk Management (8 методов, не реализовано)
```go
// tools/disk.go (8 методов)
- disk.details()
- disk.get_used()
- disk.query()
- disk.temperature_agg()
- disk.temperature_alerts()
- disk.temperatures()
- disk.update()
- disk.wipe()
```

### 13. Replication (14 методов, не реализовано)

Не путать с `pool.dataset.export_keys_for_replication()` (Приоритет 1, п.1, уже реализован) — это отдельный
вспомогательный метод экспорта ключей шифрования, а не сам домен `replication.*`/`replication.endpoint.*`
(создание и управление задачами репликации), который целиком отсутствует.

```go
// tools/replication.go (14 методов)
- replication.create()
- replication.delete()
- replication.get_instance()
- replication.query()
- replication.test()
- replication.update()
- replication.replication_schemas()
- replication.global.config()
- replication.global.update()
- replication.endpoint.create()
- replication.endpoint.delete()
- replication.endpoint.get_instance()
- replication.endpoint.query()
- replication.endpoint.update()
```

---

## Приоритет 4 (LOW)

### 14. User Management (9 методов, не реализовано)
```go
// tools/user.go (9 методов)
- user.create()
- user.delete()
- user.get_instance()
- user.query()
- user.update()
- user.user_schemas()
- user.home_directory_choices()
- user.shell_choices()
- user.update_password()
```

### 15. Cron Jobs (6 методов, не реализовано)
```go
// tools/cronjob.go (6 методов)
- cronjob.create()
- cronjob.delete()
- cronjob.get_instance()
- cronjob.query()
- cronjob.run()
- cronjob.update()
```

### 16. Device Management (1 метод, не реализовано)
```go
// tools/misc.go — вместе с DNS (п.17), т.к. оба однометодные (1 метод)
- device.get_info()
```

### 17. DNS (1 метод, не реализовано)
```go
// tools/misc.go — вместе с Device Management (п.16) (1 метод)
- dns.query()
```

---

## Итого: 224 метода остаются (из 260 отслеживаемых этим чек-листом)

### Структура файлов:
```
tools/
├── pool_management.go       # готово (28 методов); 987 строк — уже нарушает лимит 300 строк/файл,
│                             # рефакторинг существующего кода вне рамок этого документа
├── pool_dataset_choices.go  # NEW: 6 методов (Приоритет 1, п.1)
├── iscsi_auth.go            # NEW: 11 методов (Приоритет 1, п.2)
├── iscsi_extent.go          # NEW: 11 методов (Приоритет 1, п.2)
├── iscsi_portal.go          # NEW: 6 методов (Приоритет 1, п.2)
├── iscsi_target.go          # NEW: 11 методов (Приоритет 1, п.2)
├── nvmeof_host.go           # NEW: 10 методов (Приоритет 1, п.3)
├── nvmeof_subsys.go         # NEW: 10 методов (Приоритет 1, п.3)
├── nvmeof_namespace.go      # NEW: 10 методов (Приоритет 1, п.3)
├── nvmeof_port.go           # NEW: 6 методов (Приоритет 1, п.3)
├── sharing.go               # NEW: 15 методов (Приоритет 1, п.4)
├── nfs_config.go            # NEW: 6 методов (Приоритет 1, п.4)
├── app.go                   # готово (3 метода: list_apps/list_images/restart_app + доп. функциональность)
├── app_management.go        # NEW: 15 методов (Приоритет 2, п.5)
├── auth.go                  # NEW: 16 методов (Приоритет 2, п.6)
├── alert.go                 # готово (1 метод: list_alerts)
├── alert_management.go      # NEW: 12 методов (Приоритет 2, п.7)
├── cloud_backup.go          # NEW: 12 методов (Приоритет 2, п.8) — НЕ путать с tools/cloudsync.go
├── certificate.go           # NEW: 9 методов (Приоритет 3, п.9)
├── jobs.go                  # готово (2 метода: get_job/abort_job)
├── core.go                  # NEW: 12 методов (Приоритет 3, п.10)
├── directory_services.go    # NEW: 7 методов (Приоритет 3, п.11)
├── disk.go                  # NEW: 8 методов (Приоритет 3, п.12)
├── replication.go           # NEW: 14 методов (Приоритет 3, п.13)
├── user.go                  # NEW: 9 методов (Приоритет 4, п.14)
├── cronjob.go                # NEW: 6 методов (Приоритет 4, п.15)
└── misc.go                  # NEW: 2 метода — device.get_info() + dns.query() (Приоритет 4, п.16-17)
```

### Порядок реализации:
1. Pool Management, остаток (6 методов) — CRITICAL
2. iSCSI (39 методов) — CRITICAL
3. NVMe-oF (36 методов) — CRITICAL
4. Sharing NFS/SMB/WebDAV (21 метод) — HIGH
5. App Management, остаток (15 методов) — HIGH
6. Authentication (16 методов) — HIGH
7. Alert Management, остаток (12 методов) — HIGH
8. Cloud Backup (12 методов) — MEDIUM
9. Certificate Management (9 методов) — MEDIUM
10. Core System, остаток (12 методов) — MEDIUM
11. Directory Services (7 методов) — MEDIUM
12. Disk Management (8 методов) — LOW
13. Replication (14 методов) — LOW
14. User Management (9 методов) — LOW
15. Cron Jobs (6 методов) — LOW
16. Device Management (1 метод) — LOW
17. DNS (1 метод) — LOW
