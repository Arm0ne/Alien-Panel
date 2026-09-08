# Alien-Panel v1.0.0 Linux 发布包

本目录包含 Linux amd64 的中央服务、Agent 和数据库维护工具。服务器不需要
安装 Go、Node.js 或 pnpm；一键安装脚本会使用 `xpanel-central` 和
`deploy/frontend-dist/` 启动 Docker Compose 部署。

在开发机从源码重新生成发布物：

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\build-bundle.ps1
```

提交更新后的发布物后，再执行 GitHub 一键安装命令。

Agent 显示版本来自 `agent/VERSION`。发布新版本前先更新语义版本号，
再重新构建并生成 `xpanel-agent.sha256`。
