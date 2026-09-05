# Central panel deployment

The complete node onboarding, day-two operations, fault handling, gray release
and rollback procedure is documented in [`OPERATIONS_RUNBOOK.md`](../OPERATIONS_RUNBOOK.md).
The `xpanel-central.service` unit in this directory is the systemd template
used by that runbook.

Nginx deployments can use [`nginx.conf`](./nginx.conf) as the reverse-proxy
template. The existing `Caddyfile` remains available for Caddy installations;
choose one proxy, do not run both on ports 80/443.

`Caddyfile` is a production-oriented HTTPS reverse-proxy example. Replace
`panel.example.com` and the frontend root with the actual deployment paths,
then run the central Go service on loopback (`127.0.0.1:8090`) so it is not
reachable directly from the Internet.

The proxy terminates TLS, forwards `/api`, `/agent` and health checks to the
central service, serves the built Vue frontend, and adds security headers. The
Go API repeats the headers for direct API responses. HSTS is therefore only
present on the public HTTPS endpoint; local HTTP development remains usable.

## One-command Docker test deployment

The repository now contains a tested Linux release bundle under `release/` and
the matching frontend files under `deploy/frontend-dist/`. This means a test
server does not need Node.js, pnpm or Go. Push the project (including those two
directories) to a public GitHub repository, then run this single command on an
Ubuntu/Debian server:

```sh
curl -fsSL https://raw.githubusercontent.com/OWNER/REPO/main/deploy/install-docker.sh \
  | sudo bash -s -- --repo https://github.com/OWNER/REPO.git --domain panel.example.com
```

The script installs Docker if needed, downloads the bundle, generates random
administrator and Agent bootstrap secrets, starts the central API and a small
Nginx container, and keeps the database in a named Docker volume. Without a
domain, test the generated port directly:

```text
http://<server-ip>:18080
```

Without `--domain`, the temporary port binds publicly for this direct test.
With `--domain`, it binds to `127.0.0.1` and is reachable through the host
Nginx only, so port 18080 does not need to be opened in the firewall.

Credentials are written to `/opt/xpanel-central/.env` with mode `0600`.
Updating is the same command; the existing `.env` and database are preserved.
To display the generated login values on the server, run
`sudo grep -E '^(XPANEL_ADMIN_USER|XPANEL_ADMIN_PASSWORD)=' /opt/xpanel-central/.env`.
To stop the containers while keeping data:

```sh
sudo bash /opt/xpanel-central/deploy/uninstall-docker.sh
```

To completely remove an old systemd/Docker test deployment, including its
database, backups, uploaded files and Docker volume, preview and then confirm:

```sh
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/deploy/uninstall-all.sh \
  | sudo bash
curl -fsSL https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/deploy/uninstall-all.sh \
  | sudo bash -s -- --yes
```

The cleanup does not touch TLS certificates, X-Panel Agent files on node
machines, or unrelated Nginx sites.

The bundled Nginx listens on the temporary test port only. For an existing
host Nginx, proxy the HTTPS server to `127.0.0.1:18080` and set `--domain` when
installing so browser write requests pass the backend Origin check:

