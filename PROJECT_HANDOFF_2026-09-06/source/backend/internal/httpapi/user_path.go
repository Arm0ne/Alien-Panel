package httpapi

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"
)

// userPathAssignmentRequest is intentionally smaller than the legacy route
// assignment payload. The relay is derived from the user's primary inbound;
// operators only choose an optional landing node/inbound and one fixed exit IP.
type userPathAssignmentRequest struct {
	LandingNodeID    *string `json:"landingNodeId"`
	LandingInboundID *string `json:"landingInboundId"`
	ExitIPID         string  `json:"exitIpId"`
	Notes            *string `json:"notes"`
}

type userPathRecord struct {
	ID                string
	UserID            string
	RelayNodeID       string
	RelayNodeName     string
	LandingNodeID     string
	LandingNodeName   string
	LandingInboundID  string
	LandingInboundTag string
	ExitIPID          string
	ExitIPAddress     string
	ExitSourceType    string
	ExitOwnerNodeID   string
	ExitOwnerNodeName string
	Mode              string
	Notes             string
	ActiveFrom        string
	ActiveTo          string
	Valid             bool
}

func (s *Server) assignUserPath(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}
	var request userPathAssignmentRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid user path payload")
		return
	}
	exitIPID := strings.TrimSpace(request.ExitIPID)
	if exitIPID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "exitIpId is required")
		return
	}
	if err := s.validateUserExists(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
			return
		}
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user")
		return
	}

	var relayNodeID, relayType string
	var relayEnabled int
	err := s.db.QueryRow(`SELECT i.node_id, n.type, n.enabled
FROM user_inbounds ui JOIN inbounds i ON i.id = ui.inbound_id
JOIN nodes n ON n.id = i.node_id
WHERE ui.user_id = ? AND ui.is_primary = 1 AND ui.active_to IS NULL
  AND i.deleted_at IS NULL AND i.enable = 1 AND i.kind = 'user'
  AND n.type = 'relay' AND n.deleted_at IS NULL`, userID).Scan(&relayNodeID, &relayType, &relayEnabled)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusConflict, validationCode, "user has no active primary inbound")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user's relay node")
		return
	}
	if relayType != "relay" || relayEnabled != 1 {
		writeFailure(w, http.StatusConflict, validationCode, "user's primary inbound node is not an enabled relay")
		return
	}

	landingNodeID := optionalPathValue(request.LandingNodeID)
	landingInboundID := optionalPathValue(request.LandingInboundID)
	if landingInboundID != "" && landingNodeID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "landingNodeId is required when landingInboundId is set")
		return
	}
	if landingNodeID != "" {
		var nodeType string
		var enabled int
		if err := s.db.QueryRow(`SELECT type, enabled FROM nodes WHERE id = ? AND deleted_at IS NULL`, landingNodeID).Scan(&nodeType, &enabled); errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusBadRequest, validationCode, "landing node not found")
			return
		} else if err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read landing node")
			return
		} else if nodeType != "landing" || enabled != 1 {
			writeFailure(w, http.StatusConflict, validationCode, "landing node must be an enabled landing node")
			return
		}
	}
	// A landing node is an execution hop, not a user source.  Selecting it
	// always requires a concrete infrastructure inbound so the path is
	// unambiguous; a relay-direct or external S5 path must not carry landing
	// fields.
	if landingNodeID != "" && landingInboundID == "" {
		writeFailure(w, http.StatusConflict, validationCode, "landingInboundId is required when landingNodeId is set")
		return
	}
	if landingInboundID != "" {
		var inboundNodeID string
		var inboundEnabled int
		var inboundDeleted string
		if err := s.db.QueryRow(`SELECT node_id, enable, COALESCE(deleted_at, '') FROM inbounds WHERE id = ?`, landingInboundID).Scan(&inboundNodeID, &inboundEnabled, &inboundDeleted); errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusBadRequest, validationCode, "landing inbound not found")
			return
		} else if err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read landing inbound")
			return
		} else if inboundNodeID != landingNodeID {
			writeFailure(w, http.StatusBadRequest, validationCode, "landing inbound does not belong to landing node")
			return
		} else if inboundEnabled != 1 || inboundDeleted != "" {
			writeFailure(w, http.StatusConflict, validationCode, "landing inbound must be enabled and not archived")
			return
		}
	}

	var sourceType, ownerNodeID, validFrom, validTo string
	var exitEnabled, ownerEnabled int
	var ownerDeleted string
	err = s.db.QueryRow(`SELECT COALESCE(e.source_type, 'node'), COALESCE(e.owner_node_id, ''), e.enabled,
COALESCE(e.valid_from, ''), COALESCE(e.valid_to, ''), COALESCE(owner.enabled, 1), COALESCE(owner.deleted_at, '')
FROM exit_ips e LEFT JOIN nodes owner ON owner.id = e.owner_node_id
WHERE e.id = ?`, exitIPID).Scan(&sourceType, &ownerNodeID, &exitEnabled, &validFrom, &validTo, &ownerEnabled, &ownerDeleted)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "exit IP not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP")
		return
	}
	if exitEnabled != 1 || ownerEnabled != 1 || ownerDeleted != "" {
		writeFailure(w, http.StatusConflict, validationCode, "selected exit IP is disabled or its owner node is unavailable")
		return
	}
	if !pathDateActive(validFrom, validTo, time.Now().UTC()) {
		writeFailure(w, http.StatusConflict, validationCode, "selected exit IP is outside its validity period")
		return
	}

	mode := "external"
	if sourceType == "node" {
		switch {
		case ownerNodeID == relayNodeID && landingNodeID == "":
			mode = "relay"
		case landingNodeID != "" && ownerNodeID == landingNodeID:
			mode = "landing"
		default:
			writeFailure(w, http.StatusBadRequest, validationCode, "node exit IP must belong to the selected relay or landing node")
			return
		}
	} else if sourceType != "s5" {
		writeFailure(w, http.StatusBadRequest, validationCode, "unsupported exit IP source")
		return
	} else if landingNodeID != "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "external S5 exit IP cannot be combined with a landing node")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	before, _ := s.readUserPath(userID)
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin user path assignment")
		return
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE user_paths SET active_to = ?, updated_at = ? WHERE user_id = ? AND active_to IS NULL`, now, now, userID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close previous user path")
		return
	}
	pathID := newID()
	if _, err := tx.Exec(`INSERT INTO user_paths (id, user_id, relay_node_id, landing_node_id, landing_inbound_id, exit_ip_id, mode, notes, active_from, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, pathID, userID, relayNodeID, nullableRouteValue(landingNodeID), nullableRouteValue(landingInboundID), exitIPID, mode,
		nullableRouteValue(strings.TrimSpace(derefPathValue(request.Notes))), now, now, now); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save user path")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not complete user path assignment")
		return
	}
	s.writeAuditLog(r, "user.path.assign", "user", userID, userPathAuditData(before), map[string]any{"pathId": pathID, "relayNodeId": relayNodeID, "landingNodeId": nullableString(landingNodeID), "landingInboundId": nullableString(landingInboundID), "exitIpId": exitIPID, "mode": mode})
	result, err := s.readUserDetail(userID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "path saved but could not read user")
		return
	}
	writeSuccess(w, result)
}

