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
	"net/http"
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
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.HandleFunc("POST /api/auth/refreshToken", s.refreshToken)
	mux.Handle("POST /api/auth/logout", s.requireAuth(http.HandlerFunc(s.logout)))
	mux.Handle("GET /api/auth/me", s.requireAuth(http.HandlerFunc(s.getUserInfo)))
	mux.Handle("GET /api/auth/getUserInfo", s.requireAuth(http.HandlerFunc(s.getUserInfo)))
	mux.Handle("GET /api/dashboard", s.requireAuth(http.HandlerFunc(s.dashboard)))
	mux.Handle("GET /api/users", s.requireAuth(http.HandlerFunc(s.users)))
	mux.Handle("GET /api/nodes", s.requireAuth(http.HandlerFunc(s.nodes)))
	mux.Handle("GET /api/routes", s.requireAuth(http.HandlerFunc(s.routes)))
	mux.Handle("GET /api/exit-ips", s.requireAuth(http.HandlerFunc(s.exitIPs)))
	mux.Handle("GET /api/costs/summary", s.requireAuth(http.HandlerFunc(s.finance)))
	mux.Handle("GET /api/events", s.requireAuth(http.HandlerFunc(s.events)))

	return s.withRequestID(s.withCORS(mux))
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
	var response struct {
		GeneratedAt string `json:"generatedAt"`
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
	queries := []struct {
		target *int
		query  string
	}{
		{&response.Nodes.Total, `SELECT COUNT(*) FROM nodes`},
		{&response.Nodes.Online, `SELECT COUNT(*) FROM nodes WHERE health_status = 'online' AND enabled = 1`},
		{&response.Nodes.Relay, `SELECT COUNT(*) FROM nodes WHERE type = 'relay' AND enabled = 1`},
		{&response.Nodes.Landing, `SELECT COUNT(*) FROM nodes WHERE type = 'landing' AND enabled = 1`},
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
	startToday := time.Now().UTC().Truncate(24 * time.Hour).Format(time.RFC3339Nano)
	startMonth := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339Nano)
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(up + down), 0) FROM traffic_snapshots WHERE collected_at >= ?`, startToday).Scan(&response.Traffic.TodayBytes); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read traffic summary")
		return
	}
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(up + down), 0) FROM traffic_snapshots WHERE collected_at >= ?`, startMonth).Scan(&response.Traffic.MonthBytes); err != nil {
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

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)
	where := []string{"1 = 1"}
	args := make([]any, 0, 6)
	if query.keyword != "" {
		where = append(where, `(u.display_name LIKE ? OR i.tag LIKE ? OR n.name LIKE ?)`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like)
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
LEFT JOIN inbounds i ON i.id = ui.inbound_id
LEFT JOIN nodes n ON n.id = i.node_id
WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count users")
		return
	}
	rows, err := s.db.Query(`SELECT u.id, u.display_name, u.status, COALESCE(u.expiry_time, ''),
COALESCE(i.node_id, ''), COALESCE(n.name, ''), COALESCE(i.tag, ''),
COALESCE(i.client_count, 0), COALESCE(i.up, 0) + COALESCE(i.down, 0), COALESCE(i.last_seen_at, '') `+base+`
ORDER BY CASE WHEN u.expiry_time IS NULL THEN 1 ELSE 0 END, u.expiry_time ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read users")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, name, status, expiry, nodeID, nodeName, inboundTag, lastActivity string
		var clientCount int
		var traffic int64
		if err := rows.Scan(&id, &name, &status, &expiry, &nodeID, &nodeName, &inboundTag, &clientCount, &traffic, &lastActivity); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode users")
			return
		}
		items = append(items, map[string]any{
			"id": id, "name": name, "nodeId": nodeID, "nodeName": nodeName, "inboundTag": inboundTag,
			"status": status, "expiresAt": nullableString(expiry), "clientCount": clientCount,
			"trafficBytes": traffic, "lastActivityAt": nullableString(lastActivity),
		})
	}
	writeSuccess(w, pageResponse(items, total, query))
}

func (s *Server) nodes(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)
	where := []string{"1 = 1"}
	args := make([]any, 0, 4)
	if query.keyword != "" {
		where = append(where, `(n.name LIKE ? OR n.hostname LIKE ? OR n.public_ip LIKE ?)`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like)
	}
	if query.status != "" {
		where = append(where, "n.health_status = ?")
		args = append(args, query.status)
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
	rows, err := s.db.Query(`SELECT n.id, n.name, n.type, n.health_status,
COALESCE(NULLIF(n.public_ip, ''), n.hostname, ''), COALESCE(n.xpanel_version, ''), COALESCE(n.xray_version, ''),
COALESCE(n.last_seen_at, ''), COALESCE((SELECT SUM(c.monthly_amount) FROM node_costs c WHERE c.node_id = n.id AND c.currency = 'CNY'
AND c.effective_from <= date('now') AND (c.effective_to IS NULL OR c.effective_to >= date('now'))), 0)
`+base+` ORDER BY n.name ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read nodes")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, name, nodeType, status, host, xpanel, xray, lastSeen string
		var monthlyCost float64
		if err := rows.Scan(&id, &name, &nodeType, &status, &host, &xpanel, &xray, &lastSeen, &monthlyCost); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode nodes")
			return
		}
		items = append(items, map[string]any{
			"id": id, "name": name, "type": nodeType, "status": status, "host": host,
			"xpanelVersion": nullableString(xpanel), "xrayVersion": nullableString(xray),
			"lastSeenAt": nullableString(lastSeen), "monthlyCost": monthlyCost, "currency": "CNY",
		})
	}
	writeSuccess(w, pageResponse(items, total, query))
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
	rows, err := s.db.Query(`SELECT r.id, r.name, COALESCE(relay.name, ''), COALESCE(landing.name, ''),
COALESCE(r.landing_inbound_tag, ''), (SELECT COUNT(*) FROM route_exit_ips rei WHERE rei.route_id = r.id AND rei.enabled = 1),
(SELECT COUNT(DISTINCT ur.user_id) FROM user_routes ur WHERE ur.route_id = r.id AND ur.is_primary = 1 AND ur.active_to IS NULL),
CASE WHEN r.enabled = 1 THEN 'active' ELSE 'disabled' END
`+base+` ORDER BY r.name ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read routes")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, name, relay, landing, inboundTag, status string
		var exitIPCount, userCount int
		if err := rows.Scan(&id, &name, &relay, &landing, &inboundTag, &exitIPCount, &userCount, &status); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode routes")
			return
		}
		items = append(items, map[string]any{
			"id": id, "name": name, "relayNodeName": relay, "landingNodeName": landing,
			"landingInboundTag": nullableString(inboundTag), "exitIpCount": exitIPCount,
			"allocatedUserCount": userCount, "status": status,
		})
	}
	writeSuccess(w, pageResponse(items, total, query))
}

func (s *Server) exitIPs(w http.ResponseWriter, r *http.Request) {
	query := parseListQuery(r)
	where := []string{"1 = 1"}
	args := make([]any, 0, 3)
	if query.keyword != "" {
		where = append(where, `(e.ip LIKE ? OR e.provider LIKE ? OR landing.name LIKE ?)`)
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
	base := `FROM exit_ips e LEFT JOIN nodes landing ON landing.id = e.landing_node_id WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count exit IPs")
		return
	}
	rows, err := s.db.Query(`SELECT e.id, e.ip, COALESCE(landing.name, ''), COALESCE(e.provider, ''), e.enabled, e.monthly_cost, e.currency,
(SELECT COUNT(DISTINCT ur.user_id) FROM route_exit_ips rei JOIN user_routes ur ON ur.route_id = rei.route_id
WHERE rei.exit_ip_id = e.id AND rei.enabled = 1 AND ur.is_primary = 1 AND ur.active_to IS NULL), e.updated_at
`+base+` ORDER BY e.ip ASC LIMIT ? OFFSET ?`, append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IPs")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, ip, landing, provider, currency, checkedAt string
		var enabled, allocated int
		var monthlyCost float64
		if err := rows.Scan(&id, &ip, &landing, &provider, &enabled, &monthlyCost, &currency, &allocated, &checkedAt); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode exit IPs")
			return
		}
		status := "disabled"
		if enabled == 1 {
			status = "active"
		}
		items = append(items, map[string]any{
			"id": id, "address": ip, "landingNodeName": landing, "provider": nullableString(provider),
			"status": status, "monthlyCost": monthlyCost, "currency": currency,
			"allocatedUserCount": allocated, "checkedAt": nullableString(checkedAt),
		})
	}
	writeSuccess(w, pageResponse(items, total, query))
}

