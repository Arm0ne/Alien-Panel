package httpapi

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"xpanel-central/backend/internal/config"
)

const (
	successCode       = "0000"
	unauthorizedCode  = "8888"
	validationCode    = "4000"
	csrfCode          = "4403"
	notFoundCode      = "4040"
	internalErrorCode = "5000"
)

type Server struct {
	cfg     config.Config
	db      *sql.DB
	logger  *slog.Logger
	origins map[string]struct{}
}

type principal struct {
	SessionID string
	UserID    string
	UserName  string
}

type principalContextKey struct{}
type requestIDContextKey struct{}

func NewServer(cfg config.Config, database *sql.DB, logger *slog.Logger) (*Server, error) {
	if logger == nil {
		logger = slog.Default()
	}
	server := &Server{cfg: cfg, db: database, logger: logger, origins: make(map[string]struct{}, len(cfg.CorsOrigins))}
	for _, origin := range cfg.CorsOrigins {
		server.origins[origin] = struct{}{}
	}
	if err := server.ensureAdmin(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", s.healthLive)
	mux.HandleFunc("GET /health/ready", s.healthReady)
	mux.HandleFunc("POST /agent/v1/register", s.agentRegister)
	mux.HandleFunc("POST /api/agent/v1/register", s.agentRegister)
	mux.HandleFunc("POST /agent/v1/bootstrap", s.agentBootstrap)
	mux.HandleFunc("POST /api/agent/v1/bootstrap", s.agentBootstrap)
	mux.Handle("POST /agent/v1/heartbeat", s.requireAgent(http.HandlerFunc(s.agentHeartbeat)))
	mux.Handle("POST /api/agent/v1/heartbeat", s.requireAgent(http.HandlerFunc(s.agentHeartbeat)))
	mux.Handle("POST /agent/v1/sync", s.requireAgent(http.HandlerFunc(s.agentSync)))
	mux.Handle("POST /api/agent/v1/sync", s.requireAgent(http.HandlerFunc(s.agentSync)))
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.HandleFunc("POST /api/auth/refreshToken", s.refreshToken)
	mux.Handle("POST /api/auth/logout", s.requireAuth(http.HandlerFunc(s.logout)))
	mux.Handle("GET /api/auth/me", s.requireAuth(http.HandlerFunc(s.getUserInfo)))
	mux.Handle("GET /api/auth/getUserInfo", s.requireAuth(http.HandlerFunc(s.getUserInfo)))
	mux.Handle("GET /api/dashboard", s.requireAuth(http.HandlerFunc(s.dashboard)))
	mux.Handle("GET /api/users", s.requireAuth(http.HandlerFunc(s.users)))
	mux.Handle("GET /api/users/{id}", s.requireAuth(http.HandlerFunc(s.userDetail)))
	mux.Handle("GET /api/users/{id}/renewals", s.requireAuth(http.HandlerFunc(s.listUserRenewals)))
	mux.Handle("POST /api/users/{id}/renewals/{candidateId}/confirm", s.requireAuth(http.HandlerFunc(s.confirmUserRenewal)))
	mux.Handle("POST /api/users/{id}/renewals/{candidateId}/reject", s.requireAuth(http.HandlerFunc(s.rejectUserRenewal)))
	mux.Handle("GET /api/users/{id}/traffic", s.requireAuth(http.HandlerFunc(s.userTraffic)))
	mux.Handle("GET /api/users/{id}/path-assets", s.requireAuth(http.HandlerFunc(s.userPathAssets)))
	mux.Handle("PATCH /api/users/{id}", s.requireAuth(http.HandlerFunc(s.updateUser)))
	mux.Handle("PUT /api/users/{id}/path", s.requireAuth(http.HandlerFunc(s.assignUserPath)))
	mux.Handle("DELETE /api/users/{id}/path", s.requireAuth(http.HandlerFunc(s.clearUserPath)))
	mux.Handle("PUT /api/users/{id}/route", s.requireAuth(http.HandlerFunc(s.assignUserRoute)))
	mux.Handle("DELETE /api/users/{id}/route", s.requireAuth(http.HandlerFunc(s.clearUserRoute)))
	mux.Handle("GET /api/nodes", s.requireAuth(http.HandlerFunc(s.nodes)))
	mux.Handle("POST /api/nodes", s.requireAuth(http.HandlerFunc(s.createNode)))
	mux.Handle("POST /api/nodes/{id}/install-token", s.requireAuth(http.HandlerFunc(s.issueNodeInstallToken)))
	mux.Handle("PATCH /api/nodes/{id}", s.requireAuth(http.HandlerFunc(s.updateNode)))
	mux.Handle("DELETE /api/nodes/{id}", s.requireAuth(http.HandlerFunc(s.deleteNode)))
	mux.Handle("GET /api/nodes/{id}", s.requireAuth(http.HandlerFunc(s.nodeDetail)))
	mux.Handle("POST /api/nodes/{id}/sync", s.requireAuth(http.HandlerFunc(s.requestNodeSync)))
	mux.Handle("POST /api/nodes/{id}/costs", s.requireAuth(http.HandlerFunc(s.createNodeCost)))
	mux.Handle("GET /api/nodes/{id}/costs", s.requireAuth(http.HandlerFunc(s.listNodeCosts)))
	mux.Handle("PATCH /api/nodes/{id}/costs/{costId}", s.requireAuth(http.HandlerFunc(s.updateNodeCost)))
	mux.Handle("GET /api/routes", s.requireAuth(http.HandlerFunc(s.routes)))
	mux.Handle("GET /api/routes/{id}", s.requireAuth(http.HandlerFunc(s.routeDetail)))
	mux.Handle("GET /api/routes/{id}/exit-ips", s.requireAuth(http.HandlerFunc(s.routeExitIPs)))
	mux.Handle("POST /api/routes/{id}/exit-ips", s.requireAuth(http.HandlerFunc(s.bindRouteExitIP)))
	mux.Handle("PATCH /api/routes/{id}/exit-ips/{exitIpId}", s.requireAuth(http.HandlerFunc(s.updateRouteExitIP)))
	mux.Handle("DELETE /api/routes/{id}/exit-ips/{exitIpId}", s.requireAuth(http.HandlerFunc(s.unbindRouteExitIP)))
	mux.Handle("POST /api/routes", s.requireAuth(http.HandlerFunc(s.createRoute)))
	mux.Handle("PATCH /api/routes/{id}", s.requireAuth(http.HandlerFunc(s.updateRoute)))
	mux.Handle("DELETE /api/routes/{id}", s.requireAuth(http.HandlerFunc(s.deleteRoute)))
	mux.Handle("GET /api/exit-ips", s.requireAuth(http.HandlerFunc(s.exitIPs)))
	mux.Handle("GET /api/exit-ips/{id}", s.requireAuth(http.HandlerFunc(s.exitIPDetail)))
	mux.Handle("POST /api/exit-ips", s.requireAuth(http.HandlerFunc(s.createExitIP)))
	mux.Handle("PATCH /api/exit-ips/{id}", s.requireAuth(http.HandlerFunc(s.updateExitIP)))
	mux.Handle("DELETE /api/exit-ips/{id}", s.requireAuth(http.HandlerFunc(s.deleteExitIP)))
	mux.Handle("GET /api/costs/summary", s.requireAuth(http.HandlerFunc(s.finance)))
	mux.Handle("GET /api/events", s.requireAuth(http.HandlerFunc(s.events)))

	return s.withSecurityHeaders(s.withRequestID(s.withCORS(mux)))
}

func (s *Server) ensureAdmin() error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM admin_users`).Scan(&count); err != nil {
		return fmt.Errorf("count admin users: %w", err)
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(s.cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err = s.db.Exec(
		`INSERT INTO admin_users (id, username, password_hash, enabled, created_at, updated_at) VALUES (?, ?, ?, 1, ?, ?)`,
		newID(), s.cfg.AdminUsername, string(hash), now, now,
	)
	if err != nil {
		return fmt.Errorf("create initial admin user: %w", err)
	}
	s.logger.Warn("created initial administrator; change the password through deployment secrets", "username", s.cfg.AdminUsername)
	return nil
}

func (s *Server) healthLive(w http.ResponseWriter, _ *http.Request) {
	writeSuccess(w, map[string]string{"status": "ok"})
}

func (s *Server) healthReady(w http.ResponseWriter, _ *http.Request) {
	if err := s.db.Ping(); err != nil {
		writeFailure(w, http.StatusServiceUnavailable, internalErrorCode, "database is not ready")
		return
	}
	writeSuccess(w, map[string]string{"status": "ready"})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &payload); err != nil || strings.TrimSpace(payload.UserName) == "" || payload.Password == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "username and password are required")
		return
	}

	var userID, passwordHash string
	var enabled int
	err := s.db.QueryRow(`SELECT id, password_hash, enabled FROM admin_users WHERE username = ?`, strings.TrimSpace(payload.UserName)).Scan(&userID, &passwordHash, &enabled)
	if err != nil || enabled != 1 || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(payload.Password)) != nil {
		writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "invalid administrator credentials")
		return
	}

	token, refreshToken, expiresAt, err := s.createSession(userID)
	if err != nil {
		s.logger.Error("create login session", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create login session")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`UPDATE admin_users SET last_login_at = ?, updated_at = ? WHERE id = ?`, now, now, userID); err != nil {
		s.logger.Warn("update last login time", "error", err)
	}

	writeSuccess(w, map[string]any{
		"token":        token,
		"refreshToken": refreshToken,
		"expiresAt":    expiresAt,
	})
}

func (s *Server) refreshToken(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &payload); err != nil || payload.RefreshToken == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "refreshToken is required")
		return
	}

	var sessionID, userID string
	var expiresAtRaw string
	err := s.db.QueryRow(`SELECT id, admin_user_id, expires_at FROM sessions WHERE refresh_hash = ? AND revoked_at IS NULL`, hashToken(payload.RefreshToken)).Scan(&sessionID, &userID, &expiresAtRaw)
	if err != nil || isExpired(expiresAtRaw) {
		writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "refresh token is invalid or expired")
		return
	}
	if _, err := s.db.Exec(`UPDATE sessions SET revoked_at = ? WHERE id = ?`, time.Now().UTC().Format(time.RFC3339Nano), sessionID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not rotate refresh token")
		return
	}
	token, refreshTokenValue, expiresAt, err := s.createSession(userID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create refreshed session")
		return
	}
	writeSuccess(w, map[string]any{"token": token, "refreshToken": refreshTokenValue, "expiresAt": expiresAt})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if current, ok := r.Context().Value(principalContextKey{}).(principal); ok {
		_, _ = s.db.Exec(`UPDATE sessions SET revoked_at = ? WHERE id = ?`, time.Now().UTC().Format(time.RFC3339Nano), current.SessionID)
	}
	writeSuccess(w, map[string]bool{"loggedOut": true})
}

func (s *Server) getUserInfo(w http.ResponseWriter, r *http.Request) {
	current := r.Context().Value(principalContextKey{}).(principal)
	writeSuccess(w, map[string]any{
		"userId":   current.UserID,
		"userName": current.UserName,
		"roles":    []string{"R_SUPER"},
		"buttons":  []string{},
	})
}

func (s *Server) dashboard(w http.ResponseWriter, _ *http.Request) {
	s.refreshOperationalStatuses(time.Now().UTC())
	var response struct {
		GeneratedAt string  `json:"generatedAt"`
		DataAt      *string `json:"dataAt"`
		Nodes       struct {
			Total   int `json:"total"`
			Online  int `json:"online"`
			Relay   int `json:"relay"`
			Landing int `json:"landing"`
		} `json:"nodes"`
		Users struct {
			Active   int `json:"active"`
			Expiring int `json:"expiring"`
			Expired  int `json:"expired"`
		} `json:"users"`
		Traffic struct {
			TodayBytes int64 `json:"todayBytes"`
			MonthBytes int64 `json:"monthBytes"`
		} `json:"traffic"`
		Finance financeResponse `json:"finance"`
	}
	response.GeneratedAt = time.Now().UTC().Format(time.RFC3339Nano)
	response.DataAt = s.latestSuccessfulSyncAt()
	queries := []struct {
		target *int
		query  string
	}{
		{&response.Nodes.Total, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL`},
		{&response.Nodes.Online, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL AND health_status = 'online' AND enabled = 1`},
		{&response.Nodes.Relay, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL AND type = 'relay' AND enabled = 1`},
		{&response.Nodes.Landing, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL AND type = 'landing' AND enabled = 1`},
		{&response.Users.Active, `SELECT COUNT(*) FROM users WHERE status = 'active'`},
		{&response.Users.Expiring, `SELECT COUNT(*) FROM users WHERE status = 'expiring'`},
		{&response.Users.Expired, `SELECT COUNT(*) FROM users WHERE status = 'expired'`},
	}
	for _, item := range queries {
		if err := s.db.QueryRow(item.query).Scan(item.target); err != nil {
			s.logger.Error("dashboard query", "query", item.query, "error", err)
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard")
			return
		}
	}
	now := time.Now().UTC()
	startToday := now.Truncate(24 * time.Hour)
	startMonth := now.AddDate(0, 0, -30)
	var err error
	response.Traffic.TodayBytes, err = s.trafficDeltaSince(startToday, now)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read traffic summary")
		return
	}
	response.Traffic.MonthBytes, err = s.trafficDeltaSince(startMonth, now)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read traffic summary")
		return
	}
	finance, err := s.financeSummary(time.Now().UTC().Format("2006-01"))
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read finance summary")
		return
	}
	response.Finance = finance
	writeSuccess(w, response)
}

func (s *Server) trafficDeltaSince(start, end time.Time) (int64, error) {
	rows, err := s.db.Query(`SELECT t.all_time, t.reset_detected,
COALESCE((SELECT p.all_time FROM traffic_snapshots p WHERE p.inbound_id = t.inbound_id AND p.collected_at < t.collected_at ORDER BY p.collected_at DESC LIMIT 1), -1)
FROM traffic_snapshots t WHERE t.collected_at >= ? AND t.collected_at < ?`, start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var total int64
	for rows.Next() {
		var current, previous int64
		var resetDetected int
		if err := rows.Scan(&current, &resetDetected, &previous); err != nil {
			return 0, err
		}
		if resetDetected == 1 || current < previous {
			total += current
			continue
		}
		if previous < 0 {
			// The first sample is a baseline; without an earlier point its
			// period usage cannot be determined yet.
			continue
		}
		total += current - previous
	}
	return total, rows.Err()
}

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	s.refreshOperationalStatuses(time.Now().UTC())
	query := parseListQuery(r)
	where := []string{"1 = 1"}
	args := make([]any, 0, 6)
	if query.keyword != "" {
		where = append(where, `(u.display_name LIKE ? OR i.tag LIKE ? OR n.name LIKE ? OR EXISTS (SELECT 1 FROM user_routes uk JOIN routes rk ON rk.id = uk.route_id WHERE uk.user_id = u.id AND uk.active_to IS NULL AND rk.name LIKE ?))`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like, like)
	}
	if query.status != "" {
		where = append(where, "u.status = ?")
		args = append(args, query.status)
	}
	if query.nodeID != "" {
		where = append(where, "i.node_id = ?")
		args = append(args, query.nodeID)
	}
	base := `FROM users u
LEFT JOIN user_inbounds ui ON ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL
LEFT JOIN inbounds i ON i.id = ui.inbound_id AND i.kind = 'user' AND i.deleted_at IS NULL
 AND i.node_id IN (SELECT id FROM nodes WHERE type = 'relay' AND deleted_at IS NULL)
LEFT JOIN nodes n ON n.id = i.node_id
WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count users")
		return
	}
	rows, err := s.db.Query(`SELECT u.id, u.display_name, u.status, COALESCE(u.expiry_time, ''),
COALESCE(i.node_id, ''), COALESCE(n.name, ''), COALESCE(i.tag, ''),
COALESCE((SELECT r.name FROM user_routes ur JOIN routes r ON r.id = ur.route_id
 WHERE ur.user_id = u.id AND ur.is_primary = 1 AND ur.active_to IS NULL ORDER BY ur.active_from DESC LIMIT 1), ''),
COALESCE(i.client_count, 0), COALESCE(i.up, 0) + COALESCE(i.down, 0), COALESCE(i.last_seen_at, ''),
COALESCE((SELECT n2.name FROM user_paths p LEFT JOIN nodes n2 ON n2.id = p.landing_node_id WHERE p.user_id = u.id AND p.active_to IS NULL LIMIT 1), ''),
COALESCE((SELECT e2.ip FROM user_paths p JOIN exit_ips e2 ON e2.id = p.exit_ip_id WHERE p.user_id = u.id AND p.active_to IS NULL LIMIT 1), ''),
COALESCE((SELECT owner2.name FROM user_paths p JOIN exit_ips e2 ON e2.id = p.exit_ip_id LEFT JOIN nodes owner2 ON owner2.id = e2.owner_node_id WHERE p.user_id = u.id AND p.active_to IS NULL LIMIT 1), ''),
COALESCE((SELECT p.mode FROM user_paths p WHERE p.user_id = u.id AND p.active_to IS NULL LIMIT 1), ''),
COALESCE((SELECT p.id FROM user_paths p WHERE p.user_id = u.id AND p.active_to IS NULL LIMIT 1), '') `+base+`
ORDER BY CASE WHEN u.expiry_time IS NULL THEN 1 ELSE 0 END, u.expiry_time ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read users")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, name, status, expiry, nodeID, nodeName, inboundTag, routeName, lastActivity, landingName, exitIP, exitOwner, pathMode, pathID string
		var clientCount int
		var traffic int64
		if err := rows.Scan(&id, &name, &status, &expiry, &nodeID, &nodeName, &inboundTag, &routeName, &clientCount, &traffic, &lastActivity, &landingName, &exitIP, &exitOwner, &pathMode, &pathID); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode users")
			return
		}
		items = append(items, map[string]any{
			"id": id, "name": name, "nodeId": nodeID, "nodeName": nodeName, "inboundTag": inboundTag, "routeName": nullableString(routeName),
			"status": status, "expiresAt": nullableString(expiry), "clientCount": clientCount,
			"trafficBytes": traffic, "lastActivityAt": nullableString(lastActivity),
			"landingNodeName": nullableString(landingName), "exitIpAddress": nullableString(exitIP),
			"exitIpOwnerNodeName": nullableString(exitOwner), "pathMode": nullableString(pathMode),
			"pathConfigured": pathID != "",
		})
	}
	writeSuccess(w, s.pageResponse(items, total, query))
}

type updateUserRequest struct {
	DisplayName   *string  `json:"displayName"`
	MonthlyFee    *float64 `json:"monthlyFee"` // legacy monthly equivalent
	BillingCycle  *string  `json:"billingCycle"`
	BillingAmount *float64 `json:"billingAmount"`
	Currency      *string  `json:"currency"`
	Notes         *string  `json:"notes"`
}

// userRouteAssignmentRequest controls the central routing decision for a
// business user.  The route itself remains a reusable template; an optional
// routeExitIpId pins this user to one bound exit IP.  When it is omitted the
// route's enabled weighted exit-IP pool is used.
type userRouteAssignmentRequest struct {
	RouteID       string  `json:"routeId"`
	RouteExitIPID *string `json:"routeExitIpId"`
}

func (s *Server) assignUserRoute(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}
	var payload userRouteAssignmentRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid user route payload")
		return
	}
	routeID := strings.TrimSpace(payload.RouteID)
	if routeID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "routeId is required")
		return
	}
	if err := s.validateUserExists(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
			return
		}
		s.logger.Error("read user for route assignment", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
		return
	}
	var userRelayNodeID string
	if err := s.db.QueryRow(`SELECT COALESCE((SELECT i.node_id FROM user_inbounds ui
JOIN inbounds i ON i.id = ui.inbound_id
WHERE ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL LIMIT 1), '')
FROM users u WHERE u.id = ?`, userID).Scan(&userRelayNodeID); err != nil {
		s.logger.Error("read user's relay node for route assignment", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
		return
	}
	var routeEnabled int
	var routeRelayNodeID string
	if err := s.db.QueryRow(`SELECT r.enabled, r.relay_node_id FROM routes r
JOIN nodes relay ON relay.id = r.relay_node_id AND relay.deleted_at IS NULL
JOIN nodes landing ON landing.id = r.landing_node_id AND landing.deleted_at IS NULL
WHERE r.id = ?`, routeID).Scan(&routeEnabled, &routeRelayNodeID); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		return
	} else if err != nil {
		s.logger.Error("read route for user assignment", "route_id", routeID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
		return
	}
	if routeEnabled != 1 {
		writeFailure(w, http.StatusConflict, validationCode, "route is disabled")
		return
	}
	if userRelayNodeID != "" && routeRelayNodeID != userRelayNodeID {
		writeFailure(w, http.StatusBadRequest, validationCode, "route relay node must match user's primary inbound node")
		return
	}
	var routeExitIPID string
	if payload.RouteExitIPID != nil {
		routeExitIPID = strings.TrimSpace(*payload.RouteExitIPID)
	}
	if routeExitIPID != "" {
		var bindingEnabled, exitEnabled int
		if err := s.db.QueryRow(`SELECT rei.enabled, e.enabled
FROM route_exit_ips rei JOIN exit_ips e ON e.id = rei.exit_ip_id
WHERE rei.id = ? AND rei.route_id = ?`, routeExitIPID, routeID).Scan(&bindingEnabled, &exitEnabled); errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusBadRequest, validationCode, "routeExitIpId is not bound to route")
			return
		} else if err != nil {
			s.logger.Error("read route exit IP for user assignment", "route_id", routeID, "route_exit_ip_id", routeExitIPID, "error", err)
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
			return
		} else if bindingEnabled != 1 || exitEnabled != 1 {
			writeFailure(w, http.StatusConflict, validationCode, "selected exit IP is disabled")
			return
		}
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Error("begin user route assignment", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
		return
	}
	deactivate := `UPDATE user_routes SET is_primary = 0, active_to = ?, route_exit_ip_id = NULL
WHERE user_id = ? AND active_to IS NULL AND route_id <> ?`
	if _, err := tx.Exec(deactivate, now, userID, routeID); err != nil {
		_ = tx.Rollback()
		s.logger.Error("deactivate previous user route", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
		return
	}
	var assignmentID string
	err = tx.QueryRow(`SELECT id FROM user_routes WHERE user_id = ? AND route_id = ?`, userID, routeID).Scan(&assignmentID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		assignmentID = newID()
		_, err = tx.Exec(`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from, active_to, route_exit_ip_id)
