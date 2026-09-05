# X-Panel Central Backend

中央服务部署、节点接入、故障排查、备份恢复和灰度回滚请参阅项目根目录的
[`OPERATIONS_RUNBOOK.md`](../OPERATIONS_RUNBOOK.md)。

这是中央管理面板的 Go 后端骨架。它只保存中央业务数据，使用 X-Panel API 的同步数据，不修改 X-Panel/Xray 源码，也不调用 Xray 的 reset 统计接口。

## 本地启动

PowerShell：

```powershell
$env:XPANEL_ADMIN_PASSWORD = 'replace-with-a-long-secret'
$env:XPANEL_AGENT_REGISTRATION_TOKEN = 'replace-with-a-separate-bootstrap-secret'
$env:XPANEL_DATABASE = './data/panel.db'
go run ./cmd/server
```

默认监听 `:8090`。本地前端开发环境已配置为 `http://localhost:8090/api`；在 Windows 上可避免本机其他程序占用 IPv4 回环地址时的冲突。

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

数据库运维命令位于 `cmd/db-maintenance`：

```powershell
go run ./cmd/db-maintenance backup --database .\data\panel.db --backup-dir .\data\backups --retention 14
go run ./cmd/db-maintenance verify --path .\data\backups\panel-20260903T010203Z.sqlite3
go run ./cmd/db-maintenance migrate --database .\data\panel.db --backup-dir .\data\backups
go run ./cmd/db-maintenance restore --source .\data\backups\panel-20260903T010203Z.sqlite3 --database .\data\panel.db --yes
```

`backup` 使用 SQLite `VACUUM INTO` 生成不依赖 WAL/SHM 的一致性快照，默认保留
14 个版本；`verify` 执行 `integrity_check` 和 `foreign_key_check`；`migrate`
会在已有数据库升级前自动备份，然后执行内置幂等迁移。恢复前必须停止中央服务，
`restore` 会先生成 `panel-pre-restore-*.sqlite3` 安全快照并在替换后再次校验。

本地联调可写入一组可重复的模拟业务数据（不会覆盖非 `demo-*` 记录）：

```powershell
go run ./cmd/seed-demo --database .\data\panel.db
```

该命令包含模拟线路机、落地机、线路、出口 IP、一个用户、两个 Client 和流量快照，
用于验证用户详情、业务字段编辑、设备、线路和流量展示；公网或生产数据库不要执行。

## API

接口草稿见 [`openapi.yaml`](./openapi.yaml)。所有响应统一为：

```json
{
  "code": "0000",
  "msg": "ok",
  "data": {}
}
```

当前已实现健康检查、管理员登录/刷新/退出、当前管理员信息、Dashboard、各业务列表、用户详情和中央业务字段编辑、用户直接路径管理、节点详情和同步请求，以及线路关系 CRUD。Agent 端点包括：

- `POST /api/agent/v1/register`（注册节点并签发节点 Token）；
- `POST /api/agent/v1/heartbeat`（Bearer 节点 Token）；
- `POST /api/agent/v1/sync`（Bearer 节点 Token，按 `sync_id` 幂等写入 Inbound、Client 和流量快照）。

节点页面接口：

- `POST /api/nodes` 由管理员创建中央节点记录并生成一次性显示的 Agent Token；`nodeKey` 可省略，中央会生成稳定的 `node-<random>` 标识并写入 Agent 配置模板；Token 只保存哈希，X-Panel 用户名和密码只应放在 Agent 节点本机配置中；
- 节点创建同时生成 15 分钟有效、只能使用一次的在线安装 Token；`POST /api/agent/v1/bootstrap` 由 Ubuntu/Debian 一键安装脚本兑换正式 Agent Token，安装 Token 只保存哈希且兑换后立即失效；管理员可通过 `POST /api/nodes/{id}/install-token` 重新生成安装 Token；
- `PATCH /api/nodes/{id}` 修改节点元数据或启用状态。停用后该节点的 Agent heartbeat、完整同步和立即同步请求都会被拒绝，重新启用后恢复认证；兼容的 Agent register 流程不会覆盖管理员的停用状态；
- `DELETE /api/nodes/{id}` 为直接删除：节点、Agent Token、Inbound/Client、流量快照、同步记录、成本、出口 IP、用户路径以及旧线路兼容记录会在一个事务中清理；业务用户本身不会删除，删除后可重新使用原 Node Key；
- `POST /api/nodes` 的 `exitIps` 可传入最多 100 个 IPv4/IPv6 地址，创建时作为该节点的出口 IP 资产原子写入；`publicIp` 仍表示节点主公网/管理地址，不替代出口资产列表。出口 IP 的服务商、成本、有效期和备注可在出口 IP 页面补充；
- `GET /api/nodes/{id}` 返回节点元数据、Inbound、该节点拥有的公网出口 IP（线路机和落地机）、最近同步运行和状态事件；
- `POST /api/nodes/{id}/sync` 写入一次立即同步请求事件并返回 `queued`。中央服务不反向调用 X-Panel，节点 Agent 会在下一次周期同步中执行。
- `GET /api/nodes/{id}/costs` 查询节点成本记录；`POST /api/nodes/{id}/costs` 录入节点月成本；`PATCH /api/nodes/{id}/costs/{costId}` 修改成本类别、金额和备注。生效日期属于历史版本，不能在编辑时改动；日期变化应新增一条成本记录。
- 节点详情中的成本记录按生效日期展示，财务汇总按所选月份与有效区间计算节点成本。
- `GET /api/costs/summary?period=YYYY-MM` 返回有效用户数、用户月费收入、节点/出口 IP/其他成本和预计毛利润。收入按统计月份计算：用户需在该月结束前创建、统计月开始时尚未到期且未停用；不使用当前状态筛选历史月份。成本按月份半开区间与记录有效区间求交集，避免下月生效成本提前计入。
- 线路和出口 IP 列表/详情返回“配置归属用户数”：仅统计当前有效且未停用用户；出口 IP 还要求线路、出口 IP 资产和绑定均为启用状态。
- 管理员 Bearer 写请求若携带 Origin/Referer，必须匹配 `XPANEL_CORS_ORIGINS`；无浏览器来源头的 CLI 客户端仍可使用。当前认证不使用 Cookie，因此这是面向未来 Cookie 迁移的纵深 CSRF 防护。
- refresh token 为单次使用并在轮换时撤销旧会话；管理员 logout 会撤销当前会话；同一 Agent 重新注册会撤销该节点旧 Token。所有受保护接口均检查会话/节点凭据、有效期和启用状态。
- 分页接口统一返回 `dataAt`（最新成功同步时间；尚无成功同步时为 `null`），总览和财务汇总也返回同一数据时间，供前端判断数据延迟或过期。
- 用户列表回归保证按 Inbound 聚合为一行；同一 Email 出现在不同节点时仍属于不同业务用户，Client 数仅统计该 Inbound 的设备凭证。

