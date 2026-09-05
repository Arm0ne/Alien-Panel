package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	XPanelVersion string `json:"xpanel_version"`
	XrayVersion   string `json:"xray_version"`
}

type agentBootstrapRequest struct {
	InstallToken string `json:"install_token"`
	Hostname     string `json:"hostname"`
	AgentVersion string `json:"agent_version"`
}

type agentStatusPayload struct {
	XrayRunning   bool    `json:"xray_running"`
	XrayVersion   string  `json:"xray_version"`
	XPanelVersion string  `json:"xpanel_version"`
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsed    int64   `json:"memory_used"`
	MemoryTotal   int64   `json:"memory_total"`
	DiskUsed      int64   `json:"disk_used"`
	DiskTotal     int64   `json:"disk_total"`
}

type agentHeartbeatRequest struct {
	NodeKey    string             `json:"node_key"`
	ObservedAt string             `json:"observed_at"`
	Status     agentStatusPayload `json:"status"`
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
	RemoteID   int64                `json:"remote_id"`
	Tag        string               `json:"tag"`
	Remark     string               `json:"remark"`
	Protocol   string               `json:"protocol"`
	Port       int64                `json:"port"`
	Listen     string               `json:"listen"`
	Enable     bool                 `json:"enable"`
	ExpiryTime int64                `json:"expiry_time"`
	Up         int64                `json:"up"`
	Down       int64                `json:"down"`
	AllTime    int64                `json:"all_time"`
	ConfigHash string               `json:"config_hash"`
	Clients    []agentClientPayload `json:"clients"`
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
		_, err = tx.Exec(`INSERT INTO nodes (id, node_key, name, type, hostname, panel_base_path, agent_version, xpanel_version, xray_version, enabled, health_status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'unknown', ?, ?)`, nodeID, strings.TrimSpace(payload.NodeKey), strings.TrimSpace(payload.NodeName), nodeType,
			strings.TrimSpace(payload.Hostname), strings.TrimSpace(payload.PanelBasePath), strings.TrimSpace(payload.AgentVersion), strings.TrimSpace(payload.XPanelVersion), strings.TrimSpace(payload.XrayVersion), now, now)
	} else if err == nil {
		if deletedAt != "" {
			writeFailure(w, http.StatusConflict, validationCode, "node has been deleted; create a new node record before registering the Agent")
			return
		}
		// Registration refreshes Agent-owned metadata but must not override an
		// administrator's explicit disabled state. Re-enable through PATCH /api/nodes/{id}.
		_, err = tx.Exec(`UPDATE nodes SET name = ?, type = ?, hostname = ?, panel_base_path = ?, agent_version = ?, xpanel_version = ?, xray_version = ?, updated_at = ? WHERE id = ?`,
			strings.TrimSpace(payload.NodeName), nodeType, strings.TrimSpace(payload.Hostname), strings.TrimSpace(payload.PanelBasePath), strings.TrimSpace(payload.AgentVersion), strings.TrimSpace(payload.XPanelVersion), strings.TrimSpace(payload.XrayVersion), now, nodeID)
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
	if err := s.updateNodeHeartbeat(principal.NodeID, observedAt, payload.Status); err != nil {
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
		expiryText := ""
		var expiry any
		if inbound.ExpiryTime > 0 {
			expiryText = time.Unix(inbound.ExpiryTime, 0).UTC().Format(time.RFC3339Nano)
			expiry = expiryText
		}
		var previousExpiryText, previousSeenText string
		_ = tx.QueryRow(`SELECT COALESCE(expiry_time, ''), COALESCE(last_seen_at, '') FROM inbounds WHERE node_id = ? AND remote_inbound_id = ? AND deleted_at IS NULL`, principal.NodeID, remoteInboundID).
			Scan(&previousExpiryText, &previousSeenText)
		_, err := tx.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, tag, remark, protocol, port, listen, enable, expiry_time, up, down, all_time, client_count, config_hash, first_seen_at, last_seen_at, missing_since, missing_sync_count, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, 0, NULL)
ON CONFLICT(node_id, remote_inbound_id) DO UPDATE SET tag = excluded.tag, remark = excluded.remark, protocol = excluded.protocol, port = excluded.port, listen = excluded.listen, enable = excluded.enable, expiry_time = excluded.expiry_time, up = excluded.up, down = excluded.down, all_time = excluded.all_time, client_count = excluded.client_count, config_hash = excluded.config_hash, last_seen_at = excluded.last_seen_at, missing_since = NULL, missing_sync_count = 0, deleted_at = NULL`,
			newID(), principal.NodeID, remoteInboundID, inbound.Tag, inbound.Remark, inbound.Protocol, inbound.Port, inbound.Listen, boolInt(inbound.Enable), expiry, inbound.Up, inbound.Down, inbound.AllTime, len(inbound.Clients), inbound.ConfigHash, observedAt.Format(time.RFC3339Nano), observedAt.Format(time.RFC3339Nano))
		if err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("upsert inbound %s: %w", remoteInboundID, err))
			return
		}
		var inboundID string
		if err := tx.QueryRow(`SELECT id FROM inbounds WHERE node_id = ? AND remote_inbound_id = ?`, principal.NodeID, remoteInboundID).Scan(&inboundID); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("find inbound %s: %w", remoteInboundID, err))
			return
		}
		if err := s.ensureRelayInboundUserWithHistory(tx, nodeType, inboundID, remoteInboundID, inbound, expiryText, previousExpiryText, previousSeenText, observedAt); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("ensure business user for inbound %s: %w", remoteInboundID, err))
			return
		}
		resetDetected, err := detectTrafficReset(tx, inboundID, observedAt, inbound.AllTime)
		if err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("check traffic reset %s: %w", remoteInboundID, err))
			return
		}
		if resetDetected {
			_, err := tx.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES (?, ?, 'traffic_reset', 'warning', ?, ?)`, newID(), principal.NodeID, "Inbound "+remoteInboundID+" cumulative traffic moved backwards; a new baseline was recorded", now)
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
		resetFlag := 0
		if resetDetected {
			resetFlag = 1
		}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected, sync_run_id) VALUES (?, ?, ?, ?, ?, ?, ?, 'xpanel', ?, ?)`, newID(), principal.NodeID, inboundID, observedAt.Format(time.RFC3339Nano), inbound.Up, inbound.Down, inbound.AllTime, resetFlag, syncRunID); err != nil {
			s.failSync(w, tx, syncRunID, fmt.Errorf("save traffic snapshot %s: %w", remoteInboundID, err))
			return
		}
	}
	if err := s.markMissingAndArchiveInbounds(tx, principal.NodeID, observedAt); err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("mark or archive missing inbounds: %w", err))
		return
	}
	if _, err := tx.Exec(`UPDATE nodes SET health_status = ?, xpanel_version = CASE WHEN ? <> '' THEN ? ELSE xpanel_version END, xray_version = CASE WHEN ? <> '' THEN ? ELSE xray_version END,
