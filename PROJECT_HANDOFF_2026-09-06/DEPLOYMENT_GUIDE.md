# X-Panel Central 当前版本部署手册

> 适用版本：GitHub `main`，提交 `7f914fe` 及之后
> 适用服务器：Ubuntu/Debian amd64
> 推荐方式：Docker Central + Nginx 反向代理

## 1. 先区分两种包

项目根仓库的 `release/` 和 `deploy/frontend-dist/` 是可直接部署的发布产物；本交接目录的 `source/` 是干净源码快照，不包含二进制和前端构建产物。

- 只想更新测试服务器：直接使用第 2 节 GitHub 一键命令；
- 想从源码重新打包：先按第 3 节生成 Linux/amd64 发布产物，再提交到 GitHub 或复制到部署环境。

## 2. 从 GitHub 更新/安装 Docker 测试服务器

### 2.1 使用宿主机 Nginx（推荐）

`--domain` 只填写域名，不要填写 `https://`、端口或路径：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/deploy/install-docker.sh | \
sudo bash -s -- \
  --repo https://github.com/Arm0ne/Alien-Panel.git \
  --domain panel.example.com
```

脚本会：

1. 下载当前 `main` 的发布包；
2. 安装 Docker/Compose（Ubuntu/Debian 缺少时）；
3. 更新 `/opt/xpanel-central/release/xpanel-central`、前端静态文件和部署模板；
4. 保留已有 `/opt/xpanel-central/.env`、数据库和管理员凭据；
5. 构建并重启 `xpanel-central` 和内置 Web 容器。

指定域名时，内置 Web 默认监听 `127.0.0.1:18080`，宿主机 Nginx 反代到该端口。Agent 应使用域名访问中央，例如：

```yaml
central_url: https://panel.example.com/api
```

### 2.2 不使用宿主机 Nginx，直接暴露测试端口

省略 `--domain`，内置 Web 会监听 `0.0.0.0:18080`：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/deploy/install-docker.sh | \
sudo bash -s -- \
  --repo https://github.com/Arm0ne/Alien-Panel.git
```

访问：`http://服务器IP:18080`。生产环境不建议长期暴露此端口，应使用 HTTPS Nginx/Caddy。

### 2.3 更新后检查

```bash
cd /opt/xpanel-central

sudo docker compose -p xpanel-central \
  -f deploy/docker-compose.yml ps

curl -i http://127.0.0.1:18080/health/ready
```

预期：中央容器为 `healthy`，健康检查返回：

```json
{"code":"0000","msg":"ok","data":{"status":"ready"}}
```

失败时查看：

```bash
sudo docker compose -p xpanel-central \
  -f deploy/docker-compose.yml logs --tail=200 central

sudo docker inspect xpanel-central-central-1 \
  --format '{{range .State.Health.Log}}{{println .Output}}{{end}}'
```

## 3. 从干净源码重新生成发布包

### 3.1 Windows 开发机

安装 Node.js、pnpm 和 Go 1.27，进入源码根目录执行：

```powershell
cd D:\轻量Panel
powershell -ExecutionPolicy Bypass -File .\deploy\build-bundle.ps1 -GoArch amd64
```

该脚本会：

- 构建 `frontend/dist` 并复制到 `deploy/frontend-dist`；
- 使用 `GOOS=linux`、`GOARCH=amd64`、`CGO_ENABLED=0` 构建 `release/xpanel-central` 和 `release/xpanel-agent`；
- 生成 `release/xpanel-agent.sha256`。

不要直接在 Windows 下执行普通 `go build` 再把生成的 `.exe` 改名。发布前必须确认中央文件是 Linux ELF：

```powershell
Format-Hex -Path .\release\xpanel-central -Count 4
```

开头应为 `7F 45 4C 46`（ELF），不能是 `4D 5A`（Windows PE）。在 WSL/Linux 中可执行：

```bash
file release/xpanel-central release/xpanel-agent
```

应显示 `ELF 64-bit ... x86-64`。

### 3.2 Linux 构建机

```bash
export GOOS=linux GOARCH=amd64 CGO_ENABLED=0
cd backend && go build -trimpath -o ../release/xpanel-central ./cmd/server
cd ../agent && go build -trimpath -o ../release/xpanel-agent ./cmd/agent
cd ..
sha256sum release/xpanel-agent > release/xpanel-agent.sha256
```

前端仍需在 `frontend/` 执行 `pnpm build`，然后把 `frontend/dist/.` 复制到 `deploy/frontend-dist/`。