VALUES (?, ?, ?, 1, ?, NULL, ?)`, assignmentID, userID, routeID, now, nullableRouteValue(routeExitIPID))
	case err == nil:
		_, err = tx.Exec(`UPDATE user_routes SET is_primary = 1, active_from = ?, active_to = NULL, route_exit_ip_id = ? WHERE id = ?`, now, nullableRouteValue(routeExitIPID), assignmentID)
	}
	if err != nil {
		_ = tx.Rollback()
		s.logger.Error("save user route assignment", "user_id", userID, "route_id", routeID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
		return
	}
	if err := tx.Commit(); err != nil {
		s.logger.Error("commit user route assignment", "user_id", userID, "route_id", routeID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not assign user route")
		return
	}
	s.writeAuditLog(r, "user.route.assign", "user", userID, nil, map[string]any{"routeId": routeID, "routeExitIpId": nullableString(routeExitIPID)})
	result, err := s.readUserDetail(userID)
	if err != nil {
		s.logger.Error("read user after route assignment", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "route assigned but could not read user")
		return
	}
	writeSuccess(w, result)
}

func (s *Server) clearUserRoute(w http.ResponseWriter, r *http.Request) {
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
		s.logger.Error("read user for route clear", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not clear user route")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`UPDATE user_routes SET is_primary = 0, active_to = ?, route_exit_ip_id = NULL WHERE user_id = ? AND active_to IS NULL`, now, userID); err != nil {
		s.logger.Error("clear user route", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not clear user route")
		return
	}
	s.writeAuditLog(r, "user.route.clear", "user", userID, nil, map[string]any{"clearedAt": now})
	result, err := s.readUserDetail(userID)
	if err != nil {
		s.logger.Error("read user after route clear", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "route cleared but could not read user")
		return
	}
	writeSuccess(w, result)
}

func (s *Server) validateUserExists(id string) error {
	var exists int
	err := s.db.QueryRow(`SELECT 1 FROM users WHERE id = ?`, id).Scan(&exists)
	return err
}

// userDetail exposes the central business fields beside X-Panel's read-only
// Inbound and Client snapshot. It does not return X-Panel credentials or a
// complete Xray configuration.
func (s *Server) userDetail(w http.ResponseWriter, r *http.Request) {
	s.refreshOperationalStatuses(time.Now().UTC())
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}
	result, err := s.readUserDetail(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
		return
	}
	if err != nil {
		s.logger.Error("read user detail", "user_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user detail")
		return
	}
	writeSuccess(w, result)
}

// updateUser only changes central-owned business metadata. X-Panel-owned
// expiry, enable state, clients, and traffic remain sync-only fields.
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}
	var payload updateUserRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid user update payload")
		return
	}
	if payload.DisplayName == nil && payload.MonthlyFee == nil && payload.BillingCycle == nil && payload.BillingAmount == nil && payload.Currency == nil && payload.Notes == nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "at least one editable field is required")
		return
	}

	var displayName, currency, notes, billingCycle string
	var monthlyFee, billingAmount float64
	if err := s.db.QueryRow(`SELECT display_name, monthly_fee, COALESCE(billing_cycle, 'monthly'), COALESCE(billing_amount, monthly_fee), currency, COALESCE(notes, '') FROM users WHERE id = ?`, id).Scan(&displayName, &monthlyFee, &billingCycle, &billingAmount, &currency, &notes); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
		return
	} else if err != nil {
		s.logger.Error("read user for update", "user_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update user")
		return
	}
	before := map[string]any{"displayName": displayName, "monthlyFee": monthlyFee, "billingCycle": billingCycle, "billingAmount": billingAmount, "currency": currency, "notes": notes}
	if payload.DisplayName != nil {
		displayName = strings.TrimSpace(*payload.DisplayName)
		if displayName == "" || len([]rune(displayName)) > 120 {
			writeFailure(w, http.StatusBadRequest, validationCode, "displayName must contain 1 to 120 characters")
			return
		}
	}
	if payload.MonthlyFee != nil {
		if *payload.MonthlyFee < 0 || *payload.MonthlyFee > 100000000 || *payload.MonthlyFee != *payload.MonthlyFee {
			writeFailure(w, http.StatusBadRequest, validationCode, "monthlyFee must be a non-negative number")
			return
		}
		monthlyFee = *payload.MonthlyFee
		if payload.BillingAmount == nil && billingCycle == "monthly" {
			billingAmount = monthlyFee
		}
	}
	if payload.BillingCycle != nil {
		billingCycle = strings.TrimSpace(*payload.BillingCycle)
		if billingCycle != "monthly" && billingCycle != "annual" {
			writeFailure(w, http.StatusBadRequest, validationCode, "billingCycle must be monthly or annual")
			return
		}
		if payload.BillingAmount == nil {
			monthlyFee = billingAmount
			if billingCycle == "annual" {
				monthlyFee = billingAmount / 12
			}
		}
	}
	if payload.BillingAmount != nil {
		if *payload.BillingAmount < 0 || *payload.BillingAmount > 100000000 || *payload.BillingAmount != *payload.BillingAmount {
			writeFailure(w, http.StatusBadRequest, validationCode, "billingAmount must be a non-negative number")
			return
		}
		billingAmount = *payload.BillingAmount
		monthlyFee = billingAmount
		if billingCycle == "annual" {
			monthlyFee = billingAmount / 12
		}
	}
	if payload.Currency != nil {
		currency = strings.ToUpper(strings.TrimSpace(*payload.Currency))
		if currency != "CNY" {
			writeFailure(w, http.StatusBadRequest, validationCode, "currency currently supports CNY only")
			return
		}
	}
	if payload.Notes != nil {
		notes = strings.TrimSpace(*payload.Notes)
		if len([]rune(notes)) > 2000 {
			writeFailure(w, http.StatusBadRequest, validationCode, "notes must be at most 2000 characters")
			return
		}
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`UPDATE users SET display_name = ?, monthly_fee = ?, billing_cycle = ?, billing_amount = ?, currency = ?, notes = ?, updated_at = ? WHERE id = ?`, displayName, monthlyFee, billingCycle, billingAmount, currency, nullableDBString(notes), now, id); err != nil {
		s.logger.Error("update user business fields", "user_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update user")
		return
	}
	after := map[string]any{"displayName": displayName, "monthlyFee": monthlyFee, "billingCycle": billingCycle, "billingAmount": billingAmount, "currency": currency, "notes": notes}
	s.writeAuditLog(r, "user.update", "user", id, before, after)
	result, err := s.readUserDetail(id)
	if err != nil {
		s.logger.Error("read updated user", "user_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "user updated but could not read result")
		return
	}
	writeSuccess(w, result)
}

func (s *Server) readUserDetail(id string) (map[string]any, error) {
	var userID, displayName, status, currency, notes, expiry, inboundID, inboundRemoteID, inboundTag, inboundRemark, protocol, nodeID, nodeName, nodeType, lastSeen string
	var monthlyFee float64
	var port, clientCount int
	var enabled int
	var up, down, allTime int64
	err := s.db.QueryRow(`SELECT u.id, u.display_name, u.status, u.monthly_fee, u.currency, COALESCE(u.notes, ''), COALESCE(u.expiry_time, ''),
