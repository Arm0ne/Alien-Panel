# 架构与底层代码

## 1. 系统边界

Alien-Panel 是中央只读运营管理面板。它管理节点、业务用户、设备、累计流量、出口资产、成本、收费历史和事件；节点 Agent 负责从本机 X-Panel 读取数据并上报中央。中央不直接登录 X-Panel，也不执行 Xray 的 reset 统计接口。

```text
X-Panel/Xray
    │  本机 HTTP API，只读
    ▼
节点 Agent ── HTTPS/HTTP ──> 中央 Go 服务 ── SQLite WAL
                                  │
                                  ├── /api/... 管理员 API
                                  ├── /agent/v1/... Agent API
                                  ├── /health/live、/health/ready
                                  └── Vue 前端 / Nginx 或 Caddy
```

Agent 心跳和 X-Panel 采集是两条独立状态链：Agent 能向中央发送心跳时节点可以保持在线；X-Panel 采集失败只标记同步异常，不把节点误判为离线。

## 2. 组件职责

### 中央后端 `backend/`

- `cmd/server/main.go`：读取环境变量、打开 SQLite、执行迁移、启动 HTTP 服务和后台维护任务。
- `internal/db/db.go`：嵌入并按文件名顺序幂等执行 `internal/db/migrations/*.sql`，启用 WAL、外键和忙等待。
- `internal/httpapi/server.go`：路由注册、认证中间件、用户/节点/出口/财务/Dashboard 查询。
- `internal/httpapi/agent.go`：Agent 注册、bootstrap、heartbeat、完整同步、幂等 `sync_id` 和同步健康状态。
- `internal/httpapi/user_path.go`：用户路径保存、历史、有效性和多个固定出口 IP。
- `internal/httpapi/user_path_assets.go`：按用户主 Inbound 计算可选择的线路机、落地机和 S5 出口资产。
- `internal/httpapi/node_admin.go`：节点删除、关联数据清理和孤立业务用户处理。
- `internal/httpapi/orphan_users.go`：节点删除后，只删除已没有任何线路机 Inbound/路径关联的业务用户；收费记录、续费候选和财务审计历史保留。
- `internal/httpapi/node_cost.go`、`finance.go`：节点/IP/其他真实成本、服务期收入和毛利统计。
- `internal/httpapi/event.go`：同步异常、节点状态、流量回退和续费候选事件。
- `cmd/db-maintenance`：一致性备份、迁移、校验和恢复。
- `cmd/traffic-check`：只读重算累计快照的流量核对工具。

### 节点 Agent `agent/`

- `cmd/agent/main.go`：命令行入口、版本输出和运行器启动。
- `internal/config`：YAML 配置和环境变量解析，敏感文件由 systemd 用户读取。
- `internal/xpanel/client.go`：X-Panel Session 登录、Base Path 拼接、超时和重新认证。
- `internal/collector/collector.go`：兼容不同 X-Panel 响应外壳，解析 Inbound、Client、状态和累计流量；不调用 reset。
- `internal/central/client.go`：向中央发送注册后的 heartbeat/sync 请求，保留稳定请求 ID 和 `sync_id`。
- `internal/runner/runner.go`：首次立即同步、周期同步、心跳、失败重试和退避。
- `internal/buildinfo`：嵌入 `agent/VERSION`、提交号和构建时间；操作员看到的是语义版本，例如 `v1.0.3`。

### 管理前端 `frontend/`

Vue 3 + TypeScript + Naive UI，基于精简的 Soybean Admin 工程。业务页面位于 `src/views/`：

- `dashboard`：运营总览、业务流量趋势、节点和用户排行、到期与事件。
- `users`：按 Inbound 的业务用户列表、Client 设备、计费、路径和多个固定出口 IP。
- `nodes`：节点总览、Agent/X-Panel 状态、Inbound、线路机负载、成本和重命名。
- `exit-ips`：线路机/落地机/S5 出口资产、区域、成本和配置归属用户数。
- `finance`：已确认实收、成本、毛利、订单、续费率和成本明细。
- `events`：业务事件、同步异常、续费待办和已读/处理状态。

生产环境使用 `deploy/frontend-dist/`，服务器不需要 Node.js。源码包中的 `frontend/` 不包含 `node_modules` 和 `dist`，用于后续开发和重新构建。

## 3. 请求与同步流程

### Agent 注册

1. 管理员在节点页面创建节点，中央生成短期一次性的 bootstrap 安装 Token。
2. 节点执行 `agent/deploy/install-online.sh`，本地询问 X-Panel 用户名和密码。
3. 脚本调用 `/api/agent/v1/bootstrap`，换取正式节点 Token；Token 只在目标机配置文件中保存明文，中央只保存哈希。
4. Agent 发送 heartbeat 和完整 sync，中央以 `node_key` 定位节点并刷新版本、状态和数据时间。

### 完整同步

1. Agent 使用本机 X-Panel Session API 读取 Inbound、Client、累计流量和运行状态。
2. 解析结果携带 `sync_id`，中央在事务中更新节点、Inbound、Client 和流量快照。
3. 同一 `sync_id` 重试时返回之前的结果，不重复写入。
4. 只有成功的完整同步才增加 Inbound 缺失计数；连续三次成功同步均缺失才归档。
5. 相邻累计流量下降会写入 `traffic_reset` 事件并重新建立基线，绝不写入负增量。

## 4. 认证与安全

- 管理员使用短期 access token 和一次性轮换 refresh token；logout 会撤销当前会话。
- Agent 使用节点 Bearer Token；节点停用后 heartbeat、sync 和立即同步都会被拒绝。
- 浏览器写请求的 `Origin/Referer` 必须匹配 `XPANEL_CORS_ORIGINS`；无来源头的 CLI Agent 不受此限制。
- 生产服务建议只监听回环，由宿主机 Nginx/Caddy 终止 TLS；临时端口只用于测试。
- 日志禁止记录密码、管理员 Token、Agent Token 和 X-Panel Session。
