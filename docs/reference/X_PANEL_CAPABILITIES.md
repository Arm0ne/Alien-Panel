# X-Panel 完整能力说明

## 0. 文档定位与审计声明

本文档描述本地拉取的 X-Panel 免费版源码，重点说明它实际提供的 Web API、SQLite 数据、Xray 适配和后台任务，并评估它作为当前多节点运营项目节点端的能力边界。

**已固定的审计版本**

- 源码目录：`vendor-src/X-Panel`
- Git tag：`v26.6.18`
- Git commit：`2448d1bcbe2056e3014a6fef38ad095755fd1dd2`
- 提交时间：2026-06-18
- Go module 中的 Xray 依赖：`github.com/xtls/xray-core v1.260327.0`
- 运行形态：单节点 Go 程序，Xray 子进程，本机 SQLite

源码审计覆盖启动流程、Web/Session、中间件、Controller、Service、数据库模型、Xray gRPC 封装、流量和设备任务、订阅服务、服务器运维和 Telegram 可选模块。第三方依赖、发行包脚本和每个前端组件未逐行审计；生产使用前仍应按安装包和实际配置做安全验收。

## 1. 结论摘要

X-Panel 是“单节点 Xray 管理面板”，不是多节点控制平面。它的核心工作是把本地 SQLite 中的 Inbound/Client 配置生成 Xray `config.json`，启动/停止本机 Xray，并通过本机 gRPC 读取 Handler 和 Stats。

原生适合：

- 单节点 Inbound/Client 增删改查；
- Client UUID/密码、启用、限额、到期、订阅和设备限制；
- 本机 Xray 配置生成、启动、停止、重启、升级和 Geo 文件更新；
- Inbound、Client、Outbound 的本地累计流量；
- 基于 access.log 的客户端来源 IP 和设备限制；
- 订阅链接和 JSON 订阅；
- 本机系统状态、Xray 版本、日志、Telegram 运维。

不适合直接承担：

- 多节点目录、线路机/落地机关系和多出口 IP 资产；
- 中央 RBAC、审计、财务、订阅订单和利润；
- 历史映射版本和共享用户数；
- 用户到最终公网出口 IP 的精确路径或字节归属；
- 强一致、可重放、带签名的跨节点流量上报。

对当前已确认的业务口径，中央系统应手工维护“线路机 -> 落地机出口 IP”关系，并按线路机下有效 Inbound 的 `COUNT(DISTINCT inbound_id)` 统计共享用户数。X-Panel 只负责提供线路机 Inbound/Client 配置事实，不需要修改 Xray 才能完成这个指标。

## 2. 总体架构与启动流程

```text
main.go
  |
  +-- database.InitDB()       -> SQLite + GORM + AutoMigrate
  +-- XrayService              -> 生成配置、启动/停止 Xray、gRPC
  +-- InboundService            -> Inbound/Client CRUD、流量/到期
  +-- ServerService             -> 系统状态、版本、日志、备份
  +-- SettingService            -> settings 键值表
  +-- web.Server                -> Gin 管理界面/API + Session
  +-- sub.Server                -> 可选公开订阅服务
  +-- Telegram Bot              -> 可选运维通知/操作
  +-- Cron/Job                  -> 健康、流量、IP、清理和通知
```

启动顺序（`main.go:32-119`）：

1. 读取 `XUI_DEBUG`/日志级别并加载 `.env`。
2. 初始化 SQLite，自动迁移模型；空库创建 `admin/admin` 的 bcrypt 账号。
3. 创建 `XrayService`、`InboundService`、`ServerService`、`SettingService` 等实例。
4. 注入本机 `XrayAPI` 和可选 Telegram 服务。
5. 启动 Web Server 和可选订阅 Server。
6. 启动设备限制任务，监听 `SIGHUP`/`SIGTERM`。
7. `SIGHUP` 会停止并重新创建 Web/订阅服务；`SIGTERM` 停止 Xray、任务和 HTTP 服务。

Web Server 启动时初始化 Gin、域名校验、gzip、i18n、Cookie Session 和模板。配置证书时监听指定地址并启用 HTTPS；无证书时强制回环监听，适合通过 SSH 隧道访问。HTTP 读写超时为 120 秒。

