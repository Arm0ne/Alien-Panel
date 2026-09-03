# X-Panel Central 运维与故障手册

本文是 W9-08 的交接手册。公网灰度前请先在测试环境完整走一遍，命令中的域名、路径、密码和 Token 必须替换为部署环境的值。

## 1. 中央服务部署

### Linux systemd

1. 构建二进制：

   ```sh
   cd backend
   go build -o /usr/local/bin/xpanel-central ./cmd/server
   go build -o /usr/local/bin/xpanel-db-maintenance ./cmd/db-maintenance
   ```

2. 创建运行用户、目录和环境文件。环境文件权限设为 `0600`，不要提交到 Git：

   ```sh
   useradd --system --home-dir /var/lib/xpanel-central --no-create-home --shell /usr/sbin/nologin xpanel-central
   install -d -o xpanel-central -g xpanel-central -m 0750 /var/lib/xpanel-central
   install -d -o root -g xpanel-central -m 0750 /etc/xpanel-central
   install -d -m 0750 /var/backups/xpanel-central
   cat >/etc/xpanel-central/central.env <<'EOF'
   XPANEL_LISTEN=127.0.0.1:8090
   XPANEL_DATABASE=/var/lib/xpanel-central/panel.db
   XPANEL_ADMIN_USER=admin
   XPANEL_ADMIN_PASSWORD=replace-with-a-long-secret
   XPANEL_AGENT_REGISTRATION_TOKEN=replace-with-a-separate-bootstrap-secret
   XPANEL_CORS_ORIGINS=https://panel.example.com
   EOF
   chown root:xpanel-central /etc/xpanel-central
   chown root:xpanel-central /etc/xpanel-central/central.env
   chmod 640 /etc/xpanel-central/central.env
   ```

3. 安装并启动 `deploy/xpanel-central.service`：

   ```sh
   install -m 0644 deploy/xpanel-central.service /etc/systemd/system/xpanel-central.service
   systemctl daemon-reload
   systemctl enable --now xpanel-central.service
   curl --fail http://127.0.0.1:8090/health/live
   curl --fail http://127.0.0.1:8090/health/ready
   ```

   服务只监听回环地址；公网访问由 Nginx 或 `deploy/Caddyfile` 终止 HTTPS 并反向代理 `/api`、`/agent` 和健康检查。

   部署完成后执行冒烟检查：

   ```sh
   XPANEL_STAGING_URL=https://panel.example.com \
   bash deploy/staging-smoke.sh
   ```

   设置 `XPANEL_ADMIN_PASSWORD` 后，检查还会验证管理员登录和受保护的 Dashboard 接口；脚本不会输出访问 Token。

### Nginx 反向代理

如果服务器使用 Nginx，将 `deploy/nginx.conf` 复制到站点配置目录，替换域名、证书路径和前端目录后执行：

```sh
sudo cp deploy/nginx.conf /etc/nginx/sites-available/xpanel-central
sudo ln -sfn /etc/nginx/sites-available/xpanel-central /etc/nginx/sites-enabled/xpanel-central
sudo nginx -t
sudo systemctl reload nginx
```

模板中的 `proxy_pass http://127.0.0.1:8090;` 故意不带结尾斜杠，否则 Nginx 会删除 `/api` 或 `/agent` 前缀，导致中央路由返回 404。Nginx 只开放 80/443，中央服务仍保持回环监听。

#### 上传前端静态文件

在本地项目目录构建前端，并将构建目录上传到服务器的临时目录。下面的示例统一使用 `/tmp/xpanel-frontend`，不要把上传目录和 Nginx 的最终目录混用：

```powershell
Set-Location D:\轻量Panel\frontend
pnpm.cmd build
scp -P <SSH端口> -r .\dist <SSH用户>@<服务器IP>:/tmp/xpanel-frontend
```

在服务器上先确认 `index.html` 的实际位置：

```sh
find /tmp/xpanel-frontend -maxdepth 2 -type f -name index.html -print
```

若输出为 `/tmp/xpanel-frontend/index.html`，执行：

```sh
sudo install -d -m 0755 /opt/xpanel-central/frontend/dist
sudo cp -a /tmp/xpanel-frontend/. /opt/xpanel-central/frontend/dist/
```

