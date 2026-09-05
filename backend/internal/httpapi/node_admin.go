package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// nodeCreateRequest contains central metadata and optional node-owned exit
// addresses. X-Panel credentials stay on the node and are never entered into
// the central panel.
type nodeCreateRequest struct {
	NodeKey       string `json:"nodeKey"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Hostname      string `json:"hostname"`
	ManagementURL string `json:"managementUrl"`
	// PublicIP is retained for older API clients. New UI flows use
	// managementUrl and keep node exit IPs in the exit_ips asset table.
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
	ManagementURL *string `json:"managementUrl"`
	PublicIP      *string `json:"publicIp"`
	Region        *string `json:"region"`
	Provider      *string `json:"provider"`
	PanelBasePath *string `json:"panelBasePath"`
	Enabled       *bool   `json:"enabled"`
}

const nodeInstallTokenTTL = 15 * time.Minute

type nodeAdminValues struct {
	NodeKey       string
	Name          string
	Type          string
	Hostname      string
	ManagementURL string
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
	if len(values.Hostname) > 255 || len(values.ManagementURL) > 2048 || len(values.PublicIP) > 100 || len(values.Region) > 120 || len(values.Provider) > 200 {
		return errors.New("node metadata is too long")
	}
	if values.ManagementURL != "" {
		if _, err := normalizeManagementURL(values.ManagementURL); err != nil {
			return err
		}
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
		ManagementURL: strings.TrimSpace(request.ManagementURL),
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
	if values.ManagementURL != "" {
		normalized, err := normalizeManagementURL(values.ManagementURL)
		if err != nil {
			return nodeAdminValues{}, err
		}
		values.ManagementURL = normalized
		values.PanelBasePath = managementURLBasePath(normalized)
	}
	return values, validateNodeMetadata(values, true)
}

// normalizeManagementURL validates and canonicalizes the full browser/X-Panel
// management URL. It intentionally rejects credentials, queries and fragments
// so a node record can safely be displayed as a clickable maintenance link.
func normalizeManagementURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", nil
	}
	if strings.ContainsAny(value, "\r\n\t ") {
		return "", errors.New("managementUrl must not contain whitespace")
	}
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", errors.New("managementUrl must be an http or https URL with a host")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", errors.New("managementUrl must not contain credentials, query, or fragment")
	}
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	parsed.Path = "/" + strings.Trim(strings.TrimSpace(parsed.Path), "/")
	if parsed.Path == "/" {
		return parsed.String(), nil
	}
	return strings.TrimRight(parsed.String(), "/"), nil
}

func managementURLBasePath(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Path == "" {
		return "/"
	}
	path := "/" + strings.Trim(parsed.Path, "/")
	if path == "/" {
		return "/"
	}
	return path
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
		"hostname": nullableString(values.Hostname), "managementUrl": nullableString(values.ManagementURL), "publicIp": nullableString(values.PublicIP),
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
	installToken, err := randomToken()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not issue installer credential")
		return
	}
	installTokenExpiresAt := time.Now().UTC().Add(nodeInstallTokenTTL).Format(time.RFC3339Nano)
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin node creation")
		return
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO nodes (id, node_key, name, type, hostname, management_url, public_ip, region, provider, panel_base_path, enabled, health_status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 'unknown', ?, ?)`, nodeID, values.NodeKey, values.Name, values.Type, nullableRouteValue(values.Hostname), nullableRouteValue(values.ManagementURL), nullableRouteValue(values.PublicIP), nullableRouteValue(values.Region), nullableRouteValue(values.Provider), values.PanelBasePath, now, now)
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
	if _, err := tx.Exec(`INSERT INTO node_install_tokens (id, node_id, token_hash, expires_at, created_by, created_at) VALUES (?, ?, ?, ?, ?, ?)`, newID(), nodeID, hashToken(installToken), installTokenExpiresAt, currentAdminID(r), now); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save installer credential")
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
		"token": token, "installerToken": installToken, "installerTokenExpiresAt": installTokenExpiresAt,
		"enabled": true, "exitIpCount": len(exitIPs),
	})
}