证据：`main.go:32-208`、`web/web.go:202-498`、`sub/sub.go:32-209`。

## 3. 目录职责与运行文件

| 目录 | 职责 |
|---|---|
| `main.go` | 进程入口、依赖注入、信号和自定义设备任务 |
| `config/` | 版本、名称、环境变量和路径 |
| `database/` | SQLite/GORM、迁移、历史和备份 |
| `database/model/` | User、Inbound、Setting、IP、彩票等模型 |
| `web/controller/` | 登录、页面、Inbound/Server/Setting/Xray API |
| `web/service/` | 配置生成、Xray 控制、流量、状态、Telegram、设置 |
| `web/job/` | Xray 健康、流量、客户端 IP、日志和通知任务 |
| `web/session/` | Cookie Session 和当前登录用户 |
| `xray/` | Xray 进程、配置结构、gRPC Handler/Stats 封装 |
| `sub/` | 普通订阅和 JSON 订阅 |
| `util/` | 加密、系统命令、JSON、随机数等工具 |
| `web/html`、`web/assets` | 页面模板、静态资源和多语言文件 |

主要环境变量：

| 变量 | 默认/作用 |
|---|---|
| `XUI_DEBUG` | `true` 开启 Gin/GORM 调试 |
| `XUI_LOG_LEVEL` | 面板日志级别 |
| `XUI_BIN_FOLDER` | Xray 二进制和配置目录，默认 `bin` |
| `XUI_DB_FOLDER` | SQLite 目录；Linux 默认 `/etc/x-ui` |
| `XUI_LOG_FOLDER` | 日志目录；Linux 默认 `/var/log` |

Linux 默认数据库为 `/etc/x-ui/<name>.db`，Web 默认端口约为 `13688`，订阅默认端口约为 `13788`，具体以 `settings` 表为准。Xray 进程使用 `bin/xray-<os>-<arch>`、`bin/config.json`、`geoip.dat` 和 `geosite.dat`。

## 4. SQLite 数据模型

数据库由 `database/db.go:30-48` 自动迁移，当前主要表：

### 4.1 `users`

```go
type User struct {
    Id       int
    Username string
    Password string // bcrypt
}
```

只有基础登录用户字段，没有角色、权限、组织或节点范围。

### 4.2 `inbounds`

字段包括：`id`、`user_id`、`up`、`down`、`total`、`all_time`、`remark`、`enable`、`expiry_time`、`device_limit`、`listen`、`port`、`protocol`、`settings`、`stream_settings`、`tag`、`sniffing`。

Client 通常不在独立表，而在：

```text
inbounds.settings JSON -> clients[]
```

因此读取和更新 Client 需要读取、修改并重新序列化整个 Inbound 的 settings JSON。

### 4.3 `client_traffics`

```go
type ClientTraffic struct {
    Id, InboundId int
    Enable        bool
    Email         string // unique
    Up, Down      int64
    AllTime       int64
    ExpiryTime    int64
    Total         int64
    Reset         int
    LastOnline    int64
}
```

`email` 在本地数据库中全局唯一。它不适合作为中央系统跨节点主键；中央应使用 `node_id + remote_inbound_id + remote_client_id`，或另建业务用户 ID。

### 4.4 其他模型

- `outbound_traffics`：Outbound tag、up、down、total。
- `inbound_client_ips`：Client email 和来源 IP 字符串。
- `settings`：键值设置表，包含 Xray 模板、端口、证书、订阅和外部流量通知配置。
- `history_of_seeders`、`link_history`、`lottery_wins`：迁移记录、订阅历史和业务扩展表。

## 5. Xray 配置生成与进程管理

### 5.1 配置生成

`XrayService.GetXrayConfig` 的流程是：

1. 从 `settings.xrayTemplateConfig` 读取模板。
2. 查询本地所有 Inbound，只加入 `enable=true` 的 Inbound。
3. 读取 `settings.clients[]` 和 `ClientTraffic`。
4. 过滤 settings 或流量记录中禁用的 Client。
5. 根据 `SpeedLimit` 创建对应 Policy level，并给用户写入 level。
6. 组装协议用户字段，清理部分 TLS/Reality 敏感字段和 `externalProxy`。
7. 返回面板内部的 `xray.Config`，再序列化为 Xray JSON。