cpu_usage = ?, memory_used = ?, memory_total = ?, disk_used = ?, disk_total = ?, last_seen_at = ?, updated_at = ? WHERE id = ?`, healthStatus(payload.Status.XrayRunning), payload.Status.XPanelVersion, payload.Status.XPanelVersion, payload.Status.XrayVersion, payload.Status.XrayVersion,
		payload.Status.CPUUsage, nullableMetric(payload.Status.MemoryUsed), nullableMetric(payload.Status.MemoryTotal), nullableMetric(payload.Status.DiskUsed), nullableMetric(payload.Status.DiskTotal), observedAt.Format(time.RFC3339Nano), now, principal.NodeID); err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("update node after sync: %w", err))
		return
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

// ensureRelayInboundUser establishes the first-version business rule that one
// relay-node Inbound is one central business user. It deliberately never
// overwrites business fields maintained in the central panel (name, fee,
// currency, notes); X-Panel remains authoritative only for Inbound state and
// expiry. Landing-node and explicitly infrastructure-classified Inbounds stay
// outside the business-user list.
func (s *Server) ensureRelayInboundUser(tx *sql.Tx, nodeType, inboundID, remoteInboundID string, inbound agentInboundPayload, expiryText string, observedAt time.Time) error {
	return s.ensureRelayInboundUserWithHistory(tx, nodeType, inboundID, remoteInboundID, inbound, expiryText, "", "", observedAt)
}

func (s *Server) ensureRelayInboundUserWithHistory(tx *sql.Tx, nodeType, inboundID, remoteInboundID string, inbound agentInboundPayload, expiryText, previousExpiryText, previousSeenText string, observedAt time.Time) error {
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

	var userID, kind string
	if err := tx.QueryRow(`SELECT COALESCE(user_id, ''), kind FROM inbounds WHERE id = ?`, inboundID).Scan(&userID, &kind); err != nil {
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

	now := observedAt.UTC().Format(time.RFC3339Nano)
	status := statusFromInbound(inbound, observedAt)
	var expiry any
	if expiryText != "" {
		expiry = expiryText
	}
	if userID == "" {
		userID = newID()
		if _, err := tx.Exec(`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`, userID, inboundDisplayName(inbound, remoteInboundID), status, expiry, now, now); err != nil {
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
		candidateOldExpiry := previousExpiryText
		if candidateOldExpiry == "" {
			candidateOldExpiry = oldExpiry
		}
		// If the previous snapshot had no expiry (the X-Panel value was 0)
		// but this is not the first sync, retain a meaningful baseline instead
		// of silently dropping the renewal suggestion.
		if candidateOldExpiry == "" && previousSeenText != "" {
			candidateOldExpiry = previousSeenText
		}
		if candidateOldExpiry != "" && expiryText != "" {
			oldTime, oldErr := time.Parse(time.RFC3339Nano, candidateOldExpiry)
			newTime, newErr := time.Parse(time.RFC3339Nano, expiryText)
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'CNY', 'pending')`, newID(), userID, inboundID, candidateOldExpiry, expiryText, now, billingCycle, suggestedAmount); err != nil {
					return fmt.Errorf("record renewal candidate: %w", err)
				}
			}
		}
		if _, err := tx.Exec(`UPDATE users
SET expiry_time = ?, status = CASE WHEN status = 'disabled' THEN 'disabled' ELSE ? END, updated_at = ?
WHERE id = ?`, expiry, status, now, userID); err != nil {
			return fmt.Errorf("refresh business user state: %w", err)
		}
	}
	if _, err := tx.Exec(`INSERT OR IGNORE INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES (?, ?, ?, 1, ?)`, newID(), userID, inboundID, now); err != nil {
		return fmt.Errorf("save business user mapping: %w", err)
	}
	return nil
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

