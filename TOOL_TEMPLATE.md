# Шаблон для реализации TrueNAS MCP Tool

Каждый новый инструмент — это три правки в трёх слоях. Пример: `pool.create()`.

## 1. Клиентский метод — `internal/truenas/pool.go`

Оборачивает вызов TrueNAS JSON-RPC 2.0 (WebSocket, `wss://<host>/api/current`).
`c.call` сам маршалит запрос как `{"jsonrpc":"2.0","id":N,"method":...,"params":[...]}`
и ждёт ответа/job — метод передаётся ровно так, как называется в TrueNAS API
(`pool.create`, не REST-путь).

```go
// CreatePoolParams holds parameters for creating a new ZFS pool.
type CreatePoolParams struct {
    Name   string   `json:"name"`
    Disks  []string `json:"disks"`
    Layout string   `json:"layout"`
}

// CreatePool creates a new ZFS pool.
func (c *Client) CreatePool(ctx context.Context, params *CreatePoolParams) (*Pool, error) {
    if params == nil {
        return nil, errors.New("create pool: params required")
    }
    var pool Pool
    if err := c.call(ctx, "pool.create", []any{params}, &pool); err != nil {
        return nil, fmt.Errorf("create pool: %w", err)
    }
    return &pool, nil
}
```

Если метод асинхронный (TrueNAS возвращает job ID), poll через
`core.get_jobs` — см. `internal/truenas/jobs.go`.

## 2. Мок-интерфейс — `tools/client_iface.go`

Добавить сигнатуру в `truenasClient`, чтобы тесты могли мокать без живого сервера:

```go
CreatePool(ctx context.Context, params *truenas.CreatePoolParams) (*truenas.Pool, error)
```

## 3. MCP-инструмент — `tools/pool.go`

Типизированный вход через `jsonschema`-теги, регистрация через `mcp.AddTool`.

```go
type createPoolInput struct {
    Name   string   `json:"name"   jsonschema:"Pool name"`
    Disks  []string `json:"disks"  jsonschema:"List of disk devices"`
    Layout string   `json:"layout" jsonschema:"ZFS layout, e.g. STRIPE, RAIDZ1"`
}

mcp.AddTool(s, &mcp.Tool{
    Name:        "create_pool",
    Description: "Create a new ZFS pool.",
    Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
}, func(ctx context.Context, _ *mcp.CallToolRequest, p createPoolInput) (*mcp.CallToolResult, any, error) {
    if p.Name == "" || len(p.Disks) == 0 || p.Layout == "" {
        return errorResult(errors.New("create_pool: name, disks, and layout are required"))
    }
    pool, err := client.CreatePool(ctx, &truenas.CreatePoolParams{
        Name:   p.Name,
        Disks:  p.Disks,
        Layout: p.Layout,
    })
    if err != nil {
        return errorResult(fmt.Errorf("create_pool: %w", err))
    }
    return jsonResult(pool)
})
```

Регистрация вызова живёт в `register*Tools` для категории (например
`registerPoolTools`), которая уже подключена в `tools/register.go`.

## Правило 300 строк

Когда файл категории (`tools/pool_management.go`, `internal/truenas/pool.go`, …)
переваливает за ~300 строк, дробить по под-категории TrueNAS API
(например `pool.dataset.*` → отдельный файл от `pool.*`), а не по алфавиту —
границы файлов должны совпадать с границами API-namespace.

## Разрушительные операции

Инструменты, необратимо удаляющие/изменяющие данные (delete, offline, export
без confirm) регистрируются отдельно и только за флагом — см.
`tools/destructive.go` и `Config.AllowDestructive` в `tools/register.go`.