若输出为 `/tmp/xpanel-frontend/dist/index.html`，复制时将源目录改为 `/tmp/xpanel-frontend/dist/.`。复制完成后确认：

```sh
test -s /opt/xpanel-central/frontend/dist/index.html
sudo nginx -t && sudo systemctl reload nginx
```

### Windows 本地联调

```powershell
$env:XPANEL_ADMIN_PASSWORD = 'local-only-password'
$env:XPANEL_AGENT_REGISTRATION_TOKEN = 'local-bootstrap-secret'
$env:XPANEL_DATABASE = '.\data\panel.db'
go run .\cmd\server
```

前端开发服务默认是 `http://localhost:9527`，中央 API 默认是 `http://localhost:8090`。本地 HTTP 不启用 HSTS。

### Docker 一键测试部署

服务器不需要安装 Node.js 或 Go。将 `release/xpanel-central` 和
`deploy/frontend-dist` 一起推送到 GitHub 后，在服务器执行：

```sh
curl -fsSL https://raw.githubusercontent.com/OWNER/REPO/main/deploy/install-docker.sh \
  | sudo bash -s -- --repo https://github.com/OWNER/REPO.git --domain panel.example.com
```

脚本会自动安装 Docker（Ubuntu/Debian 或 Fedora）、生成随机管理员密码和
Agent 引导 Token，启动中央服务和内置 Nginx，默认通过
`http://<服务器IP>:18080` 访问。密钥保存在 `/opt/xpanel-central/.env`，更新
时不会覆盖已有密钥或数据库。未指定 `--domain` 时 18080 对公网监听；指定
`--domain` 后自动改为仅监听 `127.0.0.1`，适合已有宿主机 Nginx 的场景。此时
将域名反代到
`127.0.0.1:18080`；完整示例见 [`deploy/README.md`](deploy/README.md)。

如果需要清空旧的 systemd 或 Docker 测试部署（包含数据库、备份和上传
文件），使用部署包中的 `deploy/uninstall-all.sh`。它默认只预览，必须显式
传入 `--yes` 才会删除；TLS 证书和节点机器上的 Agent 不会被删除。

## 2. 节点接入

1. 在“节点管理”点击“接入节点”，填写节点名称、类型和管理地址；节点 Key 可以留空，由中央生成稳定 Key。
2. 需要多个公网出口时，在向导一次填写出口 IP 列表；主公网地址 `publicIp` 只作为管理地址，不替代出口资产。
3. 保存后只复制一次性显示的 Agent Token 和配置模板。X-Panel 用户名/密码只写入目标节点 `/etc/xpanel-agent/agent.yaml`，不上传中央。
4. 在节点安装 Agent：

   ```sh
   install -m 0755 xpanel-agent /usr/local/bin/xpanel-agent
   install -d -m 0750 /etc/xpanel-agent /var/lib/xpanel-agent
   install -m 0600 agent.yaml /etc/xpanel-agent/agent.yaml
   install -m 0644 agent/deploy/xpanel-agent.service /etc/systemd/system/xpanel-agent.service
   systemctl daemon-reload
   systemctl enable --now xpanel-agent.service
   ```

5. 检查 Agent 日志和中央节点详情。首次成功同步会建立线路机 Inbound 与业务用户的映射；落地机 Inbound 不会自动建立业务用户。

Agent 只读采集 X-Panel 累计计数，不调用 `QueryStats(reset=true)`。中央收到重复 `sync_id` 时返回原结果，不重复写入资源或快照。

## 3. 日常管理

- “节点管理”中的停用会立即拒绝该节点的 heartbeat、sync 和立即同步请求；重新启用后恢复认证。
- 删除节点是保留历史的软删除。若仍有关联线路或出口 IP，先解除关联；删除后旧 Token 和旧 Node Key 不能重新注册原节点。
- 用户详情中的线路分配只更新中央配置：可选择线路模板、固定出口 IP 或线路出口池；线路机必须与用户主 Inbound 匹配。
- 线路出口位置分为线路机直出、落地机出口和独立 S5。不同出口逻辑应建立不同线路，避免把不同来源混入同一出口池。
- 页面中的“配置归属用户数”是当前有效配置统计，不等同于实际产生流量的用户数。

