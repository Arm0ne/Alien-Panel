# Xray-core 完整能力说明

## 0. 文档定位与审计声明

本文档描述本地拉取的 Xray-core 源码、配置模型、运行机制和原生 API 能力，并标明对当前多节点运营项目的可用边界。

**已固定的审计版本**

- 源码目录：`vendor-src/xray-core`
- Git commit：`c1958dba04ba065cd82a05b65bfe877e2323f0cc`
- 提交时间：2026-08-27
- `core/core.go` 中的代码版本：`26.7.28`
- 许可证：Mozilla Public License 2.0

审计覆盖核心启动链路、配置解析、`app` 功能、代理协议、网络传输、统计、路由、Commander gRPC 服务和默认发行版注册清单。本文不是对 Go 标准库及全部第三方依赖的逐行审计；实际部署仍要对二进制构建选项、操作系统权限和安装包做验收。

## 1. 结论摘要

Xray-core 是一个可组合的代理运行时，不是带业务数据库的管理平台。它把入站连接解析成会话，经 Dispatcher 和 Router 选择出站，再由传输层和出站协议转发。配置、计数器和运行时 Handler 都在单个进程内存中，重启后必须重新加载配置。

对当前项目，原生能力可直接支撑：

- 入站、出站、协议和传输配置；
- 按用户、入站、出站的流量计数器；
- 用户在线会话数和连接来源 IP（开启对应 Policy 后）；
- 路由规则、Balancer、动态路由 API；
- 出站健康探测和系统/运行时指标；
- 运行时增删入站、出站、用户；
- 本机 Unix socket 或 TCP 的 Commander gRPC 管理入口。

原生能力不能直接保证：

- 用户到最终公网出口 IP 的历史字节级流量；
- 每分钟连接明细或可靠事件队列；
- 多节点身份、人工出口 IP 资产、财务、RBAC、审计和分页报表；
- API 的公网安全认证边界。

因此，线路机到落地机出口 IP 的人工映射、共享用户数和财务事实必须由中央业务系统维护；Xray 仅作为节点运行时和采集数据源。

## 2. 总体架构

```text
客户端
  |
  v
Inbound Receiver -> 入站协议认证/解码 -> Dispatcher
                                      |
                                      +-> Sniff / DNS / Policy
                                      |
                                      v
                                  Router
                              (规则 / Balancer)
                                      |
                                      v
Outbound Sender -> 出站协议/传输 -> 目标网络

Stats counters / OnlineMap / Logs / Observatory / Metrics
                                      |
                                Commander gRPC
```

Xray 的核心对象由 `core.Instance` 持有。`core.New` 先创建 `app` 中声明的 Feature，再补齐 DNS、Policy、Router、Stats 等默认 Feature，最后创建所有 Inbound/Outbound Handler；`Start` 按 Feature 启动，`Close` 逐一关闭。

## 3. 源码目录职责

| 目录 | 职责 | 关键入口 |
|---|---|---|
| `core/` | Instance、配置对象、Feature 依赖解析、生命周期 | `core/xray.go`、`core/config.proto` |
| `main/` | CLI、配置文件加载、JSON/YAML/TOML/Protobuf 注册 | `main/main.go`、`main/run.go` |
| `infra/conf/` | 面向用户的 JSON/YAML/TOML 到 protobuf 构建器 | `infra/conf/xray.go`、`router.go` |
| `app/dispatcher/` | 会话分发、嗅探、用户/系统统计包装 | `app/dispatcher/default.go` |
| `app/proxyman/` | 入站/出站 Handler Manager 和 Receiver/Sender 设置 | `app/proxyman/inbound`、`outbound` |
| `app/router/` | 规则匹配、Balancer、路由 Webhook、动态路由 | `app/router/router.go`、`webhook.go` |
| `app/stats/` | Counter、OnlineMap、统计 Channel、StatsService | `app/stats/stats.go`、`command/` |
| `app/policy/` | 超时、缓冲、用户/系统统计开关 | `app/policy/config.proto` |
| `app/dns/` | 多上游 DNS、缓存、FakeDNS、并行/回退查询 | `app/dns/` |
| `app/observatory/` | 出站健康探测和 RTT 统计 | `app/observatory/` |
| `app/metrics/` | `/debug/vars`、pprof、统计和 Observatory HTTP 输出 | `app/metrics/metrics.go` |
| `app/log/` | 访问日志、错误日志、DNS 日志、动态重启 Logger | `app/log/` |
| `app/reverse/` | Bridge/Portal 反向隧道 | `app/reverse/` |
| `proxy/` | 具体入站/出站协议实现 | `proxy/vless`、`vmess` 等 |
| `transport/internet/` | TCP、UDP、TLS、REALITY、WS、gRPC、XHTTP 等传输 | `transport/internet/` |
| `common/` | 地址、缓冲、序列化、GeoData、Session、系统工具 | `common/` |
| `features/` | Feature 接口和跨模块契约 | `features/` |

