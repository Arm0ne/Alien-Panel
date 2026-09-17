package httpapi

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type agentPrincipal struct {
	NodeID  string
	NodeKey string
}

type agentContextKey struct{}

type agentRegisterRequest struct {
	NodeKey       string `json:"node_key"`
	NodeName      string `json:"node_name"`
	NodeType      string `json:"node_type"`
	Hostname      string `json:"hostname"`
	PanelBasePath string `json:"panel_base_path"`
	AgentVersion  string `json:"agent_version"`
	XrayVersion   string `json:"xray_version"`
}

type agentBootstrapRequest struct {
	InstallToken string `json:"install_token"`
	Hostname     string `json:"hostname"`
	AgentVersion string `json:"agent_version"`
}

type agentStatusPayload struct {
	AgentVersion string  `json:"agent_version"`
	XrayRunning  bool    `json:"xray_running"`
	XrayVersion  string  `json:"xray_version"`
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsed   int64   `json:"memory_used"`
	MemoryTotal  int64   `json:"memory_total"`
	DiskUsed     int64   `json:"disk_used"`
	DiskTotal    int64   `json:"disk_total"`
}

type agentHeartbeatRequest struct {
	NodeKey         string             `json:"node_key"`
	ObservedAt      string             `json:"observed_at"`
	Status          agentStatusPayload `json:"status"`
	StatusAvailable *bool              `json:"status_available"`
	SyncState       string             `json:"sync_state"`
	SyncError       string             `json:"sync_error"`
}

type agentClientPayload struct {
	RemoteID   string `json:"remote_id"`
	Email      string `json:"email"`
	Enable     bool   `json:"enable"`
	ExpiryTime int64  `json:"expiry_time"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	AllTime    int64  `json:"all_time"`
	LastOnline int64  `json:"last_online"`
}

type agentInboundPayload struct {
	RemoteID        int64                `json:"remote_id"`
	Tag             string               `json:"tag"`
	Remark          string               `json:"remark"`
	Protocol        string               `json:"protocol"`
	Port            int64                `json:"port"`
	Listen          string               `json:"listen"`
	Enable          bool                 `json:"enable"`
	ExpiryTime      int64                `json:"expiry_time"`
	Up              int64                `json:"up"`
	Down            int64                `json:"down"`
	AllTime         int64                `json:"all_time"`
	ConfigHash      string               `json:"config_hash"`
	Clients         []agentClientPayload `json:"clients"`
	ClientsComplete bool                 `json:"clients_complete"`
}

type agentSyncRequest struct {
	NodeKey    string                `json:"node_key"`
	SyncID     string                `json:"sync_id"`
	ObservedAt string                `json:"observed_at"`
	Status     agentStatusPayload    `json:"status"`
	Inbounds   []agentInboundPayload `json:"inbounds"`
}

func (s *Server) agentRegister(w http.ResponseWriter, r *http.Request) {
	required := strings.TrimSpace(s.cfg.AgentRegistrationToken)
	if required == "" {
		writeFailure(w, http.StatusServiceUnavailable, internalErrorCode, "agent registration is disabled until XPANEL_AGENT_REGISTRATION_TOKEN is configured")
		return
	}
	if r.Header.Get("X-Agent-Registration-Token") != required {
		writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "agent registration is not authorized")
		return
	}
	var payload agentRegisterRequest
	if err := decodeAgentJSON(r, &payload); err != nil || strings.TrimSpace(payload.NodeKey) == "" || strings.TrimSpace(payload.NodeName) == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node_key and node_name are required")
		return
	}
	nodeType := strings.TrimSpace(payload.NodeType)
	if nodeType == "" {
		nodeType = "unknown"
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin node registration")
		return
	}
	defer tx.Rollback()

	var nodeID string
	var deletedAt string
	err = tx.QueryRow(`SELECT id, COALESCE(deleted_at, '') FROM nodes WHERE node_key = ?`, strings.TrimSpace(payload.NodeKey)).Scan(&nodeID, &deletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		nodeID = newID()
		_, err = tx.Exec(`INSERT INTO nodes (id, node_key, name, type, hostname, panel_base_path, agent_version, xray_version, enabled, health_status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, 'unknown', ?, ?)`, nodeID, strings.TrimSpace(payload.NodeKey), strings.TrimSpace(payload.NodeName), nodeType,
			strings.TrimSpace(payload.Hostname), strings.TrimSpace(payload.PanelBasePath), strings.TrimSpace(payload.AgentVersion), strings.TrimSpace(payload.XrayVersion), now, now)
	} else if err == nil {
		if deletedAt != "" {
			writeFailure(w, http.StatusConflict, validationCode, "node has been deleted; create a new node record before registering the Agent")
			return
		}
		// Registration refreshes Agent-owned metadata but preserves the central
		// administrator's node name. The name from the first registration is a
		// useful default; after that, PATCH /api/nodes/{id} is the authoritative
		// way to rename a node. Re-enable through the same PATCH endpoint when a
		// node has been disabled by an administrator.
		_, err = tx.Exec(`UPDATE nodes SET type = ?, hostname = ?, panel_base_path = ?, agent_version = ?, xray_version = ?, updated_at = ? WHERE id = ?`,
			nodeType, strings.TrimSpace(payload.Hostname), strings.TrimSpace(payload.PanelBasePath), strings.TrimSpace(payload.AgentVersion), strings.TrimSpace(payload.XrayVersion), now, nodeID)
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save node registration")
		return
	}
	if _, err := tx.Exec(`UPDATE node_credentials SET revoked_at = ? WHERE node_id = ? AND revoked_at IS NULL`, now, nodeID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not rotate node credentials")
		return
	}
	token, err := randomToken()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not issue node credential")
		return
	}
	if _, err := tx.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES (?, ?, ?, ?, ?)`, newID(), nodeID, hashToken(token), now, now); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save node credential")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not complete node registration")
		return
	}
	writeSuccess(w, map[string]any{"nodeId": nodeID, "nodeKey": payload.NodeKey, "token": token})
}

