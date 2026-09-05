# X-Panel Central 干净源码快照

本目录是当前项目的源码交接快照，生成日期为 2026-09-06。它包含中央后端、Agent、Vue 前端源码、数据库迁移、部署脚本和测试代码。

## 已排除内容

为保持“只含源码”的干净包，以下内容没有复制：

- 所有 `.git` 目录和 Git 历史；
- `frontend/node_modules`、`frontend/dist`、构建缓存和覆盖率目录；
- `release/` 下的二进制文件；
- `deploy/frontend-dist/` 下的生产构建文件；
- 数据库、WAL、备份、日志和临时文件；
- 真实运行环境的 `.env`、Agent 配置和任何本地凭据；快照中保留的前端 `.env` 文件和 `deploy/.env.example` 都是无秘密的构建模板；
- Windows/测试二进制（`.exe`、`.test`）和编辑器私有目录。

## 目录

```text
backend/   Go 中央服务、API、迁移和测试
agent/     Go Agent、X-Panel 采集和 systemd 安装
frontend/  Vue/TypeScript 前端源码和依赖锁文件
deploy/    Docker/Nginx/Caddy、构建、备份和冒烟脚本
```

## 从源码生成发布产物

在 Windows 开发机执行：

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\build-bundle.ps1 -GoArch amd64
```

生成后再使用项目根仓库的 Docker 部署命令。完整说明见上一级目录的 `DEPLOYMENT_GUIDE.md`。