## 4. 启动与配置加载

### 4.1 CLI 行为

`main/run.go` 支持：

- `xray run`（无参数时自动进入）；
- `-config`/`-c` 指定一个或多个配置；
- `-confdir` 读取目录中的配置文件；
- `-format=auto|json|yaml|toml|protobuf`；
- `-test` 只校验配置，不启动服务；
- `-dump` 输出合并后的配置，不启动服务。

配置查找顺序是显式参数、配置目录、当前工作目录下的 `config.json/jsonc/toml/yaml/yml`、平台默认路径，最后回退到标准输入。外部配置加载器还支持本地文件、`stdin:`、HTTP(S) GET 和 Unix socket HTTP。

证据：`main/run.go:25-229`、`main/confloader/external/external.go:20-95`。

### 4.2 格式和合并规则

- JSON、YAML、TOML 支持多文件；第一个文件作为基础，后续文件按字段覆盖，同 tag 的 Inbound/Outbound 更新，未匹配项追加。
- Protobuf 只能加载单文件。
- JSON loader 通过 `protocol` 选择具体代理类型，通过 `settings` 反序列化协议配置。
- 未使用 `main/distro/all` 或自定义 blank import 注册实现时，配置可能能解析但运行时找不到对应类型。

证据：`core/config.go:109-190`、`main/json/json.go`、`main/yaml/yaml.go`、`main/toml/toml.go`。

### 4.3 根配置模型

```json
{
  "inbounds": [],
  "outbounds": [],
  "routing": {},
  "dns": {},
  "policy": {},
  "api": {},
  "stats": {},
  "log": {},
  "metrics": {},
  "observatory": {},
  "reverse": {},
  "app": [],
  "extension": []
}
```

protobuf 根对象 `core.Config` 的必备业务部分是 `inbound`、`outbound`、`app`、`extension`；第一个 Outbound 是默认出站。每个 Inbound/Outbound 都有唯一 tag（空 tag 可表示未命名 Handler），通过 TypedMessage 组合 Receiver/Proxy 或 Sender/Proxy 配置。

证据：`core/config.proto:11-54`、`infra/conf/xray.go:376-414`。

## 5. 协议与传输能力

### 5.1 JSON loader 识别的入站协议

`infra/conf/xray.go:21-35` 识别：

- `vless`、`vmess`、`trojan`；
- `shadowsocks`、`socks`、`http`、`mixed`；
- `tunnel`/`dokodemo-door`；
- `wireguard`、`hysteria`、`tun`。

常见认证和行为包括 UUID/用户级别、密码、Fallback、嗅探、原始目的地址、TCP/UDP 开关、端口列表和 Unix socket。VLESS、VMess、Trojan、Shadowsocks 等协议都可以为用户注册统计所需的 email，但具体协议是否实现 UserManager 需要以对应实现为准。

### 5.2 JSON loader 识别的出站协议

`infra/conf/xray.go:36-52` 识别：

- `freedom`/`direct`、`blackhole`/`block`、`loopback`；
- `http`、`socks`、`shadowsocks`、`vless`、`vmess`、`trojan`；
- `hysteria`、`dns`、`wireguard`。

Shadowsocks 2022、Hysteria、TUN 等实现存在于源码，但是否被当前二进制注册取决于构建入口和平台；不能只凭目录存在就承诺所有免费版节点都支持。

### 5.3 网络传输和 Socket

`transport/internet/config.proto` 提供：

- TCP、UDP、WebSocket、gRPC、HTTPUpgrade、KCP、REALITY、XHTTP/SplitHTTP、TLS；
- TLS 证书、SNI、ALPN、会话恢复、uTLS 指纹、ECH、证书 pinning；
- REALITY server/client key、SNI、short ID、spider、Fallback 限速；
- `sendThrough`/`via`、CIDR、IPv4/IPv6 DomainStrategy、TFO、TProxy/Redirect、SO_MARK、指定 interface、dialerProxy、MPTCP、Happy Eyeballs；
- Mux、XUDP、Socket 自定义选项和 PROXY protocol；
- XHTTP/SplitHTTP 的 padding、SSE/gRPC 头、分块、连接复用和下载设置。