Client 的 `SpeedLimit` 单位为 KB/s；源码把限速映射到 Xray 用户 level，默认 level 还控制用户流量/在线统计开关。模板字段必须和目标 Xray 版本兼容，不能只依赖面板 UI 的默认值。

### 5.2 启动、重启和停止

面板把配置写入 `bin/config.json`，以子进程方式执行 Xray。Inbound/Client 写操作通常先更新 SQLite，再尝试 HandlerService 热加载；失败时设置“需要重启”标志。后台约每秒检查运行状态，每 30 秒处理重启标志；崩溃时可重新启动。

证据：`xray/process.go:174-256`、`web/service/xray.go:103-481`、`web/web.go:292-318`。

### 5.3 本机 gRPC 适配

`xray/XrayAPI.Init` 连接：

```text
127.0.0.1:<由 tag=api 的 Inbound 解析出的端口>
```

使用 insecure gRPC 的 `HandlerService` 和 `StatsService`。当前封装支持：

- `AddInbound`、`RemoveInbound`；
- 对 VMess/VLESS/Trojan/Shadowsocks/Shadowsocks 2022 执行 AddUser/RemoveUser；
- `QueryStats(reset)` 解析 user/inbound/outbound counter。

未在 `AddUser` 中处理的协议不会被动态添加，通常需要写库后重启 Xray。HandlerService 的运行时改动不会自动持久化到面板 SQLite。

证据：`xray/api.go:28-243`、`xray/client_traffic.go:3-15`。

## 6. Web 页面和 HTTP API

以下路径都受配置的 `basePath` 影响；实际 API 前缀通常是 `<basePath>/panel/api`。

### 6.1 登录和页面

| 方法 | 路径 | 功能 |
|---|---|---|
| GET | `/` | 登录页或跳转 |
| POST | `/login` | 用户名、密码、可选 TOTP 登录 |
| GET | `/logout` | 注销 |
| POST | `/getTwoFactorEnable` | 查询 2FA 状态 |
| GET | `/panel/` | 首页 |
| GET | `/panel/inbounds` | Inbound 页 |
| GET | `/panel/settings` | 设置页 |
| GET | `/panel/xray` | Xray 设置页 |
| GET | `/panel/navigation` | 导航页 |
| GET | `/panel/servers` | 服务器页 |

### 6.2 Inbound API

| 方法 | 路径 | 功能 |
|---|---|---|
| GET | `/inbounds/list` | 当前登录用户的 Inbound 列表 |
| GET | `/inbounds/get/:id` | 查询单个 Inbound |
| GET | `/inbounds/getClientTraffics/:email` | 按 email 查询 Client 流量 |
| GET | `/inbounds/getClientTrafficsById/:id` | 按 Client ID 查询流量 |
| POST | `/inbounds/add` | 新增 Inbound |
| POST | `/inbounds/del/:id` | 删除 Inbound |
| POST | `/inbounds/update/:id` | 更新 Inbound |
| POST | `/inbounds/clientIps/:email` | 查询来源 IP 记录 |
| POST | `/inbounds/clearClientIps/:email` | 清理来源 IP |
| POST | `/inbounds/addClient` | 新增 Client |
| POST | `/inbounds/:id/delClient/:clientId` | 删除 Client |
| POST | `/inbounds/updateClient/:clientId` | 更新 Client |
| POST | `/inbounds/:id/resetClientTraffic/:email` | 重置单个 Client 流量 |
| POST | `/inbounds/resetAllTraffics` | 重置所有 Inbound 流量 |
| POST | `/inbounds/resetAllClientTraffics/:id` | 重置某 Inbound 下 Client 流量 |
| POST | `/inbounds/delDepletedClients/:id` | 删除流量耗尽 Client |
| POST | `/inbounds/import` | 导入 Inbound |
| POST | `/inbounds/onlines` | 最近流量采样中的在线 Client |
| POST | `/inbounds/lastOnline` | Client 最近在线时间 |
| POST | `/inbounds/updateClientTraffic/:email` | 手工修改 Client 流量 |

证据：`web/controller/inbound.go:26-48`。

### 6.3 Server API

