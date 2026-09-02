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
| `XPANEL_MAINTENANCE_INTERVAL` | `1m` | 到期用户与节点在线状态的后台刷新间隔 |
| `XPANEL_CORS_ORIGINS` | localhost:9527 | 允许的前端来源，逗号分隔 |

## 校验

```powershell
go test ./...
go vet ./...
go build ./cmd/server
go build ./cmd/traffic-check
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

Agent Token 只以哈希形式保存在 `node_credentials`，注册响应中的明文 Token 仅返回一次。业务字段写 API 和更细的报表聚合按开发进度表继续实现。

## 节点同步和 Inbound 归档

每次成功的完整 Agent 同步都会将本次未出现的 Inbound 标记为缺失；首次缺失会写入 `inbound_missing` 事件。仅当同一个 Inbound 连续三次成功完整同步均缺失时，中央库才把它标记为归档并写入 `inbound_archived` 事件。

重新出现在后续快照中的 Inbound 会清除缺失时间、缺失计数和中央归档标记。这个策略不会删除 X-Panel 内的内容，也不会物理删除中央库中的 Inbound、Client、流量快照或业务历史；同步失败或节点离线不会增加缺失计数。

## 流量快照核对

`traffic-check` 只读打开中央 SQLite 数据库，使用和 Dashboard 相同的累计快照算法重算流量，适用于灰度部署、异常流量和账单核对。它不修改数据库，也不会访问 X-Panel/Xray。

```powershell
go run ./cmd/traffic-check --database .\data\panel.db --from 2026-09-01 --to 2026-09-02 --format json
```

- `--from` 为包含边界，`--to` 为不包含边界；可使用 RFC3339 时间或 `YYYY-MM-DD`。
- 日期形式的 `--to` 会自动扩展到次日 00:00 UTC，方便核对完整自然日。
- 可加 `--node-key <节点标识>` 仅核对一个节点；`--format` 可选 `text`（默认）或 `json`。
- 首条快照只建立基线；正常情况下累计值相减。检测到 X-Panel 流量重置或累计值回退时，以当前值重新建基线并计入该次累计值，报告会列出重置次数。