旧的 HTTP/QUIC legacy 语义不能按当前配置结构默认承诺；应使用当前版本的 gRPC、WebSocket、HTTPUpgrade、XHTTP 或 UDP/QUIC 相关实现并做实机测试。

### 5.4 平台要求

- Linux TUN/TProxy/Redirect、SO_MARK 和自动路由需要内核支持、`/dev/net/tun` 及相应 root/CAP_NET_ADMIN 权限。
- Windows TUN 依赖 Wintun；其他平台可能返回 unsupported。
- WireGuard 可使用内核或 gVisor 路径，实际能力由平台和构建决定。
- 证书、GeoIP/GeoSite 文件、Unix socket、监听端口和系统路由都需要部署前检查。

## 6. 请求处理工作原理

1. Receiver 在 IP/Unix socket 上接受 TCP 或 UDP，并按 StreamConfig 解密/解包。
2. 代理协议验证用户，建立 `session.Inbound`，包括 Inbound tag、用户 email、用户 level、源地址和目标地址。
3. Dispatcher 可执行 HTTP/TLS/QUIC/FakeDNS 嗅探，决定是否覆盖目标域名；也可将链路包装为用户、Inbound、Outbound 统计 Writer。
4. DNS Client 按域名规则、静态 hosts、上游列表、缓存、fallback 和 query strategy 解析。
5. Router 按规则顺序匹配 domain/IP/port/network/source/user/inbound/protocol/attributes/process/local 等字段；首个匹配结果指向 Outbound tag 或 Balancer。
6. Outbound Manager 创建或复用出站 Handler，经 sender socket、代理链、Mux 和传输层连接目标。
7. 正向和反向链路关闭时，Policy 超时、OnlineMap 引用计数和统计 Writer 更新状态。

`loopback` 可把出站重新送回指定 Inbound；`dokodemo-door` 可配合 Redirect/TProxy/原始目的地址做透明代理；`reverse` 用于 Bridge/Portal 反向连接，不是流量账单系统。

## 7. Policy、统计和在线语义

### 7.1 Policy 开关

```json
{
  "policy": {
    "levels": {
      "0": {
        "statsUserUplink": true,
        "statsUserDownlink": true,
        "statsUserOnline": true
      }
    },
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true,
      "statsOutboundUplink": true,
      "statsOutboundDownlink": true
    }
  }
}
```

代码层 protobuf 字段是 `Policy.Stats` 和 `SystemPolicy.Stats`。默认 Policy 的统计和在线开关关闭；只有对应 level 和 system policy 打开时，计数器才会产生。

### 7.2 Counter

计数器是进程内原子累加值，常见名称：

```text
user>>><email>>>traffic>>>uplink
user>>><email>>>traffic>>>downlink
inbound>>><tag>>>traffic>>>uplink
inbound>>><tag>>>traffic>>>downlink
outbound>>><tag>>>traffic>>>uplink
outbound>>><tag>>>traffic>>>downlink
```

`GetStats` 读取单个计数器，`QueryStats` 按字符串包含关系筛选多个计数器。`reset=true` 返回旧值并将计数器置零；读取和清零不是可重放日志，进程重启也会丢失内存计数器。

证据：`app/stats/command/command.proto:9-93`、`app/stats/command/command.go:30-205`、`app/stats/counter.go`、`app/dispatcher/default.go:140-229`。

### 7.3 OnlineMap

OnlineMap 是按来源 IP 的引用计数表：连接建立时 `AddIP`，连接上下文结束时 `RemoveIP`；同一个 IP 的多个连接只保留一个条目，并更新 Unix `last_seen`。回环来源 `127.0.0.1` 和 `[::1]` 被忽略。

原生在线 RPC：

| RPC | 返回 | 边界 |
|---|---|---|
| `GetStatsOnline` | 某个 `user>>>email>>>online` map 的 IP 数 | 必须传完整 map 名称 |
| `GetStatsOnlineIpList` | IP 到 `last_seen` | 是连接来源 IP，不是出口 IP |
| `GetAllOnlineUsers` | 当前 Count 大于 0 的 map 名称 | 不返回流量或连接详情 |
| `GetUsersStats` | 最新源码按 email 汇总在线 IP/last_seen，可选流量 | X-Panel 依赖的旧 Xray 版本可能没有 |

