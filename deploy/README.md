# Alien-Panel 部署文件

生产部署的唯一入口是包根目录的 `../install.sh`。它会下载并校验同一版本
的发布包，安装 Docker Compose，生成首次登录凭据，然后启动 `alien-panel`
Compose 项目。

```bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/Alien-Panel-v1.0.0-production/install.sh \
  | sudo bash -s -- --repo https://github.com/Arm0ne/Alien-Panel.git --ref main
```

安装后的目录为 `/opt/alien-panel`：

- `.env`：管理员密码、Agent 注册 Token 和端口配置，权限为 `0600`；
- `deploy/`：Compose、Nginx/Caddy 和备份恢复模板；
- `frontend-dist/`：已构建的 Vue 前端；
- `release/`：Linux amd64 中央服务和维护工具。

## 常用命令

```bash
cd /opt/alien-panel
sudo docker compose -p alien-panel -f deploy/docker-compose.yml ps
sudo docker compose -p alien-panel -f deploy/docker-compose.yml logs --tail=100 central
curl -fsS http://127.0.0.1:18080/health/ready
```

重复执行包根目录的安装命令即可升级。脚本会替换发布文件，保留 `.env`
和 `alien-panel_central-data` 数据卷。不要在升级时使用 `down -v`。

只停止容器并保留数据：

```bash
sudo bash /opt/alien-panel/deploy/uninstall-docker.sh
```

`install-docker.sh` 是旧入口的兼容转发脚本；本地运行时会调用包根目录
脚本，直接通过 raw URL 运行时会转发到固定的 v1.0.0 安装入口。

## 反向代理

指定 `--domain panel.example.com` 后，临时 Web 端口只监听
`127.0.0.1:18080`。宿主机 Nginx 可参考 `nginx.conf`，Caddy 可参考
`Caddyfile`；代理必须转发 `/api/`、`/agent/` 和健康检查路径，并将
`X-Forwarded-Proto` 设置为 `https`。

## 备份与恢复

`backup.sh`、`migrate.sh` 和 `restore.sh` 调用同版本的
`release/xpanel-db-maintenance`。默认数据库路径是中央容器中的
`/var/lib/xpanel-central/panel.db`，Docker 部署时应先使用维护工具或
`docker run --volumes-from` 方式挂载数据卷，再执行备份和恢复。完整步骤
见 `../docs/DATABASE.md` 与 `../docs/OPERATIONS.md`。

## 从源码重新构建

在开发机执行：

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\build-bundle.ps1 -GoArch amd64
```

构建会生成前端静态文件、中央服务、Agent、数据库维护工具、流量核对工具
及 SHA-256 校验文件。发布前应完成 `go test ./...`、`go vet ./...`、前端
typecheck/lint/build 和 `git diff --check`。