| 方法 | 路径 | 功能 |
|---|---|---|
| GET | `/server/status` | CPU、内存、磁盘、网络、连接、Xray 状态 |
| GET | `/server/getXrayVersion` | 从 GitHub 获取可安装版本 |
| GET | `/server/getConfigJson` | 返回当前生成配置 |
| GET | `/server/getDb` | 下载完整 SQLite |
| GET | `/server/getNewUUID` | 生成 UUID |
| GET | `/server/getNewX25519Cert` | 生成 X25519 密钥 |
| GET | `/server/getNewmldsa65` | 生成 ML-DSA 密钥 |
| GET | `/server/getNewmlkem768` | 生成 ML-KEM 密钥 |
| GET | `/server/getNewVlessEnc` | 生成 VLESS 加密参数 |
| POST | `/server/getNewEchCert` | 生成 ECH 证书 |
| POST | `/server/stopXrayService` | 停止 Xray |
| POST | `/server/restartXrayService` | 重启 Xray |
| POST | `/server/installXray/:version` | 下载并安装指定版本 |
| POST | `/server/updateGeofile[/:fileName]` | 更新 Geo 文件 |
| POST | `/server/logs/:count` | 获取面板日志 |
| POST | `/server/xraylogs/:count` | 获取 Xray 日志 |
| POST | `/server/importDB` | 导入 SQLite 并重启 |
| POST | `/server/history/save` | 保存订阅链接历史 |
| GET | `/server/history/load` | 查询订阅链接历史 |
| POST | `/server/install/subconverter` | 安装 Subconverter |
| POST | `/server/openPort` | 调用 UFW 放行端口 |

证据：`web/controller/server.go:43-67`。

### 6.4 Xray 和设置 API

Xray：`<basePath>/panel/xray`：

- `POST /`：Xray 模板和 Inbound tags；
- `POST /update`：保存模板；
- `GET /getXrayResult`：最后错误/状态；
- `GET /getDefaultJsonConfig`：默认配置；
- `POST /warp/:action`：WARP 数据、配置、注册、授权；
- `GET /getOutboundsTraffic`、`POST /resetOutboundsTraffic`：Outbound 流量。

设置：`<basePath>/panel/setting`：

- `POST /all`、`POST /defaultSettings`、`POST /update`；
- `POST /updateUser`：修改管理员账号密码；
- `POST /restartPanel`：重启面板；
- `GET /getDefaultJsonConfig`：默认 Xray 配置。

证据：`web/controller/xray_setting.go:24-34`、`web/controller/setting.go:34-43`。

## 7. 认证、Session 和资源归属

登录时使用 bcrypt 校验密码，可选 TOTP；成功后把 `model.User` 写入 Cookie Session。Cookie 名称为 `3x-ui`，store secret 来自设置表，Session 最大年龄来自 `sessionMaxAge`。

`/panel/api` 使用 `checkLogin`：未登录的 Ajax 返回 401，普通页面重定向到 Base Path。当前模型没有 RBAC；Inbound 列表会按 `user_id` 过滤，但多个单资源接口按 ID/email 全局查找，资源归属检查不完整。

重点风险：

- Session 明确设置了 Path、MaxAge、HttpOnly，但未显式设置 `Secure`、`SameSite`。
- 未发现通用 CSRF Token 中间件，Cookie Session 的写接口需要额外防护。
- 登录失败日志和 Telegram 通知可能包含 HTML 转义后的明文密码。
- `getDb` 可下载完整数据库，包含密码材料、Client 凭据、订阅 ID、流量和设置。
- 错误消息可能把内部错误文本返回给客户端。
- `X-Real-IP`/`X-Forwarded-For` 的信任边界需要由反向代理配置约束。
- `APIController` 创建 ServerController 时存在依赖注入不完整风险，二次开发应统一传入服务实例。

证据：`web/controller/index.go:51-96`、`web/session/session.go:21-65`、`web/controller/base.go:15-26`、`web/controller/api.go:23-34`、`web/service/inbound.go:35-43`。

## 8. 流量统计和在线状态

### 8.1 默认十秒任务

`web/web.go` 启动后约每 10 秒执行 `XrayTrafficJob`：

