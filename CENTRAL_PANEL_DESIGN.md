# 多节点 X-Panel 轻量中央管理面板设计方案

## 1. 文档信息

| 项目 | 内容 |
|---|---|
| 项目名称 | 多节点 X-Panel 轻量中央管理面板 |
| 文档版本 | v1.0 |
| 文档状态 | 开发基线 |
| 目标用户 | 跨境网络服务商的运营和运维人员 |
| 适用节点 | 线路机、落地机 |
| 核心约束 | 不修改 X-Panel/Xray 源码；不接管 Xray 流量计数器 |
| 第一版定位 | 只读同步、用户到期管理、线路关系、成本和收入汇总 |
| 前端基线 | Soybean Admin（仅保留项目需要的工程和 UI 能力） |
| 前端参考源码 | https://github.com/soybeanjs/soybean-admin |

## 2. 项目背景与目标

当前每台线路机和落地机都独立安装 X-Panel 与 Xray。线路机使用 Reality 入站，通常连接落地机的 SS 入站；部分线路也可直接使用线路机公网 IP，另有独立购买的 S5 出口。落地机仍可通过 Xray 路由和多个出站提供不同公网出口 IP。单台机器内的配置和查看并不复杂，但节点数量增长后，逐台登录面板检查用户、到期时间、流量和成本会产生明显的运营负担。

本项目建设一个轻量中央面板，集中展示和汇总现有节点数据，同时保留 X-Panel 作为每台服务器的本地配置管理者。

项目目标：

- 集中查看所有线路机用户；
- 按一个 Inbound 对应一个业务用户统计；
- 一个用户的多个 Email/Client 只作为多个设备凭证；
- 直接读取 X-Panel 已累计的 Inbound 流量，不读取或清零 Xray Counter；
- 查看用户到期时间、启用状态、客户端数量和最近活动；
- 维护线路机、落地机、落地机 SS 入站以及节点/外部出口 IP 的关系；
- 录入服务器、带宽和出口 IP 成本；
- 录入用户月费并计算收入、成本和预计毛利润；
- 支持节点离线、接口异常、流量重置和本地配置漂移提示；
- 后续可以增加中央写操作，但第一版不改变节点配置；
- X-Panel 或 Xray 独立升级时，中央历史数据和 Agent 不受影响。

## 3. 设计结论

系统采用“中央面板 + 节点 Agent + X-Panel API”的结构：

~~~text
                    中央管理面板
       Web UI / HTTP API / 业务数据库 / 同步任务
                              ^
                              | HTTPS + 节点专属 Token
                              |
              +---------------+---------------+
              |                               |
        线路机 Agent                      落地机 Agent
              |                               |
       本机 X-Panel API                 本机 X-Panel API
              |                               |
       Reality 入站/SS 出站             SS 入站/路由/出口 IP
~~~

核心责任划分：

| 领域 | 权威来源 | 中央处理方式 |
|---|---|---|
| Inbound 协议、Tag、端口 | X-Panel | 同步快照 |
| Client/Email | X-Panel | 同步设备信息 |
| Inbound 流量 | X-Panel SQLite 累计字段 | 周期性读取快照 |
| Xray/系统状态 | X-Panel Server API | 周期性读取 |
| 业务用户名称 | 中央系统 | 中央维护 |
| 用户月费 | 中央系统 | 中央维护 |
| 节点类型 | 中央系统 | 中央维护 |
| 线路关系 | 中央系统 | 人工维护 |
| 出口 IP 资产 | 中央系统 | 人工维护 |
| 节点/IP 成本 | 中央系统 | 人工维护 |

## 4. 范围与非范围

### 4.1 第一版范围

- 节点注册、启用/停用和心跳；
- Agent 登录 X-Panel 并同步 Inbound、Client、状态和累计流量；
- 线路机/落地机分类；
- 一个 Inbound 绑定一个业务用户；
- 用户到期状态和筛选；
- Client 数量和设备详情；
- 线路机到落地机关系；
- 线路机/落地机节点公网出口和独立 S5 出口资产；
- 节点、出口 IP 和用户月费录入；
- 总览、用户、节点、线路与成本四类页面；
- 日/月流量快照和趋势；
- 节点离线、同步失败、流量重置提示；
- 管理员登录、Token、审计日志和数据库备份。
- 基于 Soybean Admin 的 Vue 3 前端壳、布局、路由、主题和后台交互组件。

### 4.2 第一版明确不做

- 修改 X-Panel/Xray 源码；
- 中央直接调用 Xray 的 QueryStats(reset=true)；
- Agent 独占 Xray 流量采集；
- 中央编辑完整 Xray JSON；
- 中央自动推导 Xray 路由规则；
- 精确记录每个用户实际使用的出口 IP；
- 精确按出口 IP 归因单个用户流量；
- 支付网关、自动对账和正式财务系统；
- 复杂多租户 RBAC 和审批流；
- Telegram、多渠道消息和订阅平台重构。
- Soybean Admin 中与本项目无关的示例页面、Mock 数据、组件演示、模板 Showcase、示例权限页面和演示业务。