COALESCE(i.id, ''), COALESCE(i.remote_inbound_id, ''), COALESCE(i.tag, ''), COALESCE(i.remark, ''), COALESCE(i.protocol, ''), COALESCE(i.port, 0), COALESCE(i.enable, 0),
COALESCE(i.client_count, 0), COALESCE(i.up, 0), COALESCE(i.down, 0), COALESCE(i.all_time, 0), COALESCE(i.last_seen_at, ''),
COALESCE(n.id, ''), COALESCE(n.name, ''), COALESCE(n.type, '')
FROM users u
LEFT JOIN user_inbounds ui ON ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL
LEFT JOIN inbounds i ON i.id = ui.inbound_id AND i.kind = 'user' AND i.deleted_at IS NULL
 AND i.node_id IN (SELECT id FROM nodes WHERE type = 'relay' AND deleted_at IS NULL)
LEFT JOIN nodes n ON n.id = i.node_id
WHERE u.id = ?`, id).Scan(&userID, &displayName, &status, &monthlyFee, &currency, &notes, &expiry,
		&inboundID, &inboundRemoteID, &inboundTag, &inboundRemark, &protocol, &port, &enabled, &clientCount, &up, &down, &allTime, &lastSeen,
		&nodeID, &nodeName, &nodeType)
	if err != nil {
		return nil, err
	}
	result := map[string]any{
		"id": userID, "displayName": displayName, "status": status, "monthlyFee": monthlyFee, "currency": currency,
		"notes": nullableString(notes), "expiresAt": nullableString(expiry),
		"inbound": map[string]any{
			"id": nullableString(inboundID), "remoteId": nullableString(inboundRemoteID), "tag": nullableString(inboundTag),
			"remark": nullableString(inboundRemark), "protocol": nullableString(protocol), "port": port, "enabled": enabled == 1,
			"clientCount": clientCount, "up": up, "down": down, "allTime": allTime, "lastSeenAt": nullableString(lastSeen),
		},
		"node": map[string]any{"id": nullableString(nodeID), "name": nullableString(nodeName), "type": nullableString(nodeType)},
	}
	var billingCycle string
	var billingAmount float64
	if err := s.db.QueryRow(`SELECT COALESCE(billing_cycle, 'monthly'), COALESCE(billing_amount, monthly_fee) FROM users WHERE id = ?`, id).Scan(&billingCycle, &billingAmount); err != nil {
		return nil, err
	}
	result["billingCycle"] = billingCycle
	result["billingAmount"] = billingAmount
	candidates, err := s.readUserRenewalCandidates(id)
	if err != nil {
		return nil, err
	}
	result["renewalCandidates"] = candidates
	records, err := s.readUserBillingRecords(id)
	if err != nil {
		return nil, err
	}
	result["billingRecords"] = records
	clients := make([]map[string]any, 0)
	if inboundID != "" {
		rows, err := s.db.Query(`SELECT remote_client_id, COALESCE(email, ''), enable, COALESCE(expiry_time, ''), up, down, all_time, COALESCE(last_online, ''), COALESCE(last_seen_at, '')
FROM clients WHERE inbound_id = ? ORDER BY email ASC, remote_client_id ASC`, inboundID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var remoteID, email, clientExpiry, lastOnline, clientLastSeen string
			var clientEnabled int
			var clientUp, clientDown, clientAllTime int64
			if err := rows.Scan(&remoteID, &email, &clientEnabled, &clientExpiry, &clientUp, &clientDown, &clientAllTime, &lastOnline, &clientLastSeen); err != nil {
				_ = rows.Close()
				return nil, err
			}
			clients = append(clients, map[string]any{
				"remoteId": remoteID, "email": nullableString(email), "enabled": clientEnabled == 1,
				"expiresAt": nullableString(clientExpiry), "up": clientUp, "down": clientDown, "allTime": clientAllTime,
				"lastOnlineAt": nullableString(lastOnline), "lastSeenAt": nullableString(clientLastSeen),
			})
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return nil, err
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
	}
	result["clients"] = clients

	routes := make([]map[string]any, 0)
	rows, err := s.db.Query(`SELECT r.id, r.name, COALESCE(relay.name, ''), COALESCE(landing.name, ''), COALESCE(r.landing_inbound_id, ''), COALESCE(r.landing_inbound_tag, ''), r.enabled,
ur.is_primary, COALESCE(ur.active_from, ''), COALESCE(ur.active_to, ''),
COALESCE(rei.id, ''), COALESCE(e.id, ''), COALESCE(e.ip, ''), COALESCE(rei.scope, ''), COALESCE(e.source_type, ''), COALESCE(owner.name, '')
FROM user_routes ur
JOIN routes r ON r.id = ur.route_id
LEFT JOIN nodes relay ON relay.id = r.relay_node_id
LEFT JOIN nodes landing ON landing.id = r.landing_node_id
LEFT JOIN route_exit_ips rei ON rei.id = ur.route_exit_ip_id
LEFT JOIN exit_ips e ON e.id = rei.exit_ip_id
LEFT JOIN nodes owner ON owner.id = e.owner_node_id
WHERE ur.user_id = ? AND ur.active_to IS NULL
ORDER BY ur.is_primary DESC, r.name ASC`, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var routeID, routeName, relayName, landingName, landingInboundID, landingInboundTag string
		var activeFrom, activeTo, routeExitIPBindingID, exitIPID, exitIPAddress, exitIPScope, exitIPSourceType, exitIPOwnerName string
		var routeEnabled, isPrimary int
		if err := rows.Scan(&routeID, &routeName, &relayName, &landingName, &landingInboundID, &landingInboundTag, &routeEnabled, &isPrimary,
			&activeFrom, &activeTo, &routeExitIPBindingID, &exitIPID, &exitIPAddress, &exitIPScope, &exitIPSourceType, &exitIPOwnerName); err != nil {
			_ = rows.Close()
			return nil, err
		}
		routes = append(routes, map[string]any{
			"id": routeID, "name": routeName, "relayNodeName": nullableString(relayName), "landingNodeName": nullableString(landingName),
			"landingInboundId": nullableString(landingInboundID), "landingInboundTag": nullableString(landingInboundTag), "enabled": routeEnabled == 1, "isPrimary": isPrimary == 1,
			"activeFrom": nullableString(activeFrom), "activeTo": nullableString(activeTo),
			"routeExitIpId": nullableString(routeExitIPBindingID), "exitIpId": nullableString(exitIPID),
			"exitIpAddress": nullableString(exitIPAddress), "exitIpScope": nullableString(exitIPScope),
			"exitIpSourceType": nullableString(exitIPSourceType), "exitIpOwnerNodeName": nullableString(exitIPOwnerName),
			"assignmentMode": map[bool]string{true: "fixed", false: "pool"}[exitIPID != ""],
		})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	result["routes"] = routes
	path, err := s.readUserPath(userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	result["path"] = userPathData(path)
	pathHistory, err := s.userPathHistory(userID)
	if err != nil {
		return nil, err
	}
	result["pathHistory"] = pathHistory
	traffic, err := s.readUserTraffic(inboundID)
	if err != nil {
		return nil, err
	}
	result["traffic"] = traffic
	return result, nil
}

func (s *Server) readUserTraffic(inboundID string) ([]map[string]any, error) {
	items := make([]map[string]any, 0)
	if inboundID == "" {
		return items, nil
	}
	rows, err := s.db.Query(`SELECT collected_at, up, down, all_time, reset_detected
