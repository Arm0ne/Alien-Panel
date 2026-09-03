package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// nodeCreateRequest contains central metadata and optional node-owned exit
// addresses. X-Panel credentials stay on the node and are never entered into
// the central panel.
type nodeCreateRequest struct {
	NodeKey       string   `json:"nodeKey"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Hostname      string   `json:"hostname"`
	PublicIP      string   `json:"publicIp"`
	ExitIPs       []string `json:"exitIps"`
	Region        string   `json:"region"`
	Provider      string   `json:"provider"`
	PanelBasePath string   `json:"panelBasePath"`
}

type nodeUpdateRequest struct {
	Name          *string `json:"name"`
	Type          *string `json:"type"`
	Hostname      *string `json:"hostname"`
	PublicIP      *string `json:"publicIp"`
	Region        *string `json:"region"`
	Provider      *string `json:"provider"`
	PanelBasePath *string `json:"panelBasePath"`
	Enabled       *bool   `json:"enabled"`
}

type nodeAdminValues struct {
	NodeKey       string
	Name          string
	Type          string
	Hostname      string
	PublicIP      string
	Region        string
	Provider      string
	PanelBasePath string
	Enabled       bool
	HealthStatus  string
}

func validateNodeMetadata(values nodeAdminValues, requireType bool) error {
	if values.NodeKey == "" {
		return errors.New("nodeKey is required")
	}
	if len(values.NodeKey) > 120 {
		return errors.New("nodeKey is too long")
	}
	if values.Name == "" {
		return errors.New("name is required")
	}
	if len(values.Name) > 120 {
		return errors.New("name is too long")
	}
	if requireType && values.Type != "relay" && values.Type != "landing" {
		return errors.New("type must be relay or landing")
	}
	if values.Type != "" && values.Type != "relay" && values.Type != "landing" && values.Type != "unknown" {
		return errors.New("type must be relay, landing, or unknown")
	}
	if len(values.Hostname) > 255 || len(values.PublicIP) > 100 || len(values.Region) > 120 || len(values.Provider) > 200 {
		return errors.New("node metadata is too long")
	}
	if values.PublicIP != "" && net.ParseIP(values.PublicIP) == nil {
		return errors.New("publicIp must be a valid IP address")
	}
	if values.PanelBasePath == "" {
		return errors.New("panelBasePath is required")
	}
	if len(values.PanelBasePath) > 200 || !strings.HasPrefix(values.PanelBasePath, "/") || strings.Contains(values.PanelBasePath, "://") {
		return errors.New("panelBasePath must be a URL path")
	}
	return nil
}

func nodeCreateValues(request nodeCreateRequest) (nodeAdminValues, error) {
	values := nodeAdminValues{
		NodeKey:       strings.TrimSpace(request.NodeKey),
		Name:          strings.TrimSpace(request.Name),
		Type:          strings.TrimSpace(request.Type),
		Hostname:      strings.TrimSpace(request.Hostname),
		PublicIP:      strings.TrimSpace(request.PublicIP),
		Region:        strings.TrimSpace(request.Region),
		Provider:      strings.TrimSpace(request.Provider),
		PanelBasePath: strings.TrimSpace(request.PanelBasePath),
		Enabled:       true,
		HealthStatus:  "unknown",
	}
	if values.PanelBasePath == "" {
		values.PanelBasePath = "/"
	}
	if !strings.HasPrefix(values.PanelBasePath, "/") {
		values.PanelBasePath = "/" + values.PanelBasePath
	}
	if values.Type == "" {
		return nodeAdminValues{}, errors.New("type is required")
	}
	if values.NodeKey == "" {
		generated, err := randomToken()
		if err != nil {
			return nodeAdminValues{}, errors.New("could not generate node key")
		}
		// Node Key is an identifier, not a credential. A short random suffix
		// keeps it easy to copy while avoiding collisions across nodes.
		values.NodeKey = "node-" + generated[:16]
	}
	return values, validateNodeMetadata(values, true)
}

func nodeAdminData(id string, values nodeAdminValues) map[string]any {
	status := values.HealthStatus
	if status == "" {
		status = "unknown"
	}
	if !values.Enabled {
		status = "disabled"
	}
	return map[string]any{
		"id": id, "nodeKey": values.NodeKey, "name": values.Name, "type": values.Type,
		"hostname": nullableString(values.Hostname), "publicIp": nullableString(values.PublicIP),
		"region": nullableString(values.Region), "provider": nullableString(values.Provider),
		"panelBasePath": values.PanelBasePath, "enabled": values.Enabled, "status": status,
	}
}

func (s *Server) createNode(w http.ResponseWriter, r *http.Request) {
	var request nodeCreateRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid node payload")
		return
	}
	values, err := nodeCreateValues(request)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	exitIPs, err := normalizeNodeExitIPs(request.ExitIPs)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}

	var existingID string
	if err := s.db.QueryRow(`SELECT id FROM nodes WHERE node_key = ?`, values.NodeKey).Scan(&existingID); err == nil {
		writeFailure(w, http.StatusConflict, validationCode, "nodeKey already exists")
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check node key")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	nodeID := newID()
	token, err := randomToken()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not issue node credential")
		return
	}
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin node creation")
		return
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO nodes (id, node_key, name, type, hostname, public_ip, region, provider, panel_base_path, enabled, health_status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'unknown', ?, ?)`, nodeID, values.NodeKey, values.Name, values.Type, nullableRouteValue(values.Hostname), nullableRouteValue(values.PublicIP), nullableRouteValue(values.Region), nullableRouteValue(values.Provider), values.PanelBasePath, now, now)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			writeFailure(w, http.StatusConflict, validationCode, "nodeKey already exists")
		} else {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create node")
		}
		return
	}
	if _, err := tx.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES (?, ?, ?, ?, ?)`, newID(), nodeID, hashToken(token), now, now); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save node credential")
		return
	}
	for _, address := range exitIPs {
		family := 4
		if net.ParseIP(address).To4() == nil {
			family = 6
		}
		if _, err := tx.Exec(`INSERT INTO exit_ips (id, source_type, owner_node_id, ip, family, monthly_cost, currency, enabled, created_at, updated_at)
