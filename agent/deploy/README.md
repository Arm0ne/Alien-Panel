# Agent deployment

Build a static Linux binary from the `agent/` directory and copy it to
`/usr/local/bin/xpanel-agent` on the node. Then run `install.sh` as root.

Create `/etc/xpanel-agent/agent.yaml` before starting the service. The
installer creates the dedicated `xpanel-agent` system user and sets the
configuration owner to `xpanel-agent:xpanel-agent` with mode `0600`, so the
agent can read the X-Panel password and central Token without exposing them to
other users. Runtime data stays outside the X-Panel installation.

The uninstall script removes the systemd unit and binary but intentionally
keeps `/etc/xpanel-agent` and `/var/lib/xpanel-agent` so a later reinstall can
reuse the configuration and local state.

The running agent reads the local X-Panel Session API and sends a heartbeat and
an idempotent snapshot to `/agent/v1/heartbeat` and `/agent/v1/sync` below the
configured `central_url`. Set `central_token` to the node token issued by the
central registration endpoint. The snapshot contains X-Panel's accumulated
Inbound and Client counters; it never calls an Xray reset or stats RPC.

The collector includes compatibility mappings for common X-Panel response
shapes: `obj` and `data.list` envelopes, object or string `settings`,
`clientStats`/`clientTraffic` counters, and flat or nested (`mem`/`disk`) server
resource metrics. Sanitized v1/v2 fixtures are kept under
`internal/collector/testdata/` and run with `go test ./...`.

Remote requests use a bounded retry policy for network errors, client timeouts,
HTTP 408/425/429 and 5xx responses. A logical request is attempted at most four
times with approximately 30 seconds, 2 minutes and 10 minutes of backoff;
context cancellation interrupts the wait. The central request ID and snapshot
`sync_id` remain stable across retries, so a timeout after the server committed
the sync is safely replayed as an idempotent request. X-Panel 401 responses
trigger one session refresh; invalid credentials still enter login backoff.