// issueNodeInstallToken creates a fresh short-lived token for an existing
// node. It is intentionally separate from the long-lived node credential so
// the one-line command can be safely regenerated without exposing the Agent
// bearer token in the browser.
func (s *Server) issueNodeInstallToken(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node id is required")
		return
	}
	var name string
	if err := s.db.QueryRow(`SELECT name FROM nodes WHERE id = ? AND deleted_at IS NULL`, id).Scan(&name); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node")
		return
	}
	installToken, err := randomToken()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not issue installer credential")
		return
	}
	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	expiresAt := now.Add(nodeInstallTokenTTL).Format(time.RFC3339Nano)
	tx, err := s.db.BeginTx(r.Context(), nil)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin installer credential rotation")
		return
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(r.Context(), `UPDATE node_install_tokens SET used_at = ? WHERE node_id = ? AND used_at IS NULL`, nowText, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not revoke previous installer credentials")
		return
	}
	if _, err := tx.ExecContext(r.Context(), `INSERT INTO node_install_tokens (id, node_id, token_hash, expires_at, created_by, created_at) VALUES (?, ?, ?, ?, ?, ?)`, newID(), id, hashToken(installToken), expiresAt, currentAdminID(r), nowText); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not save installer credential")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not complete installer credential rotation")
		return
	}
	s.writeAuditLog(r, "node.install_token.issue", "node", id, nil, map[string]any{"expiresAt": expiresAt})
	writeSuccess(w, map[string]any{"nodeId": id, "nodeName": name, "installerToken": installToken, "installerTokenExpiresAt": expiresAt})
}