FROM traffic_snapshots WHERE inbound_id = ? ORDER BY collected_at DESC LIMIT 30`, inboundID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var collectedAt string
		var up, down, allTime int64
		var resetDetected int
		if err := rows.Scan(&collectedAt, &up, &down, &allTime, &resetDetected); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"collectedAt": collectedAt, "up": up, "down": down, "allTime": allTime,
			"resetDetected": resetDetected == 1,
		})
	}
	return items, rows.Err()
}

func (s *Server) nodes(w http.ResponseWriter, r *http.Request) {
	s.refreshOperationalStatuses(time.Now().UTC())
	query := parseListQuery(r)
	where := []string{"n.deleted_at IS NULL"}
	args := make([]any, 0, 4)
	if query.keyword != "" {
		where = append(where, `(n.name LIKE ? OR n.hostname LIKE ? OR n.management_url LIKE ? OR n.public_ip LIKE ?)`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like, like)
	}
	if query.status != "" && query.status != "all" {
		if query.status == "disabled" {
			where = append(where, "n.enabled = 0")
		} else {
			where = append(where, "n.enabled = 1 AND n.health_status = ?")
			args = append(args, query.status)
		}
	}
	if query.nodeType != "" {
		where = append(where, "n.type = ?")
		args = append(args, query.nodeType)
	}
	base := `FROM nodes n WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count nodes")
		return
	}
	rows, err := s.db.Query(`SELECT n.id, n.name, n.type, n.health_status, n.enabled,
COALESCE(n.management_url, ''), COALESCE(NULLIF(n.management_url, ''), NULLIF(n.public_ip, ''), n.hostname, ''), COALESCE(n.xpanel_version, ''), COALESCE(n.xray_version, ''),
COALESCE(n.cpu_usage, 0), COALESCE(n.memory_used, 0), COALESCE(n.memory_total, 0), COALESCE(n.disk_used, 0), COALESCE(n.disk_total, 0),
COALESCE(n.last_seen_at, ''), COALESCE((SELECT COALESCE(sr.finished_at, sr.started_at) FROM sync_runs sr WHERE sr.node_id = n.id AND sr.status = 'success'
ORDER BY COALESCE(sr.finished_at, sr.started_at) DESC LIMIT 1), ''), COALESCE((SELECT SUM(c.monthly_amount) FROM node_costs c WHERE c.node_id = n.id AND c.currency = 'CNY'
AND c.effective_from <= date('now') AND (c.effective_to IS NULL OR c.effective_to >= date('now'))), 0),
COALESCE((SELECT COUNT(*) FROM exit_ips e WHERE COALESCE(e.owner_node_id, e.landing_node_id) = n.id), 0)
`+base+` ORDER BY n.name ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read nodes")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, name, nodeType, status, managementURL, host, xpanel, xray, lastSeen, lastSync string
		var enabled int
		var cpuUsage, memoryUsed, memoryTotal, diskUsed, diskTotal float64
		var exitIPCount int
		var monthlyCost float64
		if err := rows.Scan(&id, &name, &nodeType, &status, &enabled, &managementURL, &host, &xpanel, &xray, &cpuUsage, &memoryUsed, &memoryTotal, &diskUsed, &diskTotal, &lastSeen, &lastSync, &monthlyCost, &exitIPCount); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode nodes")
			return
		}
		displayStatus := status
		if enabled != 1 {
			displayStatus = "disabled"
		}
		items = append(items, map[string]any{
			"id": id, "name": name, "type": nodeType, "status": displayStatus, "enabled": enabled == 1, "host": host, "managementUrl": nullableString(managementURL),
			"xpanelVersion": nullableString(xpanel), "xrayVersion": nullableString(xray),
			"cpuUsage": cpuUsage, "memoryUsed": memoryUsed, "memoryTotal": memoryTotal, "diskUsed": diskUsed, "diskTotal": diskTotal,
			"lastSeenAt": nullableString(lastSeen), "lastSyncAt": nullableString(lastSync), "monthlyCost": monthlyCost, "currency": "CNY", "exitIpCount": exitIPCount,
		})
	}
	writeSuccess(w, s.pageResponse(items, total, query))
}

// nodeDetail returns the latest central snapshot for one node. The nested
// collections are deliberately read-only mirrors of Agent/X-Panel data plus
// central operational history; no credentials or raw Xray configuration is
// exposed here.
func (s *Server) nodeDetail(w http.ResponseWriter, r *http.Request) {
	s.refreshOperationalStatuses(time.Now().UTC())
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node id is required")
		return
	}

	var nodeID, nodeKey, name, nodeType, status, hostname, managementURL, publicIP, region, provider, panelBasePath, agentVersion, xpanelVersion, xrayVersion, lastSeen, lastSync, dataAt string
	var cpuUsage sql.NullFloat64
	var memoryUsed, memoryTotal, diskUsed, diskTotal sql.NullInt64
	var enabled int
	var userCount int
	var trafficBytes int64
	err := s.db.QueryRow(`SELECT n.id, n.node_key, n.name, n.type, n.health_status,
	COALESCE(n.hostname, ''), COALESCE(n.management_url, ''), COALESCE(n.public_ip, ''), COALESCE(n.region, ''), COALESCE(n.provider, ''),
COALESCE(n.panel_base_path, ''), COALESCE(n.agent_version, ''), COALESCE(n.xpanel_version, ''), COALESCE(n.xray_version, ''),
 n.cpu_usage, n.memory_used, n.memory_total, n.disk_used, n.disk_total,
COALESCE(n.last_seen_at, ''),
COALESCE((SELECT COALESCE(sr.finished_at, sr.started_at) FROM sync_runs sr WHERE sr.node_id = n.id AND sr.status = 'success'
ORDER BY COALESCE(sr.finished_at, sr.started_at) DESC LIMIT 1), ''),
		COALESCE((SELECT MAX(COALESCE(sr.finished_at, sr.started_at)) FROM sync_runs sr WHERE sr.node_id = n.id AND sr.status = 'success'), ''),
(SELECT COUNT(DISTINCT u.id) FROM users u
 JOIN inbounds ui ON ui.user_id = u.id AND ui.node_id = n.id AND ui.kind = 'user' AND ui.deleted_at IS NULL
 WHERE n.type = 'relay' AND n.deleted_at IS NULL),
COALESCE((SELECT SUM(COALESCE(i.up, 0) + COALESCE(i.down, 0)) FROM inbounds i WHERE i.node_id = n.id AND i.deleted_at IS NULL), 0),
n.enabled
FROM nodes n WHERE n.id = ? AND n.deleted_at IS NULL`, id).Scan(&nodeID, &nodeKey, &name, &nodeType, &status,
		&hostname, &managementURL, &publicIP, &region, &provider, &panelBasePath, &agentVersion, &xpanelVersion, &xrayVersion,
		&cpuUsage, &memoryUsed, &memoryTotal, &diskUsed, &diskTotal,
		&lastSeen, &lastSync, &dataAt, &userCount, &trafficBytes, &enabled)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node not found")
		return
	}
	if err != nil {
		s.logger.Error("read node detail", "node_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node detail")
		return
	}
	if enabled != 1 {
		status = "disabled"
	}

	result := map[string]any{
		"id": nodeID, "nodeKey": nodeKey, "name": name, "type": nodeType, "status": status,
		"host": firstNonEmpty(managementURL, publicIP, hostname), "managementUrl": nullableString(managementURL), "hostname": nullableString(hostname), "publicIp": nullableString(publicIP),
		"region": nullableString(region), "provider": nullableString(provider), "panelBasePath": nullableString(panelBasePath),
		"agentVersion": nullableString(agentVersion), "xpanelVersion": nullableString(xpanelVersion), "xrayVersion": nullableString(xrayVersion),
		"cpuUsage": nullableFloat(cpuUsage), "memoryUsed": nullableInt(memoryUsed), "memoryTotal": nullableInt(memoryTotal),
		"diskUsed": nullableInt(diskUsed), "diskTotal": nullableInt(diskTotal),
		"enabled": enabled == 1, "lastSeenAt": nullableString(lastSeen), "lastSyncAt": nullableString(lastSync),
		"userCount": userCount, "trafficBytes": trafficBytes, "dataAt": nullableString(dataAt),
	}

	inbounds := make([]map[string]any, 0)
	rows, err := s.db.Query(`SELECT i.id, i.remote_inbound_id, COALESCE(i.tag, ''), COALESCE(i.remark, ''),
COALESCE(i.kind, ''), COALESCE(i.protocol, ''), COALESCE(i.port, 0), COALESCE(i.listen, ''), i.enable,
COALESCE(i.expiry_time, ''), COALESCE(i.client_count, 0), COALESCE(i.up, 0), COALESCE(i.down, 0), COALESCE(i.all_time, 0),
COALESCE(i.last_seen_at, ''), COALESCE(i.deleted_at, '')
FROM inbounds i WHERE i.node_id = ? ORDER BY (i.deleted_at IS NOT NULL) ASC, i.tag ASC, i.remote_inbound_id ASC`, id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node inbounds")
		return
	}
	for rows.Next() {
		var inboundID, remoteID, tag, remark, kind, protocol, listen, expiry, lastInboundSeen, deletedAt string
		var port, enabledInbound, clientCount int
		var up, down, allTime int64
		if err := rows.Scan(&inboundID, &remoteID, &tag, &remark, &kind, &protocol, &port, &listen, &enabledInbound,
			&expiry, &clientCount, &up, &down, &allTime, &lastInboundSeen, &deletedAt); err != nil {
			_ = rows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode node inbounds")
			return
		}
		inboundStatus := "active"
		if deletedAt != "" {
			inboundStatus = "archived"
		} else if enabledInbound != 1 {
			inboundStatus = "disabled"
		}
		inbounds = append(inbounds, map[string]any{
			"id": inboundID, "remoteId": remoteID, "tag": nullableString(tag), "remark": nullableString(remark),
			"kind": kind, "purpose": inboundPurpose(nodeType), "protocol": nullableString(protocol), "port": port, "listen": nullableString(listen),
			"enabled": enabledInbound == 1, "status": inboundStatus, "expiresAt": nullableString(expiry),
			"clientCount": clientCount, "up": up, "down": down, "allTime": allTime,
			"lastSeenAt": nullableString(lastInboundSeen), "deletedAt": nullableString(deletedAt),
		})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node inbounds")
		return
	}
	if err := rows.Close(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close node inbounds")
		return
	}
	result["inbounds"] = inbounds

	exitIPs := make([]map[string]any, 0)
	rows, err = s.db.Query(`SELECT e.id, e.ip, COALESCE(e.source_type, 'node'), COALESCE(e.owner_node_id, e.landing_node_id, ''),
COALESCE(e.provider, ''), e.family, e.enabled, e.monthly_cost, e.currency,
COALESCE(e.valid_from, ''), COALESCE(e.valid_to, ''), e.updated_at
FROM exit_ips e WHERE COALESCE(e.owner_node_id, e.landing_node_id) = ? ORDER BY e.enabled DESC, e.ip ASC`, id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node exit IPs")
		return
	}
	for rows.Next() {
		var exitID, address, sourceType, ownerNodeID, exitProvider, currency, validFrom, validTo, checkedAt string
		var family, enabledExit int
		var monthlyCost float64
		if err := rows.Scan(&exitID, &address, &sourceType, &ownerNodeID, &exitProvider, &family, &enabledExit, &monthlyCost, &currency, &validFrom, &validTo, &checkedAt); err != nil {
			_ = rows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode node exit IPs")
			return
		}
		exitStatus := "disabled"
		if enabledExit == 1 {
			exitStatus = "active"
		}
		exitIPs = append(exitIPs, map[string]any{
			"id": exitID, "address": address, "sourceType": sourceType, "ownerNodeId": nullableString(ownerNodeID),
			"ownerNodeName": nullableString(name), "ownerNodeType": nullableString(nodeType), "provider": nullableString(exitProvider), "family": family,
			"status": exitStatus, "monthlyCost": monthlyCost, "currency": currency,
			"validFrom": nullableString(validFrom), "validTo": nullableString(validTo), "checkedAt": nullableString(checkedAt),
		})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node exit IPs")
		return
	}
	if err := rows.Close(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close node exit IPs")
		return
	}
	result["exitIps"] = exitIPs
	result["exitIpCount"] = len(exitIPs)

	costs, err := s.readNodeCosts(id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node costs")
		return
	}
	costItems := make([]map[string]any, 0, len(costs))
	for _, cost := range costs {
		costItems = append(costItems, nodeCostRecordData(cost))
	}
	result["costs"] = costItems

	syncRuns := make([]map[string]any, 0)
	rows, err = s.db.Query(`SELECT id, sync_id, started_at, COALESCE(finished_at, ''), status, inbound_count, client_count, COALESCE(error_message, '')
FROM sync_runs WHERE node_id = ? ORDER BY started_at DESC LIMIT 5`, id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node sync runs")
		return
	}
	for rows.Next() {
		var runID, syncID, startedAt, finishedAt, runStatus, errorMessage string
		var inboundCount, clientCount int
		if err := rows.Scan(&runID, &syncID, &startedAt, &finishedAt, &runStatus, &inboundCount, &clientCount, &errorMessage); err != nil {
			_ = rows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode node sync runs")
			return
		}
		syncRuns = append(syncRuns, map[string]any{
			"id": runID, "syncId": syncID, "startedAt": startedAt, "finishedAt": nullableString(finishedAt),
			"status": runStatus, "inboundCount": inboundCount, "clientCount": clientCount, "errorMessage": nullableString(errorMessage),
		})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node sync runs")
		return
	}
	if err := rows.Close(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close node sync runs")
		return
	}
	result["syncRuns"] = syncRuns

	statusHistory := make([]map[string]any, 0)
	rows, err = s.db.Query(`SELECT e.id, e.event_type, e.severity, e.message, e.created_at, e.acknowledged
FROM node_events e WHERE e.node_id = ? ORDER BY e.created_at DESC LIMIT 30`, id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node status history")
		return
	}
	for rows.Next() {
		var eventID, eventType, severity, message, occurredAt string
		var acknowledged int
		if err := rows.Scan(&eventID, &eventType, &severity, &message, &occurredAt, &acknowledged); err != nil {
			_ = rows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode node status history")
			return
		}
		statusHistory = append(statusHistory, map[string]any{
			"id": eventID, "type": eventType, "severity": severity, "message": message,
			"occurredAt": occurredAt, "acknowledged": acknowledged == 1,
		})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node status history")
		return
	}
	_ = rows.Close()
	result["statusHistory"] = statusHistory

	writeSuccess(w, result)
}

func (s *Server) requestNodeSync(w http.ResponseWriter, r *http.Request) {
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
		s.logger.Error("read node for sync request", "node_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not request node sync")
		return
	}
	if enabled != 1 {
		writeFailure(w, http.StatusConflict, validationCode, "node is disabled")
		return
	}

	requestID := newID()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	message := "管理员请求立即同步；Agent 将在下一次轮询中执行（请求 " + requestID + ")"
	if _, err := s.db.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES (?, ?, 'sync_requested', 'info', ?, ?)`, requestID, id, message, now); err != nil {
		s.logger.Error("record node sync request", "node_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not request node sync")
		return
	}
	s.writeAuditLog(r, "node.sync.request", "node", id, nil, map[string]any{"requestId": requestID, "status": "queued"})
	writeSuccess(w, map[string]any{"requestId": requestID, "nodeId": id, "nodeName": name, "status": "queued", "requestedAt": now})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

type routeUpsertRequest struct {
	Name              *string `json:"name"`
	RelayNodeID       *string `json:"relayNodeId"`
	LandingNodeID     *string `json:"landingNodeId"`
	RelayOutboundTag  *string `json:"relayOutboundTag"`
	LandingInboundID  *string `json:"landingInboundId"`
	LandingInboundTag *string `json:"landingInboundTag"`
	Enabled           *bool   `json:"enabled"`
	ValidFrom         *string `json:"validFrom"`
	ValidTo           *string `json:"validTo"`
	Notes             *string `json:"notes"`
}

type routeRecord struct {
	ID                 string
	Name               string
	RelayNodeID        string
	RelayNodeName      string
	LandingNodeID      string
	LandingNodeName    string
	RelayOutboundTag   string
	LandingInboundID   string
	LandingInboundTag  string
	Enabled            bool
	ValidFrom          string
	ValidTo            string
	Notes              string
	ExitIPCount        int
	AllocatedUserCount int
}

type routeValues struct {
	Name              string
	RelayNodeID       string
	LandingNodeID     string
	RelayOutboundTag  string
	LandingInboundID  string
	LandingInboundTag string
	Enabled           bool
	ValidFrom         string
	ValidTo           string
	Notes             string
}

type sqlQueryer interface {
	QueryRow(query string, args ...any) *sql.Row
}