VALUES (?, 'node', ?, ?, ?, 0, 'CNY', 1, ?, ?)`, newID(), nodeID, address, family, now, now); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "unique") {
				writeFailure(w, http.StatusConflict, validationCode, "exit IP already exists for this node")
			} else {
				writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save node exit IP")
			}
			return
		}
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not complete node creation")
		return
	}

	s.writeAuditLog(r, "node.create", "node", nodeID, nil, nodeAdminData(nodeID, values))
	// The raw token is deliberately returned only in this response. It is not
	// included in the audit record or any later node detail/list response.
	writeSuccess(w, map[string]any{
		"nodeId": nodeID, "nodeKey": values.NodeKey, "name": values.Name, "type": values.Type,
		"token": token, "enabled": true, "exitIpCount": len(exitIPs),
	})
}

func normalizeNodeExitIPs(addresses []string) ([]string, error) {
	if len(addresses) > 100 {
		return nil, errors.New("exitIps must contain at most 100 addresses")
	}
	result := make([]string, 0, len(addresses))
	seen := make(map[string]struct{}, len(addresses))
	for _, raw := range addresses {
		address := strings.TrimSpace(raw)
		if address == "" {
			continue
		}
		parsed := net.ParseIP(address)
		if parsed == nil {
			return nil, fmt.Errorf("exit IP %q must be a valid IPv4 or IPv6 address", address)
		}
		if parsed.To4() != nil {
			address = parsed.To4().String()
		} else {
			address = parsed.String()
		}
		if _, exists := seen[address]; exists {
			return nil, fmt.Errorf("exit IP %q is duplicated", address)
		}
		seen[address] = struct{}{}
		result = append(result, address)
	}
	return result, nil
}

func (s *Server) updateNode(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node id is required")
		return
	}

	var existing nodeAdminValues
	var enabled int
	err := s.db.QueryRow(`SELECT node_key, name, type, COALESCE(hostname, ''), COALESCE(public_ip, ''), COALESCE(region, ''), COALESCE(provider, ''), COALESCE(panel_base_path, '/'), enabled, health_status