## 4. 备份、升级与恢复

每日备份（默认保留 14 个版本）：

```sh
XPANEL_DATABASE=/var/lib/xpanel-central/panel.db \
XPANEL_BACKUP_DIR=/var/backups/xpanel-central \
bash deploy/backup.sh
```

升级迁移会先备份，再执行幂等迁移：

```sh
XPANEL_DATABASE=/var/lib/xpanel-central/panel.db \
XPANEL_BACKUP_DIR=/var/backups/xpanel-central \
bash deploy/migrate.sh
```

恢复前必须停止服务并确认备份文件：

```sh
systemctl stop xpanel-central.service
bash deploy/restore.sh \
  --source /var/backups/xpanel-central/panel-YYYYMMDDTHHMMSSZ.sqlite3 --yes
systemctl start xpanel-central.service
curl --fail http://127.0.0.1:8090/health/ready
```

恢复命令会先验证源文件、创建恢复前安全快照，再原子替换数据库并复核完整性。恢复后应检查 Dashboard 数据时间、节点数量、用户线路分配和最近同步记录。

## 5. 故障排查

### 节点显示离线

1. 查看 `systemctl status xpanel-agent` 和 `journalctl -u xpanel-agent -n 100`。
2. 确认节点能访问中央 HTTPS 地址，Token 未被轮换或节点未被停用。
3. 确认 X-Panel 本机可访问 `xpanel_url`，登录凭据和 `xpanel_base_path` 正确。
4. 离线期间中央保留最后一次成功快照，不要手工删除 Inbound 或流量记录；恢复后 Agent 会继续从最新成功快照计算。

### X-Panel 接口超时或登录失败

Agent 对网络错误、超时、408/425/429/5xx 最多尝试四次，按 30 秒、2 分钟、10 分钟退避；401 只重新登录一次。连续错误凭据会进入登录失败退避。修正本机配置后重启 Agent，不要在中央清零流量。

### 中央接口超时

中央客户端会使用同一 `sync_id` 和请求 ID 重试；即使第一次请求已提交而响应丢失，后续请求也只返回幂等结果。检查中央服务状态、反向代理、DNS/TLS 和防火墙，恢复后观察 `sync_runs` 是否最终为 `success`。

### 同步失败、Inbound 缺失或归档

只有成功的完整快照才会增加缺失计数；失败同步或节点离线不会归档。连续三次成功快照仍缺失才会标记归档，并保留业务关联、Client 和历史快照。先核对 X-Panel 是否真的删除或改建了 Inbound，再决定是否恢复配置。

### 流量回退或出现重置事件

累计值下降会记录 `traffic_reset`，该次增量按新基线计算为非负值。使用只读 `traffic-check` 核对时间段，不要调用任何 Xray reset 接口：

```sh
cd /path/to/project/backend
go run ./cmd/traffic-check \
  --database /var/lib/xpanel-central/panel.db \
  --from 2026-09-01 --to 2026-09-02 --format json
```

### 页面卡顿或接口不可用

先访问 `/health/live` 和 `/health/ready`，再检查浏览器网络请求是否指向正确的中央 API。确认服务端没有反复重启、磁盘未满、数据库备份未占用异常锁；不要通过关闭认证或直接暴露 X-Panel 规避问题。

## 6. 灰度与回滚检查单

- [ ] HTTPS 域名、CSP/HSTS 和 `/health/ready` 通过；
- [ ] 已执行一次备份并验证备份文件；
- [ ] 先接入一台线路机和一台落地机，逐项核对 Inbound、Client、到期时间和累计流量；
- [ ] 验证线路机直出、落地机出口和独立 S5 三种来源；
- [ ] 验证停用节点、断开网络、恢复网络和 Agent 重启；
- [ ] 观察至少 3 天同步、离线和流量重置事件；
- [ ] 回滚时停止服务、保留当前数据库安全快照、恢复指定备份并重新执行健康检查；
- [ ] 记录版本、迁移版本、备份文件和验收结果。