func (s *Server) clearUserPath(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}
	if err := s.validateUserExists(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
			return
		}
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	before, _ := s.readUserPath(userID)
	if _, err := s.db.Exec(`UPDATE user_paths SET active_to = ?, updated_at = ? WHERE user_id = ? AND active_to IS NULL`, now, now, userID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not clear user path")
		return
	}
	s.writeAuditLog(r, "user.path.clear", "user", userID, userPathAuditData(before), map[string]any{"clearedAt": now})
	result, err := s.readUserDetail(userID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "path cleared but could not read user")
		return
	}
	writeSuccess(w, result)
}

func (s *Server) readUserPath(userID string) (*userPathRecord, error) {
	var record userPathRecord
	var relayDeleted, landingDeleted string
	var exitEnabled int
	err := s.db.QueryRow(`SELECT p.id, p.user_id, p.relay_node_id, COALESCE(relay.name, ''),
COALESCE(p.landing_node_id, ''), COALESCE(landing.name, ''), COALESCE(p.landing_inbound_id, ''), COALESCE(li.tag, ''),
p.exit_ip_id, COALESCE(e.ip, ''), COALESCE(e.source_type, 'node'), COALESCE(e.owner_node_id, ''), COALESCE(owner.name, ''),
p.mode, COALESCE(p.notes, ''), p.active_from, COALESCE(p.active_to, ''), e.enabled,
COALESCE(relay.deleted_at, ''), COALESCE(landing.deleted_at, '')
FROM user_paths p
JOIN nodes relay ON relay.id = p.relay_node_id
LEFT JOIN nodes landing ON landing.id = p.landing_node_id
LEFT JOIN inbounds li ON li.id = p.landing_inbound_id
LEFT JOIN exit_ips e ON e.id = p.exit_ip_id
LEFT JOIN nodes owner ON owner.id = e.owner_node_id
WHERE p.user_id = ? AND p.active_to IS NULL`, userID).Scan(&record.ID, &record.UserID, &record.RelayNodeID, &record.RelayNodeName,
		&record.LandingNodeID, &record.LandingNodeName, &record.LandingInboundID, &record.LandingInboundTag,
		&record.ExitIPID, &record.ExitIPAddress, &record.ExitSourceType, &record.ExitOwnerNodeID, &record.ExitOwnerNodeName,
		&record.Mode, &record.Notes, &record.ActiveFrom, &record.ActiveTo, &exitEnabled, &relayDeleted, &landingDeleted)
	if err != nil {
		return nil, err
	}
	record.Valid = exitEnabled == 1 && relayDeleted == "" && landingDeleted == ""
	return &record, nil
}