`GetUsersStats` 是本地最新 commit 新增的 RPC，启动时必须做能力探测并在不支持时降级。

## 8. 原生 Commander gRPC API

### 8.1 开启方式

```json
{
  "api": {
    "tag": "api",
    "listen": "127.0.0.1:62789",
    "services": [
      "ReflectionService",
      "HandlerService",
      "StatsService",
      "RoutingService",
      "ObservatoryService",
      "LoggerService"
    ]
  }
}
```

`infra/conf/api.go` 将 service 名称转换为 TypedMessage。Commander 使用裸 `grpc.NewServer()`，监听地址可以是 TCP 或 Unix socket；未设置 `listen` 时会通过一个出站 Handler 提供内部 gRPC 通道。默认没有 TLS、认证、ACL 或 interceptor。

**安全要求**：不要把 Xray gRPC 直接暴露到公网。节点 Agent 应在本机访问，并通过 mTLS、短期节点 Token、反向隧道或防火墙向中央服务上报。

### 8.2 HandlerService

定义于 `app/proxyman/command/command.proto`：

| RPC/Operation | 能力 |
|---|---|
| `ListInbounds` | 列出 Inbound tag，或返回 Receiver/Proxy 配置 |
| `AddInbound` / `RemoveInbound` | 运行时增加/删除 Inbound Handler |
| `AlterInbound` | 对支持 UserManager 的 Inbound 执行 Typed Operation |
| `GetInboundUsers` | 查询一个 Inbound 的全部用户或指定 email |
| `GetInboundUsersCount` | 查询 Inbound 用户数 |
| `ListOutbounds` | 列出出站配置，忽略 Commander 自身出站 |
| `AddOutbound` / `RemoveOutbound` | 运行时增加/删除 Outbound Handler |
| `AlterOutbound` | 调用已注册的 OutboundOperation；当前默认实现通常没有可用 operation |
| `AddUserOperation` | 向 Inbound UserManager 添加用户 |
| `RemoveUserOperation` | 按 email 从 Inbound UserManager 删除用户 |

运行时修改只改变 Xray 内存。它不会自动写回 X-Panel SQLite 或配置文件；在面板管理节点上应先写面板事实，再让面板生成配置/重启，或者由 Agent 负责双向一致性。删除 Outbound 的 Manager 实现还需关注 Handler 资源关闭和自定义实现差异。

### 8.3 StatsService

| RPC | 请求关键字段 | 返回/语义 |
|---|---|---|
| `GetStats` | `name`, `reset` | 单计数器 |
| `QueryStats` | `pattern`, `reset` | 包含匹配的计数器列表 |
| `GetSysStats` | 无 | Go runtime 内存、GC、goroutine、uptime |
| `GetStatsOnline` | `name` | 某在线 map 的 IP 数 |
| `GetStatsOnlineIpList` | `name` | IP/Unix last_seen |
| `GetAllOnlineUsers` | 无 | 当前在线 map 名称 |
| `GetUsersStats` | `include_traffic`, `reset` | 新版本按用户汇总在线 IP 和可选流量 |

### 8.4 RoutingService

| RPC | 能力 |
|---|---|
| `TestRoute` | 用给定 RoutingContext 测试匹配结果，可选发布到统计 Channel |
| `SubscribeRoutingStats` | 订阅路由事件流；未启用 Channel 时返回未启用 |
| `GetBalancerInfo` | 查询 override 和原则目标 |
| `OverrideBalancerTarget` | 临时覆盖 Balancer 目标 |
| `AddRule` / `RemoveRule` / `ListRule` | 动态增删和列出路由规则 |

当前 `app/router/command` 注册时传入的 routing stats channel 为 `nil`，因此不能假设 `SubscribeRoutingStats` 在默认发行版配置中可用。

### 8.5 ObservatoryService、LoggerService、ReflectionService

- `ObservatoryService.GetOutboundStatus` 返回出站 alive、delay、last error、last seen 等健康结果。
- `LoggerService.RestartLogger` 关闭并重新启动 Logger。
- `ReflectionService` 暴露 gRPC reflection，便于 Agent 做服务发现；reflection 本身不提供认证。

### 8.6 官方 CLI API 包装

`main/commands/all/api` 对同一套 gRPC 服务提供命令行入口，常用命令包括：