// agentBootstrap exchanges a short-lived installer token for the normal node
// bearer token. The installer token is bound to a pre-created node, consumed
// atomically, and never stored in plaintext.
func (s *Server) agentBootstrap(w http.ResponseWriter, r *http.Request) {
	var payload agentBootstrapRequest
	if err := decodeAgentJSON(r, &payload); err != nil || strings.TrimSpace(payload.InstallToken) == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "install_token is required")
		return
	}
	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	tx, err := s.db.BeginTx(r.Context(), nil)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin Agent bootstrap")
		return
	}
	defer tx.Rollback()

	var nodeID, nodeKey, nodeName, nodeType string
	err = tx.QueryRowContext(r.Context(), `SELECT n.id, n.node_key, n.name, n.type
FROM node_install_tokens it JOIN nodes n ON n.id = it.node_id
WHERE it.token_hash = ? AND it.used_at IS NULL AND it.expires_at > ?
  AND n.enabled = 1 AND n.deleted_at IS NULL`, hashToken(strings.TrimSpace(payload.InstallToken)), nowText).
		Scan(&nodeID, &nodeKey, &nodeName, &nodeType)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "installer token is invalid, expired, or already used")
		return
	}
	if err != nil {
		s.logger.Error("read Agent installer token", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read Agent installer token")
		return
	}
	if _, err := tx.ExecContext(r.Context(), `UPDATE node_install_tokens SET used_at = ? WHERE token_hash = ? AND used_at IS NULL`, nowText, hashToken(strings.TrimSpace(payload.InstallToken))); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not consume Agent installer token")
		return
	}
	if _, err := tx.ExecContext(r.Context(), `UPDATE node_credentials SET revoked_at = ? WHERE node_id = ? AND revoked_at IS NULL`, nowText, nodeID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not rotate node credentials")
		return
	}
	centralToken, err := randomToken()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not issue node credential")
		return
	}
	if _, err := tx.ExecContext(r.Context(), `INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES (?, ?, ?, ?, ?)`, newID(), nodeID, hashToken(centralToken), nowText, nowText); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save node credential")
		return
	}
	if _, err := tx.ExecContext(r.Context(), `UPDATE nodes SET hostname = CASE WHEN ? <> '' THEN ? ELSE hostname END,
agent_version = CASE WHEN ? <> '' THEN ? ELSE agent_version END, updated_at = ? WHERE id = ?`,
		strings.TrimSpace(payload.Hostname), strings.TrimSpace(payload.Hostname), strings.TrimSpace(payload.AgentVersion), strings.TrimSpace(payload.AgentVersion), nowText, nodeID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update node bootstrap metadata")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not complete Agent bootstrap")
		return
	}
	writeSuccess(w, map[string]any{
		"node_id": nodeID, "node_key": nodeKey, "node_name": nodeName, "node_type": nodeType,
		"central_token": centralToken,
	})
}

func (s *Server) agentHeartbeat(w http.ResponseWriter, r *http.Request) {
	principal := r.Context().Value(agentContextKey{}).(agentPrincipal)
	var payload agentHeartbeatRequest
	if err := decodeAgentJSON(r, &payload); err != nil || strings.TrimSpace(payload.NodeKey) == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node_key is required")
		return
	}
	if payload.NodeKey != principal.NodeKey {
		writeFailure(w, http.StatusForbidden, unauthorizedCode, "node credential does not match node_key")
		return
	}
	observedAt, err := parseObservedAt(payload.ObservedAt)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "observed_at must be RFC3339")
		return
	}
	syncState := strings.TrimSpace(payload.SyncState)
	if syncState != "" && syncState != "failed" {
		writeFailure(w, http.StatusBadRequest, validationCode, "sync_state must be failed when present")
		return
	}
	// Older Agents did not send status_available and always included an
	// X-Panel snapshot in their heartbeat. Keep that protocol working while
	// allowing newer Agents to send a liveness-only heartbeat.
	statusAvailable := payload.StatusAvailable == nil || *payload.StatusAvailable
	if err := s.updateNodeHeartbeat(principal.NodeID, observedAt, payload.Status, statusAvailable, syncState, payload.SyncError); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update node heartbeat")
		return
	}
	writeSuccess(w, map[string]any{"accepted": true, "nodeId": principal.NodeID, "observedAt": observedAt.Format(time.RFC3339Nano)})
}

func (s *Server) agentSync(w http.ResponseWriter, r *http.Request) {
	principal := r.Context().Value(agentContextKey{}).(agentPrincipal)
	var payload agentSyncRequest
	if err := decodeAgentJSON(r, &payload); err != nil || strings.TrimSpace(payload.NodeKey) == "" || strings.TrimSpace(payload.SyncID) == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node_key and sync_id are required")
		return
	}
	if payload.NodeKey != principal.NodeKey {
		writeFailure(w, http.StatusForbidden, unauthorizedCode, "node credential does not match node_key")
		return
	}
	observedAt, err := parseObservedAt(payload.ObservedAt)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "observed_at must be RFC3339")
		return
	}

	var existingStatus string
	var existingInboundCount, existingClientCount int
	err = s.db.QueryRow(`SELECT status, inbound_count, client_count FROM sync_runs WHERE sync_id = ? AND node_id = ?`, payload.SyncID, principal.NodeID).Scan(&existingStatus, &existingInboundCount, &existingClientCount)
	if err == nil {
		writeSuccess(w, map[string]any{"sync_id": payload.SyncID, "status": existingStatus, "inboundCount": existingInboundCount, "clientCount": existingClientCount, "idempotent": true})
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check sync id")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin sync")
		return
	}
	defer tx.Rollback()
	syncRunID := newID()
	if _, err := tx.Exec(`INSERT INTO sync_runs (id, node_id, sync_id, started_at, status) VALUES (?, ?, ?, ?, 'running')`, syncRunID, principal.NodeID, payload.SyncID, now); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			writeSuccess(w, map[string]any{"sync_id": payload.SyncID, "status": "running", "idempotent": true})
			return
		}
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create sync run")
		return
	}
	var nodeType string
	if err := tx.QueryRow(`SELECT type FROM nodes WHERE id = ?`, principal.NodeID).Scan(&nodeType); err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("read node type: %w", err))
		return
	}

	clientCount := 0
	for _, inbound := range payload.Inbounds {
		remoteInboundID := strconv.FormatInt(inbound.RemoteID, 10)
		incomingClientIDs := normalizeClientIDs(inbound.Clients)
		expiryText := ""
		var expiry any
		if inbound.ExpiryTime > 0 {
			expiryText = time.Unix(inbound.ExpiryTime, 0).UTC().Format(time.RFC3339Nano)
			expiry = expiryText
		}
		_, err := tx.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, tag, remark, protocol, port, listen, enable, expiry_time, up, down, all_time, client_count, config_hash, first_seen_at, last_seen_at, missing_since, missing_sync_count, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, 0, NULL)
ON CONFLICT(node_id, remote_inbound_id) DO UPDATE SET tag = excluded.tag, remark = excluded.remark, protocol = excluded.protocol, port = excluded.port, listen = excluded.listen, enable = excluded.enable, expiry_time = excluded.expiry_time, up = excluded.up, down = excluded.down, all_time = excluded.all_time, client_count = excluded.client_count, config_hash = excluded.config_hash, last_seen_at = excluded.last_seen_at, missing_since = NULL, missing_sync_count = 0, deleted_at = NULL`,
			newID(), principal.NodeID, remoteInboundID, inbound.Tag, inbound.Remark, inbound.Protocol, inbound.Port, inbound.Listen, boolInt(inbound.Enable), expiry, inbound.Up, inbound.Down, inbound.AllTime, len(incomingClientIDs), inbound.ConfigHash, observedAt.Format(time.RFC3339Nano), observedAt.Format(time.RFC3339Nano))
		if err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("upsert inbound %s: %w", remoteInboundID, err))
			return
		}
		var inboundID string
		if err := tx.QueryRow(`SELECT id FROM inbounds WHERE node_id = ? AND remote_inbound_id = ?`, principal.NodeID, remoteInboundID).Scan(&inboundID); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("find inbound %s: %w", remoteInboundID, err))
			return
		}
		var inboundUserID, inboundKind, inboundUserDeletedAt, storedClientSet string
		if err := tx.QueryRow(`SELECT COALESCE(i.user_id, ''), COALESCE(i.kind, 'unknown'), COALESCE(u.deleted_at, ''), COALESCE(i.client_set_json, '[]')