func (s *Server) readRoute(id string) (routeRecord, error) {
	var route routeRecord
	var enabled int
	err := s.db.QueryRow(`SELECT r.id, r.name, r.relay_node_id, COALESCE(relay.name, ''), r.landing_node_id,
COALESCE(landing.name, ''), COALESCE(r.relay_outbound_tag, ''), COALESCE(r.landing_inbound_id, ''),
COALESCE(r.landing_inbound_tag, ''), r.enabled, COALESCE(r.valid_from, ''), COALESCE(r.valid_to, ''), COALESCE(r.notes, '')
FROM routes r
LEFT JOIN nodes relay ON relay.id = r.relay_node_id
LEFT JOIN nodes landing ON landing.id = r.landing_node_id
WHERE r.id = ?`, id).Scan(&route.ID, &route.Name, &route.RelayNodeID, &route.RelayNodeName,
		&route.LandingNodeID, &route.LandingNodeName, &route.RelayOutboundTag, &route.LandingInboundID,
		&route.LandingInboundTag, &enabled, &route.ValidFrom, &route.ValidTo, &route.Notes)
	if err != nil {
		return routeRecord{}, err
	}
	route.Enabled = enabled == 1
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM route_exit_ips rei JOIN exit_ips e ON e.id = rei.exit_ip_id
WHERE rei.route_id = ? AND rei.enabled = 1 AND e.enabled = 1`, id).Scan(&route.ExitIPCount); err != nil {
		return routeRecord{}, err
	}
	if err := s.db.QueryRow(`SELECT COUNT(DISTINCT ur.user_id) FROM user_routes ur
JOIN users u ON u.id = ur.user_id
WHERE ur.route_id = ? AND ur.is_primary = 1 AND ur.active_to IS NULL
AND u.status <> 'disabled' AND datetime(u.created_at) <= datetime('now')
AND (u.expiry_time IS NULL OR datetime(u.expiry_time) >= datetime('now'))`, id).Scan(&route.AllocatedUserCount); err != nil {
		return routeRecord{}, err
	}
	return route, nil
}

func routeRecordData(route routeRecord) map[string]any {
	status := "disabled"
	if route.Enabled {
		status = "active"
	}
	return map[string]any{
		"id": route.ID, "name": route.Name,
		"relayNodeId": route.RelayNodeID, "relayNodeName": nullableString(route.RelayNodeName),
		"landingNodeId": route.LandingNodeID, "landingNodeName": nullableString(route.LandingNodeName),
		"relayOutboundTag": nullableString(route.RelayOutboundTag),
		"landingInboundId": nullableString(route.LandingInboundID), "landingInboundTag": nullableString(route.LandingInboundTag),
		"enabled": route.Enabled, "validFrom": nullableString(route.ValidFrom), "validTo": nullableString(route.ValidTo),
		"notes": nullableString(route.Notes), "exitIpCount": route.ExitIPCount,
		"allocatedUserCount": route.AllocatedUserCount, "status": status,
	}
}

func routeValuesFromRequest(request routeUpsertRequest, existing *routeRecord) (routeValues, error) {
	values := routeValues{Enabled: true}
	if existing != nil {
		values = routeValues{
			Name: existing.Name, RelayNodeID: existing.RelayNodeID, LandingNodeID: existing.LandingNodeID,
			RelayOutboundTag: existing.RelayOutboundTag, LandingInboundID: existing.LandingInboundID,
			LandingInboundTag: existing.LandingInboundTag, Enabled: existing.Enabled,
			ValidFrom: existing.ValidFrom, ValidTo: existing.ValidTo, Notes: existing.Notes,
		}
	}
	if request.Name != nil {
		values.Name = strings.TrimSpace(*request.Name)
	}
	if request.RelayNodeID != nil {
		values.RelayNodeID = strings.TrimSpace(*request.RelayNodeID)
	}
	if request.LandingNodeID != nil {
		values.LandingNodeID = strings.TrimSpace(*request.LandingNodeID)
	}
	if request.RelayOutboundTag != nil {
		values.RelayOutboundTag = strings.TrimSpace(*request.RelayOutboundTag)
	}
	if request.LandingInboundID != nil {
		values.LandingInboundID = strings.TrimSpace(*request.LandingInboundID)
	}
	if request.LandingInboundTag != nil {
		values.LandingInboundTag = strings.TrimSpace(*request.LandingInboundTag)
	}
	if request.Enabled != nil {
		values.Enabled = *request.Enabled
	}
	if request.ValidFrom != nil {
		values.ValidFrom = strings.TrimSpace(*request.ValidFrom)
	}
	if request.ValidTo != nil {
		values.ValidTo = strings.TrimSpace(*request.ValidTo)
	}
	if request.Notes != nil {
		values.Notes = strings.TrimSpace(*request.Notes)
	}
	return values, validateRouteValues(values)
}

func validateRouteValues(values routeValues) error {
	if values.Name == "" {
		return errors.New("route name is required")
	}
	if len(values.Name) > 120 {
		return errors.New("route name is too long")
	}
	if values.RelayNodeID == "" || values.LandingNodeID == "" {
		return errors.New("relayNodeId and landingNodeId are required")
	}
	if values.RelayNodeID == values.LandingNodeID {
		return errors.New("relay and landing nodes must be different")
	}
	if len(values.RelayOutboundTag) > 200 || len(values.LandingInboundTag) > 200 {
		return errors.New("route tag is too long")
	}
	if len(values.LandingInboundID) > 120 {
		return errors.New("landingInboundId is too long")
	}
	if len(values.Notes) > 2000 {
		return errors.New("route notes are too long")
	}
	from, fromSet, err := parseRouteDate(values.ValidFrom)
	if err != nil {
		return errors.New("validFrom must be RFC3339 or YYYY-MM-DD")
	}
	to, toSet, err := parseRouteDate(values.ValidTo)
	if err != nil {
		return errors.New("validTo must be RFC3339 or YYYY-MM-DD")
	}
	if fromSet && toSet && to.Before(from) {
		return errors.New("validTo must not be earlier than validFrom")
	}
	return nil
}

func parseRouteDate(value string) (time.Time, bool, error) {
	if value == "" {
		return time.Time{}, false, nil
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, true, nil
	}
	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return parsed, true, nil
	}
	return time.Time{}, false, errors.New("invalid route date")
}

func validateRouteRelations(queryer sqlQueryer, values routeValues) error {
	var relayType, landingType string
	if err := queryer.QueryRow(`SELECT type FROM nodes WHERE id = ? AND deleted_at IS NULL`, values.RelayNodeID).Scan(&relayType); errors.Is(err, sql.ErrNoRows) {
		return errors.New("relay node not found")
	} else if err != nil {
		return fmt.Errorf("read relay node: %w", err)
	}
	if relayType != "relay" {
		return errors.New("relayNodeId must refer to a relay node")
	}
	if err := queryer.QueryRow(`SELECT type FROM nodes WHERE id = ? AND deleted_at IS NULL`, values.LandingNodeID).Scan(&landingType); errors.Is(err, sql.ErrNoRows) {
		return errors.New("landing node not found")
	} else if err != nil {
		return fmt.Errorf("read landing node: %w", err)
	}
	if landingType != "landing" {
		return errors.New("landingNodeId must refer to a landing node")
	}
	if values.LandingInboundID != "" {
		var count int
		if err := queryer.QueryRow(`SELECT COUNT(*) FROM inbounds WHERE id = ? AND node_id = ?`, values.LandingInboundID, values.LandingNodeID).Scan(&count); err != nil {
			return fmt.Errorf("read landing inbound: %w", err)
		}
		if count == 0 {
			return errors.New("landingInboundId does not belong to landingNodeId")
		}
	}
	return nil
}

func (s *Server) routeDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "route id is required")
		return
	}
	route, err := s.readRoute(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		return
	}
	if err != nil {
		s.logger.Error("read route detail", "route_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route detail")
		return
	}
	writeSuccess(w, routeRecordData(route))
}

func (s *Server) createRoute(w http.ResponseWriter, r *http.Request) {
	var request routeUpsertRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid route payload")
		return
	}
	values, err := routeValuesFromRequest(request, nil)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if err := validateRouteRelations(s.db, values); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	id := newID()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`INSERT INTO routes (id, name, relay_node_id, landing_node_id, relay_outbound_tag, landing_inbound_id,
landing_inbound_tag, enabled, valid_from, valid_to, notes, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, values.Name, values.RelayNodeID, values.LandingNodeID,
		values.RelayOutboundTag, nullableRouteValue(values.LandingInboundID), nullableRouteValue(values.LandingInboundTag), boolInt(values.Enabled),
		nullableRouteValue(values.ValidFrom), nullableRouteValue(values.ValidTo), nullableRouteValue(values.Notes), now, now); err != nil {
		s.logger.Error("create route", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create route")
		return
	}
	route, err := s.readRoute(id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read created route")
		return
	}
	s.writeAuditLog(r, "route.create", "route", id, nil, routeRecordData(route))
	writeSuccess(w, routeRecordData(route))
}

func (s *Server) updateRoute(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "route id is required")
		return
	}
	existing, err := s.readRoute(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route")
		return
	}
	var request routeUpsertRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid route payload")
		return
	}
	values, err := routeValuesFromRequest(request, &existing)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if err := validateRouteRelations(s.db, values); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`UPDATE routes SET name = ?, relay_node_id = ?, landing_node_id = ?, relay_outbound_tag = ?, landing_inbound_id = ?,
landing_inbound_tag = ?, enabled = ?, valid_from = ?, valid_to = ?, notes = ?, updated_at = ? WHERE id = ?`, values.Name,
		values.RelayNodeID, values.LandingNodeID, nullableRouteValue(values.RelayOutboundTag), nullableRouteValue(values.LandingInboundID),
		nullableRouteValue(values.LandingInboundTag), boolInt(values.Enabled), nullableRouteValue(values.ValidFrom), nullableRouteValue(values.ValidTo),
		nullableRouteValue(values.Notes), now, id); err != nil {
		s.logger.Error("update route", "route_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update route")
		return
	}
	updated, err := s.readRoute(id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read updated route")
		return
	}
	s.writeAuditLog(r, "route.update", "route", id, routeRecordData(existing), routeRecordData(updated))
	writeSuccess(w, routeRecordData(updated))
}

func (s *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "route id is required")
		return
	}
	route, err := s.readRoute(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route")
		return
	}
	var userBindings, exitBindings int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM user_routes WHERE route_id = ?`, id).Scan(&userBindings); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check route bindings")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM route_exit_ips WHERE route_id = ?`, id).Scan(&exitBindings); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check route bindings")
		return
	}
	if userBindings > 0 || exitBindings > 0 {
		writeFailure(w, http.StatusConflict, validationCode, "route has bindings; disable it instead of deleting")
		return
	}
	if _, err := s.db.Exec(`DELETE FROM routes WHERE id = ?`, id); err != nil {
		s.logger.Error("delete route", "route_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete route")
		return
	}
	s.writeAuditLog(r, "route.delete", "route", id, routeRecordData(route), map[string]any{"deleted": true})
	writeSuccess(w, map[string]any{"id": id, "deleted": true})
}

func nullableRouteValue(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func (s *Server) routes(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)
	where := []string{"1 = 1"}
	args := make([]any, 0, 3)
	if query.keyword != "" {
		where = append(where, `(r.name LIKE ? OR relay.name LIKE ? OR landing.name LIKE ? OR r.landing_inbound_tag LIKE ?)`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like, like)
	}
	if query.status != "" {
		if query.status == "active" {
			where = append(where, "r.enabled = 1")
		} else if query.status == "disabled" {
			where = append(where, "r.enabled = 0")
		}
	}
	base := `FROM routes r
LEFT JOIN nodes relay ON relay.id = r.relay_node_id
LEFT JOIN nodes landing ON landing.id = r.landing_node_id
WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count routes")
		return
	}
	rows, err := s.db.Query(`SELECT r.id, r.name, r.relay_node_id, r.landing_node_id, COALESCE(relay.name, ''), COALESCE(landing.name, ''),
COALESCE(r.landing_inbound_tag, ''), (SELECT COUNT(*) FROM route_exit_ips rei JOIN exit_ips e ON e.id = rei.exit_ip_id
 WHERE rei.route_id = r.id AND rei.enabled = 1 AND e.enabled = 1),
(SELECT COUNT(DISTINCT ur.user_id) FROM user_routes ur JOIN users u ON u.id = ur.user_id
 WHERE ur.route_id = r.id AND ur.is_primary = 1 AND ur.active_to IS NULL
 AND u.status <> 'disabled' AND datetime(u.created_at) <= datetime('now')
 AND (u.expiry_time IS NULL OR datetime(u.expiry_time) >= datetime('now'))),
CASE WHEN r.enabled = 1 THEN 'active' ELSE 'disabled' END
`+base+` ORDER BY r.name ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read routes")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, name, relayNodeID, landingNodeID, relay, landing, inboundTag, status string
		var exitIPCount, userCount int
		if err := rows.Scan(&id, &name, &relayNodeID, &landingNodeID, &relay, &landing, &inboundTag, &exitIPCount, &userCount, &status); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode routes")
			return
		}
		items = append(items, map[string]any{
			"id": id, "name": name, "relayNodeId": relayNodeID, "landingNodeId": landingNodeID,
			"relayNodeName": relay, "landingNodeName": landing,
			"landingInboundTag": nullableString(inboundTag), "exitIpCount": exitIPCount,
			"allocatedUserCount": userCount, "status": status,
		})
	}
	writeSuccess(w, s.pageResponse(items, total, query))
}

type routeExitIPUpsertRequest struct {
	ExitIPID         *string  `json:"exitIpId"`
	Scope            *string  `json:"scope"`
	AllocationWeight *float64 `json:"allocationWeight"`
	Enabled          *bool    `json:"enabled"`
}

type routeExitIPRecord struct {
	ID               string
	RouteID          string
	ExitIPID         string
	Address          string
	Family           int
	SourceType       string
	OwnerNodeID      string
	OwnerNodeName    string
	OwnerNodeType    string
	LandingNodeID    string
	LandingNodeName  string
	Scope            string
	AllocationWeight float64
	Enabled          bool
}

func (s *Server) readRouteExitIPs(routeID string) ([]routeExitIPRecord, error) {
	rows, err := s.db.Query(`SELECT rei.id, rei.route_id, rei.exit_ip_id, e.ip, e.family,
COALESCE(e.source_type, 'node'), COALESCE(e.owner_node_id, e.landing_node_id, ''),
COALESCE(owner.name, ''), COALESCE(owner.type, ''), COALESCE(e.landing_node_id, ''),
COALESCE(landing.name, ''), COALESCE(rei.scope, 'landing'), rei.allocation_weight, rei.enabled
FROM route_exit_ips rei
JOIN exit_ips e ON e.id = rei.exit_ip_id
LEFT JOIN nodes owner ON owner.id = COALESCE(e.owner_node_id, e.landing_node_id)
LEFT JOIN nodes landing ON landing.id = e.landing_node_id
WHERE rei.route_id = ?
ORDER BY e.ip ASC`, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]routeExitIPRecord, 0)
	for rows.Next() {
		var item routeExitIPRecord
		var enabled int
		if err := rows.Scan(&item.ID, &item.RouteID, &item.ExitIPID, &item.Address, &item.Family,
			&item.SourceType, &item.OwnerNodeID, &item.OwnerNodeName, &item.OwnerNodeType,
			&item.LandingNodeID, &item.LandingNodeName, &item.Scope, &item.AllocationWeight, &enabled); err != nil {
			return nil, err
		}
		item.Enabled = enabled == 1
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func routeExitIPRecordData(item routeExitIPRecord) map[string]any {
	return map[string]any{
		"id": item.ID, "routeId": item.RouteID, "exitIpId": item.ExitIPID, "address": item.Address,
		"family": item.Family, "sourceType": item.SourceType, "ownerNodeId": nullableString(item.OwnerNodeID),
		"ownerNodeName": nullableString(item.OwnerNodeName), "ownerNodeType": nullableString(item.OwnerNodeType),
		"landingNodeId": nullableString(item.LandingNodeID), "landingNodeName": nullableString(item.LandingNodeName),
		"scope":            item.Scope,
		"allocationWeight": item.AllocationWeight, "enabled": item.Enabled,
	}
}

func validateAllocationWeight(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 || value > 1000000 {
		return errors.New("allocationWeight must be greater than 0")
	}
	return nil
}

func (s *Server) routeLandingNodeID(routeID string) (string, error) {
	var landingNodeID string
	err := s.db.QueryRow(`SELECT landing_node_id FROM routes WHERE id = ?`, routeID).Scan(&landingNodeID)
	return landingNodeID, err
}

func (s *Server) routeNodeIDs(routeID string) (relayNodeID, landingNodeID string, err error) {
	err = s.db.QueryRow(`SELECT relay_node_id, landing_node_id FROM routes WHERE id = ?`, routeID).Scan(&relayNodeID, &landingNodeID)
	return
}

func inferExitIPScope(sourceType, ownerNodeID, relayNodeID, landingNodeID string) string {
	if sourceType == "s5" {
		return "external"
	}
	if ownerNodeID == relayNodeID {
		return "relay"
	}
	return "landing"
}

func validateExitIPBinding(scope, sourceType, ownerNodeID, relayNodeID, landingNodeID string) error {
	switch scope {
	case "relay":
		if sourceType != "node" || ownerNodeID != relayNodeID {
			return errors.New("relay scope requires an exit IP owned by this route relay node")
		}
	case "landing":
		if sourceType != "node" || ownerNodeID != landingNodeID {
			return errors.New("landing scope requires an exit IP owned by this route landing node")
		}
	case "external":
		if sourceType != "s5" {
			return errors.New("external scope requires an S5 exit IP")
		}
	default:
		return errors.New("scope must be relay, landing, or external")
	}
	return nil
}

func (s *Server) routeExitIPs(w http.ResponseWriter, r *http.Request) {
	routeID := strings.TrimSpace(r.PathValue("id"))
	if routeID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "route id is required")
		return
	}
	if _, err := s.routeLandingNodeID(routeID); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route")
		return
	}
	items, err := s.readRouteExitIPs(routeID)
	if err != nil {
		s.logger.Error("read route exit IP bindings", "route_id", routeID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route exit IP bindings")
		return
	}
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, routeExitIPRecordData(item))
	}
	writeSuccess(w, result)
}

