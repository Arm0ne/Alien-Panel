# 生产部署、一键安装与升级

## 1. 支持范围

本交付包的预构建二进制目标为 Linux amd64，适用于 Ubuntu/Debian 和已安装 Docker Compose 的 Fedora/RHEL 系统。脚本会尝试安装 Docker；如果系统不提供受支持的包管理器，请先手动安装 Docker Engine 和 Compose 插件。

中央 Docker 服务使用 Alpine 容器运行 `release/xpanel-central`，前端使用内置 Nginx 容器，数据保存在名为 `alien-panel_central-data` 的 Docker volume。公网生产建议在宿主机使用 Nginx 或 Caddy 终止 TLS。

## 2. 一行安装

脚本入口是本目录的 `install.sh`，它可以从 GitHub raw URL 直接执行。脚本本身会下载同一版本目录的压缩包，校验 SHA-256，再复制运行所需文件到安装目录。

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh \
  | sudo bash -s -- --repo https://github.com/Arm0ne/Alien-Panel.git --ref main --domain panel.example.com
```

使用默认仓库和 `main` 分支时，命令可以简化为：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh | sudo bash -s -- --domain panel.example.com
```

临时测试端口：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh \
  | sudo bash -s -- --repo https://github.com/Arm0ne/Alien-Panel.git --ref main --port 18080
```

常用参数：

| 参数 | 默认值 | 用途 |
|---|---|---|
| `--repo` | 必填 | GitHub 仓库 URL |
| `--ref` | `main` | 分支或 tag |
| `--domain` | 空 | HTTPS 域名；设置后临时端口只监听回环 |
| `--port` | `18080` | 内置 Web 端口 |
| `--bind` | 无 | `0.0.0.0` 或 `127.0.0.1` |
| `--dir` | `/opt/alien-panel` | 安装目录 |
| `--admin-user` | `admin` | 首次管理员用户名 |

首次安装创建：

- `/opt/alien-panel/.env`，包含随机 `XPANEL_ADMIN_PASSWORD`、`XPANEL_AGENT_REGISTRATION_TOKEN` 和 CORS 来源；
- `/opt/alien-panel/deploy/`，Docker Compose 和 Nginx 模板；
- `/opt/alien-panel/frontend-dist/`，前端静态文件；
- `/opt/alien-panel/release/xpanel-central`，中央 Linux amd64 二进制。

查看凭据：

```bash
sudo grep -E '^(XPANEL_ADMIN_USER|XPANEL_ADMIN_PASSWORD|XPANEL_AGENT_REGISTRATION_TOKEN)=' /opt/alien-panel/.env
```

## 3. HTTPS 反向代理

指定 `--domain` 后，内置 Web 监听 `127.0.0.1:18080`。宿主机 Nginx 的最小反代：

```nginx
server {
    listen 443 ssl;
    server_name panel.example.com;

    location / {
        proxy_pass http://127.0.0.1:18080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
    }
}
```

也可以直接参考 `deploy/nginx.conf` 或 `deploy/Caddyfile`。Agent 的 `central_url` 应填写公开 HTTPS API 基址，例如 `https://panel.example.com/api`。

## 4. 安装后检查

```bash
cd /opt/alien-panel
sudo docker compose -p alien-panel -f deploy/docker-compose.yml ps
curl -fsS http://127.0.0.1:18080/health/live
curl -fsS http://127.0.0.1:18080/health/ready
sudo docker compose -p alien-panel -f deploy/docker-compose.yml logs --tail=100 central
```

浏览器登录后，在节点管理创建节点并复制 Agent 一行安装命令。Agent 安装脚本会在目标节点本地询问 X-Panel 用户名和密码，凭据不会发送到中央。

## 5. 升级

重复执行相同命令即可：

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh \
  | sudo bash -s -- --repo https://github.com/Arm0ne/Alien-Panel.git --ref main --domain panel.example.com
```

脚本会替换中央二进制、前端静态文件和部署模板，保留 `.env` 与数据库 volume。服务启动时自动执行尚未应用的 SQLite 迁移，包括历史账单导入的 `021_billing_history_import.sql`。升级前建议先执行备份。

## 6. 手动源码构建

在开发机需要 Go 1.27、Node.js 20.19+ 和 pnpm 10.5+：

```powershell
cd D:\轻量Panel
powershell -ExecutionPolicy Bypass -File .\deploy\build-bundle.ps1 -GoArch amd64
```

这会构建 `frontend/dist`、复制 `deploy/frontend-dist`，并生成 Linux/amd64 的 central、Agent 和校验文件。发布前执行后端/Agent 测试、`go vet`、前端 typecheck/lint/build 和 `git diff --check`。

## 7. 停止和卸载

保留数据库但停止容器：

```bash
sudo bash /opt/alien-panel/deploy/uninstall-docker.sh
```

删除容器、安装目录和 volume 前必须确认备份。不要把 `docker compose down -v` 用作普通升级命令。
