# X-Panel Central Backend

这是中央管理面板的 Go 后端骨架。它只保存中央业务数据，使用 X-Panel API 的同步数据，不修改 X-Panel/Xray 源码，也不调用 Xray 的 reset 统计接口。

## 本地启动

PowerShell：

```powershell
$env:XPANEL_ADMIN_PASSWORD = 'replace-with-a-long-secret'
$env:XPANEL_AGENT_REGISTRATION_TOKEN = 'replace-with-a-separate-bootstrap-secret'
$env:XPANEL_DATABASE = './data/panel.db'
go run ./cmd/server
```

默认监听 `:8090`。前端开发环境已配置为 `http://127.0.0.1:8090/api`。

常用配置：

| 环境变量 | 默认值 | 说明 |
|---|---|---|
| `XPANEL_LISTEN` | `:8090` | HTTP 监听地址 |
| `XPANEL_DATABASE` | `./data/panel.db` | SQLite 数据库路径 |
| `XPANEL_ADMIN_USER` | `admin` | 首次启动创建的管理员用户名 |
| `XPANEL_ADMIN_PASSWORD` | 无默认值 | 必填，首次启动创建管理员密码 |
| `XPANEL_AGENT_REGISTRATION_TOKEN` | 无默认值 | Agent 首次注册接口的引导密钥；未配置时注册接口保持禁用 |
| `XPANEL_SESSION_TTL` | `24h` | Session 有效期 |
| `XPANEL_CORS_ORIGINS` | localhost:9527 | 允许的前端来源，逗号分隔 |

## 校验

```powershell
go test ./...
go build ./cmd/server
```

数据库使用 SQLite WAL，迁移文件位于 `internal/db/migrations/`，服务启动时自动执行且按版本幂等。

## API

接口草稿见 [`openapi.yaml`](./openapi.yaml)。所有响应统一为：

```json
{
  "code": "0000",
  "msg": "ok",
  "data": {}
}
```

当前已实现健康检查、管理员登录/刷新/退出、当前管理员信息、Dashboard 和各业务只读列表。Agent 端点包括：

- `POST /api/agent/v1/register`（注册节点并签发节点 Token）；
- `POST /api/agent/v1/heartbeat`（Bearer 节点 Token）；
- `POST /api/agent/v1/sync`（Bearer 节点 Token，按 `sync_id` 幂等写入 Inbound、Client 和流量快照）。

Agent Token 只以哈希形式保存在 `node_credentials`，注册响应中的明文 Token 仅返回一次。业务字段写 API、流量重置检测和更细的聚合任务按开发进度表继续实现。