func (s *Server) bindRouteExitIP(w http.ResponseWriter, r *http.Request) {
	routeID := strings.TrimSpace(r.PathValue("id"))
	if routeID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "route id is required")
		return
	}
	relayNodeID, landingNodeID, err := s.routeNodeIDs(routeID)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route")
		return
	}
	var request routeExitIPUpsertRequest
	if err := decodeJSON(r, &request); err != nil || request.ExitIPID == nil || strings.TrimSpace(*request.ExitIPID) == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "exitIpId is required")
		return
	}
	exitIPID := strings.TrimSpace(*request.ExitIPID)
	var sourceType, ownerNodeID string
	if err := s.db.QueryRow(`SELECT COALESCE(source_type, 'node'), COALESCE(owner_node_id, landing_node_id, '') FROM exit_ips WHERE id = ?`, exitIPID).Scan(&sourceType, &ownerNodeID); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "exit IP not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP")
		return
	}
	scope := inferExitIPScope(sourceType, ownerNodeID, relayNodeID, landingNodeID)
	if request.Scope != nil && strings.TrimSpace(*request.Scope) != "" {
		scope = strings.TrimSpace(*request.Scope)
	}
	if err := validateExitIPBinding(scope, sourceType, ownerNodeID, relayNodeID, landingNodeID); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	weight := 1.0
	if request.AllocationWeight != nil {
		weight = *request.AllocationWeight
	}
	if err := validateAllocationWeight(weight); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	enabled := true
	if request.Enabled != nil {
		enabled = *request.Enabled
	}
	var existingID string
	if err := s.db.QueryRow(`SELECT id FROM route_exit_ips WHERE route_id = ? AND exit_ip_id = ?`, routeID, exitIPID).Scan(&existingID); err == nil {
		writeFailure(w, http.StatusConflict, validationCode, "exit IP is already bound to this route")
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check route exit IP binding")
		return
	}
	bindingID := newID()
	if _, err := s.db.Exec(`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, scope, allocation_weight, enabled) VALUES (?, ?, ?, ?, ?, ?)`,
		bindingID, routeID, exitIPID, scope, weight, boolInt(enabled)); err != nil {
		s.logger.Error("bind route exit IP", "route_id", routeID, "exit_ip_id", exitIPID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not bind exit IP to route")
		return
	}
	items, err := s.readRouteExitIPs(routeID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read created route exit IP binding")
		return
	}
	var created routeExitIPRecord
	for _, item := range items {
		if item.ID == bindingID {
			created = item
			break
		}
	}
	s.writeAuditLog(r, "route.exit_ip.bind", "route_exit_ip", bindingID, nil, routeExitIPRecordData(created))
	writeSuccess(w, routeExitIPRecordData(created))
}

func (s *Server) updateRouteExitIP(w http.ResponseWriter, r *http.Request) {
	routeID := strings.TrimSpace(r.PathValue("id"))
	exitIPID := strings.TrimSpace(r.PathValue("exitIpId"))
	if routeID == "" || exitIPID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "route id and exit IP id are required")
		return
	}
	if _, err := s.routeLandingNodeID(routeID); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route")
		return
	}
	items, err := s.readRouteExitIPs(routeID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route exit IP bindings")
		return
	}
	var existing routeExitIPRecord
	found := false
	for _, item := range items {
		if item.ExitIPID == exitIPID {
			existing = item
			found = true
			break
		}
	}
	if !found {
		writeFailure(w, http.StatusNotFound, notFoundCode, "route exit IP binding not found")
		return
	}
	var request routeExitIPUpsertRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid route exit IP payload")
		return
	}
	weight := existing.AllocationWeight
	if request.AllocationWeight != nil {
		weight = *request.AllocationWeight
	}
	if err := validateAllocationWeight(weight); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	enabled := existing.Enabled
	if request.Enabled != nil {
		enabled = *request.Enabled
	}
	if !enabled && existing.Enabled {
		var fixedUserCount int
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM user_routes WHERE route_exit_ip_id = ? AND active_to IS NULL`, existing.ID).Scan(&fixedUserCount); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check user route assignments")
			return
		}
		if fixedUserCount > 0 {
			writeFailure(w, http.StatusConflict, validationCode, "exit IP is fixed to active users; reassign them before disabling")
			return
		}
	}
	scope := existing.Scope
	if request.Scope != nil && strings.TrimSpace(*request.Scope) != "" {
		scope = strings.TrimSpace(*request.Scope)
		relayNodeID, landingNodeID, routeErr := s.routeNodeIDs(routeID)
		if routeErr != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route")
			return
		}
		var sourceType, ownerNodeID string
		if err := s.db.QueryRow(`SELECT COALESCE(source_type, 'node'), COALESCE(owner_node_id, landing_node_id, '') FROM exit_ips WHERE id = ?`, existing.ExitIPID).Scan(&sourceType, &ownerNodeID); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP")
			return
		}
		if err := validateExitIPBinding(scope, sourceType, ownerNodeID, relayNodeID, landingNodeID); err != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
			return
		}
	}
	if _, err := s.db.Exec(`UPDATE route_exit_ips SET scope = ?, allocation_weight = ?, enabled = ? WHERE id = ?`, scope, weight, boolInt(enabled), existing.ID); err != nil {
		s.logger.Error("update route exit IP binding", "binding_id", existing.ID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update route exit IP binding")
		return
	}
	updatedItems, err := s.readRouteExitIPs(routeID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read updated route exit IP binding")
		return
	}
	var updated routeExitIPRecord
	for _, item := range updatedItems {
		if item.ID == existing.ID {
			updated = item
			break
		}
	}
	s.writeAuditLog(r, "route.exit_ip.update", "route_exit_ip", existing.ID, routeExitIPRecordData(existing), routeExitIPRecordData(updated))
	writeSuccess(w, routeExitIPRecordData(updated))
}

func (s *Server) unbindRouteExitIP(w http.ResponseWriter, r *http.Request) {
	routeID := strings.TrimSpace(r.PathValue("id"))
	exitIPID := strings.TrimSpace(r.PathValue("exitIpId"))
	if routeID == "" || exitIPID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "route id and exit IP id are required")
		return
	}
	items, err := s.readRouteExitIPs(routeID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read route exit IP bindings")
		return
	}
	var existing routeExitIPRecord
	found := false
	for _, item := range items {
		if item.ExitIPID == exitIPID {
			existing = item
			found = true
			break
		}
	}
	if !found {
		if _, routeErr := s.routeLandingNodeID(routeID); errors.Is(routeErr, sql.ErrNoRows) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "route not found")
		} else {
			writeFailure(w, http.StatusNotFound, notFoundCode, "route exit IP binding not found")
		}
		return
	}
	var fixedUserCount int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM user_routes WHERE route_exit_ip_id = ? AND active_to IS NULL`, existing.ID).Scan(&fixedUserCount); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check user route assignments")
		return
	}
	if fixedUserCount > 0 {
		writeFailure(w, http.StatusConflict, validationCode, "exit IP is fixed to active users; reassign them before unbinding")
		return
	}
	if _, err := s.db.Exec(`DELETE FROM route_exit_ips WHERE id = ?`, existing.ID); err != nil {
		s.logger.Error("unbind route exit IP", "binding_id", existing.ID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not unbind route exit IP")
		return
	}
	s.writeAuditLog(r, "route.exit_ip.unbind", "route_exit_ip", existing.ID, routeExitIPRecordData(existing), map[string]any{"deleted": true})
	writeSuccess(w, map[string]any{"id": existing.ID, "routeId": routeID, "exitIpId": exitIPID, "deleted": true})
}

type exitIPUpsertRequest struct {
	Address       *string  `json:"address"`
	SourceType    *string  `json:"sourceType"`
	OwnerNodeID   *string  `json:"ownerNodeId"`
	LandingNodeID *string  `json:"landingNodeId"`
	Family        *int     `json:"family"`
	Provider      *string  `json:"provider"`
	MonthlyCost   *float64 `json:"monthlyCost"`
	Currency      *string  `json:"currency"`
	Enabled       *bool    `json:"enabled"`
	ValidFrom     *string  `json:"validFrom"`
	ValidTo       *string  `json:"validTo"`
	Notes         *string  `json:"notes"`
}

type exitIPRecord struct {
	ID                 string
	Address            string
	SourceType         string
	OwnerNodeID        string
	OwnerNodeName      string
	OwnerNodeType      string
	LandingNodeID      string
	LandingNodeName    string
	Family             int
	Provider           string
	MonthlyCost        float64
	Currency           string
	Enabled            bool
	ValidFrom          string
	ValidTo            string
	Notes              string
	AllocatedUserCount int
	CheckedAt          string
}

type exitIPValues struct {
	Address       string
	SourceType    string
	OwnerNodeID   string
	LandingNodeID string
	Family        int
	Provider      string
	MonthlyCost   float64
	Currency      string
	Enabled       bool
	ValidFrom     string
	ValidTo       string
	Notes         string
}

func (s *Server) readExitIP(id string) (exitIPRecord, error) {
	var record exitIPRecord
	var enabled int
	err := s.db.QueryRow(`SELECT e.id, e.ip, COALESCE(e.source_type, 'node'), COALESCE(e.owner_node_id, e.landing_node_id, ''),
COALESCE(owner.name, ''), COALESCE(owner.type, ''), COALESCE(e.landing_node_id, ''), COALESCE(landing.name, ''), e.family,
COALESCE(e.provider, ''), e.monthly_cost, e.currency, e.enabled, COALESCE(e.valid_from, ''),
COALESCE(e.valid_to, ''), COALESCE(e.notes, ''), e.updated_at
FROM exit_ips e
LEFT JOIN nodes owner ON owner.id = COALESCE(e.owner_node_id, e.landing_node_id)
LEFT JOIN nodes landing ON landing.id = e.landing_node_id
	WHERE e.id = ?`, id).Scan(&record.ID, &record.Address, &record.SourceType, &record.OwnerNodeID, &record.OwnerNodeName,
		&record.OwnerNodeType, &record.LandingNodeID, &record.LandingNodeName, &record.Family, &record.Provider,
		&record.MonthlyCost, &record.Currency, &enabled, &record.ValidFrom,
		&record.ValidTo, &record.Notes, &record.CheckedAt)
	if err != nil {
		return exitIPRecord{}, err
	}
	record.Enabled = enabled == 1
	if err := s.db.QueryRow(`SELECT COUNT(DISTINCT p.user_id) FROM user_paths p
JOIN users u ON u.id = p.user_id
WHERE p.exit_ip_id = ? AND p.active_to IS NULL
AND u.status <> 'disabled' AND datetime(u.created_at) <= datetime('now')
AND (u.expiry_time IS NULL OR datetime(u.expiry_time) >= datetime('now'))`, id).Scan(&record.AllocatedUserCount); err != nil {
		return exitIPRecord{}, err
	}
	return record, nil
}

func exitIPRecordData(record exitIPRecord) map[string]any {
	status := "disabled"
	if record.Enabled {
		status = "active"
	}
	return map[string]any{
		"id": record.ID, "address": record.Address,
		"sourceType": record.SourceType, "ownerNodeId": nullableString(record.OwnerNodeID),
		"ownerNodeName": nullableString(record.OwnerNodeName), "ownerNodeType": nullableString(record.OwnerNodeType),
		"landingNodeId": nullableString(record.LandingNodeID), "landingNodeName": nullableString(record.LandingNodeName),
		"family": record.Family, "provider": nullableString(record.Provider),
		"monthlyCost": record.MonthlyCost, "currency": record.Currency, "enabled": record.Enabled,
		"validFrom": nullableString(record.ValidFrom), "validTo": nullableString(record.ValidTo),
		"notes": nullableString(record.Notes), "allocatedUserCount": record.AllocatedUserCount,
		"checkedAt": nullableString(record.CheckedAt), "status": status,
	}
}