| 命令 | 对应能力 |
|---|---|
| `xray api lsi` / `rmi` / `adi` | 列出、删除、增加 Inbound |
| `xray api lso` / `rmo` / `ado` | 列出、删除、增加 Outbound |
| `xray api inbounduser` / `inboundusercount` | 查询用户和用户数 |
| `xray api adu` / `rmu` | 增删 Inbound 用户 |
| `xray api stats` / `statsquery` / `statssys` | 读取单个、批量和系统统计 |
| `xray api statsonline` / `statsonlineiplist` / `statsgetallonlineusers` | 查询在线数、来源 IP 和在线用户 |
| `xray api lsrules` / `adrules` / `rmrules` | 列出、增加、删除路由规则 |
| `xray api bi` / `bo` | 查询或覆盖 Balancer 目标 |
| `xray api restartlogger` | 重启 Logger |
| `xray api sib` | 按来源 IP 动态增加阻断规则 |

CLI 默认用 `127.0.0.1:8080` 和 insecure gRPC；应通过 `--server`、隧道或本机执行，不能把 CLI 默认值误解为远程安全 API。

## 9. Routing Webhook 与路由能力

路由规则可按以下条件组合：domain、IP、端口、network、source IP/port、user email、Inbound tag、protocol、attributes、local IP/port、process、local OS、VLESS route。字段在同一规则内是 AND 关系，规则按列表顺序匹配；Balancer 支持 random、least load、least ping、round robin 和 fallback。

`routing.rules[].webhook` 是路由事件通知，不是流量回调。JSON 事件可能包含：

```json
{
  "email": "user@example",
  "protocol": "vless",
  "network": "tcp",
  "source": "198.51.100.2:1234",
  "destination": "example.com:443",
  "originalTarget": "...",
  "routeTarget": "...",
  "inboundTag": "relay-in",
  "outboundTag": "direct",
  "ts": 1760000000
}
```

实现支持 URL、Unix socket、静态 headers 和按 email 的去重窗口；异步 POST 超时 5 秒，失败只记录日志，没有签名、可靠队列、持久化、重试或重放。事件不含 uplink/downlink 字节数，也不含最终 NAT 后公网出口 IP，因此只能做路径辅助观测。

证据：`app/router/config.proto:14-113`、`app/router/webhook.go:22-225`。

## 10. DNS、GeoData、日志、指标和健康

### DNS

DNS Client 的 nameserver scheme 分支支持 localhost、普通 UDP、TCP、DoH（`https`/`https+local`）、h2c DoH（`h2c`/`h2c+local`）、DoQ（`quic+local`）、FakeDNS，以及多上游、缓存、stale、expected/unexpected IP、fallback、并行查询、A/AAAA 策略、EDNS client IP、静态 hosts 和按 tag 选择。源码没有独立的通用 `dot` scheme；不要把它概括为完整 DoT 接入。`proxy/dns` 还可按 QTYPE/domain 执行 Direct、Drop、Return、Hijack 和 server 重写。

### GeoData

`geoip:`、`geosite:`、CIDR、regexp、domain、full、keyword、dotless、属性过滤和反向匹配均可用于路由。`app/geodata` 可按 cron 下载和原子替换 asset；缺失、损坏或不兼容的 dat 文件会导致规则加载失败。GeoData 文件来源和校验必须纳入部署供应链。

### 日志

`app/log` 支持 None、Console、File、Event，分别配置 error/access/DNS 日志、级别和地址掩码。访问日志对在线来源 IP、用户和故障排查有价值，但不是完整连接账单；日志格式和轮转由部署层负责。

### Metrics

`app/metrics` 可监听 TCP 或挂载为 Outbound tag，提供 `/debug/vars`、pprof、inbound/outbound/user counter 和 Observatory 状态。它没有身份认证和多租户隔离，生产环境只应绑定回环/管理网并加反向代理。

### Observatory

按 `subjectSelector` 选择出站，经指定出站发 HTTPS 探测，返回 alive、delay、last error、last seen。Burst health ping 还可统计 average/min/max/deviation/fail。它证明出站可用性和 RTT，不证明某个用户实际使用了哪个公网 IP。

## 11. 特殊能力