## 5. 业务口径

### 5.1 用户口径

~~~text
一个线路机上的业务 Inbound = 一个业务用户
一个 Inbound 下的 Client/Email = 一个设备凭证
~~~

同一用户在手机、电脑和平板上有三个 Client 时，中央仍只统计一个用户，设备数量显示为 3。

已登记为 `relay` 的线路机完成成功快照后，中央会为每个新发现的业务 Inbound 自动创建一个业务用户和一条一对一关联。初始名称依次取 X-Panel `remark`、`tag`、`Inbound <remote_id>`；后续同步只更新 X-Panel 权威字段（启用状态、到期时间、流量和 Client），不会覆盖中央维护的名称、月费、币种和备注。

X-Panel 的 email 在本地数据库中具有唯一性，但中央不使用 email 作为主键。中央远程资源唯一键为：

~~~text
node_id + remote_inbound_id
node_id + remote_inbound_id + remote_client_id
~~~

### 5.2 Inbound 分类

中央将同步发现的 Inbound 分为：

- user：线路机上的业务用户 Inbound；
- infrastructure：落地机 SS 入站、管理入站、API 入站等基础设施；
- unknown：新发现但尚未人工确认的 Inbound。

线路机的新 Inbound 默认归为 `user` 并建立业务用户；落地机 Inbound 按节点类型归为基础设施资源，不会建立业务用户。API 同时返回 `purpose=business|infrastructure`，前端不得把落地节点 Inbound 放入业务用户列表。已标记为 `infrastructure` 的线路机 Inbound 也不会被后续同步自动创建成用户。

第一版不自动从 Xray 路由规则推断用户归属。

### 5.3 用户流量口径

第一版用户流量以线路机 Reality 入站的 X-Panel Inbound 累计值为准：

~~~text
用户总流量 = Inbound.up + Inbound.down
~~~

X-Panel 内部约每 10 秒从 Xray 读取并清零 Counter，再累计写入 SQLite。中央只读 X-Panel API 返回的累计字段，不清零 Counter，因此不会与 X-Panel 任务争抢统计。

### 5.4 在线口径

在线展示分为：

- 节点在线：Agent 最近心跳成功；
- 最近有流量：最近同步周期内累计值增加；
- 最近活跃：使用 X-Panel 的 last_online 或最近增量时间；
- 节点离线：超过离线阈值没有心跳。

X-Panel 的 onlines 基于采样增量，不是严格实时连接列表。页面使用“最近有流量”或“最近活跃”，不使用“实时连接列表”。

### 5.5 出口 IP 口径

中央可统计“配置归属用户数”：当前 `user_paths` 直接指定该出口 IP 的有效用户数；旧线路绑定统计仅用于兼容历史页面。

现有共享 SS 入站和随机 Xray 路由无法证明某个用户在某个时刻实际使用了哪一个出口 IP。因此页面必须使用“配置归属用户数”措辞，不得宣称“实际使用用户数”。

出口资产不再假设全部属于落地机，统一按来源建模：

- `source_type=node`：节点公网出口，`owner_node_id` 可指向线路机或落地机；线路绑定时用 `scope=relay` 表示线路机直出，用 `scope=landing` 表示落地机出口；
- `source_type=s5`：独立购买的 S5 出口，不属于任何 VPS，`owner_node_id` 为空，只能以 `scope=external` 绑定线路；
- `landing_node_id` 保留为旧 API/数据兼容字段，仅对落地机资产回填，不作为新的归属判断依据；
- 绑定接口按线路两端节点和来源类型做强校验，禁止跨节点或将 S5 当作节点出口绑定；
- `scope` 仍是旧线路出口池的配置维度；新用户路径不经过出口池，直接校验出口资产归属，避免把两种出口混在同一池中随机使用；
- 用户实际路径通过 `user_paths` 完成：线路机从用户主 Inbound 自动确定；线路机直出只指定线路机出口 IP，经落地机必须同时指定落地机和该节点基础设施 Inbound，独立 S5 不得与落地机混用。路径保存会关闭上一条当前路径并保留历史，配置属于中央运营配置，不会直接改写 X-Panel/Xray；
- 旧 `routes`、`user_routes`、`route_exit_ips` 数据和接口保留用于历史兼容，但不再作为新用户日常分配的必填前置步骤；用户列表直接展示线路机、落地机、出口 IP 及其归属节点。
- 中央只维护 S5 资产、节点/出口成本和用户路径（旧线路关系仅兼容），不保存明文账号密码；实际认证凭据应使用加密 Secret 或节点本地 Secret 引用。

### 5.6 财务口径