func exitIPValuesFromRequest(request exitIPUpsertRequest, existing *exitIPRecord) (exitIPValues, error) {
	values := exitIPValues{Family: 4, Currency: "CNY", Enabled: true, SourceType: "node"}
	if existing != nil {
		values = exitIPValues{
			Address: existing.Address, SourceType: existing.SourceType, OwnerNodeID: existing.OwnerNodeID,
			LandingNodeID: existing.LandingNodeID, Family: existing.Family,
			Provider: existing.Provider, MonthlyCost: existing.MonthlyCost, Currency: existing.Currency,
			Enabled: existing.Enabled, ValidFrom: existing.ValidFrom, ValidTo: existing.ValidTo, Notes: existing.Notes,
		}
		if values.SourceType == "" {
			values.SourceType = "node"
		}
	}
	if request.Address != nil {
		values.Address = strings.TrimSpace(*request.Address)
	}
	if request.SourceType != nil {
		values.SourceType = strings.ToLower(strings.TrimSpace(*request.SourceType))
		if values.SourceType == "s5" {
			// Explicitly switching an existing node asset to S5 clears its
			// previous owner even when JSON carries ownerNodeId: null.
			values.OwnerNodeID = ""
		}
	}
	if request.OwnerNodeID != nil {
		values.OwnerNodeID = strings.TrimSpace(*request.OwnerNodeID)
	}
	// landingNodeId is retained as a compatibility alias for older clients.
	if request.LandingNodeID != nil {
		values.OwnerNodeID = strings.TrimSpace(*request.LandingNodeID)
	}
	if request.Family != nil {
		values.Family = *request.Family
	}
	if request.Provider != nil {
		values.Provider = strings.TrimSpace(*request.Provider)
	}
	if request.MonthlyCost != nil {
		values.MonthlyCost = *request.MonthlyCost
	}
	if request.Currency != nil {
		values.Currency = strings.TrimSpace(*request.Currency)
	}
	if request.Enabled != nil {
		values.Enabled = *request.Enabled
	}
	if request.ValidFrom != nil {
		values.ValidFrom = strings.TrimSpace(*request.ValidFrom)
	}
	if request.ValidTo != nil {
		values.ValidTo = strings.TrimSpace(*request.ValidTo)
	}
	if request.Notes != nil {
		values.Notes = strings.TrimSpace(*request.Notes)
	}
	if values.SourceType == "node" {
		// The legacy column is only populated for landing-node assets.  Relay
		// assets still have a null landingNodeId but remain addressable by owner.
		values.LandingNodeID = values.OwnerNodeID
	} else {
		values.OwnerNodeID = ""
		values.LandingNodeID = ""
	}
	if request.Family == nil && request.Address != nil {
		if parsed := net.ParseIP(values.Address); parsed != nil {
			if parsed.To4() == nil {
				values.Family = 6
			} else {
				values.Family = 4
			}
		}
	}
	return values, validateExitIPValues(values)
}

func validateExitIPValues(values exitIPValues) error {
	parsed := net.ParseIP(values.Address)
	if parsed == nil {
		return errors.New("address must be a valid IPv4 or IPv6 address")
	}
	if values.Family != 4 && values.Family != 6 {
		return errors.New("family must be 4 or 6")
	}
	if values.Family == 4 && parsed.To4() == nil {
		return errors.New("family 4 requires an IPv4 address")
	}
	if values.Family == 6 && parsed.To4() != nil {
		return errors.New("family 6 requires an IPv6 address")
	}
	if values.SourceType != "node" && values.SourceType != "s5" {
		return errors.New("sourceType must be node or s5")
	}
	if values.SourceType == "node" && values.OwnerNodeID == "" {
		return errors.New("ownerNodeId is required for node exit IPs")
	}
	if values.SourceType == "s5" && values.OwnerNodeID != "" {
		return errors.New("ownerNodeId must be empty for S5 exit IPs")
	}
	if len(values.Provider) > 200 {
		return errors.New("provider is too long")
	}
	if math.IsNaN(values.MonthlyCost) || math.IsInf(values.MonthlyCost, 0) || values.MonthlyCost < 0 || values.MonthlyCost > 100000000 {
		return errors.New("monthlyCost must be a non-negative number")
	}
	if values.Currency != "CNY" {
		return errors.New("currency currently supports CNY only")
	}
	if len(values.Notes) > 2000 {
		return errors.New("exit IP notes are too long")
	}
	from, fromSet, err := parseRouteDate(values.ValidFrom)
	if err != nil {
		return errors.New("validFrom must be RFC3339 or YYYY-MM-DD")
	}
	to, toSet, err := parseRouteDate(values.ValidTo)
	if err != nil {
		return errors.New("validTo must be RFC3339 or YYYY-MM-DD")
	}
	if fromSet && toSet && to.Before(from) {
		return errors.New("validTo must not be earlier than validFrom")
	}
	return nil
}

func validateExitIPRelations(queryer sqlQueryer, sourceType, ownerNodeID string) error {
	if sourceType == "s5" {
		return nil
	}
	var nodeType string
	if err := queryer.QueryRow(`SELECT type FROM nodes WHERE id = ? AND deleted_at IS NULL`, ownerNodeID).Scan(&nodeType); errors.Is(err, sql.ErrNoRows) {
		return errors.New("owner node not found")
	} else if err != nil {
		return fmt.Errorf("read owner node: %w", err)
	}
	if nodeType != "landing" && nodeType != "relay" {
		return errors.New("ownerNodeId must refer to a relay or landing node")
	}
	return nil
}

// Older clients used landingNodeId as a landing-only selector. Keep that
// contract when the new sourceType/ownerNodeId fields are absent, while
// allowing the new API to explicitly create relay-owned node assets.
func validateLegacyLandingAlias(queryer sqlQueryer, request exitIPUpsertRequest, ownerNodeID string) error {
	if request.SourceType != nil || request.OwnerNodeID != nil || request.LandingNodeID == nil || ownerNodeID == "" {
		return nil
	}
	var nodeType string
	if err := queryer.QueryRow(`SELECT type FROM nodes WHERE id = ? AND deleted_at IS NULL`, ownerNodeID).Scan(&nodeType); errors.Is(err, sql.ErrNoRows) {
		return errors.New("landing node not found")
	} else if err != nil {
		return fmt.Errorf("read landing node: %w", err)
	}
	if nodeType != "landing" {
		return errors.New("landingNodeId must refer to a landing node")
	}
	return nil
}

func (s *Server) exitIPDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "exit IP id is required")
		return
	}
	record, err := s.readExitIP(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "exit IP not found")
		return
	}
	if err != nil {
		s.logger.Error("read exit IP detail", "exit_ip_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP detail")
		return
	}
	writeSuccess(w, exitIPRecordData(record))
}

func (s *Server) createExitIP(w http.ResponseWriter, r *http.Request) {
	var request exitIPUpsertRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid exit IP payload")
		return
	}
	values, err := exitIPValuesFromRequest(request, nil)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if err := validateExitIPRelations(s.db, values.SourceType, values.OwnerNodeID); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if err := validateLegacyLandingAlias(s.db, request, values.OwnerNodeID); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	var duplicateID string
	err = s.db.QueryRow(`SELECT id FROM exit_ips WHERE ip = ? AND ((source_type = 's5' AND ? = 's5') OR (source_type = 'node' AND ? = 'node' AND owner_node_id = ?))`, values.Address, values.SourceType, values.SourceType, values.OwnerNodeID).Scan(&duplicateID)
	if err == nil {
		writeFailure(w, http.StatusConflict, validationCode, "exit IP already exists for this source")
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check duplicate exit IP")
		return
	}
	id := newID()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	legacyLandingNodeID := ""
	if values.SourceType == "node" {
		var ownerType string
		if err := s.db.QueryRow(`SELECT type FROM nodes WHERE id = ? AND deleted_at IS NULL`, values.OwnerNodeID).Scan(&ownerType); err != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, "owner node not found")
			return
		}
		if ownerType == "landing" {
			legacyLandingNodeID = values.OwnerNodeID
		}
	}
	if _, err := s.db.Exec(`INSERT INTO exit_ips (id, landing_node_id, source_type, owner_node_id, ip, family, provider, monthly_cost, currency, enabled, valid_from, valid_to, notes, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, nullableDBString(legacyLandingNodeID), values.SourceType, nullableDBString(values.OwnerNodeID), values.Address, values.Family,
		nullableRouteValue(values.Provider), values.MonthlyCost, values.Currency, boolInt(values.Enabled), nullableRouteValue(values.ValidFrom),
		nullableRouteValue(values.ValidTo), nullableRouteValue(values.Notes), now, now); err != nil {
		s.logger.Error("create exit IP", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create exit IP")
		return
	}
	record, err := s.readExitIP(id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read created exit IP")
		return
	}
	s.writeAuditLog(r, "exit_ip.create", "exit_ip", id, nil, exitIPRecordData(record))
	writeSuccess(w, exitIPRecordData(record))
}

func (s *Server) updateExitIP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "exit IP id is required")
		return
	}
	existing, err := s.readExitIP(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "exit IP not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP")
		return
	}
	var request exitIPUpsertRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid exit IP payload")
		return
	}
	values, err := exitIPValuesFromRequest(request, &existing)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if err := validateExitIPRelations(s.db, values.SourceType, values.OwnerNodeID); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if err := validateLegacyLandingAlias(s.db, request, values.OwnerNodeID); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if values.SourceType != existing.SourceType || values.OwnerNodeID != existing.OwnerNodeID {
		var bindingCount, pathCount int
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM route_exit_ips WHERE exit_ip_id = ?`, id).Scan(&bindingCount); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check exit IP bindings")
			return
		}
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM user_paths WHERE exit_ip_id = ? AND active_to IS NULL`, id).Scan(&pathCount); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check user path assignments")
			return
		}
		if bindingCount > 0 || pathCount > 0 {
			writeFailure(w, http.StatusConflict, validationCode, "bound exit IP cannot move to another landing node")
			return
		}
	}
	var duplicateID string
	err = s.db.QueryRow(`SELECT id FROM exit_ips WHERE ip = ? AND id <> ? AND ((source_type = 's5' AND ? = 's5') OR (source_type = 'node' AND ? = 'node' AND owner_node_id = ?))`, values.Address, id, values.SourceType, values.SourceType, values.OwnerNodeID).Scan(&duplicateID)
	if err == nil {
		writeFailure(w, http.StatusConflict, validationCode, "exit IP already exists for this source")
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check duplicate exit IP")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	legacyLandingNodeID := ""
	if values.SourceType == "node" {
		var ownerType string
		if err := s.db.QueryRow(`SELECT type FROM nodes WHERE id = ? AND deleted_at IS NULL`, values.OwnerNodeID).Scan(&ownerType); err != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, "owner node not found")
			return
		}
		if ownerType == "landing" {
			legacyLandingNodeID = values.OwnerNodeID
		}
	}
	if _, err := s.db.Exec(`UPDATE exit_ips SET landing_node_id = ?, source_type = ?, owner_node_id = ?, ip = ?, family = ?, provider = ?, monthly_cost = ?, currency = ?, enabled = ?, valid_from = ?, valid_to = ?, notes = ?, updated_at = ? WHERE id = ?`,
		nullableDBString(legacyLandingNodeID), values.SourceType, nullableDBString(values.OwnerNodeID), values.Address, values.Family, nullableRouteValue(values.Provider), values.MonthlyCost, values.Currency,
		boolInt(values.Enabled), nullableRouteValue(values.ValidFrom), nullableRouteValue(values.ValidTo), nullableRouteValue(values.Notes), now, id); err != nil {
		s.logger.Error("update exit IP", "exit_ip_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update exit IP")
		return
	}
	updated, err := s.readExitIP(id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read updated exit IP")
		return
	}
	s.writeAuditLog(r, "exit_ip.update", "exit_ip", id, exitIPRecordData(existing), exitIPRecordData(updated))
	writeSuccess(w, exitIPRecordData(updated))
}

func (s *Server) deleteExitIP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "exit IP id is required")
		return
	}
	record, err := s.readExitIP(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "exit IP not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP")
		return
	}
	var bindingCount, pathCount int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM route_exit_ips WHERE exit_ip_id = ?`, id).Scan(&bindingCount); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check exit IP bindings")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM user_paths WHERE exit_ip_id = ? AND active_to IS NULL`, id).Scan(&pathCount); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check user path assignments")
		return
	}
	if bindingCount > 0 || pathCount > 0 {
		writeFailure(w, http.StatusConflict, validationCode, "exit IP has active assignments; disable it or change the users first")
		return
	}
	if _, err := s.db.Exec(`DELETE FROM exit_ips WHERE id = ?`, id); err != nil {
		s.logger.Error("delete exit IP", "exit_ip_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not delete exit IP")
		return
	}
	s.writeAuditLog(r, "exit_ip.delete", "exit_ip", id, exitIPRecordData(record), map[string]any{"deleted": true})
	writeSuccess(w, map[string]any{"id": id, "deleted": true})
}

func (s *Server) exitIPs(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)
	where := []string{"1 = 1"}
	args := make([]any, 0, 3)
	if query.keyword != "" {
		where = append(where, `(e.ip LIKE ? OR e.provider LIKE ? OR owner.name LIKE ?)`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like)
	}
	if query.status != "" {
		if query.status == "active" {
			where = append(where, "e.enabled = 1")
		} else if query.status == "disabled" {
			where = append(where, "e.enabled = 0")
		}
	}
	base := `FROM exit_ips e
LEFT JOIN nodes owner ON owner.id = COALESCE(e.owner_node_id, e.landing_node_id)
LEFT JOIN nodes landing ON landing.id = e.landing_node_id
WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count exit IPs")
		return
	}
	rows, err := s.db.Query(`SELECT e.id, e.ip, COALESCE(e.source_type, 'node'), COALESCE(e.owner_node_id, e.landing_node_id, ''), COALESCE(owner.name, ''), COALESCE(owner.type, ''),
