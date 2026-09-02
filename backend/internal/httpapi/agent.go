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

type agentStatusPayload struct {
	XrayRunning   bool   `json:"xray_running"`
	XrayVersion   string `json:"xray_version"`
	XPanelVersion string `json:"xpanel_version"`
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
	err = tx.QueryRow(`SELECT id FROM nodes WHERE node_key = ?`, strings.TrimSpace(payload.NodeKey)).Scan(&nodeID)
	if errors.Is(err, sql.ErrNoRows) {
		nodeID = newID()
		_, err = tx.Exec(`INSERT INTO nodes (id, node_key, name, type, hostname, panel_base_path, agent_version, xpanel_version, xray_version, enabled, health_status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'unknown', ?, ?)`, nodeID, strings.TrimSpace(payload.NodeKey), strings.TrimSpace(payload.NodeName), nodeType,
			strings.TrimSpace(payload.Hostname), strings.TrimSpace(payload.PanelBasePath), strings.TrimSpace(payload.AgentVersion), strings.TrimSpace(payload.XPanelVersion), strings.TrimSpace(payload.XrayVersion), now, now)
	} else if err == nil {
		_, err = tx.Exec(`UPDATE nodes SET name = ?, type = ?, hostname = ?, panel_base_path = ?, agent_version = ?, xpanel_version = ?, xray_version = ?, enabled = 1, updated_at = ? WHERE id = ?`,
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

	clientCount := 0
	for _, inbound := range payload.Inbounds {
		remoteInboundID := strconv.FormatInt(inbound.RemoteID, 10)
		expiry := epochString(inbound.ExpiryTime)
		_, err := tx.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, tag, remark, protocol, port, listen, enable, expiry_time, up, down, all_time, client_count, config_hash, first_seen_at, last_seen_at, missing_since, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, NULL)
ON CONFLICT(node_id, remote_inbound_id) DO UPDATE SET tag = excluded.tag, remark = excluded.remark, protocol = excluded.protocol, port = excluded.port, listen = excluded.listen, enable = excluded.enable, expiry_time = excluded.expiry_time, up = excluded.up, down = excluded.down, all_time = excluded.all_time, client_count = excluded.client_count, config_hash = excluded.config_hash, last_seen_at = excluded.last_seen_at, missing_since = NULL, deleted_at = NULL`,
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
	if _, err := tx.Exec(`UPDATE inbounds SET missing_since = COALESCE(missing_since, ?) WHERE node_id = ? AND deleted_at IS NULL AND (last_seen_at IS NULL OR last_seen_at < ?)`, observedAt.Format(time.RFC3339Nano), principal.NodeID, observedAt.Format(time.RFC3339Nano)); err != nil {
		s.failSync(w, tx, syncRunID, fmt.Errorf("mark missing inbounds: %w", err))
		return
	}
	if _, err := tx.Exec(`UPDATE nodes SET health_status = ?, xpanel_version = CASE WHEN ? <> '' THEN ? ELSE xpanel_version END, xray_version = CASE WHEN ? <> '' THEN ? ELSE xray_version END, last_seen_at = ?, updated_at = ? WHERE id = ?`, healthStatus(payload.Status.XrayRunning), payload.Status.XPanelVersion, payload.Status.XPanelVersion, payload.Status.XrayVersion, payload.Status.XrayVersion, observedAt.Format(time.RFC3339Nano), now, principal.NodeID); err != nil {
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
	_, err := s.db.Exec(`UPDATE nodes SET health_status = ?, xpanel_version = CASE WHEN ? <> '' THEN ? ELSE xpanel_version END, xray_version = CASE WHEN ? <> '' THEN ? ELSE xray_version END, last_seen_at = ?, updated_at = ? WHERE id = ?`, healthStatus(status.XrayRunning), status.XPanelVersion, status.XPanelVersion, status.XrayVersion, status.XrayVersion, observedAt.Format(time.RFC3339Nano), now, nodeID)
	return err
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
		err := s.db.QueryRow(`SELECT n.id, n.node_key, c.revoked_at FROM node_credentials c JOIN nodes n ON n.id = c.node_id WHERE c.token_hash = ? AND n.enabled = 1`, hashToken(token)).Scan(&principal.NodeID, &principal.NodeKey, &revokedAt)
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
