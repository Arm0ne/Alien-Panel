package httpapi

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"
)

// userPathAssets returns the resources that can be used to configure one
// business user's active path.  A user's relay and primary inbound are
// derived from the user itself; landing inbounds are infrastructure assets
// and are deliberately never returned as business users.
func (s *Server) userPathAssets(w http.ResponseWriter, r *http.Request) {
	s.refreshOperationalStatuses(time.Now().UTC())
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}

	var exists int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = ?`, userID).Scan(&exists); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user")
		return
	}
	if exists == 0 {
		writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
		return
	}

	var relayID, relayName, relayType, relayStatus, relayHost, relayLastSeen, relayLastSync string
	var relayEnabled int
	err := s.db.QueryRow(`SELECT n.id, n.name, n.type, n.health_status, n.enabled,
COALESCE(NULLIF(n.public_ip, ''), n.hostname, ''), COALESCE(n.last_seen_at, ''),
COALESCE((SELECT COALESCE(sr.finished_at, sr.started_at) FROM sync_runs sr
 WHERE sr.node_id = n.id AND sr.status = 'success'
 ORDER BY COALESCE(sr.finished_at, sr.started_at) DESC LIMIT 1), '')
FROM user_inbounds ui
JOIN inbounds i ON i.id = ui.inbound_id
JOIN nodes n ON n.id = i.node_id
WHERE ui.user_id = ? AND ui.is_primary = 1 AND ui.active_to IS NULL
  AND i.deleted_at IS NULL AND i.enable = 1 AND i.kind = 'user'
  AND n.type = 'relay' AND n.deleted_at IS NULL`, userID).Scan(
		&relayID, &relayName, &relayType, &relayStatus, &relayEnabled, &relayHost, &relayLastSeen, &relayLastSync,
	)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusConflict, validationCode, "user has no active primary inbound")
		return
	}
	if err != nil {
		s.logger.Error("read user path relay", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user's relay node")
		return
	}
	if relayType != "relay" || relayEnabled != 1 {
		writeFailure(w, http.StatusConflict, validationCode, "user's primary inbound node is not an enabled relay")
		return
	}
	if relayStatus == "" {
		relayStatus = "unknown"
	}

	landingNodes, err := s.readLandingPathAssets()
	if err != nil {
		s.logger.Error("read landing path assets", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read landing path assets")
		return
	}
	relayExitIPs, err := s.readSelectablePathExitIPs(relayID)
	if err != nil {
		s.logger.Error("read relay path exit IPs", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read relay exit IP assets")
		return
	}
	externalExitIPs, err := s.readSelectablePathExitIPs("")
	if err != nil {
		s.logger.Error("read external path exit IPs", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read external exit IP assets")
		return
	}

	generatedAt := time.Now().UTC().Format(time.RFC3339Nano)
	dataAt := s.latestSuccessfulSyncAt()
	relayData := map[string]any{
		"id": relayID, "name": relayName, "type": relayType, "status": relayStatus,
		"enabled": relayEnabled == 1, "host": relayHost,
		"lastSeenAt": nullableString(relayLastSeen), "lastSyncAt": nullableString(relayLastSync),
	}
	writeSuccess(w, map[string]any{
		"generatedAt":     generatedAt,
		"dataAt":          dataAt,
		"relay":           relayData,
		"landingNodes":    landingNodes,
		"relayExitIps":    relayExitIPs,
		"externalExitIps": externalExitIPs,
	})
}

func (s *Server) readLandingPathAssets() ([]map[string]any, error) {
	rows, err := s.db.Query(`SELECT n.id, n.name, n.type, n.health_status, n.enabled,
COALESCE(NULLIF(n.public_ip, ''), n.hostname, ''), COALESCE(n.last_seen_at, ''),
COALESCE((SELECT COALESCE(sr.finished_at, sr.started_at) FROM sync_runs sr
 WHERE sr.node_id = n.id AND sr.status = 'success'
 ORDER BY COALESCE(sr.finished_at, sr.started_at) DESC LIMIT 1), '')
FROM nodes n WHERE n.type = 'landing' AND n.deleted_at IS NULL ORDER BY n.name ASC, n.id ASC`)
	if err != nil {
		return nil, err
	}
	type landingNodeRow struct {
		id, name, nodeType, status, host, lastSeen, lastSync string
		enabled                                              int
	}
	nodeRows := make([]landingNodeRow, 0)
	for rows.Next() {
		var id, name, nodeType, status, host, lastSeen, lastSync string
		var enabled int
		if err := rows.Scan(&id, &name, &nodeType, &status, &enabled, &host, &lastSeen, &lastSync); err != nil {
			_ = rows.Close()
			return nil, err
		}
		nodeRows = append(nodeRows, landingNodeRow{id: id, name: name, nodeType: nodeType, status: status, enabled: enabled, host: host, lastSeen: lastSeen, lastSync: lastSync})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	items := make([]map[string]any, 0, len(nodeRows))
	for _, node := range nodeRows {
		status := node.status
		if status == "" {
			status = "unknown"
		}
		id, name, nodeType, enabled, host, lastSeen, lastSync := node.id, node.name, node.nodeType, node.enabled, node.host, node.lastSeen, node.lastSync
		inbounds, inboundState, err := s.readLandingPathInbounds(id, lastSync)
		if err != nil {
			return nil, err
		}
		exitIPs, err := s.readSelectablePathExitIPs(id)
		if err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"id": id, "name": name, "type": nodeType, "status": status, "enabled": enabled == 1, "host": host,
			"lastSeenAt": nullableString(lastSeen), "lastSyncAt": nullableString(lastSync),
			"inboundState": inboundState, "inbounds": inbounds, "exitIps": exitIPs,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Server) readLandingPathInbounds(nodeID, lastSync string) ([]map[string]any, string, error) {
	rows, err := s.db.Query(`SELECT i.id, i.remote_inbound_id, COALESCE(i.tag, ''), COALESCE(i.remark, ''),