构建后建议执行：

```bash
cd backend && go test ./... && go vet ./...
cd ../frontend && pnpm typecheck && pnpm build
```

## 4. 宿主机 Nginx 配置

参考项目根目录 `deploy/nginx.conf`。必须反代以下路径到 `127.0.0.1:18080`：

- `/api/`：管理员 API；
- `/agent/`：兼容 Agent API；
- `/health/live`、`/health/ready`：健康检查；
- 其余路径：前端静态文件或 `/index.html` fallback。

如果使用 HTTPS，Agent 的 `central_url` 必须使用有效证书的域名。不要把包含用户名、密码或节点 Token 的配置提交到 GitHub。

## 5. Agent 安装（Ubuntu/Debian amd64）

### 5.1 一行在线安装

在中央面板“节点管理 → 接入节点”创建或选择节点，复制页面生成的一行命令，在目标节点以 root 执行。命令会：

1. 通过一次性安装 Token 获取节点凭据；
2. 下载并校验 Linux/amd64 Agent；
3. 在目标机本地询问 X-Panel 用户名和密码；
4. 写入 `/etc/xpanel-agent/agent.yaml`（权限 0600）；
5. 安装并启动 `/etc/systemd/system/xpanel-agent.service`。

X-Panel 凭据只写入目标节点，不上传中央。安装 Token 15 分钟后失效且只能使用一次。

### 5.2 手动安装兜底

```bash
sudo install -m 0755 xpanel-agent /usr/local/bin/xpanel-agent
sudo install -d -m 0750 /etc/xpanel-agent /var/lib/xpanel-agent
sudo install -m 0600 agent.yaml /etc/xpanel-agent/agent.yaml
sudo install -m 0644 xpanel-agent.service /etc/systemd/system/xpanel-agent.service
sudo systemctl daemon-reload
sudo systemctl enable --now xpanel-agent.service
```

检查：

```bash
sudo systemctl status xpanel-agent --no-pager
sudo journalctl -u xpanel-agent -f -o cat
```

正常时会周期性出现 `xpanel snapshot synchronized`。Agent 默认首次立即同步，之后按 `sync_interval`（默认 60 秒）轮询。

## 6. 升级、回滚和数据库安全

升级前建议备份：

```bash
cd /opt/xpanel-central
sudo ./deploy/backup.sh
```

更新时重新执行第 2 节命令即可。不要执行 `docker compose down -v`，否则会删除数据库 volume。升级脚本不会覆盖已有 `.env`。

出现问题时：

```bash
sudo docker compose -p xpanel-central \
  -f /opt/xpanel-central/deploy/docker-compose.yml down

# 恢复备份前，先创建当前数据库的安全副本，再使用 deploy/restore.sh。
```

恢复、迁移和全量清理请严格参考根目录 `OPERATIONS_RUNBOOK.md`；全量清理会删除数据库、备份和上传文件，只有明确确认后才执行。

## 7. Agent 离线排查

1. 中央服务：

```bash
curl -fsS http://127.0.0.1:18080/health/ready
sudo docker compose -p xpanel-central -f /opt/xpanel-central/deploy/docker-compose.yml ps
```

2. Agent：

```bash
sudo systemctl is-active xpanel-agent
sudo journalctl -u xpanel-agent --since '15 minutes ago' --no-pager -o cat
```

3. 核对配置：

```bash
sudo grep -E '^(central_url|node_key):' /etc/xpanel-agent/agent.yaml
```

4. 如果 Agent 使用 `http://服务器IP:18080`，但服务器 `.env` 中是 `XPANEL_BIND_ADDRESS=127.0.0.1`，该地址对远程 Agent 不可达。应改用宿主机 HTTPS 域名，或在确认防火墙策略后将绑定地址调整为 `0.0.0.0` 并重建 Web 容器。

5. 如果日志出现 `exec format error`，检查中央二进制是否 ELF x86-64；不要重新安装 Agent，先更新中央发布包。

## 8. 安全底线

- 不提交 `.env`、`agent.yaml`、真实账号密码、Token、私钥和备份数据库；
- 生产使用 HTTPS 和有效证书；
- Agent Token、管理员密码和数据库备份按秘密材料管理；
- 只有在备份成功并确认可恢复后，才进行全量节点升级；
- 任何中央写入 X-Panel 的功能都必须经过新的权限、审计和回滚设计，本版本不提供该能力。