type financeResponse struct {
	Period      string           `json:"period"`
	Currency    string           `json:"currency"`
	MonthIncome float64          `json:"monthIncome"`
	MonthCost   float64          `json:"monthCost"`
	GrossProfit float64          `json:"grossProfit"`
	Breakdown   []map[string]any `json:"breakdown,omitempty"`
	DataAt      string           `json:"dataAt,omitempty"`
}

func (s *Server) finance(w http.ResponseWriter, r *http.Request) {
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
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(monthly_fee), 0) FROM users WHERE status IN ('active', 'expiring') AND currency = 'CNY' AND (expiry_time IS NULL OR expiry_time >= ?)`, start).Scan(&result.MonthIncome); err != nil {
		return financeResponse{}, err
	}
	var nodeCost, otherCost, exitCost float64
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(monthly_amount), 0) FROM node_costs WHERE currency = 'CNY' AND effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)`, end, start).Scan(&nodeCost); err != nil {
		return financeResponse{}, err
	}
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(monthly_amount), 0) FROM other_costs WHERE currency = 'CNY' AND effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)`, end, start).Scan(&otherCost); err != nil {
		return financeResponse{}, err
	}
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(monthly_cost), 0) FROM exit_ips WHERE currency = 'CNY' AND enabled = 1 AND (valid_from IS NULL OR valid_from <= ?) AND (valid_to IS NULL OR valid_to >= ?)`, end, start).Scan(&exitCost); err != nil {
		return financeResponse{}, err
	}
	result.Breakdown = append(result.Breakdown,
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
	writeSuccess(w, pageResponse(items, total, query))
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
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalContextKey{}, current)))
	})
}

func (s *Server) withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if requestID == "" {
			requestID = newID()
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
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

func pageResponse(items []map[string]any, total int, query listQuery) map[string]any {
	return map[string]any{"items": items, "total": total, "page": query.page, "pageSize": query.pageSize}
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
