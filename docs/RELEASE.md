# 版本发布与回滚

## 当前版本

- 产品：Alien-Panel `v1.0.0`。
- Agent：`v1.0.3`。
- 目标：Linux amd64。
- 数据库迁移：`001`–`020`。
- 源码提交：`4bc4ade`。

## 发布前门禁

```text
[ ] 空库和已有数据库迁移通过
[ ] 后端 go test ./... 和 go vet ./... 通过
[ ] Agent go test ./... 通过
[ ] 前端 typecheck、lint、build 通过
[ ] 生成物为 Linux amd64 ELF
[ ] Agent SHA-256 与下载文件匹配
[ ] 一台线路机和一台落地机真实同步核对
[ ] 三种出口路径和多个固定出口 IP 已验证
[ ] 数据库备份和恢复演练完成
[ ] HTTPS、健康检查和 CORS 来源校验通过
```

## 构建发布包

在主源码根目录执行：

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\build-bundle.ps1 -GoArch amd64
```

脚本读取 `agent/VERSION`，将版本、提交号和构建时间写入 Agent；生成 `release/xpanel-agent.sha256`。构建后更新 `RELEASE_MANIFEST.json` 的 `sourceCommit` 和版本字段，再提交 `release/`、`deploy/frontend-dist/` 和本目录。

## 回滚

1. 先保存当前日志和数据库安全副本。
2. 将安装命令的 `--ref` 改为已验证的旧 tag/commit 对应发布目录。
3. 停止并重新启动容器，确认旧二进制和前端资源来自同一版本。
4. 如果迁移不可逆或应用无法启动，停止服务后使用同版本维护工具恢复备份。
5. 检查 `/health/ready`、节点心跳、用户列表、路径历史、财务记录和事件中心。

回滚不会自动把数据库结构降级。涉及新迁移的版本必须使用迁移前数据库备份恢复，再启动旧版本；仅替换二进制不能绕过已执行迁移。

## 变更记录

每次发布记录：产品版本、Agent 版本、源码提交、数据库迁移版本、构建时间、发布人、备份文件、灰度节点、验证结果和回滚点。不要把密码和 Token 写入变更记录。