1. 检查 Xray 是否运行。
2. 调用 `XrayService.GetXrayTraffic()`。
3. 内部执行 `StatsService.QueryStats(reset=true)`。
4. 解析 user/inbound/outbound counter。
5. 写入本地 Inbound、ClientTraffic、OutboundTraffics。
6. 可选 POST 到 `externalTrafficInformURI`。
7. 触发到期、周期流量、自动恢复/禁用等处理。

因此 Xray counter 在面板读取后会被清零。本次返回值是“上一次任务到本次任务之间的区间增量”，不是可重放快照。

### 8.2 累计规则

```text
Inbound:  up += delta_up; down += delta_down; all_time += up + down
Client:   up += delta_up; down += delta_down; all_time += up + down
Outbound: up += delta_up; down += delta_down; total = up + down
```

### 8.3 `/onlines` 和 `lastOnline`

`/onlines` 不是 Xray 当前连接列表。面板只把本次采样中 `up + down > 0` 的 email 放入内存在线列表，并更新 `last_online`。长连接但采样区间没有字节的用户可能不出现在 `/onlines`；面板重启后内存在线列表也会丢失。

Xray 原生 OnlineMap（如果目标版本和 Policy 启用）记录的是连接来源 IP，不是落地机最终公网出口 IP；X-Panel 当前依赖版本没有把最新 `GetUsersStats` RPC 作为稳定契约。

### 8.4 外部流量通知

设置项：`externalTrafficInformEnable`、`externalTrafficInformURI`。Payload 类似：

```json
{
  "clientTraffics": [],
  "inboundTraffics": []
}
```

它没有节点 ID、出口 IP、签名、稳定 batch/sample ID、可靠重试、持久化队列和重放接口，也不包含 Outbound 计数。中央接收端只能把它当作“面板清零后的不可重放区间”，不能作为严格财务对账。

不能让中央服务同时对同一节点执行 `QueryStats(reset=true)`，否则两个消费者会互相清零并丢失流量。

证据：`web/job/xray_traffic_job.go:23-71`、`web/service/xray.go:403-419`、`web/service/inbound.go:956-1316`。

## 9. 客户端来源 IP 和设备限制

`CheckClientIpJob` 每 10 秒扫描 Xray access.log，解析用户 email 和 `from` 后的来源 IP：

- 保存到 `InboundClientIps`；
- 维护 email -> IP -> lastSeen 的内存状态；
- 约 3 分钟没有新日志则视为不活跃；
- Client 的 `limitIp`/Inbound 的 `device_limit` 超限时可通知、封禁或通过 HandlerService 删除并恢复临时用户；
- Linux 还可检查 Fail2Ban；
- 日志会按条件复制到持久日志并截断。

这些 IP 是“访问 Xray 入站的客户端来源 IP”，不是落地机对公网连接使用的出口 IP，也不能证明用户走过哪一个 NAT 出口。它们适合设备限制和安全审计，不适合当前出口 IP 共享用户数统计。

证据：`web/job/check_client_ip_job.go:83-265`、`428-660`。

## 10. 订阅服务

订阅是独立 HTTP Server，`subEnable=true` 时启动；可配置 HTTP/HTTPS、监听地址、域名、路径和标题。主要路由：

```text
/<subPath>/<subid>
/<subJsonPath>/<subid>
```

普通订阅：

- 从 `settings.clients[]` 中筛选 `enable=true`、`subId=subid`；
- 仅支持 VMess、VLESS、Trojan、Shadowsocks 链接生成；
- 处理 TLS/Reality、WebSocket、gRPC、HTTPUpgrade、XHTTP/外部代理参数；
- 返回纯文本或 Base64。

JSON 订阅：

- 生成客户端配置数组；
- 支持预置 fragment、noise、mux、rules 和默认 outbounds；
- 同样聚合 Client 流量并返回 `Subscription-Userinfo`。

响应头包括 `Subscription-Userinfo`、`Profile-Update-Interval`、`Profile-Title`。订阅接口不经过面板登录 Session，`subid` 本身相当于访问凭证；泄露后可能暴露配置、流量和有效期信息。应使用高熵 token、HTTPS、限流和可撤销机制。

证据：`sub/sub.go:40-209`、`sub/subController.go:51-125`、`sub/subService.go:37-121`、`sub/subJsonService.go`。