~~~text
月收入 = 有效用户月费之和
月成本 = 线路机成本 + 落地机成本 + 出口 IP 成本 + 其他成本
预计毛利润 = 月收入 - 月成本
~~~

第一版只做运营估算，不替代会计凭证。所有成本记录必须有生效日期，便于按月份查询。

## 6. 系统组件

### 6.1 中央后端

职责：

- 管理员认证和 Session；
- 接收 Agent 注册、心跳和同步快照；
- 维护业务用户、节点、线路、出口 IP、成本和收入数据；
- 计算到期状态、流量增量、收入、成本和毛利润；
- 提供前端查询 API；
- 记录同步、操作和异常事件；
- 执行定时任务。

### 6.2 节点 Agent

职责：

- 使用本机 X-Panel Session API 登录；
- 读取 Inbound、Client、流量和状态；
- 将本机版本差异转换为统一结构；
- 向中央主动发起 HTTPS 请求；
- 返回同步成功或失败结果；
- 不修改 X-Panel/Xray；
- 不调用 Xray reset 统计接口。

Agent 设计为独立静态 Go 二进制，安装、升级和回滚不影响 X-Panel/Xray。

### 6.3 中央前端

中央前端以 Soybean Admin 源码作为工程和 UI 基础，采用其 Vue 3、Vite、TypeScript、Pinia、路由布局、主题切换、请求封装和后台组件组织方式。Soybean Admin 只提供技术底座，不提供本项目的业务模型和数据。

前端必须遵循以下原则：

- 使用中央后端真实 API，不保留模板 Mock 数据；
- 保留桌面后台的侧栏、顶部栏、面包屑、主题和响应式布局；
- 页面以表格、筛选、状态标签、抽屉详情和确认弹窗为主；
- 用户、节点、线路、出口 IP、成本、同步事件是唯一业务导航；
- 不把模板示例 Dashboard、图表或权限演示直接当作产品页面；
- 不做复杂营销首页、拓扑动画或与运营无关的可视化；
- 所有业务字段、状态名称和统计口径以本设计方案为准。

推荐前端目录边界：

~~~text
src/
  layouts/                 保留后台布局和响应式壳
  router/                  只保留本项目路由
  store/                   登录、节点筛选和全局 UI 状态
  service/api/             中央后端 API client
  views/dashboard/         总览
  views/users/             用户和 Client 设备
  views/nodes/             节点
  views/routes/            线路和出口 IP
  views/finance/           成本、月费和毛利润
  views/events/            同步和异常事件
  components/              项目专用表格、状态、流量和详情组件
  constants/               状态、字段、权限和格式化规则
~~~

### 6.4 Soybean Admin 保留、改造和删除范围

#### 保留

- Vue 3、Vite、TypeScript 工程配置；
- 路由、布局、菜单、面包屑和页面容器；
- Pinia 状态管理和登录状态守卫；
- HTTP 请求封装、错误处理和请求拦截；
- 现有组件库的表格、表单、抽屉、弹窗、通知和分页；
- 主题、暗色模式、响应式布局和基础图标；
- 通用加载、空状态、错误状态和权限守卫。

#### 改造

- 菜单改为总览、用户、节点、线路与出口 IP、成本财务、同步事件；
- Dashboard 改为真实节点、用户、到期、流量和财务聚合；
- 表格列改为 Inbound 用户口径，Client 仅显示为设备；
- 请求类型和 API client 改为中央后端接口；
- 主题和颜色用于区分在线、离线、到期、异常和数据过期状态；
- 权限先保留管理员单角色结构，为后续 RBAC 留下路由元数据。

#### 删除

- 模板自带的示例 Dashboard 和假数据；
- 组件、图表、表单、表格 Showcase 页面；
- 示例权限、角色、部门、菜单管理页面；
- 示例多语言内容和无关业务字典；
- 演示用 keep-alive 页面、测试页面和占位页面；
- 与节点运营无关的营销、工作台和项目管理模块。

删除模板功能时不能删除其被项目使用的基础依赖；必须先确认没有业务页面引用，再清理路由、菜单、权限和资源文件。

前端初始化时必须记录 Soybean Admin 的具体 Git commit、Node 包管理器和依赖锁文件。后续窗口统一使用该固定版本，禁止一边开发一边无记录地拉取上游最新代码。上线前保留上游 LICENSE/版权要求，并在项目 README 中注明基于该仓库改造。

### 6.5 数据库

MVP 使用 SQLite WAL 单机部署。数据库规模和并发增长后迁移 PostgreSQL。业务表使用稳定中央 ID，远端 ID 仅作为来源字段。

## 7. 数据模型

以下为逻辑模型，实际字段可按 ORM 规范调整。

### 7.1 管理员和节点

~~~text
admin_users
- id
- username
- password_hash
- totp_secret
- enabled
- last_login_at
- created_at
- updated_at