func userPathData(record *userPathRecord) map[string]any {
	if record == nil {
		return nil
	}
	return map[string]any{
		"id": record.ID, "relayNodeId": record.RelayNodeID, "relayNodeName": nullableString(record.RelayNodeName),
		"landingNodeId": nullableString(record.LandingNodeID), "landingNodeName": nullableString(record.LandingNodeName),
		"landingInboundId": nullableString(record.LandingInboundID), "landingInboundTag": nullableString(record.LandingInboundTag),
		"exitIpId": record.ExitIPID, "exitIpAddress": nullableString(record.ExitIPAddress), "exitIpSourceType": record.ExitSourceType,
		"exitIpOwnerNodeId": nullableString(record.ExitOwnerNodeID), "exitIpOwnerNodeName": nullableString(record.ExitOwnerNodeName),
		"mode": record.Mode, "notes": nullableString(record.Notes), "activeFrom": nullableString(record.ActiveFrom),
		"activeTo": nullableString(record.ActiveTo), "valid": record.Valid,
	}
}

func userPathAuditData(record *userPathRecord) any {
	return userPathData(record)
}

func optionalPathValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func derefPathValue(value *string) string {
	return optionalPathValue(value)
}

func pathDateActive(fromValue, toValue string, now time.Time) bool {
	from, fromSet, err := parseRouteDate(fromValue)
	if err != nil {
		return false
	}
	to, toSet, err := parseRouteDate(toValue)
	if err != nil {
		return false
	}
	today := now.UTC().Truncate(24 * time.Hour)
	if fromSet && today.Before(from.UTC().Truncate(24*time.Hour)) {
		return false
	}
	if toSet && today.After(to.UTC().Truncate(24*time.Hour)) {
		return false
	}
	return true
}

func (s *Server) userPathHistory(userID string) ([]map[string]any, error) {
	rows, err := s.db.Query(`SELECT p.id, p.user_id, p.relay_node_id, COALESCE(relay.name, ''), COALESCE(p.landing_node_id, ''), COALESCE(landing.name, ''),
	COALESCE(p.landing_inbound_id, ''), COALESCE(li.tag, ''), p.exit_ip_id, COALESCE(e.ip, ''), COALESCE(e.source_type, 'node'), COALESCE(owner.id, ''), COALESCE(owner.name, ''),
p.mode, COALESCE(p.notes, ''), p.active_from, COALESCE(p.active_to, '')
FROM user_paths p JOIN nodes relay ON relay.id = p.relay_node_id
LEFT JOIN nodes landing ON landing.id = p.landing_node_id LEFT JOIN inbounds li ON li.id = p.landing_inbound_id
LEFT JOIN exit_ips e ON e.id = p.exit_ip_id LEFT JOIN nodes owner ON owner.id = e.owner_node_id
WHERE p.user_id = ? ORDER BY p.active_from DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var record userPathRecord
		if err := rows.Scan(&record.ID, &record.UserID, &record.RelayNodeID, &record.RelayNodeName, &record.LandingNodeID, &record.LandingNodeName,
			&record.LandingInboundID, &record.LandingInboundTag, &record.ExitIPID, &record.ExitIPAddress, &record.ExitSourceType, &record.ExitOwnerNodeID, &record.ExitOwnerNodeName,
			&record.Mode, &record.Notes, &record.ActiveFrom, &record.ActiveTo); err != nil {
			return nil, err
		}
		items = append(items, userPathData(&record))
	}
	return items, rows.Err()
}