FROM inbounds i LEFT JOIN users u ON u.id = i.user_id WHERE i.id = ?`, inboundID).
			Scan(&inboundUserID, &inboundKind, &inboundUserDeletedAt, &storedClientSet); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("read inbound client state %s: %w", remoteInboundID, err))
			return
		}
		previousClientIDs, err := decodeClientSet(storedClientSet)
		if err != nil {
			// A malformed value should never make a sync fail.  Rebuild the
			// baseline from the live rows and let this sync repair the JSON.
			previousClientIDs = nil
		}
		existingClientIDs, err := readInboundClientIDsTx(tx, inboundID)
		if err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("read existing clients %s: %w", remoteInboundID, err))
			return
		}
		if len(previousClientIDs) == 0 && len(existingClientIDs) > 0 {
			previousClientIDs = existingClientIDs
		}
		// A deleted central user means this Inbound belonged to a customer that
		// was already retired (usually after the Inbound was archived).  Its
		// return must start a new customer lifecycle instead of reusing the old
		// soft-deleted user.  Do not create a second replacement-confirmation
		// event for that case; the archive/retirement boundary already confirms
		// that the old operational state is no longer current.
		deletedUserReappeared := inboundUserID != "" && inboundUserDeletedAt != ""
		replacementDetected := inbound.ClientsComplete && nodeType == "relay" && inboundUserID != "" &&
			!deletedUserReappeared &&
			inboundKind != "infrastructure" && len(previousClientIDs) > 0 && len(incomingClientIDs) > 0 &&
			clientSetsDisjoint(previousClientIDs, incomingClientIDs)
		resetDetected, err := detectTrafficReset(tx, inboundID, observedAt, inbound.AllTime)
		if err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("check traffic reset %s: %w", remoteInboundID, err))
			return
		}
		if resetDetected {
			if _, err := tx.Exec(`UPDATE inbounds SET traffic_baseline_up = ?, traffic_baseline_down = ?, traffic_baseline_all_time = ?, traffic_baseline_at = ? WHERE id = ?`, inbound.Up, inbound.Down, inbound.AllTime, observedAt.Format(time.RFC3339Nano), inboundID); err != nil {
				s.failSync(w, tx, syncRunID, fmt.Errorf("reset inbound traffic baseline %s: %w", remoteInboundID, err))
				return
			}
			err := insertNodeEventTx(tx, nodeEventSpec{
				NodeID: principal.NodeID, EventType: "traffic_reset", Category: "node", Severity: "warning",
				Title: "流量累计值发生回退", Message: "Inbound " + remoteInboundID + " 累计流量回退，已建立新基线",
				ResourceType: "inbound", ResourceID: inboundID, ActionType: "inspect_traffic", Source: "agent", CorrelationID: syncRunID, OccurredAt: observedAt,
			})
			if err != nil {
				s.failSync(w, tx, syncRunID, fmt.Errorf("record traffic reset %s: %w", remoteInboundID, err))
				return
			}
		}
		for _, client := range inbound.Clients {
			if strings.TrimSpace(client.RemoteID) == "" {
				continue
			}
			clientCount++
			_, err := tx.Exec(`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, expiry_time, up, down, all_time, last_online, last_seen_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(node_id, inbound_id, remote_client_id) DO UPDATE SET email = excluded.email, enable = excluded.enable, expiry_time = excluded.expiry_time, up = excluded.up, down = excluded.down, all_time = excluded.all_time, last_online = excluded.last_online, last_seen_at = excluded.last_seen_at`,
				newID(), principal.NodeID, inboundID, client.RemoteID, client.Email, boolInt(client.Enable), epochString(client.ExpiryTime), client.Up, client.Down, client.AllTime, epochString(client.LastOnline), observedAt.Format(time.RFC3339Nano))
			if err != nil {
				s.failSync(w, tx, syncRunID, fmt.Errorf("upsert client %s: %w", client.RemoteID, err))
				return
			}
		}
		if inbound.ClientsComplete {
			if _, err := deleteStaleClientsTx(tx, inboundID, incomingClientIDs); err != nil {
				s.failSync(w, tx, syncRunID, fmt.Errorf("remove deleted clients %s: %w", remoteInboundID, err))
				return
			}
			clientSet := incomingClientIDs
			if len(clientSet) == 0 {
				// Keep the last non-empty identity set through an empty sync so a
				// later replacement can still be detected after the old rows were
				// purged.  This also repairs pre-024 databases on first use.
				clientSet = previousClientIDs
			}
			encoded, encodeErr := encodeClientSet(clientSet)
			if encodeErr != nil {
				s.failSync(w, tx, syncRunID, fmt.Errorf("encode client state %s: %w", remoteInboundID, encodeErr))
				return
			}
			if _, err := tx.Exec(`UPDATE inbounds SET client_set_json = ? WHERE id = ?`, encoded, inboundID); err != nil {
				s.failSync(w, tx, syncRunID, fmt.Errorf("save client state %s: %w", remoteInboundID, err))
				return
			}
		}
		if err := s.ensureRelayInboundUserWithOptions(tx, nodeType, inboundID, remoteInboundID, inbound, observedAt, replacementDetected); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("ensure business user for inbound %s: %w", remoteInboundID, err))
			return
		}
		if replacementDetected {
			if err := insertClientReplacementEventTx(tx, principal.NodeID, inboundID, inboundUserID, previousClientIDs, incomingClientIDs, observedAt); err != nil {
				s.failSync(w, tx, syncRunID, fmt.Errorf("record client replacement %s: %w", remoteInboundID, err))
				return
			}
		}
		resetFlag := 0
		if resetDetected {
			resetFlag = 1
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected, sync_run_id) VALUES (?, ?, ?, ?, ?, ?, ?, 'xpanel', ?, ?)`, newID(), principal.NodeID, inboundID, observedAt.Format(time.RFC3339Nano), inbound.Up, inbound.Down, inbound.AllTime, resetFlag, syncRunID); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("save traffic snapshot %s: %w", remoteInboundID, err))
			return
		}
	}
	err = s.markMissingAndArchiveInbounds(tx, principal.NodeID, observedAt)
	if err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("mark or archive missing inbounds: %w", err))
		return
	}
	orphanedUserIDs, err := orphanedUserIDsTx(tx)
	if err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("find orphaned users after sync: %w", err))
		return
	}
	if len(orphanedUserIDs) > 0 {
		var nodeName string
		_ = tx.QueryRow(`SELECT name FROM nodes WHERE id = ?`, principal.NodeID).Scan(&nodeName)
		if strings.TrimSpace(nodeName) == "" {
			nodeName = principal.NodeKey
		}
		if err := s.autoDeleteOrphanedUsersTx(tx, r, orphanedUserIDs, principal.NodeID, nodeName, time.Now().UTC()); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("clean up orphaned users after sync: %w", err))
			return
		}
	}
	var previousHealth, previousSyncState string
	_ = tx.QueryRow(`SELECT health_status, sync_status FROM nodes WHERE id = ?`, principal.NodeID).Scan(&previousHealth, &previousSyncState)
	newHealth := healthStatus(payload.Status.XrayRunning)
	if _, err := tx.Exec(`UPDATE nodes SET health_status = ?, agent_version = CASE WHEN ? <> '' THEN ? ELSE agent_version END, xray_version = CASE WHEN ? <> '' THEN ? ELSE xray_version END,
cpu_usage = ?, memory_used = ?, memory_total = ?, disk_used = ?, disk_total = ?, sync_status = 'success', last_sync_error = NULL, last_seen_at = ?, updated_at = ? WHERE id = ?`, healthStatus(payload.Status.XrayRunning), strings.TrimSpace(payload.Status.AgentVersion), strings.TrimSpace(payload.Status.AgentVersion), payload.Status.XrayVersion, payload.Status.XrayVersion,
		payload.Status.CPUUsage, nullableMetric(payload.Status.MemoryUsed), nullableMetric(payload.Status.MemoryTotal), nullableMetric(payload.Status.DiskUsed), nullableMetric(payload.Status.DiskTotal), observedAt.Format(time.RFC3339Nano), now, principal.NodeID); err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("update node after sync: %w", err))
		return
	}
	if previousHealth == "offline" && newHealth != "offline" {
		if err := insertNodeEventTx(tx, nodeEventSpec{
			NodeID: principal.NodeID, EventType: "node_recovered", Category: "node", Severity: "info", Title: "节点已恢复在线",
			Message: "节点已重新完成 Agent 同步", ResourceType: "node", ResourceID: principal.NodeID, Source: "agent", CorrelationID: syncRunID, OccurredAt: observedAt,
		}); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("record node recovery: %w", err))
			return
		}
	}
	if previousSyncState == "failed" {
		if _, err := tx.Exec(`UPDATE node_events SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?), dedupe_key = NULL
WHERE dedupe_key = ? AND event_status NOT IN ('resolved', 'dismissed')`, now, now, "xpanel-sync-failed:"+principal.NodeID); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("resolve X-Panel sync failure event: %w", err))
			return
		}
		if err := insertNodeEventTx(tx, nodeEventSpec{
			NodeID: principal.NodeID, EventType: "xpanel_sync_recovered", Category: "sync", Severity: "info", Title: "X-Panel 采集已恢复",
			Message: "Agent 已重新完成 X-Panel 数据同步", ResourceType: "node", ResourceID: principal.NodeID,
			ActionType: "inspect_sync_run", Source: "agent", CorrelationID: syncRunID, OccurredAt: observedAt,
		}); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("record X-Panel sync recovery: %w", err))
			return
		}
	}
	if _, err := tx.Exec(`UPDATE sync_runs SET finished_at = ?, status = 'success', inbound_count = ?, client_count = ? WHERE id = ?`, now, len(payload.Inbounds), clientCount, syncRunID); err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("finish sync run: %w", err))
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not commit sync")
		return
	}
	writeSuccess(w, map[string]any{"sync_id": payload.SyncID, "status": "success", "inboundCount": len(payload.Inbounds), "clientCount": clientCount, "idempotent": false})
}