func statusFromInbound(inbound agentInboundPayload, observedAt time.Time) string {
	if !inbound.Enable {
		return "disabled"
	}
	if inbound.ExpiryTime <= 0 {
		return "active"
	}
	expiresAt := time.Unix(inbound.ExpiryTime, 0).UTC()
	if !expiresAt.After(observedAt.UTC()) {
		return "expired"
	}
	if !expiresAt.After(observedAt.UTC().Add(userExpiringWindow)) {
		return "expiring"
	}
	return "active"
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
			if _, err := tx.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES (?, ?, 'inbound_missing', 'warning', ?, ?)`,
				newID(), nodeID, "Inbound "+label+" was absent from a successful snapshot; it will be archived only after three consecutive missing snapshots", observedAtText); err != nil {
				return fmt.Errorf("record inbound missing %s: %w", inbound.remoteID, err)
			}
		}
		if inbound.oldMissingCount+1 < missingInboundArchiveAfter {
			continue
		}
		if _, err := tx.Exec(`UPDATE inbounds SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`, observedAtText, inbound.id); err != nil {
			return fmt.Errorf("archive inbound %s: %w", inbound.remoteID, err)
		}
		if _, err := tx.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES (?, ?, 'inbound_archived', 'warning', ?, ?)`,
			newID(), nodeID, "Inbound "+label+" was absent from three consecutive successful snapshots and was archived; historical data was retained", observedAtText); err != nil {
			return fmt.Errorf("record inbound archive %s: %w", inbound.remoteID, err)
		}
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
	_, _ = tx.Exec(`UPDATE sync_runs SET finished_at = ?, status = 'failed', error_message = ? WHERE id = ?`, time.Now().UTC().Format(time.RFC3339Nano), message, syncRunID)
	_ = tx.Commit()
	writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not apply node sync")
	return err
}

func (s *Server) updateNodeHeartbeat(nodeID string, observedAt time.Time, status agentStatusPayload) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := s.db.Exec(`UPDATE nodes SET health_status = ?, xpanel_version = CASE WHEN ? <> '' THEN ? ELSE xpanel_version END, xray_version = CASE WHEN ? <> '' THEN ? ELSE xray_version END,
cpu_usage = ?, memory_used = ?, memory_total = ?, disk_used = ?, disk_total = ?, last_seen_at = ?, updated_at = ? WHERE id = ?`, healthStatus(status.XrayRunning), status.XPanelVersion, status.XPanelVersion, status.XrayVersion, status.XrayVersion,
		status.CPUUsage, nullableMetric(status.MemoryUsed), nullableMetric(status.MemoryTotal), nullableMetric(status.DiskUsed), nullableMetric(status.DiskTotal), observedAt.Format(time.RFC3339Nano), now, nodeID)
	return err
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