## 11. 状态、日志、升级和系统运维

### 11.1 状态

`/server/status` 聚合：CPU/核心/频率、内存、Swap、磁盘、系统 uptime、Load Average、TCP/UDP 数、网卡累计字节和速率、Xray 状态/版本/错误、Go goroutine/内存/uptime，以及通过 ipify 等外部服务探测到的 IPv4/IPv6。

外部 IP 是节点/进程级观测，不能代表每个用户或每条出站的实际公网 IP。

### 11.2 日志

支持面板日志和 Xray access/error 日志查询，可按数量、级别、过滤项、direct/blocked/proxy 分类。日志是诊断来源，不是可靠流量账本。

### 11.3 Xray 升级和 Geo 文件

面板通过 GitHub Releases API 获取版本，下载 ZIP、解压并替换二进制后重启；可更新 geoip/geosite 等文件。审计代码未形成完整的签名或 SHA256 供应链校验，生产应在反向代理/发布层补充固定来源和完整性验证。

### 11.4 数据库备份恢复

- `getDb` 先 checkpoint，再直接返回完整 SQLite 文件。
- `importDB` 检查 SQLite 文件头，写临时文件，验证迁移，停止 Xray，备份旧库，替换新库，再重启。

这是整库导入导出，不是多节点同步协议；中心系统不应以下载 SQLite 作为长期同步方式。

### 11.5 高权限系统命令

Web API 可触发 Xray 安装/更新、Geo 下载、Subconverter 安装、UFW 放行端口、面板重启和密钥生成。这些操作需要最小权限、审计、CSRF 防护和网络边界。

证据：`web/service/server.go:318-509`、`628-757`、`789-929`、`1111-1391`。

## 12. Telegram Bot（可选）

启用后可提供：登录通知、流量日报、CPU 告警、Client/Inbound/来源 IP 查询、Client 增删改、备份发送、在线和到期通知，以及部分面板操作。Bot 使用的权限和节点范围同样不构成中央 RBAC；接入中央系统时应把它视为本机运维适配器。

## 13. 对当前项目的能力映射

| 当前需求 | X-Panel 原生状态 | 实现结论 |
|---|---|---|
| 节点状态、Xray 版本 | `/server/status`、`getXrayVersion` | Agent 可读取并上报 |
| Inbound/Client 同步 | `/inbounds/list/get`、生成配置 | 可行；需版本化 Session adapter |
| Inbound/Client 写操作 | CRUD + 本机 Handler 热加载/重启 | 可行；中央必须严格授权和审计 |
| 用户数 | 面板可读 Client，Xray 可读 Inbound users | 当前业务按 Inbound 计数，中央去重 |
| 落地机多个出口 IP | 面板无节点/出口资产模型 | 中央独立 `exit_ips` 表 |
| 线路机 -> 出口 IP 映射 | 面板无此概念 | 中央人工维护映射及有效期 |
| 出口 IP 共享用户数 | 面板不能计算跨节点指标 | 中央按映射线路机的有效 Inbound `COUNT(DISTINCT inbound_id)` |
| Client 流量 | 本地十秒 reset 区间和累计字段 | 可观测；严格对账需独占采集器/WAL |
| Inbound/Outbound 流量 | 本地可累计，外部通知不含 Outbound | 可选，必须定义唯一计费边界 |
| 在线用户 | `/onlines` 是正增量采样 | 仅展示“最近有流量”，不要称严格在线 |
| 来源 IP/设备限制 | access.log + 内存 TTL | 可做安全/设备功能，不是出口 IP |
| 精确公网出口 IP | 节点级外部 IP 探测 | 不能推导每用户实际出口 |
| 订阅 | 普通/JSON 公开 token | 可复用，但需 token 安全和撤销 |
| 财务/成本/利润 | 无模型 | 中央系统实现 |

## 14. 推荐中央接入方式

### 14.1 Agent 只做节点适配

Agent 部署在每台节点本机：

1. 使用 X-Panel Session API 读取 Inbound/Client/状态/配置。
2. 使用本机 Xray gRPC 做能力探测和可选统计读取。
3. 将远端 ID、tag、email、配置 hash、版本和采样时间标准化上报。
4. 中央写入 `nodes`、`inbounds`、`clients`、`sync_runs`，远端删除先标记 `missing_since`。
5. 不把人工映射随节点同步覆盖或删除。

