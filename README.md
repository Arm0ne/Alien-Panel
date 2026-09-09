# Alien-Panel v1.0.0

这是 Alien-Panel 的第一个可部署生产版本。它包含最新的中央服务、节点 Agent、Vue 管理前端、数据库迁移、Linux amd64 发布包和部署脚本。

本目录是独立的交付包，不依赖当前开发工作区中的 `node_modules`、本地数据库、日志、预览 HTML 或 Git 元数据。生产服务器只需要 Linux amd64、Docker Compose 和网络访问；服务器不需要安装 Go、Node.js 或 pnpm。

## 快速部署

推荐使用 HTTPS 域名和宿主机 Nginx/Caddy。将本目录提交到 GitHub 后，在 Ubuntu/Debian amd64 服务器执行：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh \
  | sudo bash -s -- --repo https://github.com/Arm0ne/Alien-Panel.git --ref main --domain panel.example.com
```

仓库使用默认值时可以进一步缩短为：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh | sudo bash -s -- --domain panel.example.com
```

如果只做临时公网测试，不使用域名：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh \
  | sudo bash -s -- --repo https://github.com/Arm0ne/Alien-Panel.git --ref main
```

安装脚本会下载并校验本目录的发布包，安装 Docker，生成随机管理员密码和 Agent 引导 Token，启动中央服务与内置 Nginx。首次安装的凭据保存在 `/opt/alien-panel/.env`，权限为 `0600`。升级时会保留该文件和 Docker 数据卷。

临时测试地址为 `http://服务器IP:18080`。指定 `--domain` 后，内置 Web 只监听 `127.0.0.1:18080`，需要由宿主机 Nginx 或 Caddy 终止 HTTPS 并反向代理到这个地址。

## 目录

```text
backend/       中央 Go 服务、SQLite 迁移、API、测试和运维命令
agent/         节点 Go Agent、X-Panel 只读采集、systemd 安装文件
frontend/      Vue 3 + TypeScript 管理前端源码，不含 node_modules/dist
deploy/        Docker、Nginx/Caddy、构建、备份、迁移、恢复和冒烟脚本
release/       Linux amd64 central/agent 二进制及 Agent SHA-256 校验值
docs/          当前版本说明；历史设计和原始计划位于 docs/reference/
checksums/     发布文件校验清单
```

## 当前版本口径

- 产品版本：`v1.0.0`。
- Agent 显示版本：`v1.0.3`，由 `agent/VERSION` 控制；Agent 构建提交号只作为诊断信息。
- 数据库迁移：`001_initial.sql` 到 `021_billing_history_import.sql`，服务启动时自动幂等执行。
- 用户出口路径支持多个固定出口 IP；旧的 `exitIpId` API 字段仍兼容，第一项作为兼容值。
- 业务用户按线路机用户 Inbound 建模，Client/Email 是设备；不同 Inbound 或节点的相同 Email 不合并。

## 文档入口

- [架构与底层代码](docs/ARCHITECTURE.md)
- [业务逻辑与数据口径](docs/BUSINESS-LOGIC.md)
- [数据库、迁移与备份](docs/DATABASE.md)
- [生产部署、一键安装与升级](docs/DEPLOYMENT.md)
- [日常运维与故障排查](docs/OPERATIONS.md)
- [版本发布与回滚](docs/RELEASE.md)
- [发布清单](RELEASE_MANIFEST.json)

## 安全边界

中央服务只保存 Agent 上报的 X-Panel 数据，不向 X-Panel/Xray 写入用户、Client、到期或流量字段。生产环境必须使用 HTTPS；管理员密码、Agent 注册 Token、节点 Token、X-Panel 密码和数据库备份不能提交到 GitHub。删除数据库卷会删除业务数据，升级时不要执行 `docker compose down -v`。

本版本已通过后端和 Agent Go 测试、`go vet`、前端类型检查、Lint 和生产构建；正式接入全部节点前仍应完成公网灰度、备份恢复演练和真实数据核对。
