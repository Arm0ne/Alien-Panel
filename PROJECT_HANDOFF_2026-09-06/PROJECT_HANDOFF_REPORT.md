# X-Panel Central 项目交接报告

> 交接日期：2026-09-06
> 当前代码基线：根仓库 `main`，GitHub 提交 `7f914fe`
> 仓库：[Arm0ne/Alien-Panel](https://github.com/Arm0ne/Alien-Panel)

## 1. 项目结论

本项目已经完成从“读取多个 X-Panel 节点”的技术验证，演进为一个可用于日常运营的轻量中央管理面板：中央服务保存节点、业务用户、流量、出口资产、成本、收费和事件数据；节点 Agent 只读采集本机 X-Panel，并通过 HTTPS/HTTP API 将快照上报中央。

当前版本适合进入公网测试和小规模灰度。核心代码、数据库迁移、前端生产构建、Docker/Nginx 部署文件和自动化测试均已在仓库中。正式生产前仍需要完成真实节点灰度、连续运行观察、数据逐项核对和回滚演练。

## 2. 总体架构

| 组件 | 技术 | 职责 |
|---|---|---|
| 中央后端 | Go 1.27、SQLite WAL | 管理员认证、节点 API、同步接收、业务查询、财务、事件和维护任务 |
| 节点 Agent | Go 1.27、systemd | 登录本机 X-Panel、只读采集 Inbound/Client/流量/状态、心跳和完整同步 |
| 管理前端 | Vue 3 + TypeScript + Soybean Admin 精简基座 | 总览、用户、节点、出口 IP、线路兼容入口、财务、事件中心 |
| 反向代理 | Docker 内置 Nginx；宿主机可用 Nginx/Caddy | 静态前端、API/Agent 反代、HTTPS 和安全响应头 |
| 数据库 | SQLite，迁移 001–014 | 单机轻量部署；WAL、外键、幂等迁移、备份恢复 |

请求路径约定：

- 管理员 API：`/api/...`；
- Agent API：`/api/agent/v1/...`，同时保留 `/agent/v1/...` 兼容路径；
- 健康检查：`/health/live`、`/health/ready`；
- Agent 配置中的 `central_url` 推荐填写到 `/api`，例如 `https://panel.example.com/api`。

## 3. 已确定的业务口径

### 3.1 用户和 Inbound

- 业务用户以线路机（relay）上的用户 Inbound 为来源，一个线路机主 Inbound 对应一个中央业务用户。
- Inbound 下的 Client/Email 是设备，不单独计为用户；多个 Client 不增加用户数。
- 落地机 Inbound 是基础设施入口，不进入业务用户列表，也不会自动创建用户。
- 同名 Email 在不同节点不合并，节点和 Inbound 是数据边界。

### 3.2 用户路径和出口

日常分配使用 `user_paths`，不要求先建立线路模板：

1. `relay`：线路机直出，出口 IP 归属于线路机；
2. `landing`：线路机 → 落地机 Inbound → 落地机出口 IP；
3. `external`：线路机 → 独立购买的 S5 出口 IP。

旧的 `routes`、`user_routes`、`route_exit_ips` 保留用于历史数据和兼容接口，但不再是新用户分配的前置步骤。

### 3.3 节点生命周期

- 节点管理地址使用完整 URL，可含协议、端口和面板路径；列表和详情可直接点击运维。
- Node Key 可留空，由中央生成稳定 Key；Agent Token 为一次性安装引导 Token，正式节点 Token 存哈希。
- 停用节点拒绝新的心跳、同步和立即同步请求；历史数据保留。
- 删除节点会事务性清理该节点及其 Inbound、Client、同步记录、快照、成本、出口 IP、用户路径和旧线路兼容记录；业务用户本身不自动删除，避免误删预留用户。

### 3.4 流量和状态

- Agent 只读 X-Panel 累计计数，不调用 `QueryStats(reset=true)`。
- 中央保存快照，按相邻快照计算增量；计数回退产生 `traffic_reset` 事件，绝不产生负流量。
- 用户详情显示累计总量增长曲线，支持 1 小时、6 小时、1 天、7 天；悬停可显示采样速率，但曲线主体是累计使用总量。
- 连续三次成功同步都缺失的 Inbound 才归档；一次或两次缺失只记录事件并保留历史。
- 节点超过 5 分钟无心跳标记离线，历史快照和业务数据不删除。

### 3.5 计费和事件

- 用户支持月付和年付，周期金额独立保存；财务按服务周期和统计月份计算有效用户收入。
- Agent 发现线路机 Inbound 到期时间向后延长时，只生成待确认续费候选，不直接记为收费。
- 管理员确认后写入不可变收费记录；也可以标记为非收费变更。年费按服务区间分摊到对应月份。
- “事件中心”展示业务事件、节点告警和同步失败；立即同步仅作为内部调度记录。
- 事件区分未读和未处理；顶部消息按钮显示红点，登录后和在线期间对待确认续费进行强提醒，并通过浏览器会话避免重复弹窗。

## 4. 源码结构

```text
backend/                         中央 Go 服务、迁移、API、测试
agent/                           节点 Go Agent、X-Panel 适配、systemd 安装
frontend/                        Vue/TypeScript 管理前端（独立上游 Git 历史）
deploy/                          Docker、Nginx/Caddy、备份恢复、构建和冒烟脚本
release/                         可部署二进制和校验文件（源码快照中排除）
deploy/frontend-dist/            前端生产构建（源码快照中排除）
CENTRAL_PANEL_DESIGN.md          原始总体设计
CENTRAL_PANEL_DEVELOPMENT_PLAN.md 原始开发计划和检查点
OPERATIONS_RUNBOOK.md            运维与故障手册
```

核心数据库迁移：

- `001_initial.sql`：管理员、节点、用户、Inbound、Client、线路、出口、成本、流量、审计等基础表；
- `004_exit_ip_sources.sql`：线路机/落地机出口和独立 S5 来源；
- `008_user_paths.sql`：直接用户路径和旧线路数据迁移；
- `010_node_install_tokens.sql`、`011_node_management_url.sql`：一键接入和管理地址；
- `012_traffic_snapshot_inbound_time.sql`：趋势查询索引；
- `013_user_billing.sql`：月/年付、收费记录和续费候选；
- `014_event_center.sql`：事件分类、已读/处理状态、去重和关联资源。

## 5. 当前功能清单

- 管理员登录、刷新 Token、退出、Bearer 认证和来源校验；
- 节点接入、启用/停用、永久删除、管理地址、成本和多出口 IP；
- Ubuntu/Debian amd64 Agent 一行安装命令；手动 YAML/systemd 方式仍保留；
- X-Panel Base Path、Session 过期重登录、v1/v2 响应字段兼容；
- 心跳、完整同步、幂等 `sync_id`、有界重试、节点离线和同步失败事件；
- 用户按 Inbound 聚合、Client 设备、到期状态、路径分配和路径历史；
- 线路机直出、落地机出口、独立 S5 三类出口资产；
- 累计流量快照、趋势、回退检测和数据新鲜度；
- 月费/年费、Agent 续费候选、收费确认和财务历史月份统计；
- 事件中心、待办数量、已读/全部已读、处理和顶部强提醒；
- Docker 一键测试、宿主机 Nginx/Caddy、备份/迁移/恢复和本地冒烟脚本。

## 6. 测试与发布状态

已完成：

- `backend`: `go test ./...`、`go vet ./...` 通过；
- `frontend`: `pnpm.cmd typecheck`、`pnpm.cmd build` 通过；
- Agent 解析夹具、认证、重试、幂等同步、流量趋势、计费和事件集成测试已加入；
- Linux/Windows staging smoke 脚本和安全响应头检查已加入；
- 最新发布包已修正为 Linux/amd64 ELF，避免 Docker 执行 Windows 二进制导致 `exec format error`。

尚未视为正式完成：

- 公网测试服务器上的真实节点逐项数据核对；
- 线路机直出、落地机出口、独立 S5 三种路径的完整业务验收；
- 至少 20%–30% 节点的 3–5 天灰度观察；
- 真实续费、年费分摊、事件提醒和财务账目核对；
- 数据库恢复、Agent 升级回滚和管理员培训记录；
- 100 台节点/单节点 500 Inbound 的性能基线验证。

## 7. 交接时的注意事项

1. 不要把 `frontend/` 推送到当前的 Soybean 上游远程；该目录在主仓库中被忽略，源码快照已单独复制。
2. Docker 镜像中的 `release/xpanel-central` 必须是 Linux/amd64 ELF，不要直接把 Windows 下默认构建的 `.exe` 改名使用。
3. 使用 `--domain` 时，内置 Web 默认绑定 `127.0.0.1:18080`，宿主机 Nginx 应反代到该地址；Agent 应访问域名而不是未开放的公网 18080。
4. 更新中央服务时不要删除 Docker volume 或 `/opt/xpanel-central/.env`；其中包含数据库、管理员密码和 Agent 引导 Token。
5. 删除节点不会自动删除业务用户；预留的线路机 Inbound 和业务用户必须由管理员按业务需要单独处理。

## 8. 配套文档

- [设计与进度差异交接](./DESIGN_PROGRESS_HANDOFF.md)
- [部署手册](./DEPLOYMENT_GUIDE.md)
- [源码快照说明](./source/README.md)
- 原始设计：[../CENTRAL_PANEL_DESIGN.md](../CENTRAL_PANEL_DESIGN.md)
- 原始计划：[../CENTRAL_PANEL_DEVELOPMENT_PLAN.md](../CENTRAL_PANEL_DEVELOPMENT_PLAN.md)
- 运维手册：[../OPERATIONS_RUNBOOK.md](../OPERATIONS_RUNBOOK.md)
