# Agent deployment

Build a static Linux binary from the `agent/` directory and copy it to
`/usr/local/bin/xpanel-agent` on the node. Then run `install.sh` as root.

Create `/etc/xpanel-agent/agent.yaml` with mode `0600` before starting the
service. The installer creates the dedicated `xpanel-agent` system user and
keeps configuration and runtime data outside the X-Panel installation.

The uninstall script removes the systemd unit and binary but intentionally
keeps `/etc/xpanel-agent` and `/var/lib/xpanel-agent` so a later reinstall can
reuse the configuration and local state.

The running agent reads the local X-Panel Session API and sends a heartbeat and
an idempotent snapshot to `/agent/v1/heartbeat` and `/agent/v1/sync` below the
configured `central_url`. Set `central_token` to the node token issued by the
central registration endpoint. The snapshot contains X-Panel's accumulated
Inbound and Client counters; it never calls an Xray reset or stats RPC.
