# Linux release bundle

`xpanel-central` is the Linux amd64 central-panel binary consumed by
`deploy/install-docker.sh`. The frontend bundle is kept in
`deploy/frontend-dist/` so a server does not need a Node.js toolchain.

Regenerate both artifacts from the project root with:

```powershell
powershell -ExecutionPolicy Bypass -File .\deploy\build-bundle.ps1
```

Commit the updated files before running the GitHub one-command installer.

The Agent's displayed version comes from `agent/VERSION`. Update it to the
next semantic version before rebuilding a release bundle.