nodes
- id
- node_key
- name
- type                 -- relay / landing
- deleted_at           -- soft-deletion timestamp; hidden from active management
- hostname
- public_ip
- region
- provider
- panel_base_path
- agent_version
- xpanel_version
- xray_version
- enabled
- health_status        -- unknown / online / offline / error
- last_seen_at
- created_at
- updated_at

node_credentials
- id
- node_id
- token_hash
- last_rotated_at
- revoked_at
- created_at
~~~

中央只保存 Node Token 的哈希。X-Panel 用户名和密码只保存在节点 Agent 本地配置中。

节点的 `public_ip` 仅表示节点主公网/管理地址；节点可拥有多个出口资产，统一记录在 `exit_ips`，通过 `owner_node_id` 归属线路机或落地机。接入向导允许一次提交多个 IPv4/IPv6 地址，创建节点与出口资产在同一事务中完成；出口资产的服务商、成本、有效期和备注在出口 IP 页面维护。

### 7.2 用户、Inbound 和 Client

~~~text
users
- id
- display_name
- status               -- active / expiring / expired / disabled / error
- monthly_fee
- currency
- expiry_time
- notes
- created_at
- updated_at

inbounds
- id
- node_id
- remote_inbound_id
- user_id
- kind                 -- user / infrastructure / unknown
- tag
- remark
- protocol
- port
- listen
- enable
- expiry_time
- up
- down
- all_time
- client_count
- config_hash
- first_seen_at
- last_seen_at
- missing_since
- deleted_at

clients
- id
- node_id
- inbound_id
- remote_client_id
- email
- enable
- expiry_time
- up
- down
- all_time
- last_online
- limit_ip
- last_seen_at

user_inbounds
- id
- user_id
- inbound_id
- is_primary
- active_from
- active_to
~~~

约束：

- 一个 Inbound 最多绑定一个业务用户；
- 第一版一个用户只有一个主 Inbound；
- 用户更换线路时保留历史关系；
- 不使用 email 作为中央用户主键。

### 7.3 线路、出口 IP 和成本

~~~text
routes
- id
- name
- relay_node_id
- landing_node_id
- relay_outbound_tag
- landing_inbound_id
- landing_inbound_tag
- enabled
- valid_from
- valid_to
- notes

user_routes
- id
- user_id
- route_id
- route_exit_ip_id (nullable; fixed route-exit binding, null means route pool)
- is_primary
- active_from
- active_to