func readInboundClientIDsTx(tx *sql.Tx, inboundID string) ([]string, error) {
	rows, err := tx.Query(`SELECT remote_client_id FROM clients WHERE inbound_id = ? AND TRIM(remote_client_id) <> '' ORDER BY remote_client_id`, inboundID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return normalizeClientIDStrings(ids), nil
}

func normalizeClientIDs(clients []agentClientPayload) []string {
	ids := make([]string, 0, len(clients))
	for _, client := range clients {
		if id := strings.TrimSpace(client.RemoteID); id != "" {
			ids = append(ids, id)
		}
	}
	return normalizeClientIDStrings(ids)
}

func normalizeClientIDStrings(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	result := make([]string, 0, len(ids))
	for _, value := range ids {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func decodeClientSet(value string) ([]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	var ids []string
	if err := json.Unmarshal([]byte(value), &ids); err != nil {
		return nil, err
	}
	return normalizeClientIDStrings(ids), nil
}

func encodeClientSet(ids []string) (string, error) {
	encoded, err := json.Marshal(normalizeClientIDStrings(ids))
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func clientSetsDisjoint(left, right []string) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(left))
	for _, value := range left {
		seen[value] = struct{}{}
	}
	for _, value := range right {
		if _, exists := seen[value]; exists {
			return false
		}
	}
	return true
}

func deleteStaleClientsTx(tx *sql.Tx, inboundID string, incomingIDs []string) (int64, error) {
	incoming := make(map[string]struct{}, len(incomingIDs))
	for _, id := range incomingIDs {
		incoming[id] = struct{}{}
	}
	rows, err := tx.Query(`SELECT id, remote_client_id FROM clients WHERE inbound_id = ?`, inboundID)
	if err != nil {
		return 0, err
	}
	staleIDs := make([]string, 0)
	for rows.Next() {
		var id, remoteID string
		if err := rows.Scan(&id, &remoteID); err != nil {
			return 0, err
		}
		if _, exists := incoming[strings.TrimSpace(remoteID)]; !exists {
			staleIDs = append(staleIDs, id)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}
	var removed int64
	for _, id := range staleIDs {
		result, err := tx.Exec(`DELETE FROM clients WHERE id = ?`, id)
		if err != nil {
			return removed, err
		}
		if count, err := result.RowsAffected(); err == nil {
			removed += count
		}
	}
	return removed, nil
}

func clientSetHash(ids []string) string {
	encoded, _ := encodeClientSet(ids)
	hash := sha256.Sum256([]byte(encoded))
	return hex.EncodeToString(hash[:])
}

func insertClientReplacementEventTx(tx *sql.Tx, nodeID, inboundID, userID string, oldClientIDs, newClientIDs []string, observedAt time.Time) error {
	var inboundName, userName string
	if err := tx.QueryRow(`SELECT COALESCE(NULLIF(i.remark, ''), NULLIF(i.tag, ''), i.remote_inbound_id), COALESCE(u.display_name, '')
FROM inbounds i LEFT JOIN users u ON u.id = i.user_id WHERE i.id = ?`, inboundID).Scan(&inboundName, &userName); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("read replacement event context: %w", err)
	}
	if userName == "" {
		userName = userID
	}
	if inboundName == "" {
		inboundName = inboundID
	}
	payload := map[string]any{
		"inboundId": inboundID, "userId": userID, "inboundName": inboundName,
		"userName":       userName,
		"oldClientCount": len(oldClientIDs), "newClientCount": len(newClientIDs),
		"oldClientSetHash": clientSetHash(oldClientIDs), "clientSetHash": clientSetHash(newClientIDs),
	}
	return insertNodeEventTx(tx, nodeEventSpec{
		NodeID: nodeID, EventType: "client_set_replacement_detected", Category: "business", Severity: "warning",
		Title: "检测到客户更换，待确认", Message: fmt.Sprintf("Inbound「%s」的 Client 已全部更换。确认后将关闭旧用户的路径和流量，仅保留账务记录，并创建新用户。", inboundName),
		RequiresAction: true, EventStatus: "open", ResourceType: "inbound", ResourceID: inboundID, ActionType: "reset_user",
		Payload: payload, DedupeKey: "client-replacement:" + inboundID + ":" + clientSetHash(newClientIDs), Source: "agent", CorrelationID: inboundID, OccurredAt: observedAt,
	})
}

// ensureRelayInboundUser keeps the business rule that one relay-node Inbound
// is one central business user. The Inbound identifies the user and owns its
// route and traffic, while its enabled Clients determine the user's expiry.
// This deliberately does not merge equal Email values across Inbounds/nodes.
func (s *Server) ensureRelayInboundUser(tx *sql.Tx, nodeType, inboundID, remoteInboundID string, inbound agentInboundPayload, observedAt time.Time) error {
	return s.ensureRelayInboundUserWithOptions(tx, nodeType, inboundID, remoteInboundID, inbound, observedAt, false)
}

func (s *Server) ensureRelayInboundUserWithOptions(tx *sql.Tx, nodeType, inboundID, remoteInboundID string, inbound agentInboundPayload, observedAt time.Time, suppressRenewalCandidate bool) error {
	if nodeType != "relay" {
		// A landing/unknown node may have been synchronized by an older
		// version that incorrectly attached its Inbound to a business user.
		// Remove only that invalid association; the user record itself is
		// central data and must remain available for later reuse.
		if _, err := tx.Exec(`DELETE FROM user_inbounds WHERE inbound_id = ?`, inboundID); err != nil {
			return fmt.Errorf("clear non-relay business mapping: %w", err)
		}
		if _, err := tx.Exec(`UPDATE inbounds SET user_id = NULL, kind = 'infrastructure' WHERE id = ?`, inboundID); err != nil {
			return fmt.Errorf("classify non-relay inbound: %w", err)
		}
		return nil
	}

	var userID, kind, userDeletedAt string
	if err := tx.QueryRow(`SELECT COALESCE(i.user_id, ''), COALESCE(i.kind, 'unknown'), COALESCE(u.deleted_at, '')
FROM inbounds i LEFT JOIN users u ON u.id = i.user_id WHERE i.id = ?`, inboundID).Scan(&userID, &kind, &userDeletedAt); err != nil {
		return fmt.Errorf("read inbound classification: %w", err)
	}
	if kind == "infrastructure" {
		return nil
	}
	if kind == "unknown" {
		if _, err := tx.Exec(`UPDATE inbounds SET kind = 'user' WHERE id = ?`, inboundID); err != nil {
			return fmt.Errorf("classify relay inbound as user: %w", err)
		}
	}
	if userID != "" && userDeletedAt != "" {
		newUserID, err := s.replaceDeletedInboundUserTx(tx, inboundID, userID, inbound, remoteInboundID, observedAt)
		if err != nil {
			return err
		}
		userID = newUserID
	}

	now := observedAt.UTC().Format(time.RFC3339Nano)
	state := userStateFromClients(inbound, observedAt)
	expiry := nullableDBString(state.ExpiryText)
	newUserCreated := false
	if userID == "" {
		userID = newID()
		newUserCreated = true
		if _, err := tx.Exec(`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`, userID, inboundDisplayName(inbound, remoteInboundID), state.Status, expiry, now, now); err != nil {
			return fmt.Errorf("create business user: %w", err)
		}
		if _, err := tx.Exec(`UPDATE inbounds SET user_id = ?, kind = 'user' WHERE id = ?`, userID, inboundID); err != nil {
			return fmt.Errorf("link inbound to business user: %w", err)
		}
	} else {
		// Read the previous central snapshot before replacing it.  An expiry
		// extension is only a renewal suggestion; it is not automatically
		// treated as a payment because operators may grant time for free.
		var oldExpiry, billingCycle string
		var billingAmount, monthlyFee float64
		if err := tx.QueryRow(`SELECT COALESCE(expiry_time, ''), COALESCE(billing_cycle, 'monthly'), COALESCE(billing_amount, 0), COALESCE(monthly_fee, 0) FROM users WHERE id = ?`, userID).
			Scan(&oldExpiry, &billingCycle, &billingAmount, &monthlyFee); err != nil {
			return fmt.Errorf("read previous user billing state: %w", err)
		}
		if !suppressRenewalCandidate && oldExpiry != "" && state.ExpiryText != "" {
			oldTime, oldErr := time.Parse(time.RFC3339Nano, oldExpiry)
			newTime, newErr := time.Parse(time.RFC3339Nano, state.ExpiryText)
			if oldErr == nil && newErr == nil && newTime.After(oldTime) {
				if billingCycle != "annual" {
					billingCycle = "monthly"
				}
				suggestedAmount := billingAmount
				if suggestedAmount <= 0 {
					suggestedAmount = monthlyFee
					if billingCycle == "annual" {
						suggestedAmount *= 12
					}
				}
				if _, err := tx.Exec(`INSERT OR IGNORE INTO user_renewal_candidates
(id, user_id, inbound_id, old_expiry_at, new_expiry_at, detected_at, suggested_cycle, suggested_amount, currency, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'CNY', 'pending')`, newID(), userID, inboundID, oldExpiry, state.ExpiryText, now, billingCycle, suggestedAmount); err != nil {
					return fmt.Errorf("record renewal candidate: %w", err)
				}
				var candidateID string
				if err := tx.QueryRow(`SELECT id FROM user_renewal_candidates WHERE user_id = ? AND inbound_id = ? AND new_expiry_at = ?`, userID, inboundID, state.ExpiryText).Scan(&candidateID); err != nil {
					return fmt.Errorf("read renewal candidate: %w", err)
				}
				payload := map[string]any{
					"candidateId": candidateID, "userId": userID, "inboundId": inboundID,
					"oldExpiryAt": oldExpiry, "newExpiryAt": state.ExpiryText,
					"billingCycle": billingCycle, "suggestedAmount": suggestedAmount, "currency": "CNY",
				}
				var eventUserName string
				_ = tx.QueryRow(`SELECT COALESCE(display_name, '') FROM users WHERE id = ?`, userID).Scan(&eventUserName)
				if eventUserName == "" {
					eventUserName = userID
				}
				payload["userName"] = eventUserName
				var eventNodeID string
				_ = tx.QueryRow(`SELECT node_id FROM inbounds WHERE id = ?`, inboundID).Scan(&eventNodeID)
				if err := insertNodeEventTx(tx, nodeEventSpec{
					NodeID: eventNodeID, EventType: "renewal_candidate_detected", Category: "business", Severity: "warning",
					Title: "检测到续费变更，待确认", Message: "用户「" + eventUserName + "」到期时间由 " + oldExpiry + " 延长至 " + state.ExpiryText + "，请确认是否计入收费",
					RequiresAction: true, EventStatus: "open", ResourceType: "renewal", ResourceID: candidateID,
					ActionType: "confirm_renewal", Payload: payload, DedupeKey: "renewal-candidate:" + candidateID,
					Source: "agent", CorrelationID: candidateID, OccurredAt: observedAt,
				}); err != nil {
					return fmt.Errorf("record renewal event: %w", err)
				}
			}
		}
		if _, err := tx.Exec(`UPDATE users
SET expiry_time = ?, status = ?, updated_at = ?
WHERE id = ?`, expiry, state.Status, now, userID); err != nil {
			return fmt.Errorf("refresh business user state: %w", err)
		}
	}
	if _, err := tx.Exec(`INSERT OR IGNORE INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES (?, ?, ?, 1, ?)`, newID(), userID, inboundID, now); err != nil {
		return fmt.Errorf("save business user mapping: %w", err)
	}
	if newUserCreated {
		if err := insertNewUserProfileEventTx(tx, inboundID, userID, inboundDisplayName(inbound, remoteInboundID), observedAt); err != nil {
			return fmt.Errorf("record new user profile event: %w", err)
		}
	}
	if err := reconcileClientExpiryMismatchEvent(tx, inboundID, userID, inboundDisplayName(inbound, remoteInboundID), state, observedAt); err != nil {
		return err
	}
	return nil
}

// replaceDeletedInboundUserTx starts a new operational customer lifecycle
// when a previously archived Inbound reappears while still pointing at a
// soft-deleted central user.  Billing rows remain attached to the old user;
// paths, routes and traffic belong to the retired customer and are closed or
// removed before the new user is linked by ensureRelayInboundUser.
func (s *Server) replaceDeletedInboundUserTx(tx *sql.Tx, inboundID, oldUserID string, inbound agentInboundPayload, remoteInboundID string, observedAt time.Time) (string, error) {
	nowText := observedAt.UTC().Format(time.RFC3339Nano)

	if _, err := tx.Exec(`UPDATE user_inbounds SET active_to = ? WHERE user_id = ? AND inbound_id = ? AND active_to IS NULL`, nowText, oldUserID, inboundID); err != nil {
		return "", fmt.Errorf("close deleted user inbound mapping: %w", err)
	}

	var remainingAssociations int
	if err := tx.QueryRow(`SELECT COUNT(*)
FROM user_inbounds ui JOIN inbounds i ON i.id = ui.inbound_id
JOIN nodes n ON n.id = i.node_id
WHERE ui.user_id = ? AND ui.active_to IS NULL AND i.id <> ? AND i.deleted_at IS NULL
  AND i.kind = 'user' AND n.type = 'relay' AND n.deleted_at IS NULL`, oldUserID, inboundID).Scan(&remainingAssociations); err != nil {
		return "", fmt.Errorf("count deleted user inbound mappings: %w", err)
	}
	if remainingAssociations == 0 {
		if _, err := tx.Exec(`UPDATE user_paths SET active_to = ?, updated_at = ? WHERE user_id = ? AND active_to IS NULL`, nowText, nowText, oldUserID); err != nil {
			return "", fmt.Errorf("close deleted user paths: %w", err)
		}
		if _, err := tx.Exec(`UPDATE user_routes SET is_primary = 0, active_to = ?, route_exit_ip_id = NULL WHERE user_id = ? AND active_to IS NULL`, nowText, oldUserID); err != nil {
			return "", fmt.Errorf("close deleted user routes: %w", err)
		}
	}
	if err := closeRenewalCandidatesForInboundTx(tx, oldUserID, inboundID, nowText); err != nil {
		return "", err
	}
	if _, err := deleteInboundTrafficTx(tx, inboundID); err != nil {
		return "", fmt.Errorf("purge deleted user traffic: %w", err)
	}
	if _, err := tx.Exec(`UPDATE inbounds
SET traffic_baseline_up = ?, traffic_baseline_down = ?, traffic_baseline_all_time = ?, traffic_baseline_at = ?
WHERE id = ?`, inbound.Up, inbound.Down, inbound.AllTime, nowText, inboundID); err != nil {
		return "", fmt.Errorf("save reactivated inbound traffic baseline: %w", err)
	}

	state := userStateFromClients(inbound, observedAt)
	newUserID := newID()
	if _, err := tx.Exec(`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`, newUserID, inboundDisplayName(inbound, remoteInboundID), state.Status, nullableDBString(state.ExpiryText), nowText, nowText); err != nil {
		return "", fmt.Errorf("create user for reactivated inbound: %w", err)
	}
	if _, err := tx.Exec(`UPDATE inbounds SET user_id = ?, kind = 'user' WHERE id = ?`, newUserID, inboundID); err != nil {
		return "", fmt.Errorf("link user for reactivated inbound: %w", err)
	}
	if err := insertNewUserProfileEventTx(tx, inboundID, newUserID, inboundDisplayName(inbound, remoteInboundID), observedAt); err != nil {
		return "", fmt.Errorf("record reactivated user profile event: %w", err)
	}

	if err := s.writeAuditLogTx(tx, nil, "user.replace_from_deleted_inbound", "user", oldUserID,
		map[string]any{"inboundId": inboundID, "displayName": inboundDisplayName(inbound, remoteInboundID)},
		map[string]any{"newUserId": newUserID, "reason": "archived_inbound_reappeared", "trafficSamplesRemoved": true}, observedAt); err != nil {
		return "", fmt.Errorf("write reactivation audit: %w", err)
	}
	return newUserID, nil
}

type businessUserState struct {
	Status         string
	ExpiryText     string
	ExpiryMismatch bool
	ExpiryValues   []int64
}

func userStateFromClients(inbound agentInboundPayload, observedAt time.Time) businessUserState {
	if !inbound.Enable {
		return businessUserState{Status: "disabled"}
	}
	uniqueExpiry := make(map[int64]struct{})
	activeClients := 0
	var earliestExpiry int64
	for _, client := range inbound.Clients {
		if !client.Enable || strings.TrimSpace(client.RemoteID) == "" {
			continue
		}
		activeClients++
		uniqueExpiry[client.ExpiryTime] = struct{}{}
		if client.ExpiryTime > 0 && (earliestExpiry == 0 || client.ExpiryTime < earliestExpiry) {
			earliestExpiry = client.ExpiryTime
		}
	}
	if activeClients == 0 {
		return businessUserState{Status: "disabled"}
	}
	values := make([]int64, 0, len(uniqueExpiry))
	for value := range uniqueExpiry {
		values = append(values, value)
	}
	sort.Slice(values, func(left, right int) bool { return values[left] < values[right] })
	state := businessUserState{Status: "active", ExpiryMismatch: len(values) > 1, ExpiryValues: values}
	if earliestExpiry <= 0 {
		return state
	}
	state.ExpiryText = time.Unix(earliestExpiry, 0).UTC().Format(time.RFC3339Nano)
	expiresAt := time.Unix(earliestExpiry, 0).UTC()
	if !expiresAt.After(observedAt.UTC()) {
		state.Status = "expired"
	} else if !expiresAt.After(observedAt.UTC().Add(userExpiringWindow)) {
		state.Status = "expiring"
	}
	return state
}

func reconcileClientExpiryMismatchEvent(tx *sql.Tx, inboundID, userID, fallbackName string, state businessUserState, observedAt time.Time) error {
	const eventType = "client_expiry_mismatch"
	dedupeKey := eventType + ":" + inboundID
	now := observedAt.UTC().Format(time.RFC3339Nano)
	if !state.ExpiryMismatch {
		_, err := tx.Exec(`UPDATE node_events SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?), dedupe_key = NULL
WHERE dedupe_key = ? AND event_status NOT IN ('resolved', 'dismissed')`, now, now, dedupeKey)
		return err
	}
	userName := fallbackName
	var nodeID string
	_ = tx.QueryRow(`SELECT COALESCE(u.display_name, ''), COALESCE(i.node_id, '') FROM users u JOIN inbounds i ON i.user_id = u.id WHERE u.id = ? AND i.id = ?`, userID, inboundID).Scan(&userName, &nodeID)
	if strings.TrimSpace(userName) == "" {
		userName = userID
	}
	message := "用户「" + userName + "」的已启用 Client 到期时间不一致（" + formatClientExpiryValues(state.ExpiryValues) + "）；中央已按最早有限到期时间计算用户状态。"
	result, err := tx.Exec(`UPDATE node_events SET severity = 'warning', title = '多设备到期时间不一致', message = ?, created_at = ?,
visibility = 'public', requires_action = 0, event_status = 'open', acknowledged = 0, read_at = NULL, resolved_at = NULL, resolved_by = NULL,
resource_type = 'inbound', resource_id = ?, action_type = 'inspect_sync_run', source = 'agent'
WHERE dedupe_key = ?`, message, now, inboundID, dedupeKey)
	if err != nil {
		return err
	}
	updated, err := result.RowsAffected()
	if err != nil || updated > 0 {
		return err
	}
	return insertNodeEventTx(tx, nodeEventSpec{
		NodeID: nodeID, EventType: eventType, Category: "business", Severity: "warning", Title: "多设备到期时间不一致",
		Message: message, ResourceType: "inbound", ResourceID: inboundID, ActionType: "inspect_sync_run",
		DedupeKey: dedupeKey, Source: "agent", OccurredAt: observedAt,
	})
}

func formatClientExpiryValues(values []int64) string {
	formatted := make([]string, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			formatted = append(formatted, "长期")
			continue
		}
		formatted = append(formatted, time.Unix(value, 0).UTC().Format("2006-01-02 15:04 UTC"))
	}
	return strings.Join(formatted, "、")
}

func inboundDisplayName(inbound agentInboundPayload, remoteInboundID string) string {
	if value := strings.TrimSpace(inbound.Remark); value != "" {
		return value
	}
	if value := strings.TrimSpace(inbound.Tag); value != "" {
		return value
	}
	return "Inbound " + remoteInboundID
}

func (s *Server) markMissingAndArchiveInbounds(tx *sql.Tx, nodeID string, observedAt time.Time) error {
	observedAtText := observedAt.UTC().Format(time.RFC3339Nano)
	rows, err := tx.Query(`SELECT id, remote_inbound_id, COALESCE(tag, ''), missing_sync_count
FROM inbounds
WHERE node_id = ? AND deleted_at IS NULL AND (last_seen_at IS NULL OR last_seen_at < ?)`, nodeID, observedAtText)
	if err != nil {
		return fmt.Errorf("find missing inbounds: %w", err)
	}

	type missingInbound struct {
		id              string
		remoteID        string
		tag             string
		oldMissingCount int
	}
	missing := make([]missingInbound, 0)
	for rows.Next() {
		var inbound missingInbound
		if err := rows.Scan(&inbound.id, &inbound.remoteID, &inbound.tag, &inbound.oldMissingCount); err != nil {
			_ = rows.Close()
			return fmt.Errorf("read missing inbound: %w", err)
		}
		missing = append(missing, inbound)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("iterate missing inbounds: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close missing inbound query: %w", err)
	}

	if _, err := tx.Exec(`UPDATE inbounds
SET missing_since = COALESCE(missing_since, ?), missing_sync_count = missing_sync_count + 1
WHERE node_id = ? AND deleted_at IS NULL AND (last_seen_at IS NULL OR last_seen_at < ?)`, observedAtText, nodeID, observedAtText); err != nil {
		return fmt.Errorf("mark missing inbounds: %w", err)
	}

	for _, inbound := range missing {
		label := inbound.remoteID
		if inbound.tag != "" {
			label += " (" + inbound.tag + ")"
		}
		if inbound.oldMissingCount == 0 {
			if err := insertNodeEventTx(tx, nodeEventSpec{
				NodeID: nodeID, EventType: "inbound_missing", Category: "node", Severity: "warning",
				Title: "Inbound 暂时缺失", Message: "Inbound " + label + " 未出现在成功同步快照中，连续三次缺失后才会归档",
				ResourceType: "inbound", ResourceID: inbound.id, ActionType: "inspect_inbound", Source: "agent", OccurredAt: observedAt,
			}); err != nil {
				return fmt.Errorf("record inbound missing %s: %w", inbound.remoteID, err)
			}
		}
		if inbound.oldMissingCount+1 < missingInboundArchiveAfter {
			continue
		}
		if _, err := tx.Exec(`UPDATE inbounds SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`, observedAtText, inbound.id); err != nil {
			return fmt.Errorf("archive inbound %s: %w", inbound.remoteID, err)
		}
		if err := s.cleanupArchivedInboundTx(tx, inbound.id, observedAt); err != nil {
			return fmt.Errorf("clean up archived inbound %s: %w", inbound.remoteID, err)
		}
		if err := insertNodeEventTx(tx, nodeEventSpec{
			NodeID: nodeID, EventType: "inbound_archived", Category: "node", Severity: "warning",
			Title: "Inbound 已归档", Message: "Inbound " + label + " 连续三次同步缺失，已归档并清理运营数据，账务记录继续保留",
			ResourceType: "inbound", ResourceID: inbound.id, ActionType: "inspect_inbound", Source: "agent", OccurredAt: observedAt,
		}); err != nil {
			return fmt.Errorf("record inbound archive %s: %w", inbound.remoteID, err)
		}
	}
	return nil
}

// cleanupArchivedInboundTx removes operational state that belongs exclusively
// to an Inbound which has disappeared from a complete Agent snapshot. Billing
// records are kept on the user row. If the user still has another live relay
// Inbound, its shared path remains active; otherwise the path is closed so the
// orphan cleanup that follows this sync can retire the user safely.
func (s *Server) cleanupArchivedInboundTx(tx *sql.Tx, inboundID string, observedAt time.Time) error {
	nowText := observedAt.UTC().Format(time.RFC3339Nano)
	var userID string
	if err := tx.QueryRow(`SELECT COALESCE(user_id, '') FROM inbounds WHERE id = ?`, inboundID).Scan(&userID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM clients WHERE inbound_id = ?`, inboundID); err != nil {
		return fmt.Errorf("delete clients: %w", err)
	}
	if _, err := deleteInboundTrafficTx(tx, inboundID); err != nil {
		return fmt.Errorf("delete traffic snapshots: %w", err)
	}
	if userID == "" {
		return nil
	}
	if _, err := tx.Exec(`UPDATE user_inbounds SET active_to = ? WHERE user_id = ? AND inbound_id = ? AND active_to IS NULL`, nowText, userID, inboundID); err != nil {
		return fmt.Errorf("close inbound mapping: %w", err)
	}
	var remainingAssociations int
	if err := tx.QueryRow(`SELECT COUNT(*)
FROM user_inbounds ui JOIN inbounds i ON i.id = ui.inbound_id
WHERE ui.user_id = ? AND ui.active_to IS NULL AND i.deleted_at IS NULL AND i.kind = 'user'`, userID).Scan(&remainingAssociations); err != nil {
		return fmt.Errorf("count remaining user inbounds: %w", err)
	}
	if remainingAssociations > 0 {
		if err := closeRenewalCandidatesForInboundTx(tx, userID, inboundID, nowText); err != nil {
			return err
		}
		return nil
	}
	if _, err := tx.Exec(`UPDATE user_paths SET active_to = ?, updated_at = ? WHERE user_id = ? AND active_to IS NULL`, nowText, nowText, userID); err != nil {
		return fmt.Errorf("close user paths: %w", err)
	}
	if _, err := tx.Exec(`UPDATE user_routes SET is_primary = 0, active_to = ?, route_exit_ip_id = NULL WHERE user_id = ? AND active_to IS NULL`, nowText, userID); err != nil {
		return fmt.Errorf("close user routes: %w", err)
	}
	return nil
}

func closeRenewalCandidatesForInboundTx(tx *sql.Tx, userID, inboundID, nowText string) error {
	if _, err := tx.Exec(`UPDATE user_renewal_candidates SET status = 'rejected', processed_at = ?,
notes = CASE WHEN COALESCE(notes, '') = '' THEN '自动关闭：客户已更换' ELSE notes || char(10) || '自动关闭：客户已更换' END
WHERE user_id = ? AND inbound_id = ? AND status = 'pending'`, nowText, userID, inboundID); err != nil {
		return fmt.Errorf("close renewal candidates: %w", err)
	}
	if _, err := tx.Exec(`UPDATE node_events SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?)
WHERE event_type = 'renewal_candidate_detected' AND resource_type = 'renewal' AND event_status NOT IN ('resolved', 'dismissed')
  AND resource_id IN (SELECT id FROM user_renewal_candidates WHERE user_id = ? AND inbound_id = ? AND status = 'rejected')`, nowText, nowText, userID, inboundID); err != nil {
		return fmt.Errorf("close renewal events: %w", err)
	}
	return nil
}

func detectTrafficReset(tx *sql.Tx, inboundID string, observedAt time.Time, allTime int64) (bool, error) {
	var previous int64
	err := tx.QueryRow(`SELECT all_time FROM traffic_snapshots WHERE inbound_id = ? AND collected_at < ? ORDER BY collected_at DESC LIMIT 1`, inboundID, observedAt.Format(time.RFC3339Nano)).Scan(&previous)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return allTime < previous, nil
}

func (s *Server) failSync(w http.ResponseWriter, tx *sql.Tx, syncRunID string, err error) error {
	message := err.Error()
	now := time.Now().UTC()
	finishedAt := now.Format(time.RFC3339Nano)
	_, _ = tx.Exec(`UPDATE sync_runs SET finished_at = ?, status = 'failed', error_message = ? WHERE id = ?`, finishedAt, message, syncRunID)
	var nodeID string
	if scanErr := tx.QueryRow(`SELECT node_id FROM sync_runs WHERE id = ?`, syncRunID).Scan(&nodeID); scanErr == nil {
		if eventErr := insertNodeEventTx(tx, nodeEventSpec{
			NodeID: nodeID, EventType: "sync_failed", Category: "sync", Severity: "error", Title: "节点同步失败",
			Message: message, ResourceType: "sync_run", ResourceID: syncRunID, ActionType: "retry_sync",
			Source: "agent", CorrelationID: syncRunID, OccurredAt: now,
		}); eventErr != nil {
			s.logger.Warn("record sync failure event", "sync_id", syncRunID, "error", eventErr)
		}
	}
	_ = tx.Commit()
	writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not apply node sync")
	return err
}

func (s *Server) updateNodeHeartbeat(nodeID string, observedAt time.Time, status agentStatusPayload, statusAvailable bool, syncState, syncError string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var previousHealth, previousSyncState string
	_ = s.db.QueryRow(`SELECT health_status, sync_status FROM nodes WHERE id = ?`, nodeID).Scan(&previousHealth, &previousSyncState)
	newHealth := "online"
	if statusAvailable {
		newHealth = healthStatus(status.XrayRunning)
	}
	syncError = truncateAgentError(syncError)
	var err error
	if statusAvailable {
		_, err = s.db.Exec(`UPDATE nodes SET health_status = ?, agent_version = CASE WHEN ? <> '' THEN ? ELSE agent_version END, xray_version = CASE WHEN ? <> '' THEN ? ELSE xray_version END,
cpu_usage = ?, memory_used = ?, memory_total = ?, disk_used = ?, disk_total = ?, sync_status = CASE WHEN ? = 'failed' THEN 'failed' ELSE sync_status END, last_sync_error = CASE WHEN ? = 'failed' THEN ? ELSE last_sync_error END, last_seen_at = ?, updated_at = ? WHERE id = ?`, newHealth, strings.TrimSpace(status.AgentVersion), strings.TrimSpace(status.AgentVersion), status.XrayVersion, status.XrayVersion,
			status.CPUUsage, nullableMetric(status.MemoryUsed), nullableMetric(status.MemoryTotal), nullableMetric(status.DiskUsed), nullableMetric(status.DiskTotal), syncState, syncState, nullableDBString(syncError), observedAt.Format(time.RFC3339Nano), now, nodeID)
	} else {
		_, err = s.db.Exec(`UPDATE nodes SET health_status = ?, agent_version = CASE WHEN ? <> '' THEN ? ELSE agent_version END,
sync_status = CASE WHEN ? = 'failed' THEN 'failed' ELSE sync_status END, last_sync_error = CASE WHEN ? = 'failed' THEN ? ELSE last_sync_error END, last_seen_at = ?, updated_at = ? WHERE id = ?`, newHealth, strings.TrimSpace(status.AgentVersion), strings.TrimSpace(status.AgentVersion), syncState, syncState, nullableDBString(syncError), observedAt.Format(time.RFC3339Nano), now, nodeID)
	}
	if err == nil && previousHealth == "offline" && newHealth != "offline" {
		eventErr := insertNodeEvent(s.db, nodeEventSpec{
			NodeID: nodeID, EventType: "node_recovered", Category: "node", Severity: "info", Title: "节点已恢复在线",
			Message: "节点已重新收到 Agent 心跳", ResourceType: "node", ResourceID: nodeID, Source: "agent", OccurredAt: observedAt,
		})
		if eventErr != nil {
			s.logger.Warn("record node recovery event", "node_id", nodeID, "error", eventErr)
		}
	}
	if err == nil && syncState == "failed" && previousSyncState != "failed" {
		eventErr := insertNodeEvent(s.db, nodeEventSpec{
			NodeID: nodeID, EventType: "xpanel_sync_failed", Category: "sync", Severity: "warning", Title: "X-Panel 采集异常",
			Message: "Agent 在线，但 X-Panel 数据采集失败：" + firstNonEmpty(syncError, "未提供错误详情"), ResourceType: "node", ResourceID: nodeID,
			ActionType: "inspect_sync_run", DedupeKey: "xpanel-sync-failed:" + nodeID, Source: "agent", OccurredAt: observedAt,
		})
		if eventErr != nil {
			s.logger.Warn("record X-Panel sync failure event", "node_id", nodeID, "error", eventErr)
		}
	}
	return err
}

func truncateAgentError(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 500 {
		return value[:500]
	}
	return value
}

func nullableMetric(value int64) any {
	if value <= 0 {
		return nil
	}
	return value
}

func (s *Server) requireAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "node authentication required")
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		var principal agentPrincipal
		var revokedAt sql.NullString
		err := s.db.QueryRow(`SELECT n.id, n.node_key, c.revoked_at FROM node_credentials c JOIN nodes n ON n.id = c.node_id WHERE c.token_hash = ? AND n.enabled = 1 AND n.deleted_at IS NULL`, hashToken(token)).Scan(&principal.NodeID, &principal.NodeKey, &revokedAt)
		if err != nil || revokedAt.Valid {
			writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "node authentication required")
			return
		}
		next.ServeHTTP(w, r.WithContext(contextWithAgent(r, principal)))
	})
}

func contextWithAgent(r *http.Request, principal agentPrincipal) context.Context {
	return context.WithValue(r.Context(), agentContextKey{}, principal)
}

func decodeAgentJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 4<<20))
	return decoder.Decode(target)
}

func parseObservedAt(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, errors.New("observed_at is required")
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func epochString(value int64) any {
	if value <= 0 {
		return nil
	}
	return time.Unix(value, 0).UTC().Format(time.RFC3339Nano)
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func healthStatus(xrayRunning bool) string {
	if xrayRunning {
		return "online"
	}
	return "degraded"
}