func currentAdminID(r *http.Request) any {
	if current, ok := r.Context().Value(principalContextKey{}).(principal); ok && current.UserID != "" {
		return current.UserID
	}
	return nil
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
	err := s.db.QueryRow(`SELECT node_key, name, type, COALESCE(hostname, ''), COALESCE(management_url, ''), COALESCE(public_ip, ''), COALESCE(region, ''), COALESCE(provider, ''), COALESCE(panel_base_path, '/'), enabled, health_status
FROM nodes WHERE id = ? AND deleted_at IS NULL`, id).Scan(&existing.NodeKey, &existing.Name, &existing.Type, &existing.Hostname, &existing.ManagementURL, &existing.PublicIP, &existing.Region, &existing.Provider, &existing.PanelBasePath, &enabled, &existing.HealthStatus)
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
	if request.Name == nil && request.Type == nil && request.Hostname == nil && request.ManagementURL == nil && request.PublicIP == nil && request.Region == nil && request.Provider == nil && request.PanelBasePath == nil && request.Enabled == nil {
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
	if request.ManagementURL != nil {
		updated.ManagementURL = strings.TrimSpace(*request.ManagementURL)
		if updated.ManagementURL != "" {
			normalized, normalizeErr := normalizeManagementURL(updated.ManagementURL)
			if normalizeErr != nil {
				writeFailure(w, http.StatusBadRequest, validationCode, normalizeErr.Error())
				return
			}
			updated.ManagementURL = normalized
			updated.PanelBasePath = managementURLBasePath(normalized)
		}
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
	if _, err := s.db.Exec(`UPDATE nodes SET name = ?, type = ?, hostname = ?, management_url = ?, public_ip = ?, region = ?, provider = ?, panel_base_path = ?, enabled = ?, updated_at = ? WHERE id = ?`,
		updated.Name, updated.Type, nullableRouteValue(updated.Hostname), nullableRouteValue(updated.ManagementURL), nullableRouteValue(updated.PublicIP), nullableRouteValue(updated.Region), nullableRouteValue(updated.Provider), updated.PanelBasePath, boolInt(updated.Enabled), now, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update node")
		return
	}
	if existing.Enabled != updated.Enabled {
		eventType := "node_disabled"
		title := "节点已停用"
		message := "管理员停用了节点"
		if updated.Enabled {
			eventType = "node_enabled"
			title = "节点已启用"
			message = "管理员启用了节点"
		}
		if err := insertNodeEvent(s.db, nodeEventSpec{
			NodeID: id, EventType: eventType, Category: "node", Severity: "info", Title: title, Message: message,
			ResourceType: "node", ResourceID: id, Source: "admin", OccurredAt: time.Now().UTC(),
		}); err != nil {
			s.logger.Warn("record node status change", "node_id", id, "error", err)
		}
	}
	s.writeAuditLog(r, "node.update", "node", id, nodeAdminData(id, existing), nodeAdminData(id, updated))
	writeSuccess(w, nodeAdminData(id, updated))
}

// deleteNode permanently removes a node and the node-owned operational data.
// The old routes tables are compatibility data only; they are cleaned up in
// the same transaction and must not block deletion of a node in the current
// direct user-path model.
func (s *Server) deleteNode(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node id is required")
		return
	}
	ctx := r.Context()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin node deletion")
		return
	}
	defer tx.Rollback()
	var name string
	var enabled int
	if err := tx.QueryRowContext(ctx, `SELECT name, enabled FROM nodes WHERE id = ?`, id).Scan(&name, &enabled); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node")
		return
	}

	// Remove references first so SQLite foreign-key enforcement remains on and
	// the whole operation is atomic. Users themselves are intentionally kept;
	// deleting a node only clears the node-owned connection/path records.
	execDelete := func(label, query string, args ...any) error {
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			s.logger.Error("delete node data", "node_id", id, "resource", label, "error", err)
			return err
		}
		return nil
	}
	if err := execDelete("user paths", `DELETE FROM user_paths WHERE relay_node_id = ? OR landing_node_id = ? OR landing_inbound_id IN (SELECT id FROM inbounds WHERE node_id = ?) OR exit_ip_id IN (SELECT id FROM exit_ips WHERE owner_node_id = ? OR landing_node_id = ?)`, id, id, id, id, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node user paths")
		return
	}
	if err := execDelete("legacy user routes", `DELETE FROM user_routes WHERE route_id IN (SELECT id FROM routes WHERE relay_node_id = ? OR landing_node_id = ?)`, id, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node legacy user routes")
		return
	}
	if err := execDelete("legacy route exit IP bindings", `DELETE FROM route_exit_ips WHERE route_id IN (SELECT id FROM routes WHERE relay_node_id = ? OR landing_node_id = ?) OR exit_ip_id IN (SELECT id FROM exit_ips WHERE owner_node_id = ? OR landing_node_id = ?)`, id, id, id, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node legacy route exit IP bindings")
		return
	}
	if err := execDelete("legacy routes", `DELETE FROM routes WHERE relay_node_id = ? OR landing_node_id = ?`, id, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node legacy routes")
		return
	}
	if err := execDelete("traffic snapshots", `DELETE FROM traffic_snapshots WHERE node_id = ? OR inbound_id IN (SELECT id FROM inbounds WHERE node_id = ?)`, id, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node traffic snapshots")
		return
	}
	if err := execDelete("user inbound bindings", `DELETE FROM user_inbounds WHERE inbound_id IN (SELECT id FROM inbounds WHERE node_id = ?)`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node user inbound bindings")
		return
	}
	if err := execDelete("renewal candidate inbound references", `UPDATE user_renewal_candidates SET inbound_id = NULL WHERE inbound_id IN (SELECT id FROM inbounds WHERE node_id = ?)`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not preserve renewal history")
		return
	}
	if err := execDelete("clients", `DELETE FROM clients WHERE node_id = ? OR inbound_id IN (SELECT id FROM inbounds WHERE node_id = ?)`, id, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node clients")
		return
	}
	if err := execDelete("inbounds", `DELETE FROM inbounds WHERE node_id = ?`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node inbounds")
		return
	}
	if err := execDelete("exit IPs", `DELETE FROM exit_ips WHERE owner_node_id = ? OR landing_node_id = ?`, id, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node exit IPs")
		return
	}
	if err := execDelete("node costs", `DELETE FROM node_costs WHERE node_id = ?`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node costs")
		return
	}
	if err := execDelete("sync runs", `DELETE FROM sync_runs WHERE node_id = ?`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node sync runs")
		return
	}
	if err := execDelete("node events", `DELETE FROM node_events WHERE node_id = ?`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node events")
		return
	}
	if err := execDelete("Agent credentials", `DELETE FROM node_credentials WHERE node_id = ?`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node Agent credentials")
		return
	}
	if err := execDelete("Agent installer credentials", `DELETE FROM node_install_tokens WHERE node_id = ?`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node installer credentials")
		return
	}
	if err := execDelete("node", `DELETE FROM nodes WHERE id = ?`, id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete node")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not complete node deletion")
		return
	}
	s.writeAuditLog(r, "node.delete", "node", id, map[string]any{"name": name, "enabled": enabled == 1}, map[string]any{"deleted": true, "hardDelete": true})
	writeSuccess(w, map[string]any{"id": id, "deleted": true})
}