- `freedom`：直连，支持目标解析策略、`sendThrough`、fragment/noise、PROXY protocol 等；当前构建对命名协议访问私网目标有安全限制。
- `blackhole`：丢弃连接，可返回 None 或 HTTP 响应，也可吸收 UDP。
- `loopback`：把流量重新送入指定 Inbound，可再次嗅探。
- `dokodemo-door`：透明代理、Redirect/TProxy、原始目的地址和端口映射。
- `reverse`：Bridge/Portal 通过 mux worker 维护反向连接，可 DRAIN，但不提供业务用户映射。
- `WireGuard`/`TUN`：系统级入口/出口、peer 和路由能力，受平台权限限制。

## 12. 当前项目能力映射

| 项目需求 | Xray 原生状态 | 推荐接入 |
|---|---|---|
| 节点健康/版本 | 部分原生：SysStats、Observatory；版本由二进制/Agent读取 | Agent 本机探测并上报 |
| Inbound/Outbound/用户配置 | 原生 HandlerService + 配置文件 | Agent 读取；写操作先落面板再同步 |
| 线路机 Inbound 用户数 | 可读 Handler 用户数/面板配置 | 中央按 Inbound 事实表统计 |
| 落地机多出口 IP | Xray 只有 `sendThrough/via` 配置语义 | 中央 `exit_ips` 和人工映射表 |
| 出口 IP 共享用户数 | 不是 Xray 运行时指标 | `COUNT(DISTINCT inbound_id)`，人工映射为事实 |
| 用户/Inbound/Outbound 流量 | 可提供计数器 | 选择唯一采集者，快照转增量 |
| 用户到出口 IP 精确流量 | 原生不能保证 | 固定 outbound、跨节点身份和额外埋点后再做 |
| 路由路径事件 | Routing Webhook 可提供辅助事件 | 标为 `observed`，不可作为计费凭据 |
| 分钟历史/强一致对账 | 原生没有事件队列 | Agent WAL、采样桶和缺口标记 |

## 13. 版本兼容和安全边界

X-Panel `v26.6.18` 编译依赖 Xray `v1.260327.0`，本地最新 clone 的 `GetUsersStats` 不一定存在于该依赖版本。Agent 应：

1. 读取版本和生成配置；
2. 探测 Reflection 或逐个调用 RPC；
3. 遇到 `Unimplemented`/`Unavailable` 时记录能力集合并降级；
4. 不把同一 Stats counter 的 `reset=true` 读取交给两个消费者；
5. 记录节点 Xray generation、配置 hash、采集时间和缺口。

Xray Commander、Metrics 和 Reflection 默认没有认证；公网暴露会允许远程修改代理和用户。生产部署应使用本机 Unix socket 或回环地址、主机防火墙、mTLS/隧道、最小服务集合和审计代理。

## 14. 关键源码索引

| 主题 | 文件 |
|---|---|
| CLI/启动 | `main/main.go`、`main/run.go` |
| 配置加载 | `core/config.go`、`main/json/json.go`、`main/yaml/yaml.go`、`main/toml/toml.go` |
| 根配置/生命周期 | `core/config.proto`、`core/xray.go` |
| 默认注册 | `main/distro/all/all.go` |
| JSON 协议映射 | `infra/conf/xray.go` |
| 路由解析 | `infra/conf/router.go`、`app/router/router.go` |
| Dispatcher/统计包装 | `app/dispatcher/default.go` |
| 入站/出站 Manager | `app/proxyman/inbound/inbound.go`、`app/proxyman/outbound/outbound.go` |
| Commander | `app/commander/commander.go`、`app/commander/service.go` |
| Handler API | `app/proxyman/command/command.proto`、`command.go` |
| Stats API | `app/stats/command/command.proto`、`command.go` |
| 在线语义 | `app/stats/online_map.go` |
| 路由 API/Webhook | `app/router/command/command.proto`、`command.go`、`app/router/webhook.go` |
| Policy | `app/policy/config.proto`、`app/policy/config.go` |
| DNS/GeoData | `app/dns/config.proto`、`app/dns/dns.go`、`app/geodata/geodata.go` |
| Observatory/Metrics | `app/observatory/`、`app/metrics/metrics.go` |
| 平台 TUN | `proxy/tun/` |

## 15. 最终判断

Xray-core 原生 API 足以作为当前项目的节点运行时数据源，但它不是中央运营系统。对于已经确认的“线路机 -> 落地机出口 IP”人工映射和出口共享用户数，最稳妥的实现是不改 Xray：中央系统维护映射，Agent 同步线路机 Inbound，按有效期和启用状态计算唯一 Inbound 数；流量、在线和路由数据作为独立的可选观测事实，并明确其来源和置信度。