Agent 到中央使用 mTLS 或节点专属签名 Token；中央不要直接保存面板管理员密码并向公网暴露 Xray gRPC。

### 14.2 共享用户数口径

中央维护：

```text
nodes(node_type = relay|landing)
exit_ips(landing_node_id, ip, family, enabled)
relay_exit_ip_mappings(relay_node_id, exit_ip_id, valid_from, valid_to, enabled)
inbounds(node_id, remote_inbound_id, enable, expiry_at, deleted_at)
```

统计某个出口 IP：

```sql
SELECT COUNT(DISTINCT i.id) AS shared_user_count
FROM relay_exit_ip_mappings m
JOIN nodes n ON n.id = m.relay_node_id AND n.node_type = 'relay'
JOIN inbounds i ON i.node_id = n.id
WHERE m.exit_ip_id = :exit_ip_id
  AND m.enabled = TRUE
  AND i.enabled = TRUE
  AND i.deleted_at IS NULL
  AND (i.expiry_at IS NULL OR i.expiry_at > :as_of);
```

一个 Inbound 下多个 Client 仍只算一个业务用户；一台线路机映射多个出口 IP 时，该 Inbound 在每个有效 IP 下各计一次。这是配置归属指标，不代表此刻在线或发生了流量。

### 14.3 流量模式选择

- **兼容模式**：接收 X-Panel `externalTrafficInformURI`，接受不可重放和缺口，不能用作严格对账。
- **独占模式**：停用/替换面板的 `QueryStats(reset=true)` 任务，由 Agent 独占 `reset=false` 快照并计算 delta；需要改造部署和本地 WAL。

同一节点不能让 X-Panel 和中央 Agent 同时 reset 同一组 Xray counter。

## 15. 版本和实现风险

- X-Panel 编译依赖 Xray `v1.260327.0`，与本地最新 Xray clone 的 protobuf/API 可能不同；Agent 必须能力探测。
- Client 在 JSON 中，单个 Client 更新存在并发覆盖 settings 的风险，中央同步应使用版本号/hash 或串行锁。
- `email` 全局唯一与多节点相同 email 冲突，不能作为中央主键。
- Xray 运行时变更不自动写回 SQLite，重启后会丢失；写操作必须明确持久化所有者。
- 外部流量通知和 `/onlines` 都不是可靠事件流。
- `getConfigJson` 可能返回敏感凭据，中央只读同步要做字段脱敏和密钥存储。

## 16. 关键源码索引

| 主题 | 文件 |
|---|---|
| 入口和信号 | `main.go` |
| 路径/环境配置 | `config/config.go` |
| SQLite/迁移 | `database/db.go` |
| 数据模型 | `database/model/model.go`、`xray/client_traffic.go` |
| Web 启动/任务 | `web/web.go` |
| 登录/Session | `web/controller/index.go`、`web/controller/base.go`、`web/session/session.go` |
| API 路由 | `web/controller/api.go`、`inbound.go`、`server.go`、`setting.go`、`xray_setting.go` |
| 配置生成/重启 | `web/service/xray.go`、`xray/process.go` |
| Xray gRPC | `xray/api.go` |
| 流量任务 | `web/job/xray_traffic_job.go`、`web/service/inbound.go`、`outbound.go` |
| 来源 IP/设备 | `web/job/check_client_ip_job.go` |
| 状态/备份/升级 | `web/service/server.go` |
| 订阅 | `sub/sub.go`、`sub/subController.go`、`sub/subService.go`、`sub/subJsonService.go` |
| Telegram | `web/service/tgbot.go` |

## 17. 最终判断

X-Panel 免费版可以作为每台线路机/落地机的本地控制和数据适配层，但不能直接变成中央多节点运营平台。当前最小可行方案是：中央维护节点、多个出口 IP 和人工线路机映射；Agent 只读同步 X-Panel 的 Inbound/Client 和健康状态；统计页按有效映射下唯一 Inbound 数显示共享用户数。流量、在线、路径、财务和高权限写操作应作为后续模块，以独立数据模型、版本兼容层、认证和审计补齐。
