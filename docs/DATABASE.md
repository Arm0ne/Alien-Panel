# 数据库、迁移与备份

## 1. 运行方式

中央使用 SQLite WAL，默认数据库路径为 `/var/lib/xpanel-central/panel.db`（Docker 中位于 `central-data` volume）。外键和 busy timeout 已在 `backend/internal/db/db.go` 打开。服务启动时读取嵌入的迁移文件并按文件名顺序执行，每个文件只记录一次。

当前迁移到 `021_billing_history_import.sql`，空库和已有库都必须通过迁移测试。生产包不携带任何真实 `panel.db`、WAL/SHM 文件或备份。

## 2. 核心表关系

```text
nodes
  ├── inbounds ── user_inbounds ── users
  │                    └── clients
  ├── exit_ips
  └── node_costs

users ── user_paths ── user_path_exit_ips ── exit_ips
users ── user_billing_records
users ── user_renewal_candidates ── node_events
nodes ── sync_runs ── traffic_snapshots
```

关键约束：

- `users` 是中央业务对象，不等同于 Client。
- `user_paths.active_to IS NULL` 表示当前生效路径；旧路径保留历史。
- `user_path_exit_ips` 的 `(user_path_id, exit_ip_id)` 唯一，`position` 保留选择顺序，路径删除时关联记录级联删除。
- `user_paths.exit_ip_id` 是第一出口 IP兼容列，迁移 020 会从它回填一条关联。
- 收费记录和审计日志是业务历史，不随节点删除而物理删除。
- `user_billing_records.origin=historical_import` 表示一次性历史账单导入；`verification_status=unverified` 的待核实记录保持在历史中，但不计入现金或服务期收入，核验后才进入财务统计。

## 3. 备份

推荐使用同版本的 `xpanel-db-maintenance`。Docker 部署时让 Compose 临时
启动 `central` 服务并复用 `central-data` 卷，避免把容器内路径误当成宿主机路径：

```bash
sudo install -d -m 0750 /var/backups/alien-panel
sudo docker compose -p alien-panel -f /opt/alien-panel/deploy/docker-compose.yml run --rm --no-deps \
  -v /opt/alien-panel/release/xpanel-db-maintenance:/usr/local/bin/xpanel-db-maintenance:ro \
  -v /var/backups/alien-panel:/backups \
  --entrypoint /usr/local/bin/xpanel-db-maintenance central backup \
  --database /var/lib/xpanel-central/panel.db --backup-dir /backups --retention 14
```

`backup` 使用 SQLite `VACUUM INTO`，包含已提交的 WAL 数据，不需要直接复制 `.db-wal` 或 `.db-shm`。`verify` 会执行 `integrity_check` 和 `foreign_key_check`。

## 4. 升级迁移

Docker 升级重新执行一键脚本即可；它不会覆盖 `/opt/alien-panel/.env` 或 `central-data` volume。手动 systemd 部署时先备份，再执行：

```bash
XPANEL_DATABASE=/var/lib/xpanel-central/panel.db \
XPANEL_BACKUP_DIR=/var/backups/alien-panel \
XPANEL_DB_MAINTENANCE_BIN=/usr/local/bin/xpanel-db-maintenance \
bash deploy/migrate.sh
```

迁移失败时保留迁移前快照。不要手工编辑 `schema_migrations`，也不要跳过迁移文件。

## 5. 管理页面备份与恢复

登录管理页面后打开“数据备份”，点击“创建并下载备份”即可生成一致性 SQLite 文件。下载接口只允许管理员访问，备份生成在临时目录，响应完成后自动删除临时文件。

恢复时选择本系统生成的 `.sqlite3` 文件并确认。服务会先校验 SQLite 完整性和外键，再创建恢复前安全快照，原子替换数据库并重新打开连接；任何失败都会保留当前数据。恢复成功后所有旧管理员会话失效，需要重新登录。备份文件包含用户、节点、入站、出口 IP、账单、流量、事件和审计历史，请离线妥善保管。

## 6. 恢复

恢复前停止中央服务并确认备份来源：

```bash
sudo docker compose -p alien-panel -f /opt/alien-panel/deploy/docker-compose.yml down
sudo docker compose -p alien-panel -f /opt/alien-panel/deploy/docker-compose.yml run --rm --no-deps \
  -v /opt/alien-panel/release/xpanel-db-maintenance:/usr/local/bin/xpanel-db-maintenance:ro \
  -v /var/backups/alien-panel:/backups:ro \
  --entrypoint /usr/local/bin/xpanel-db-maintenance central restore \
  --source /backups/panel-YYYYMMDDTHHMMSSZ.sqlite3 \
  --database /var/lib/xpanel-central/panel.db --yes
sudo docker compose -p alien-panel -f /opt/alien-panel/deploy/docker-compose.yml up -d
sudo docker compose -p alien-panel -f /opt/alien-panel/deploy/docker-compose.yml up -d
```

恢复命令会先验证源文件，并在替换前创建安全副本。恢复后检查 `/health/ready`、Dashboard 的数据时间、节点数量、用户路径和最近同步记录。