COALESCE(i.kind, ''), COALESCE(i.protocol, ''), COALESCE(i.port, 0), COALESCE(i.listen, ''), i.enable,
COALESCE(i.expiry_time, ''), COALESCE(i.client_count, 0), COALESCE(i.up, 0), COALESCE(i.down, 0), COALESCE(i.all_time, 0),
COALESCE(i.last_seen_at, ''), COALESCE(i.deleted_at, '')
FROM inbounds i WHERE i.node_id = ? AND i.deleted_at IS NULL
ORDER BY i.enable DESC, i.tag ASC, i.remote_inbound_id ASC`, nodeID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	items := make([]map[string]any, 0)
	activeCount := 0
	for rows.Next() {
		var id, remoteID, tag, remark, kind, protocol, listen, expiry, lastSeen string
		var port, enabled, clientCount int
		var up, down, allTime int64
		if err := rows.Scan(&id, &remoteID, &tag, &remark, &kind, &protocol, &port, &listen, &enabled,
			&expiry, &clientCount, &up, &down, &allTime, &lastSeen, new(string)); err != nil {
			return nil, "", err
		}
		if enabled == 1 {
			activeCount++
		}
		items = append(items, map[string]any{
			"id": id, "remoteId": remoteID, "tag": nullableString(tag), "remark": nullableString(remark),
			"kind": kind, "purpose": "infrastructure", "protocol": nullableString(protocol), "port": port,
			"listen": nullableString(listen), "enabled": enabled == 1, "status": mapInboundStatus(enabled),
			"expiresAt": nullableString(expiry), "clientCount": clientCount, "up": up, "down": down, "allTime": allTime,
			"lastSeenAt": nullableString(lastSeen), "deletedAt": nil,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}
	state := "pending"
	if activeCount > 0 {
		state = "ready"
	} else if strings.TrimSpace(lastSync) != "" {
		state = "empty"
	}
	return items, state, nil
}

func mapInboundStatus(enabled int) string {
	if enabled == 1 {
		return "active"
	}
	return "disabled"
}

func inboundPurpose(nodeType string) string {
	if nodeType == "landing" {
		return "infrastructure"
	}
	return "business"
}

func (s *Server) readSelectablePathExitIPs(ownerNodeID string) ([]map[string]any, error) {
	args := []any{}
	where := `e.enabled = 1
AND (e.valid_from IS NULL OR date(e.valid_from) <= date('now'))
AND (e.valid_to IS NULL OR date(e.valid_to) >= date('now'))`
	if ownerNodeID == "" {
		where += ` AND e.source_type = 's5'`
	} else {
		where += ` AND e.source_type = 'node' AND COALESCE(e.owner_node_id, e.landing_node_id) = ?`
		args = append(args, ownerNodeID)
	}
	rows, err := s.db.Query(`SELECT e.id, e.ip, COALESCE(e.source_type, 'node'),
COALESCE(e.owner_node_id, e.landing_node_id, ''), COALESCE(owner.name, ''), COALESCE(owner.type, ''),
COALESCE(e.landing_node_id, ''), COALESCE(landing.name, ''), e.family, COALESCE(e.provider, ''),
e.enabled, e.monthly_cost, e.currency, COALESCE(e.valid_from, ''), COALESCE(e.valid_to, ''), e.updated_at,
(SELECT COUNT(DISTINCT p.user_id) FROM user_paths p JOIN users u ON u.id = p.user_id
 WHERE p.exit_ip_id = e.id AND p.active_to IS NULL AND u.status <> 'disabled'
 AND (u.expiry_time IS NULL OR datetime(u.expiry_time) >= datetime('now')))
FROM exit_ips e
LEFT JOIN nodes owner ON owner.id = COALESCE(e.owner_node_id, e.landing_node_id)
LEFT JOIN nodes landing ON landing.id = e.landing_node_id
WHERE `+where+` ORDER BY e.ip ASC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, address, sourceType, ownerID, ownerName, ownerType, landingID, landingName, provider, currency, validFrom, validTo, checkedAt string
		var family, enabled, allocated int
		var monthlyCost float64
		if err := rows.Scan(&id, &address, &sourceType, &ownerID, &ownerName, &ownerType, &landingID, &landingName,
			&family, &provider, &enabled, &monthlyCost, &currency, &validFrom, &validTo, &checkedAt, &allocated); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"id": id, "address": address, "sourceType": sourceType, "ownerNodeId": nullableString(ownerID),
			"ownerNodeName": nullableString(ownerName), "ownerNodeType": nullableString(ownerType),
			"landingNodeId": nullableString(landingID), "landingNodeName": nullableString(landingName),
			"family": family, "provider": nullableString(provider), "status": "active", "monthlyCost": monthlyCost,
			"currency": currency, "validFrom": nullableString(validFrom), "validTo": nullableString(validTo),
			"checkedAt": nullableString(checkedAt), "allocatedUserCount": allocated,
		})
	}
	return items, rows.Err()
}