COALESCE(e.landing_node_id, ''), COALESCE(landing.name, ''), e.family, COALESCE(e.provider, ''), e.enabled, e.monthly_cost, e.currency,
(SELECT COUNT(DISTINCT p.user_id) FROM user_paths p
JOIN users u ON u.id = p.user_id
WHERE p.exit_ip_id = e.id AND p.active_to IS NULL
AND u.status <> 'disabled' AND datetime(u.created_at) <= datetime('now')
AND (u.expiry_time IS NULL OR datetime(u.expiry_time) >= datetime('now'))), e.updated_at
`+base+` ORDER BY e.ip ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IPs")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, ip, sourceType, ownerNodeID, ownerNodeName, ownerNodeType, landingNodeID, landing, provider, currency, checkedAt string
		var family, enabled, allocated int
		var monthlyCost float64
		if err := rows.Scan(&id, &ip, &sourceType, &ownerNodeID, &ownerNodeName, &ownerNodeType, &landingNodeID, &landing, &family, &provider, &enabled, &monthlyCost, &currency, &allocated, &checkedAt); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode exit IPs")
			return
		}
		status := "disabled"
		if enabled == 1 {
			status = "active"
		}
		items = append(items, map[string]any{
			"id": id, "address": ip, "sourceType": sourceType, "ownerNodeId": nullableString(ownerNodeID), "ownerNodeName": nullableString(ownerNodeName), "ownerNodeType": nullableString(ownerNodeType),
			"landingNodeId": nullableString(landingNodeID), "landingNodeName": nullableString(landing), "family": family, "provider": nullableString(provider),
			"status": status, "monthlyCost": monthlyCost, "currency": currency,
			"allocatedUserCount": allocated, "checkedAt": nullableString(checkedAt),
		})
	}
	writeSuccess(w, s.pageResponse(items, total, query))
}

type financeResponse struct {
	Period             string           `json:"period"`
	Currency           string           `json:"currency"`
	EffectiveUserCount int              `json:"effectiveUserCount"`
	MonthIncome        float64          `json:"monthIncome"`
	CashIncome         float64          `json:"cashIncome"`
	MonthCost          float64          `json:"monthCost"`
	GrossProfit        float64          `json:"grossProfit"`
	Breakdown          []map[string]any `json:"breakdown,omitempty"`
	DataAt             *string          `json:"dataAt"`
}

func (s *Server) finance(w http.ResponseWriter, r *http.Request) {
	s.refreshOperationalStatuses(time.Now().UTC())
	period := r.URL.Query().Get("period")
	if period == "" {
		period = time.Now().UTC().Format("2006-01")
	}
	result, err := s.financeSummary(period)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	writeSuccess(w, result)
}

func (s *Server) financeSummary(period string) (financeResponse, error) {
	start, end, err := monthBounds(period)
	if err != nil {
		return financeResponse{}, errors.New("period must use YYYY-MM format")
	}
	result := financeResponse{Period: period, Currency: "CNY", Breakdown: make([]map[string]any, 0)}
	result.DataAt = s.latestSuccessfulSyncAt()
	// Revenue is a temporal estimate. The current operational status is
	// intentionally not used here: a user that is expired today can still be
	// valid for a historical month. A user is effective when it was created
	// before the month ended, had not expired at the month start, and was not
	// explicitly disabled. This preserves the central/X-Panel disabled state
	// while keeping historical periods stable as time moves forward.
	if err := s.db.QueryRow(`SELECT COUNT(*)
FROM users
WHERE currency = 'CNY'
  AND status <> 'disabled'
  AND datetime(created_at) < datetime(?)
	  AND (expiry_time IS NULL OR datetime(expiry_time) >= datetime(?))`, end, start).Scan(&result.EffectiveUserCount); err != nil {
		return financeResponse{}, err
	}
	result.MonthIncome, err = s.accruedUserIncome(start, end)
	if err != nil {
		return financeResponse{}, err
	}
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM user_billing_records
WHERE currency = 'CNY' AND status = 'confirmed' AND paid_at IS NOT NULL
  AND datetime(paid_at) >= datetime(?) AND datetime(paid_at) < datetime(?)`, start, end).Scan(&result.CashIncome); err != nil {
		return financeResponse{}, err
	}
	var nodeCost, otherCost, exitCost float64
	// Cost dates are stored as calendar dates with inclusive effective_to.
	// Compare them as dates against the half-open month [start, end) so a
	// record beginning on the first day of the next month is not included.
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(monthly_amount), 0) FROM node_costs WHERE currency = 'CNY' AND date(effective_from) < date(?) AND (effective_to IS NULL OR date(effective_to) >= date(?))`, end, start).Scan(&nodeCost); err != nil {
		return financeResponse{}, err
	}
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(monthly_amount), 0) FROM other_costs WHERE currency = 'CNY' AND date(effective_from) < date(?) AND (effective_to IS NULL OR date(effective_to) >= date(?))`, end, start).Scan(&otherCost); err != nil {
		return financeResponse{}, err
	}
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(monthly_cost), 0) FROM exit_ips WHERE currency = 'CNY' AND enabled = 1 AND (valid_from IS NULL OR date(valid_from) < date(?)) AND (valid_to IS NULL OR date(valid_to) >= date(?))`, end, start).Scan(&exitCost); err != nil {
		return financeResponse{}, err
	}
	result.Breakdown = append(result.Breakdown,
		map[string]any{"label": "用户月费收入", "amount": result.MonthIncome},
		map[string]any{"label": "节点成本", "amount": nodeCost},
		map[string]any{"label": "出口 IP 成本", "amount": exitCost},
		map[string]any{"label": "其他成本", "amount": otherCost},
	)
	result.MonthCost = nodeCost + exitCost + otherCost
	result.GrossProfit = result.MonthIncome - result.MonthCost
	return result, nil
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)
	where := []string{"1 = 1"}
	args := make([]any, 0, 3)
	if query.keyword != "" {
		where = append(where, `(e.event_type LIKE ? OR e.message LIKE ? OR n.name LIKE ?)`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like)
	}
	if query.severity != "" {
		where = append(where, "e.severity = ?")
		args = append(args, query.severity)
	}
	base := `FROM node_events e LEFT JOIN nodes n ON n.id = e.node_id WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count events")
		return
	}
	rows, err := s.db.Query(`SELECT e.id, e.event_type, e.severity, COALESCE(n.name, ''), e.message, e.created_at, e.acknowledged
`+base+` ORDER BY e.created_at DESC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read events")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, eventType, severity, nodeName, message, occurredAt string
		var acknowledged int
		if err := rows.Scan(&id, &eventType, &severity, &nodeName, &message, &occurredAt, &acknowledged); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode events")
			return
		}
		items = append(items, map[string]any{
			"id": id, "type": eventType, "severity": severity, "nodeName": nullableString(nodeName),
			"message": message, "occurredAt": occurredAt, "acknowledged": acknowledged == 1,
		})
	}
	writeSuccess(w, s.pageResponse(items, total, query))
}

func (s *Server) createSession(userID string) (string, string, string, error) {
	token, err := randomToken()
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err := randomToken()
	if err != nil {
		return "", "", "", err
	}
	now := time.Now().UTC()
	expiresAt := now.Add(s.cfg.SessionTTL).Format(time.RFC3339Nano)
	_, err = s.db.Exec(`INSERT INTO sessions (id, admin_user_id, token_hash, refresh_hash, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`, newID(), userID, hashToken(token), hashToken(refreshToken), expiresAt, now.Format(time.RFC3339Nano))
	return token, refreshToken, expiresAt, err
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "authentication required")
			return
		}
		rawToken := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		var current principal
		var expiresAt string
		err := s.db.QueryRow(`SELECT s.id, a.id, a.username, s.expires_at FROM sessions s JOIN admin_users a ON a.id = s.admin_user_id WHERE s.token_hash = ? AND s.revoked_at IS NULL AND a.enabled = 1`, hashToken(rawToken)).Scan(&current.SessionID, &current.UserID, &current.UserName, &expiresAt)
		if err != nil {
			writeFailure(w, http.StatusUnauthorized, unauthorizedCode, "authentication required")
			return
		}
		if isExpired(expiresAt) {
			writeFailure(w, http.StatusUnauthorized, "9999", "session expired")
			return
		}
		if isStateChangingMethod(r.Method) && !s.isAllowedWriteOrigin(r) {
			writeFailure(w, http.StatusForbidden, csrfCode, "request origin is not allowed")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalContextKey{}, current)))
	})
}

func isStateChangingMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// isAllowedWriteOrigin is a defense-in-depth CSRF check for administrator
// Bearer requests. The browser does not attach the Bearer token by itself,
// but validating Origin/Referer prevents a future cookie-auth migration or a
// misconfigured client from turning a cross-site request into a write. CLI
// and Agent-style clients with no browser provenance headers remain supported.
func (s *Server) isAllowedWriteOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		referer := strings.TrimSpace(r.Header.Get("Referer"))
		if referer == "" {
			return true
		}
		parsed, err := url.Parse(referer)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return false
		}
		origin = parsed.Scheme + "://" + parsed.Host
	}
	_, allowed := s.origins[origin]
	return allowed
}

func (s *Server) withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if requestID == "" {
			requestID = newID()
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDContextKey{}, requestID)))
	})
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if _, allowed := s.origins[origin]; allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// withSecurityHeaders applies headers that are safe for the JSON API and do
// not depend on a particular reverse proxy. HSTS is only emitted when the
// request is already HTTPS (directly or as reported by a trusted TLS
// terminator); emitting it on local HTTP would make browser development
// sessions unexpectedly sticky.
func (s *Server) withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		w.Header().Set("Cache-Control", "no-store")
		if r.TLS != nil || strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")), "https") {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

type listQuery struct {
	page     int
	pageSize int
	offset   int
	keyword  string
	status   string
	nodeID   string
	nodeType string
	severity string
}

func parseListQuery(r *http.Request) listQuery {
	page := positiveInt(r.URL.Query().Get("page"), 1)
	pageSize := positiveInt(r.URL.Query().Get("page_size"), 50)
	if pageSize > 200 {
		pageSize = 200
	}
	return listQuery{
		page: page, pageSize: pageSize, offset: (page - 1) * pageSize,
		keyword: strings.TrimSpace(r.URL.Query().Get("keyword")), status: strings.TrimSpace(r.URL.Query().Get("status")),
		nodeID: strings.TrimSpace(r.URL.Query().Get("node_id")), nodeType: strings.TrimSpace(r.URL.Query().Get("node_type")),
		severity: strings.TrimSpace(r.URL.Query().Get("severity")),
	}
}

func (s *Server) pageResponse(items []map[string]any, total int, query listQuery) map[string]any {
	return map[string]any{
		"items": items, "total": total, "page": query.page, "pageSize": query.pageSize,
		"dataAt": s.latestSuccessfulSyncAt(),
	}
}

// latestSuccessfulSyncAt is the timestamp of the newest complete snapshot
// represented in central storage. A null result is intentional when no
// successful Agent sync has been received yet; callers must surface that
// state instead of inventing a current timestamp.
func (s *Server) latestSuccessfulSyncAt() *string {
	var value sql.NullString
	if err := s.db.QueryRow(`SELECT MAX(COALESCE(finished_at, started_at)) FROM sync_runs WHERE status = 'success'`).Scan(&value); err != nil {
		s.logger.Warn("read latest successful sync time", "error", err)
		return nil
	}
	if !value.Valid || strings.TrimSpace(value.String) == "" {
		return nil
	}
	return &value.String
}

func positiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func monthBounds(period string) (string, string, error) {
	start, err := time.Parse("2006-01", period)
	if err != nil {
		return "", "", err
	}
	return start.UTC().Format(time.RFC3339Nano), start.AddDate(0, 1, 0).UTC().Format(time.RFC3339Nano), nil
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nullableFloat(value sql.NullFloat64) any {
	if !value.Valid {
		return nil
	}
	return value.Float64
}

func nullableInt(value sql.NullInt64) any {
	if !value.Valid {
		return nil
	}
	return value.Int64
}

// nullableDBString preserves the distinction between an empty optional value
// and a populated one in SQLite. API consumers still receive null for either
// an unset or intentionally cleared note.
func nullableDBString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

// writeAuditLog is deliberately best-effort: a completed business-field
// update must not be rolled back merely because its auxiliary audit write
// fails. It is only called after the central database update succeeds.
func (s *Server) writeAuditLog(r *http.Request, action, resourceType, resourceID string, before, after any) {
	beforeJSON, err := marshalAuditState(before)
	if err != nil {
		s.logger.Warn("encode audit before state", "action", action, "resource_id", resourceID, "error", err)
		return
	}
	afterJSON, err := marshalAuditState(after)
	if err != nil {
		s.logger.Warn("encode audit after state", "action", action, "resource_id", resourceID, "error", err)
		return
	}

	var adminUserID any
	if current, ok := r.Context().Value(principalContextKey{}).(principal); ok && current.UserID != "" {
		adminUserID = current.UserID
	}
	var requestID any
	if value, ok := r.Context().Value(requestIDContextKey{}).(string); ok && value != "" {
		requestID = value
	}
	var ip any
	if value := clientIP(r.RemoteAddr); value != "" {
		ip = value
	}

	if _, err := s.db.Exec(`INSERT INTO audit_logs (id, admin_user_id, action, resource_type, resource_id, request_id, before_json, after_json, ip, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		newID(), adminUserID, action, resourceType, resourceID, requestID, string(beforeJSON), string(afterJSON), ip, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		s.logger.Warn("write audit log", "action", action, "resource_type", resourceType, "resource_id", resourceID, "error", err)
	}
}

// marshalAuditState round-trips through JSON so maps, slices and structs all
// receive the same recursive redaction. Audit records are for operational
// history, not a second secret store; values under credential-like keys are
// replaced before they can reach SQLite or logs.
func marshalAuditState(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}
	redactAuditValue(decoded)
	return json.Marshal(decoded)
}

func redactAuditValue(value any) {
	switch typed := value.(type) {
	case map[string]any:
		for key, item := range typed {
			if isSensitiveAuditKey(key) {
				typed[key] = "[REDACTED]"
				continue
			}
			redactAuditValue(item)
		}
	case []any:
		for _, item := range typed {
			redactAuditValue(item)
		}
	}
}

func isSensitiveAuditKey(key string) bool {
	lower := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(key, "-", ""), "_", ""))
	for _, marker := range []string{"password", "token", "secret", "privatekey", "credential", "subscription"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func clientIP(remoteAddr string) string {
	remoteAddr = strings.TrimSpace(remoteAddr)
	if remoteAddr == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	if net.ParseIP(remoteAddr) != nil {
		return remoteAddr
	}
	return ""
}

func isExpired(value string) bool {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	return err != nil || !parsed.After(time.Now().UTC())
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, envelope{Code: successCode, Msg: "ok", Data: data})
}

func writeFailure(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, envelope{Code: code, Msg: message, Data: nil})
}

type envelope struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func newID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(bytes)
}

func randomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