```nginx
location / {
    proxy_pass http://127.0.0.1:18080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

When frontend or backend code changes, rebuild the bundle from the project
root before pushing:

```powershell
Set-Location frontend
pnpm.cmd build
Set-Location ..
Copy-Item -Path frontend/dist/* -Destination deploy/frontend-dist -Recurse -Force
$env:GOOS='linux'; $env:GOARCH='amd64'; $env:CGO_ENABLED='0'
go build -trimpath -o release/xpanel-central ./backend/cmd/server
Remove-Item Env:GOOS,Env:GOARCH,Env:CGO_ENABLED
```

Before exposing the service:

1. Set `XPANEL_ADMIN_PASSWORD`, `XPANEL_AGENT_REGISTRATION_TOKEN`, and
   `XPANEL_DATABASE` through the service manager rather than committing them.
2. Restrict the backend listener to loopback and allow only ports 80/443 in the
   firewall.
3. Verify `GET /health/live` and `GET /health/ready` through the HTTPS domain,
   and confirm the response headers before registering an Agent.

### Ubuntu/Debian one-line Agent installation

The node onboarding dialog generates a short-lived, single-use installer
command for Ubuntu/Debian amd64. It downloads the Agent from GitHub, exchanges
the installer token for the node credential, verifies the Agent checksum, asks
for the target node's local X-Panel credentials, and installs the systemd
service. The password is read from the target terminal and is never sent to
the central panel. The existing manual YAML installation remains available as
a fallback.

### Uploading the Vue frontend

Build on the local Windows workstation and upload to a temporary directory:

```powershell
Set-Location D:\轻量Panel\frontend
pnpm.cmd build
scp -P <SSH端口> -r .\dist <SSH用户>@<服务器IP>:/tmp/xpanel-frontend
```

On the server, locate `index.html` before copying. If it is directly under
`/tmp/xpanel-frontend`, use:

```sh
sudo install -d -m 0755 /opt/xpanel-central/frontend/dist
sudo cp -a /tmp/xpanel-frontend/. /opt/xpanel-central/frontend/dist/
```

If it is under `/tmp/xpanel-frontend/dist`, use that directory as the source
instead. Never use `/tmp/dist` unless that is the path reported by `find`.

The staging smoke checks automate these health and security-header assertions:

```sh
XPANEL_STAGING_URL=https://panel.example.com \
bash deploy/staging-smoke.sh
```

Set `XPANEL_ADMIN_PASSWORD` (and optionally `XPANEL_ADMIN_USER`) to also verify
administrator login and the protected Dashboard endpoint. The script never
prints the returned access token. On Windows, run
`deploy/staging-smoke.ps1` with the equivalent environment variables.

## Database backup and upgrades

Build and install `backend/cmd/db-maintenance` beside the central service:

```sh
cd backend
go build -o /usr/local/bin/xpanel-db-maintenance ./cmd/db-maintenance
```

The maintenance binary uses SQLite `VACUUM INTO`, so backups include committed
WAL data and do not require copying `panel.db-wal` or `panel.db-shm`. It verifies
both SQLite integrity and foreign-key consistency before reporting success.

Daily backup (14 snapshots by default):

```sh
XPANEL_DATABASE=/var/lib/xpanel-central/panel.db \
XPANEL_BACKUP_DIR=/var/backups/xpanel-central \
bash /path/to/project/deploy/backup.sh
```

`backup.sh` is suitable for a systemd timer or cron job. Set
`XPANEL_BACKUP_RETENTION` to keep between 14 and 30 snapshots (the command
accepts any positive value when used directly). Backups are mode `0600`.

Before an application upgrade, run `deploy/migrate.sh`. When systemd is
available it stops and restarts an active service automatically. It creates a
pre-migration snapshot, then runs all embedded migrations
idempotently. A failed migration leaves the pre-migration snapshot untouched.

To restore, first stop the service and use an explicit confirmation. The
restore command validates the source, creates a `panel-pre-restore-*.sqlite3`
safety snapshot of the current database, atomically installs the backup, and
validates the result:

```sh
XPANEL_DATABASE=/var/lib/xpanel-central/panel.db \
bash /path/to/project/deploy/restore.sh \
  --source /var/backups/xpanel-central/panel-20260903T010203Z.sqlite3 --yes
```

The scripts automatically stop and restart an active systemd service (default
`xpanel-central.service`). Set
`XPANEL_SKIP_SERVICE_STOP=1` when the service is managed elsewhere. On Windows,
the equivalent `backup.ps1`, `migrate.ps1` and `restore.ps1` wrappers call
`backend/db-maintenance.exe`; restore requires `-Yes`.