user_paths (new primary assignment model)
- id
- user_id
- relay_node_id (copied from the user's primary Inbound)
- landing_node_id (nullable)
- landing_inbound_id (nullable)
- exit_ip_id (fixed, required)
- mode (relay / landing / external)
- notes
- active_from
- active_to

路径模式约束：`mode=relay` 不得有落地字段且出口必须属于线路机；`mode=landing` 必须同时有落地节点、落地基础设施 Inbound 且出口属于落地机；`mode=external` 只能使用独立 S5，落地字段必须为空。

exit_ips
- id
- source_type (node / s5)
- owner_node_id (nullable; relay or landing)
- landing_node_id (legacy nullable alias)
- ip
- family
- provider
- monthly_cost
- currency
- enabled
- valid_from
- valid_to
- notes

route_exit_ips
- id
- route_id
- exit_ip_id
- scope (relay / landing / external)
- allocation_weight
- enabled

node_costs
- id
- node_id
- category
- monthly_amount
- currency
- effective_from
- effective_to
- notes

other_costs
- id
- name
- category
- monthly_amount
- currency
- effective_from
- effective_to
- notes
~~~

### 7.4 流量、同步和审计

~~~text
traffic_snapshots
- id
- node_id
- inbound_id
- collected_at
- up
- down
- all_time
- source
- reset_detected
- sync_run_id

sync_runs
- id
- node_id
- sync_id
- started_at
- finished_at
- status
- inbound_count
- client_count
- error_message

node_events
- id
- node_id
- event_type
- message
- created_at

audit_logs
- id
- admin_user_id
- action
- resource_type
- resource_id
- request_id
- before_json
- after_json
- ip
- created_at
~~~

审计日志中不记录密码、Token、Reality 私钥、完整订阅链接或完整 Xray 配置。

## 8. Agent 同步协议

### 8.1 注册

~~~text
POST /agent/v1/register
~~~

请求包含 Node Key、名称、节点类型、主机名、Agent 版本和本机 X-Panel 地址。中央返回中央 Node ID、正式 Token 和同步策略。正式 Token 只显示一次。

生产环境的注册请求必须额外携带部署侧保存的 `X-Agent-Registration-Token` 引导密钥；引导密钥只用于首次注册/轮换，不作为节点同步 Bearer Token。中央只在注册响应中返回一次正式节点 Token，后续仅保存其哈希。

### 8.2 心跳

~~~text
POST /agent/v1/heartbeat
~~~

心跳上报：

- 节点 Key；
- 观测时间；
- Xray 是否运行；
- Xray 版本；
- X-Panel 版本；
- CPU、内存、磁盘；
- Agent 版本。

### 8.3 完整同步

~~~text
POST /agent/v1/sync
~~~

一次同步上传一个完整节点快照。核心字段：

~~~json
{
  "node_key": "relay-001",
  "sync_id": "relay-001-20260902T120000-abc123",
  "observed_at": "2026-09-02T12:00:00+08:00",
  "status": {
    "xray_running": true,
    "xray_version": "26.6.27",
    "xpanel_version": "26.4.25"
  },
  "inbounds": [
    {
      "remote_id": 15,
      "tag": "user-15",
      "remark": "客户 A",
      "protocol": "vless",
      "enable": true,
      "expiry_time": 1792022400,
      "up": 123456,
      "down": 987654,
      "all_time": 1111110,
      "clients": [
        {
          "remote_id": 31,
          "email": "phone",
          "enable": true,
          "all_time": 500000,
          "last_online": 1791000000
        }
      ]
    }
  ]
}
~~~

### 8.4 幂等和重试

- sync_id 全局唯一；
- 重复上传同一 sync_id 返回原同步结果；
- Agent 请求超时后按 30 秒、2 分钟、10 分钟退避重试；
- 中央按 node_key + remote_id 更新 Inbound；
- Client 按 node_key + remote_inbound_id + remote_client_id 更新；
- 采样按 inbound_id + collected_at 去重；
- 同步成功后才更新 last_seen_at；
- 部分失败时保留上一次有效快照并标记 partial。

### 8.5 缺失资源处理

一次同步未发现 Inbound 时：

1. 仅在成功的完整快照中设置 `missing_since`、增加 `missing_sync_count`；首次写入 `inbound_missing` 事件；
2. 不立即删除业务用户和历史数据；
3. 连续三次成功完整快照仍未发现时，写入 `inbound_archived` 事件并在中央库标记为 archived；
4. 同步失败、部分失败或节点离线均不增加缺失计数；重新出现时清除缺失时间、计数和中央归档标记；
5. 归档仅影响中央展示状态，不会修改 X-Panel，也不会物理删除中央的业务关联、Client 或历史流量快照。

## 9. X-Panel 采集适配

### 9.1 采集接口

Agent 使用 X-Panel 已有 HTTP API：

~~~text
POST <basePath>/login
GET  <basePath>/panel/api/inbounds/list
GET  <basePath>/panel/api/inbounds/get/:id
GET  <basePath>/panel/api/inbounds/getClientTraffics/:email
GET  <basePath>/panel/api/server/status
GET  <basePath>/panel/api/server/getXrayVersion
~~~

不同 X-Panel 版本的 Base Path、响应字段和登录行为可能不同，Agent 需要实现版本探测和兼容适配，不应假设所有节点版本完全一致。

### 9.2 登录会话

Agent 启动时登录 X-Panel，缓存 Session Cookie。遇到 401 时重新登录一次；连续失败后进入退避，避免频繁触发面板防护。日志中不得输出密码和完整 Cookie。

### 9.3 采集策略

- 首次同步读取完整 Inbound 和 Client；
- 后续同步优先使用列表接口；
- 对列表中缺少详细字段的 Inbound，再请求详情接口；
- Server Status 作为节点状态来源；
- 不下载完整 SQLite；
- 不上传完整 Xray 配置；
- 不访问 Xray gRPC，除非后续健康探测明确需要。

## 10. 流量处理

### 10.1 读取原则

中央只读取：

~~~text
Inbound.up
Inbound.down
Inbound.all_time
~~~

中央绝不调用：

~~~text
StatsService.QueryStats(reset=true)
~~~

这样 X-Panel 仍然是唯一的 Xray Counter 清零者。

### 10.2 快照和增量

每次完整同步保存 traffic_snapshots。展示累计值时使用最近快照；展示日/月流量时使用同一 Inbound 相邻快照的差值汇总。

~~~text
delta = current.all_time - previous.all_time
~~~

若 current.all_time 小于 previous.all_time：

- delta 记为 0；
- reset_detected=true；
- 写入 traffic_reset 事件；
- 开始新的统计基线；
- 不产生负流量。

### 10.3 流量重置和数据库恢复

常见原因：X-Panel 手动重置、SQLite 恢复、节点重装、Inbound 删除重建。中央保留历史快照和重置事件，报表显示“统计基线已重置”，避免把恢复后的流量和旧周期错误相加。

### 10.4 数据新鲜度

每条流量数据显示 last_collected_at：

- 超过 10 分钟未同步：数据延迟；
- 超过 30 分钟未同步：数据过期或节点离线；
- 节点恢复后：继续从最新成功快照计算。

## 11. 用户生命周期

### 11.1 第一版只读模式

第一版用户来源为线路机 X-Panel Inbound。成功同步时会自动建立一对一业务用户关联；中央可编辑：

- 业务名称；
- 月费；
- 备注；
- Inbound 类型；
- 线路绑定。

中央不修改 X-Panel Inbound 和 Client 配置。

### 11.2 到期状态

按每天或每小时任务计算：

~~~text
active:    expiry_time 为空或大于提醒窗口
expiring:  expiry_time 在提醒窗口内
expired:   expiry_time 小于当前时间
disabled:  Inbound.enable=false 或中央手工禁用
error:     连续同步失败或资源状态不确定
~~~

提醒窗口默认为 7 天，可配置为 3、7、15 天。

### 11.3 后续写操作

第二阶段可通过 X-Panel API 支持创建、续期、禁用和删除：

~~~text
POST /inbounds/add
POST /inbounds/update/:id
POST /inbounds/del/:id
POST /inbounds/addClient
POST /inbounds/:id/delClient/:clientId
POST /inbounds/updateClient/:clientId
~~~

写操作必须采用任务表、幂等 request_id、执行结果和审计日志。中央修改成功不代表 Xray 已热加载成功，必须显示“中央已提交”和“节点已执行”两个状态。

## 12. 页面设计

### 12.1 总览页面

指标：

- 节点总数、在线节点、异常节点；
- 线路机数、落地机数；
- 有效用户数、7 天内到期数；
- 今日流量、本月流量；
- 本月用户收入、本月节点成本、本月出口 IP 成本；
- 预计毛利润；
- 最近同步失败和流量重置事件。

### 12.2 用户页面

每行一个业务 Inbound：

| 列 | 内容 |
|---|---|
| 用户 | 中央业务名称和远程 Remark |
| 节点 | 线路机名称 |
| Inbound | Tag、协议、端口 |
| 设备数 | 启用 Client 数 |
| 流量 | 上行、下行、总量 |
| 到期 | 到期时间和剩余天数 |
| 月费 | 金额和币种 |
| 线路 | 落地机和出口 IP 配置 |
| 状态 | 正常、即将到期、已到期、异常 |
| 同步 | 最近成功时间 |
| 操作 | 查看详情、编辑业务字段、立即同步 |

筛选：节点、线路、状态、到期范围、流量范围、是否最近活跃、同步状态。

用户详情抽屉显示 Client/Email 设备列表、Inbound 流量趋势、最近同步、线路和出口 IP 配置归属。

当前实现补充：用户详情接口已先提供最近 30 条累计流量快照和 `resetDetected` 标识，页面以轻量列表展示；业务字段编辑仅写中央数据库并记录审计日志，X-Panel/Xray 字段保持只读。

### 12.3 节点页面

显示节点类型、区域、IP、在线状态、CPU、内存、磁盘、Xray/X-Panel 版本、用户数、总流量、月成本和最后同步时间。

节点详情显示：

- Inbound 列表；
- 节点拥有的公网出口 IP（线路机直出或落地机出口）；
- 最近同步错误；
- 最近状态历史。

### 12.4 线路与成本页面

使用表格维护：

- 线路机；
- 落地机；
- 线路机 Outbound Tag；
- 落地机 SS Inbound；
- 出口 IP 来源（线路机、落地机或独立 S5）；
- 出口位置（线路机直出、落地机出口或外部 S5）；
- 出口 IP 月成本；
- 线路有效用户数；
- 线路备注。

### 12.5 事件和同步页面

显示同步失败、节点离线、流量重置、Inbound 缺失和版本兼容异常。支持按节点、事件类型和日期筛选。

### 12.6 前端路由和菜单

Soybean Admin 的原始菜单只作为参考，最终菜单固定为：

| 路由 | 页面 | 数据来源 |
|---|---|---|
| `/dashboard` | 总览 | `/api/dashboard` |
| `/users` | 用户列表 | `/api/users` |
| `/users/:id` | 用户详情 | `/api/users/:id` |
| `/nodes` | 节点列表 | `/api/nodes` |
| `/nodes/:id` | 节点详情 | `/api/nodes/:id` |
| `/routes` | 线路关系 | `/api/routes` |
| `/exit-ips` | 出口 IP | `/api/exit-ips` |
| `/finance` | 成本与收入 | `/api/costs/summary` |
| `/events` | 同步和异常事件 | `/api/events`、`/api/sync-runs` |

前端实现不得保留模板示例路由，也不得在菜单中出现角色、部门、字典、组件示例、项目管理或无关工作台入口。

### 12.7 UI 验收重点

- 首屏明确显示节点、用户、到期、流量和财务五类核心信息；
- 用户表格按 Inbound 聚合，不能把多个 Client 渲染成多条用户记录；
- 设备详情在用户抽屉内展示，不抢占主导航；
- 数据延迟、节点离线、流量重置和同步失败必须有明显但克制的状态标识；
- 统计卡片、表格和图表均显示数据时间或统计周期；
- 删除和高风险操作使用确认弹窗；
- 用户路径明确区分线路机直出、经落地机和独立 S5；经落地机必须选择落地 Inbound，落地 Inbound 只能在路径配置中出现，不能出现在用户来源列表；
- 落地节点未同步时显示等待同步状态，不把基础设施资源静默显示成空用户列表；
- 空数据、加载中、权限不足和接口错误都有专门状态；
- 不使用模板中的 Mock 数字、随机图表或静态占位用户。

## 13. 中央 HTTP API

### 13.1 管理员 API

~~~text
POST /api/auth/login
POST /api/auth/logout
GET  /api/auth/me

GET  /api/dashboard
GET  /api/users
GET  /api/users/:id
GET  /api/users/:id/path-assets
PATCH /api/users/:id
PUT  /api/users/:id/route
DELETE /api/users/:id/route
GET  /api/nodes
GET  /api/nodes/:id
PATCH /api/nodes/:id
POST /api/nodes/:id/sync
GET  /api/routes
POST /api/routes
PATCH /api/routes/:id
DELETE /api/routes/:id
GET  /api/exit-ips
POST /api/exit-ips
PATCH /api/exit-ips/:id
DELETE /api/exit-ips/:id
GET  /api/costs/summary
GET  /api/events
GET  /api/sync-runs
~~~

### 13.2 分页和筛选

列表 API 统一支持 page、page_size、sort、order、keyword、status、node_id。默认 page_size=50，最大 200。总览接口返回聚合值和数据时间，不在前端重复计算财务指标。

前端基座使用以下稳定的 JSON 字段约定，后端可以在数据库或 Agent 字段不同的情况下做映射：

- 分页响应为 `{ items, total, page, pageSize, dataAt }`；查询参数保持 `page_size`、`node_id` 等 API 命名；
- 用户摘要至少包含 `id`、`name`、`nodeId`、`nodeName`、`inboundTag`、`status`、`expiresAt`、`clientCount`、`trafficBytes`、`lastActivityAt`；
- 节点摘要至少包含 `id`、`name`、`type`、`status`、`host`、`xpanelVersion`、`xrayVersion`、`lastSeenAt`、`lastSyncAt`、`monthlyCost`；
- 线路、出口 IP、财务和事件字段以中央前端 `src/typings/api/central.d.ts` 为准；线路同时区分线路机 Reality 入站和落地机 `landingInboundTag`（SS 入站）；
- 时间字段统一使用 RFC 3339 字符串，金额使用数字加 `currency`，流量使用字节数；
- `dataAt` 表示本次查询使用的最新同步时间，缺失或接口失败时页面必须显示明确状态，不得补造统计数字。

## 14. 一致性和冲突处理

### 14.1 字段权威

- X-Panel 字段同步覆盖中央镜像字段；
- 中央业务字段不被节点同步覆盖；
- 远端资源缺失先标记，不立即删除；
- 中央手工分类和线路关系始终保留。

### 14.2 配置漂移

Agent 计算 Inbound 的标准化 config_hash。Hash 变化时写入事件，但第一版不自动覆盖节点或中央业务字段。管理员查看差异后决定是否接受。

### 14.3 并发

同一节点同步串行处理。中央手工修改用户业务字段时使用 updated_at 或版本号，避免后台页面覆盖他人修改。

## 15. 安全设计

### 15.1 Agent 通信

- HTTPS 是最低要求；
- 每个节点独立 Token；
- Token 使用哈希存储；
- 支持撤销和轮换；
- 可选 mTLS；
- 请求包含时间戳和随机 request_id；
- 服务端限制时钟偏差和重复请求；
- 节点停用后立即拒绝其同步请求。

### 15.2 X-Panel 和 Xray

- Agent 默认访问 127.0.0.1；
- X-Panel 不要求开放公网管理端口；
- Xray Commander 保持本机或 Unix Socket；
- 不把 Xray gRPC 暴露到公网；
- 不上传完整配置、私钥、密码和订阅 Token。

### 15.3 中央面板

- 管理员密码使用 bcrypt 或 Argon2 哈希；
- 支持 TOTP；
- Cookie 设置 HttpOnly、Secure、SameSite；
- 所有写接口启用 CSRF 防护；
- 登录限速和失败审计；
- 重要删除操作二次确认；
- 日志脱敏；
- 数据库定期备份和恢复演练。

## 16. 可观测性和运维

中央服务记录：

- Agent 注册和心跳；
- 同步耗时、数量和错误；
- X-Panel 登录失败；
- API 响应格式异常；
- 节点离线和恢复；
- 流量重置；
- 管理员操作。

健康检查：

~~~text
GET /health/live
GET /health/ready
~~~

备份策略：

- 每日 SQLite 备份；
- 保留 14 至 30 个版本；
- 升级前执行手工快照；
- 定期在隔离目录验证恢复。

## 17. 技术选型和部署

### 17.1 推荐技术

| 层 | 技术 |
|---|---|
| Agent | Go 静态二进制 |
| 中央后端 | Go + Gin/Echo |
| 数据库 MVP | SQLite WAL |
| 数据库扩展 | PostgreSQL |
| 前端基座 | Soybean Admin（Vue 3 + Vite + TypeScript + Pinia + UnoCSS + 现有组件库） |
| 前端业务层 | 本项目专属 views、components、API 类型和状态 |
| 反向代理 | Caddy 或 Nginx |
| 节点服务管理 | systemd |
| 日志 | JSON 日志 + logrotate |

### 17.2 中央部署

~~~text
Caddy/Nginx :443
       |
       v
central-panel
       |
       v
panel.db (SQLite WAL)
~~~

### 17.3 Agent 部署

~~~text
/usr/local/bin/xpanel-agent
/etc/xpanel-agent/config.yaml
/var/log/xpanel-agent/
~~~

Agent 使用独立 systemd 服务，X-Panel/Xray 升级和重启不改变 Agent 数据目录。

## 18. 测试方案

### 18.1 单元测试

- X-Panel 响应字段映射；
- 到期状态计算；
- 流量增量和重置判断；
- 费用汇总；
- 远程资源幂等更新；
- 缺失 Inbound 延迟归档；
- Token 校验和时间戳校验。

### 18.2 集成测试

- 不同 Base Path 的 X-Panel；
- 登录过期和重新登录；
- 节点重复同步；
- 节点暂时离线；
- X-Panel 手动重置流量；
- Inbound 增删改；
- 同 Email 出现在不同节点；
- Agent 版本升级和回滚。

前端还必须验证：

- Soybean Admin 示例路由和 Mock 数据已移除；
- 登录、路由守卫、刷新 Session 和退出流程正常；
- 用户列表按 Inbound 聚合，Client 只显示为设备；
- 筛选、分页、抽屉详情和错误状态使用真实 API；
- 主题切换和窄屏布局不遮挡核心数据；
- 前端删除模板页面后没有死链、菜单空项或未使用资源引用。

### 18.3 验收测试

1. 一个 Inbound 下三个 Client 只显示一个业务用户；
2. 中央流量与 X-Panel 页面累计值一致；
3. 中央同步不改变 Xray/X-Panel 的统计值；
4. 相同 Email 在不同节点不会覆盖；
5. 流量重置不出现负数；
6. 节点离线保留最后快照；
7. 节点恢复后自动继续同步；
8. 能显示 7 天内到期用户；
9. 能录入节点和出口 IP 成本；
10. 能录入用户月费并计算预计毛利润；
11. X-Panel/Xray 升级不会丢失中央历史数据。

## 19. 风险和应对

| 风险 | 影响 | 应对 |
|---|---|---|
| X-Panel 版本不同 | 字段和登录行为差异 | Agent 版本探测、适配器和原始响应测试夹具 |
| X-Panel 流量被手动重置 | 历史增量断点 | 保存快照、检测回退、重新建立基线 |
| 节点离线 | 数据不新鲜 | 显示最后同步时间和离线状态，不删除历史 |
| 本地 Inbound 被删除 | 中央用户误删 | 缺失三次后归档，人工确认 |
| 共享 SS 无法归因出口 IP | 出口 IP 报表误解 | 页面明确“配置归属用户数” |
| Agent 凭据泄露 | 节点数据暴露 | 节点独立 Token、轮换、HTTPS、最小权限 |
| SQLite 损坏 | 中央数据不可用 | WAL、每日备份、恢复演练 |
| 中央和节点写操作冲突 | 配置不一致 | 第一版只读；后续使用任务、版本和审计 |

## 20. 版本演进

### v1.0

只读同步、用户/节点/线路/成本展示、到期提醒、流量快照。

### v1.1

更多历史报表、通知、导出、配置漂移差异和成本分摊规则。

### v2.0

中央创建用户、续期、禁用、添加设备、任务队列、执行结果回传和完整审计。

### v3.0（可选）

成本按用户/线路/流量分摊、支付记录、用户实际出口归因、多管理员角色和多租户。

## 21. 最终判断

本项目的最小可靠闭环为：

~~~text
X-Panel 继续管理节点
Agent 只读取 X-Panel
中央按 Inbound 管理业务用户
Client 只统计设备数
流量直接读取 X-Panel 累计值
线路、出口 IP、成本和月费由中央维护
~~~

Soybean Admin 仅作为前端工程和 UI 基础，项目必须删除无关模板功能并重新实现所有业务页面。这样既能复用成熟的后台交互，又不会让模板的示例数据或业务结构污染中央管理口径。

这套方案不会修改现有 X-Panel/Xray，兼容后续独立升级，能够先解决多节点查看、到期管理和成本收入汇总问题，同时明确不承诺现有链路无法提供的“实际出口 IP 精确归因”。