FROM nodes WHERE id = ? AND deleted_at IS NULL`, id).Scan(&existing.NodeKey, &existing.Name, &existing.Type, &existing.Hostname, &existing.PublicIP, &existing.Region, &existing.Provider, &existing.PanelBasePath, &enabled, &existing.HealthStatus)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node")
		return
	}
	existing.Enabled = enabled == 1

	var request nodeUpdateRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid node payload")
		return
	}
	if request.Name == nil && request.Type == nil && request.Hostname == nil && request.PublicIP == nil && request.Region == nil && request.Provider == nil && request.PanelBasePath == nil && request.Enabled == nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "node update payload is empty")
		return
	}
	updated := existing
	if request.Name != nil {
		updated.Name = strings.TrimSpace(*request.Name)
	}
	if request.Type != nil {
		updated.Type = strings.TrimSpace(*request.Type)
	}
	if request.Hostname != nil {
		updated.Hostname = strings.TrimSpace(*request.Hostname)
	}
	if request.PublicIP != nil {
		updated.PublicIP = strings.TrimSpace(*request.PublicIP)
	}
	if request.Region != nil {
		updated.Region = strings.TrimSpace(*request.Region)
	}
	if request.Provider != nil {
		updated.Provider = strings.TrimSpace(*request.Provider)
	}
	if request.PanelBasePath != nil {
		updated.PanelBasePath = strings.TrimSpace(*request.PanelBasePath)
		if updated.PanelBasePath == "" {
			updated.PanelBasePath = "/"
		}
		if !strings.HasPrefix(updated.PanelBasePath, "/") {
			updated.PanelBasePath = "/" + updated.PanelBasePath
		}
	}
	if request.Enabled != nil {
		updated.Enabled = *request.Enabled
	}
	if err := validateNodeMetadata(updated, false); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`UPDATE nodes SET name = ?, type = ?, hostname = ?, public_ip = ?, region = ?, provider = ?, panel_base_path = ?, enabled = ?, updated_at = ? WHERE id = ?`,
		updated.Name, updated.Type, nullableRouteValue(updated.Hostname), nullableRouteValue(updated.PublicIP), nullableRouteValue(updated.Region), nullableRouteValue(updated.Provider), updated.PanelBasePath, boolInt(updated.Enabled), now, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update node")
		return
	}
	if existing.Enabled != updated.Enabled {
		eventType := "node_disabled"
		message := "Node disabled by administrator"
		if updated.Enabled {
			eventType = "node_enabled"
			message = "Node enabled by administrator"
		}
		if _, err := s.db.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES (?, ?, ?, 'info', ?, ?)`, newID(), id, eventType, message, now); err != nil {
			s.logger.Warn("record node status change", "node_id", id, "error", err)
		}
	}
	s.writeAuditLog(r, "node.update", "node", id, nodeAdminData(id, existing), nodeAdminData(id, updated))
	writeSuccess(w, nodeAdminData(id, updated))
}

// deleteNode hides a node from active management and revokes all of its Agent
// credentials. Operational history remains intact so audits and historical
// reports do not lose their context. Routes and exit IP assets are protected
// from becoming orphaned and must be removed or unbound first.
func (s *Server) deleteNode(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node id is required")
		return
	}
	var name string
	var enabled int
	if err := s.db.QueryRow(`SELECT name, enabled FROM nodes WHERE id = ? AND deleted_at IS NULL`, id).Scan(&name, &enabled); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node")
		return
	}
	var routeCount, exitIPCount int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM routes WHERE relay_node_id = ? OR landing_node_id = ?`, id, id).Scan(&routeCount); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check node routes")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM exit_ips WHERE COALESCE(owner_node_id, landing_node_id) = ?`, id).Scan(&exitIPCount); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check node exit IPs")
		return
	}
	if routeCount > 0 || exitIPCount > 0 {
		parts := make([]string, 0, 2)
		if routeCount > 0 {
			parts = append(parts, fmt.Sprintf("%d 条线路", routeCount))
		}
		if exitIPCount > 0 {
			parts = append(parts, fmt.Sprintf("%d 个出口 IP", exitIPCount))
		}
		writeFailure(w, http.StatusConflict, validationCode, "节点仍有关联的"+strings.Join(parts, "和")+"，请先移除或解除绑定")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin node deletion")
		return
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE nodes SET enabled = 0, deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`, now, now, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node")
		return
	}
	if _, err := tx.Exec(`UPDATE node_credentials SET revoked_at = ? WHERE node_id = ? AND revoked_at IS NULL`, now, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not revoke node credentials")
		return
	}
	if _, err := tx.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES (?, ?, 'node_deleted', 'warning', ?, ?)`, newID(), id, "节点已删除，Agent 凭据已撤销，历史数据保留", now); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not record node deletion")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not complete node deletion")
		return
	}
	s.writeAuditLog(r, "node.delete", "node", id, map[string]any{"name": name, "enabled": enabled == 1}, map[string]any{"deleted": true, "deletedAt": now, "historyPreserved": true})
	writeSuccess(w, map[string]any{"id": id, "deleted": true, "deletedAt": now, "historyPreserved": true})
}