用户详情接口：

- `GET /api/users/{id}` 返回业务用户、主 Inbound、Client/Email 设备、线路机、当前用户路径、路径历史和已分配线路兼容快照；
- `GET /api/users/{id}/path-assets` 返回该用户可用的路径资源：用户主 Inbound 所属线路机、启用的落地节点及其基础设施 Inbound、线路机出口 IP 和独立 S5 出口 IP。落地 Inbound 仅作为路径资源，不会创建或展示为业务用户；没有成功同步时会返回 `inboundState=pending`，而不是静默空列表；
- `PATCH /api/users/{id}` 仅更新中央维护的 `displayName`、`monthlyFee`、`currency`（当前仅支持 `CNY`）和 `notes`，并写入 `audit_logs`；
- `PUT /api/users/{id}/route` 是旧线路模板分配兼容接口；可选 `routeExitIpId` 固定到该线路已绑定的某个出口 IP，省略时按线路出口池权重分配；`DELETE /api/users/{id}/route` 解除当前线路并保留历史关系。新日常流程应使用下方的 `user_paths` 接口；这些操作只更新中央配置，不会直接改写 X-Panel/Xray；
- `PUT /api/users/{id}/path` 保存新的直接路径。线路机由用户主 Inbound 自动确定；线路机直出只能绑定线路机出口 IP，经落地机时必须同时指定启用的落地节点和该节点的基础设施 Inbound，并绑定该落地机出口 IP；独立 S5 不得与落地节点混用。保存新路径会关闭旧路径并保留历史；`DELETE /api/users/{id}/path` 解除当前路径但不删除历史；
- 用户详情页不会修改 X-Panel/Xray，也不会写入 X-Panel 的到期、启用、Client 或流量字段。

线路关系接口：

- `GET /api/routes/{id}` 返回线路机、落地机、Outbound/Inbound 标签、有效期和绑定统计；
- `POST /api/routes`、`PATCH /api/routes/{id}` 创建或更新线路关系，并校验节点类型、落地 Inbound 归属及日期范围；
- `DELETE /api/routes/{id}` 仅允许删除没有用户或出口 IP 绑定的线路；有绑定时返回 `409`，应先停用线路。
- `GET /api/routes/{id}/exit-ips` 查询线路已绑定的出口 IP；
- `POST /api/routes/{id}/exit-ips` 按 `scope` 绑定出口：`relay` 只能绑定该线路线路机直出、`landing` 只能绑定该线路落地机、`external` 只能绑定独立 S5；省略 `scope` 时按资产归属兼容推断；重复或跨节点绑定会被拒绝；
- `scope` 属于旧线路出口池；新用户路径不依赖线路出口池，直接按出口资产归属校验。
- `PATCH /api/routes/{id}/exit-ips/{exitIpId}` 更新绑定权重和启用状态，`DELETE` 解绑；绑定变更会刷新线路和出口 IP 的配置归属统计并写入审计日志。

出口 IP 接口：

- `GET /api/exit-ips/{id}` 返回地址、协议族、来源类型、所属节点/独立 S5、成本、有效期和配置归属用户数；旧 `landingNodeId` 字段保留兼容；
- `POST /api/exit-ips`、`PATCH /api/exit-ips/{id}` 创建或更新 IPv4/IPv6 出口 IP 资产。`sourceType=node` 时用 `ownerNodeId` 指向线路机或落地机，`sourceType=s5` 表示独立购买的 S5（不绑定节点）；旧 `landingNodeId` 请求仍兼容为落地机资产。接口校验来源、地址族、成本和日期范围；
- `DELETE /api/exit-ips/{id}` 仅允许删除没有线路绑定或当前用户路径的资产；有绑定时返回 `409`，应先停用出口 IP 或调整用户路径。

Agent Token 只以哈希形式保存在 `node_credentials`，注册响应中的明文 Token 仅返回一次。业务字段写 API 和更细的报表聚合按开发进度表继续实现。

已登记为 `relay` 的节点成功同步后，中央会将每个新发现的 Inbound 自动建立为一个业务用户，并关联其 Client/Email 设备。同步只更新来自 X-Panel 的 Inbound、到期、流量和设备字段，不会覆盖中央维护的业务名称、月费、币种和备注。落地节点的 Inbound 不会自动创建业务用户。

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
