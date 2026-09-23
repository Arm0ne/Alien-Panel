package httpapi

import (
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"
)

type dashboardRangeSpec struct {
	name     string
	duration time.Duration
	bucket   time.Duration
}

var dashboardRanges = map[string]dashboardRangeSpec{
	"today": {name: "today", duration: 24 * time.Hour, bucket: time.Hour},
	"7d":    {name: "7d", duration: 7 * 24 * time.Hour, bucket: 24 * time.Hour},
	"30d":   {name: "30d", duration: 30 * 24 * time.Hour, bucket: 24 * time.Hour},
}

// Dashboard ranges represent business calendar days. The panel is operated
// in China, so a "today" range starts at midnight in Asia/Shanghai and is
// converted to UTC only when querying the database (timestamps are stored in
// UTC). Keeping this boundary explicit avoids using the container's timezone
// or accidentally making the business day start at 08:00 in the UI.
var dashboardBusinessLocation = time.FixedZone("Asia/Shanghai", 8*60*60)

type dashboardTrafficPoint struct {
	Time          string `json:"time"`
	UploadBytes   int64  `json:"uploadBytes"`
	DownloadBytes int64  `json:"downloadBytes"`
	TotalBytes    int64  `json:"totalBytes"`
	SampleCount   int    `json:"sampleCount"`
	ResetDetected bool   `json:"resetDetected"`
	HasGap        bool   `json:"hasGap"`
}

type dashboardTrafficSummary struct {
	UploadBytes   int64   `json:"uploadBytes"`
	DownloadBytes int64   `json:"downloadBytes"`
	TotalBytes    int64   `json:"totalBytes"`
	SampleCount   int     `json:"sampleCount"`
	Coverage      float64 `json:"coverage"`
}

type dashboardTrafficTrend struct {
	Range   string                  `json:"range"`
	From    string                  `json:"from"`
	To      string                  `json:"to"`
	Bucket  string                  `json:"bucket"`
	DataAt  *string                 `json:"dataAt"`
	Points  []dashboardTrafficPoint `json:"points"`
	Summary dashboardTrafficSummary `json:"summary"`
}

type dashboardNodeTrafficItem struct {
	NodeID        string  `json:"nodeId"`
	NodeName      string  `json:"nodeName"`
	NodeType      string  `json:"nodeType"`
	Status        string  `json:"status"`
	UploadBytes   int64   `json:"uploadBytes"`
	DownloadBytes int64   `json:"downloadBytes"`
	TotalBytes    int64   `json:"totalBytes"`
	Share         float64 `json:"share"`
	DataAt        *string `json:"dataAt"`
}

type dashboardUserTrafficItem struct {
	UserID         string  `json:"userId"`
	UserName       string  `json:"userName"`
	NodeID         string  `json:"nodeId"`
	NodeName       string  `json:"nodeName"`
	InboundID      string  `json:"inboundId"`
	InboundTag     string  `json:"inboundTag"`
	ClientCount    int     `json:"clientCount"`
	Status         string  `json:"status"`
	ExpiresAt      *string `json:"expiresAt"`
	UploadBytes    int64   `json:"uploadBytes"`
	DownloadBytes  int64   `json:"downloadBytes"`
	TotalBytes     int64   `json:"totalBytes"`
	LastActivityAt *string `json:"lastActivityAt"`
	DataAt         *string `json:"dataAt"`
}

type dashboardEventItem struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"`
	Category       string  `json:"category"`
	Severity       string  `json:"severity"`
	Title          string  `json:"title"`
	NodeID         *string `json:"nodeId"`
	NodeName       string  `json:"nodeName"`
	Message        string  `json:"message"`
	OccurredAt     string  `json:"occurredAt"`
	Acknowledged   bool    `json:"acknowledged"`
	RequiresAction bool    `json:"requiresAction"`
	Status         string  `json:"status"`
}

type dashboardInbound struct {
	id           string
	nodeID       string
	nodeName     string
	nodeStatus   string
	userID       string
	userName     string
	inboundTag   string
	status       string
	expiresAt    *string
	clientCount  int
	lastActivity *string
}

type dashboardSnapshot struct {
	inboundID string
	at        time.Time
	up        int64
	down      int64
	reset     bool
}

type dashboardInboundTraffic struct {
	uploadBytes   int64
	downloadBytes int64
	totalBytes    int64
	lastActivity  *string
}

type dashboardNodeTraffic struct {
	nodeID        string
	nodeName      string
	nodeType      string
	status        string
	uploadBytes   int64
	downloadBytes int64
	totalBytes    int64
}

type dashboardTrafficAggregate struct {
	trend         dashboardTrafficTrend
	byInbound     map[string]*dashboardInboundTraffic
	byNode        map[string]*dashboardNodeTraffic
	eligibleCount int
	trendPoints   map[int64]*dashboardTrafficPoint
	coverageSecs  float64
}

func parseDashboardRange(r *http.Request, now time.Time) (dashboardRangeSpec, time.Time, time.Time, error) {
	name := strings.TrimSpace(r.URL.Query().Get("range"))
	if name == "" {
		name = "today"
	}
	if name == "custom" {
		fromText := strings.TrimSpace(r.URL.Query().Get("from"))
		toText := strings.TrimSpace(r.URL.Query().Get("to"))
		if fromText == "" || toText == "" {
			return dashboardRangeSpec{}, time.Time{}, time.Time{}, errors.New("custom range requires from and to")
		}
		from, err := parseDashboardTime(fromText, false)
		if err != nil {
			return dashboardRangeSpec{}, time.Time{}, time.Time{}, errors.New("invalid custom range from")
		}
		to, err := parseDashboardTime(toText, true)
		if err != nil {
			return dashboardRangeSpec{}, time.Time{}, time.Time{}, errors.New("invalid custom range to")
		}
		if !to.After(from) || to.Sub(from) > 90*24*time.Hour {
			return dashboardRangeSpec{}, time.Time{}, time.Time{}, errors.New("custom range must be between 1 and 90 days")
		}
		return dashboardRangeSpec{name: "custom", duration: to.Sub(from), bucket: dashboardBucketForDuration(to.Sub(from))}, from, to, nil
	}
	spec, ok := dashboardRanges[name]
	if !ok {
		return dashboardRangeSpec{}, time.Time{}, time.Time{}, errors.New("range must be today, 7d, 30d, or custom")
	}
	var from time.Time
	if name == "today" {
		from = dashboardDayStart(now)
	} else {
		from = now.Add(-spec.duration)
	}
	return spec, from, now, nil
}

func dashboardDayStart(now time.Time) time.Time {
	localNow := now.In(dashboardBusinessLocation)
	return time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, dashboardBusinessLocation).UTC()
}

func dashboardMonthStart(now time.Time) time.Time {
	localNow := now.In(dashboardBusinessLocation)
	return time.Date(localNow.Year(), localNow.Month(), 1, 0, 0, 0, 0, dashboardBusinessLocation).UTC()
}

func parseDashboardTime(value string, endOfDay bool) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.UTC(), nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, dashboardBusinessLocation)
	if err != nil {
		return time.Time{}, err
	}
	if endOfDay {
		return parsed.AddDate(0, 0, 1).UTC(), nil
	}
	return parsed.UTC(), nil
}

func dashboardBucketForDuration(duration time.Duration) time.Duration {
	if duration <= 48*time.Hour {
		return time.Hour
	}
	return 24 * time.Hour
}

func dashboardBucketLabel(bucket time.Duration) string {
	if bucket == time.Hour {
		return "1h"
	}
	return "1d"
}

func (s *Server) beginDashboardFlight(key string) (*dashboardFlight, bool) {
	s.dashboardFlightMu.Lock()
	defer s.dashboardFlightMu.Unlock()
	if flight, ok := s.dashboardFlights[key]; ok {
		return flight, false
	}
	flight := &dashboardFlight{done: make(chan struct{})}
	s.dashboardFlights[key] = flight
	return flight, true
}

func (s *Server) finishDashboardFlight(key string, flight *dashboardFlight, data map[string]any) {
	s.dashboardFlightMu.Lock()
	if current, ok := s.dashboardFlights[key]; ok && current == flight {
		flight.data = data
		delete(s.dashboardFlights, key)
		close(flight.done)
	}
	s.dashboardFlightMu.Unlock()
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	spec, from, to, err := parseDashboardRange(r, now)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	cacheKey := r.URL.Query().Encode()
	s.dashboardCacheMu.Lock()
	for key, cached := range s.dashboardCache {
		if !now.Before(cached.expiresAt) {
			delete(s.dashboardCache, key)
		}
	}
	cached, cacheHit := s.dashboardCache[cacheKey]
	s.dashboardCacheMu.Unlock()
	// Cache misses perform several SQLite reads; do not serialize those reads
	// behind the cache mutex.
	if cacheHit && now.Before(cached.expiresAt) {
		writeSuccess(w, cached.data)
		return
	}
	flight, leader := s.beginDashboardFlight(cacheKey)
	if !leader {
		<-flight.done
		if flight.data != nil {
			writeSuccess(w, flight.data)
			return
		}
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard")
		return
	}
	var dashboardResult map[string]any
	defer func() { s.finishDashboardFlight(cacheKey, flight, dashboardResult) }()
	inbounds, err := s.dashboardInbounds()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard users")
		return
	}
	todaySpec := dashboardRanges["today"]
	todayStart := dashboardDayStart(now)
	windowSpecs := []dashboardTrafficWindowSpec{{from: from, to: to, spec: spec}}
	todayIndex := 0
	if spec.name == todaySpec.name && from.Equal(todayStart) && to.Equal(now) {
		// Reuse the selected range when the dashboard is showing today.
	} else {
		todayIndex = len(windowSpecs)
		windowSpecs = append(windowSpecs, dashboardTrafficWindowSpec{from: todayStart, to: now, spec: todaySpec})
	}
	monthStart := dashboardMonthStart(now)
	monthSpec := dashboardRangeSpec{name: "month", duration: now.Sub(monthStart), bucket: 24 * time.Hour}
	trafficWindows, err := s.dashboardTrafficWindows(windowSpecs, inbounds)
	if err != nil {
		s.logger.Error("dashboard traffic query", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard traffic")
		return
	}
	traffic := trafficWindows[0]
	todayTraffic := trafficWindows[todayIndex]
	monthTraffic, err := s.dashboardTrafficSummary(monthStart, now, monthSpec, inbounds)
	if err != nil {
		s.logger.Error("dashboard month traffic summary", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read month traffic")
		return
	}

	var nodes struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Relay   int `json:"relay"`
		Landing int `json:"landing"`
	}
	for _, item := range []struct {
		target *int
		query  string
	}{
		{&nodes.Total, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL`},
		{&nodes.Online, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL AND health_status = 'online' AND enabled = 1`},
		{&nodes.Relay, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL AND type = 'relay' AND enabled = 1`},
		{&nodes.Landing, `SELECT COUNT(*) FROM nodes WHERE deleted_at IS NULL AND type = 'landing' AND enabled = 1`},
	} {
		if err := s.db.QueryRow(item.query).Scan(item.target); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard nodes")
			return
		}
	}
	expiry, err := s.dashboardExpiry(now)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard users")
		return
	}
	events, err := s.dashboardEvents()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard events")
		return
	}
	finance, err := s.financeSummary(now.Format("2006-01"))
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read finance summary")
		return
	}
	userRanking, err := s.dashboardUserRanking(traffic, inbounds)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard users")
		return
	}
	dataAt := s.latestSuccessfulSyncAt()

	response := map[string]any{
		"generatedAt": time.Now().UTC().Format(time.RFC3339Nano),
		"dataAt":      dataAt,
		"range": map[string]any{
			"name": spec.name, "from": from.Format(time.RFC3339Nano), "to": to.Format(time.RFC3339Nano), "bucket": dashboardBucketLabel(spec.bucket),
		},
		"nodes": nodes,
		"users": map[string]int{"active": expiry.Active, "expiring": expiry.Expiring7d, "expired": expiry.Expired},
		"traffic": map[string]any{
			"todayBytes": todayTraffic.trend.Summary.TotalBytes, "monthBytes": monthTraffic.trend.Summary.TotalBytes,
			"uploadBytes": traffic.trend.Summary.UploadBytes, "downloadBytes": traffic.trend.Summary.DownloadBytes,
			"totalBytes": traffic.trend.Summary.TotalBytes, "businessScope": "relay-user-inbound",
		},
		"trafficTrend":       traffic.trend,
		"nodeTrafficRanking": s.dashboardNodeRanking(traffic),
		"userTrafficRanking": userRanking,
		"expiry":             expiry,
		"events":             events,
		"finance":            finance,
	}
	s.dashboardCacheMu.Lock()
	s.dashboardCache[cacheKey] = dashboardCacheEntry{data: response, expiresAt: now.Add(dashboardCacheTTL)}
	s.dashboardCacheMu.Unlock()
	dashboardResult = response
	writeSuccess(w, response)
}

type dashboardExpiry struct {
	Active         int `json:"active"`
	NormalAfter30d int `json:"normalAfter30d"`
	Expiring7d     int `json:"expiring7d"`
	Expiring30d    int `json:"expiring30d"`
	Expired        int `json:"expired"`
}

func (s *Server) dashboardExpiry(now time.Time) (dashboardExpiry, error) {
	var result dashboardExpiry
	err := s.db.QueryRow(`SELECT
COALESCE(SUM(CASE WHEN u.status = 'active' THEN 1 ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN u.expiry_time IS NULL OR datetime(u.expiry_time) > datetime(?, '+30 days') THEN 1 ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN u.status = 'expiring' THEN 1 ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN u.expiry_time IS NOT NULL AND datetime(u.expiry_time) > datetime(?, '+7 days') AND datetime(u.expiry_time) <= datetime(?, '+30 days') THEN 1 ELSE 0 END), 0),
COALESCE(SUM(CASE WHEN u.status = 'expired' THEN 1 ELSE 0 END), 0)
FROM users u
WHERE u.deleted_at IS NULL AND u.status <> 'disabled'
  AND EXISTS (SELECT 1 FROM user_inbounds ui JOIN inbounds i ON i.id = ui.inbound_id JOIN nodes n ON n.id = i.node_id
    WHERE ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL AND i.kind = 'user' AND i.deleted_at IS NULL
      AND n.type = 'relay' AND n.deleted_at IS NULL)`,
		now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)).Scan(&result.Active, &result.NormalAfter30d, &result.Expiring7d, &result.Expiring30d, &result.Expired)
	return result, err
}

type dashboardEventsResponse struct {
	PendingCount int                  `json:"pendingCount"`
	Items        []dashboardEventItem `json:"items"`
}

func (s *Server) dashboardEvents() (dashboardEventsResponse, error) {
	result := dashboardEventsResponse{Items: make([]dashboardEventItem, 0)}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events WHERE visibility = 'public' AND ` + pendingEventFilter("")).Scan(&result.PendingCount); err != nil {
		return result, err
	}
	rows, err := s.db.Query(`SELECT e.id, e.event_type, e.event_category, e.severity, COALESCE(e.title, ''),
COALESCE(e.node_id, ''), COALESCE(n.name, ''), e.message, e.created_at, e.acknowledged, e.requires_action, e.event_status
FROM node_events e LEFT JOIN nodes n ON n.id = e.node_id
WHERE e.visibility = 'public' AND ` + pendingEventFilter("e") + `
ORDER BY CASE e.severity WHEN 'error' THEN 0 WHEN 'warning' THEN 1 ELSE 2 END, e.created_at DESC LIMIT 5`)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var item dashboardEventItem
		var nodeID string
		var acknowledged, requiresAction int
		if err := rows.Scan(&item.ID, &item.Type, &item.Category, &item.Severity, &item.Title, &nodeID, &item.NodeName, &item.Message, &item.OccurredAt, &acknowledged, &requiresAction, &item.Status); err != nil {
			return result, err
		}
		item.NodeID = dashboardNullableString(nodeID)
		item.Acknowledged = acknowledged == 1
		item.RequiresAction = requiresAction == 1
		result.Items = append(result.Items, item)
	}
	return result, rows.Err()
}

func (s *Server) dashboardInbounds() (map[string]dashboardInbound, error) {
	rows, err := s.db.Query(`SELECT i.id, i.node_id, COALESCE(n.name, ''), COALESCE(n.health_status, 'unknown'),
COALESCE(i.user_id, ''), COALESCE(u.display_name, ''), COALESCE(i.tag, ''), COALESCE(u.status, 'unknown'),
COALESCE(u.expiry_time, ''), COALESCE((SELECT COUNT(*) FROM clients c WHERE c.inbound_id = i.id), 0),
COALESCE((SELECT MAX(NULLIF(c.last_online, '')) FROM clients c WHERE c.inbound_id = i.id), '')
FROM inbounds i JOIN nodes n ON n.id = i.node_id
JOIN users u ON u.id = i.user_id AND u.deleted_at IS NULL
WHERE i.kind = 'user' AND i.deleted_at IS NULL AND n.type = 'relay' AND n.deleted_at IS NULL
  AND EXISTS (SELECT 1 FROM user_inbounds ui WHERE ui.user_id = i.user_id AND ui.inbound_id = i.id AND ui.is_primary = 1 AND ui.active_to IS NULL)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]dashboardInbound)
	for rows.Next() {
		var item dashboardInbound
		var expiry, activity string
		if err := rows.Scan(&item.id, &item.nodeID, &item.nodeName, &item.nodeStatus, &item.userID, &item.userName, &item.inboundTag, &item.status, &expiry, &item.clientCount, &activity); err != nil {
			return nil, err
		}
		item.expiresAt = dashboardNullableString(expiry)
		item.lastActivity = dashboardNullableString(activity)
		result[item.id] = item
	}
	return result, rows.Err()
}

type dashboardTrafficWindowSpec struct {
	from time.Time
	to   time.Time
	spec dashboardRangeSpec
}

func newDashboardTrafficAggregate(from, to time.Time, spec dashboardRangeSpec, inbounds map[string]dashboardInbound) dashboardTrafficAggregate {
	result := dashboardTrafficAggregate{
		trend:     dashboardTrafficTrend{Range: spec.name, From: from.Format(time.RFC3339Nano), To: to.Format(time.RFC3339Nano), Bucket: dashboardBucketLabel(spec.bucket), Points: make([]dashboardTrafficPoint, 0)},
		byInbound: make(map[string]*dashboardInboundTraffic), byNode: make(map[string]*dashboardNodeTraffic),
		eligibleCount: len(inbounds), trendPoints: make(map[int64]*dashboardTrafficPoint),
	}
	for _, inbound := range inbounds {
		result.byInbound[inbound.id] = &dashboardInboundTraffic{lastActivity: inbound.lastActivity}
		result.byNode[inbound.nodeID] = &dashboardNodeTraffic{nodeID: inbound.nodeID, nodeName: inbound.nodeName, nodeType: "relay", status: inbound.nodeStatus}
	}
	return result
}

func (result *dashboardTrafficAggregate) trendPoint(key int64, start time.Time) *dashboardTrafficPoint {
	point := result.trendPoints[key]
	if point == nil {
		point = &dashboardTrafficPoint{Time: start.Format(time.RFC3339Nano)}
		result.trendPoints[key] = point
	}
	return point
}

func (s *Server) dashboardTraffic(from, to time.Time, spec dashboardRangeSpec, inbounds map[string]dashboardInbound) (dashboardTrafficAggregate, error) {
	results, err := s.dashboardTrafficWindows([]dashboardTrafficWindowSpec{{from: from, to: to, spec: spec}}, inbounds)
	if err != nil {
		return dashboardTrafficAggregate{}, err
	}
	return results[0], nil
}

// dashboardTrafficWindows reads a single broad snapshot window and computes
// every dashboard traffic range during the same ordered pass. The dashboard
// needs the selected range, today, and month totals; issuing one SQLite scan
// per range made the cost grow linearly with the number of ranges as history
// accumulated.
func (s *Server) dashboardTrafficWindows(windowSpecs []dashboardTrafficWindowSpec, inbounds map[string]dashboardInbound) ([]dashboardTrafficAggregate, error) {
	results := make([]dashboardTrafficAggregate, len(windowSpecs))
	if len(windowSpecs) == 0 {
		return results, nil
	}
	broadFrom := windowSpecs[0].from
	broadTo := windowSpecs[0].to
	for index, window := range windowSpecs {
		results[index] = newDashboardTrafficAggregate(window.from, window.to, window.spec, inbounds)
		if window.from.Before(broadFrom) {
			broadFrom = window.from
		}
		if window.to.After(broadTo) {
			broadTo = window.to
		}
	}
	if len(inbounds) == 0 {
		return results, nil
	}

	// Read the broad range plus one baseline sample per Inbound. The baseline
	// is shared by all requested ranges and lets the ordered pass below compute
	// deltas when a range starts between two stored snapshots.
	rows, err := s.db.Query(`WITH eligible AS MATERIALIZED (
  SELECT i.id AS inbound_id
  FROM inbounds i
  JOIN nodes n ON n.id = i.node_id
  JOIN users u ON u.id = i.user_id AND u.deleted_at IS NULL
  WHERE i.kind = 'user' AND i.deleted_at IS NULL
    AND n.type = 'relay' AND n.deleted_at IS NULL
    AND EXISTS (
      SELECT 1 FROM user_inbounds ui
      WHERE ui.user_id = i.user_id AND ui.inbound_id = i.id
        AND ui.is_primary = 1 AND ui.active_to IS NULL
    )
), range_rows AS MATERIALIZED (
  SELECT t.inbound_id, t.collected_at, t.up, t.down, t.reset_detected
  FROM traffic_snapshots t
  JOIN eligible e ON e.inbound_id = t.inbound_id
  WHERE t.collected_at >= ? AND t.collected_at <= ?
), baseline_times AS MATERIALIZED (
  SELECT e.inbound_id,
    (SELECT MAX(t.collected_at)
     FROM traffic_snapshots t
     WHERE t.inbound_id = e.inbound_id AND t.collected_at < ?) AS collected_at
  FROM eligible e
), baseline_rows AS MATERIALIZED (
  SELECT t.inbound_id, t.collected_at, t.up, t.down, t.reset_detected
  FROM traffic_snapshots t
  JOIN baseline_times b ON b.inbound_id = t.inbound_id AND b.collected_at = t.collected_at
)
SELECT inbound_id, collected_at, up, down, reset_detected
FROM (
  SELECT inbound_id, collected_at, up, down, reset_detected FROM baseline_rows
  UNION ALL
  SELECT inbound_id, collected_at, up, down, reset_detected FROM range_rows
)
ORDER BY inbound_id, collected_at`, broadFrom.Format(time.RFC3339Nano), broadTo.Format(time.RFC3339Nano), broadFrom.Format(time.RFC3339Nano))
	if err != nil {
		return results, err
	}
	defer rows.Close()

	var currentInbound string
	var previous dashboardSnapshot
	var hasPrevious bool
	for rows.Next() {
		var snapshot dashboardSnapshot
		var collected string
		var reset int
		if err := rows.Scan(&snapshot.inboundID, &collected, &snapshot.up, &snapshot.down, &reset); err != nil {
			return results, err
		}
		parsed, err := time.Parse(time.RFC3339Nano, collected)
		if err != nil {
			continue
		}
		snapshot.at = parsed.UTC()
		snapshot.reset = reset == 1
		if snapshot.inboundID != currentInbound {
			currentInbound = snapshot.inboundID
			hasPrevious = false
		}
		if !hasPrevious {
			previous = snapshot
			hasPrevious = true
			for index, window := range windowSpecs {
				if !snapshot.at.After(window.to) {
					results[index].trend.DataAt = dashboardNullableString(snapshot.at.Format(time.RFC3339Nano))
				}
			}
			continue
		}

		resetDetected := snapshot.reset || snapshot.up < previous.up || snapshot.down < previous.down
		uploadDelta := snapshot.up - previous.up
		downDelta := snapshot.down - previous.down
		if resetDetected {
			uploadDelta, downDelta = snapshot.up, snapshot.down
		}
		if uploadDelta < 0 {
			uploadDelta = 0
		}
		if downDelta < 0 {
			downDelta = 0
		}
		interval := snapshot.at.Sub(previous.at)
		if interval > 0 {
			if inbound, ok := inbounds[snapshot.inboundID]; ok {
				for index, window := range windowSpecs {
					if snapshot.at.Before(window.from) || snapshot.at.After(window.to) {
						continue
					}
					result := &results[index]
					result.trend.DataAt = dashboardNullableString(snapshot.at.Format(time.RFC3339Nano))
					inboundTraffic := result.byInbound[snapshot.inboundID]
					if inboundTraffic == nil {
						continue
					}
					inboundTraffic.uploadBytes += uploadDelta
					inboundTraffic.downloadBytes += downDelta
					inboundTraffic.totalBytes += uploadDelta + downDelta
					nodeTraffic := result.byNode[inbound.nodeID]
					if nodeTraffic != nil {
						nodeTraffic.uploadBytes += uploadDelta
						nodeTraffic.downloadBytes += downDelta
						nodeTraffic.totalBytes += uploadDelta + downDelta
					}
					bucketStart := time.Unix((snapshot.at.Unix()/int64(window.spec.bucket.Seconds()))*int64(window.spec.bucket.Seconds()), 0).UTC()
					bucketKey := bucketStart.Unix()
					point := result.trendPoint(bucketKey, bucketStart)
					point.UploadBytes += uploadDelta
					point.DownloadBytes += downDelta
					point.TotalBytes += uploadDelta + downDelta
					point.SampleCount++
					point.ResetDetected = point.ResetDetected || resetDetected
					point.HasGap = point.HasGap || interval > window.spec.bucket*3
					result.trend.Summary.UploadBytes += uploadDelta
					result.trend.Summary.DownloadBytes += downDelta
					result.trend.Summary.TotalBytes += uploadDelta + downDelta
					result.trend.Summary.SampleCount++
					result.coverageSecs += interval.Seconds()
				}
			}
		}
		previous = snapshot
		for index, window := range windowSpecs {
			if !snapshot.at.After(window.to) {
				results[index].trend.DataAt = dashboardNullableString(snapshot.at.Format(time.RFC3339Nano))
			}
		}
	}
	if err := rows.Err(); err != nil {
		return results, err
	}

	for index, window := range windowSpecs {
		result := &results[index]
		keys := make([]int64, 0, len(result.trendPoints))
		for key := range result.trendPoints {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		var cumulativeUpload, cumulativeDownload int64
		for _, key := range keys {
			point := *result.trendPoints[key]
			cumulativeUpload += point.UploadBytes
			cumulativeDownload += point.DownloadBytes
			point.UploadBytes = cumulativeUpload
			point.DownloadBytes = cumulativeDownload
			point.TotalBytes = cumulativeUpload + cumulativeDownload
			result.trend.Points = append(result.trend.Points, point)
		}
		result.trend.Summary.Coverage = 0
		if result.eligibleCount > 0 && window.spec.duration > 0 {
			result.trend.Summary.Coverage = result.coverageSecs / (window.spec.duration.Seconds() * float64(result.eligibleCount))
			if result.trend.Summary.Coverage > 1 {
				result.trend.Summary.Coverage = 1
			}
		}
	}
	return results, nil
}

// dashboardTrafficSummary calculates a range that is only displayed as a
// total (currently the month card) inside SQLite. Returning one aggregate row
// avoids transferring every minute-level snapshot to the Go process just to
// add the same deltas together.
func (s *Server) dashboardTrafficSummary(from, to time.Time, spec dashboardRangeSpec, inbounds map[string]dashboardInbound) (dashboardTrafficAggregate, error) {
	result := newDashboardTrafficAggregate(from, to, spec, inbounds)
	if len(inbounds) == 0 {
		return result, nil
	}
	var uploadBytes, downloadBytes, totalBytes, sampleCount int64
	var observedSeconds float64
	err := s.db.QueryRow(`WITH eligible AS MATERIALIZED (
  SELECT i.id AS inbound_id
  FROM inbounds i
  JOIN nodes n ON n.id = i.node_id
  JOIN users u ON u.id = i.user_id AND u.deleted_at IS NULL
  WHERE i.kind = 'user' AND i.deleted_at IS NULL
    AND n.type = 'relay' AND n.deleted_at IS NULL
    AND EXISTS (
      SELECT 1 FROM user_inbounds ui
      WHERE ui.user_id = i.user_id AND ui.inbound_id = i.id
        AND ui.is_primary = 1 AND ui.active_to IS NULL
    )
), range_rows AS MATERIALIZED (
  SELECT t.inbound_id, t.collected_at, t.up, t.down, t.reset_detected
  FROM traffic_snapshots t
  JOIN eligible e ON e.inbound_id = t.inbound_id
  WHERE t.collected_at >= ? AND t.collected_at <= ?
), baseline_times AS MATERIALIZED (
  SELECT e.inbound_id,
    (SELECT MAX(t.collected_at)
     FROM traffic_snapshots t
     WHERE t.inbound_id = e.inbound_id AND t.collected_at < ?) AS collected_at
  FROM eligible e
), baseline_rows AS MATERIALIZED (
  SELECT t.inbound_id, t.collected_at, t.up, t.down, t.reset_detected
  FROM traffic_snapshots t
  JOIN baseline_times b ON b.inbound_id = t.inbound_id AND b.collected_at = t.collected_at
), snapshots AS MATERIALIZED (
  SELECT inbound_id, collected_at, up, down, reset_detected FROM baseline_rows
  UNION ALL
  SELECT inbound_id, collected_at, up, down, reset_detected FROM range_rows
), ordered AS (
  SELECT inbound_id, collected_at, up, down, reset_detected,
    LAG(collected_at) OVER (PARTITION BY inbound_id ORDER BY collected_at) AS previous_at,
    LAG(up) OVER (PARTITION BY inbound_id ORDER BY collected_at) AS previous_up,
    LAG(down) OVER (PARTITION BY inbound_id ORDER BY collected_at) AS previous_down
  FROM snapshots
), deltas AS (
  SELECT collected_at,
    CASE WHEN reset_detected = 1 OR up < previous_up THEN up ELSE MAX(up - previous_up, 0) END AS upload_delta,
    CASE WHEN reset_detected = 1 OR down < previous_down THEN down ELSE MAX(down - previous_down, 0) END AS download_delta,
    (julianday(collected_at) - julianday(previous_at)) * 86400.0 AS interval_seconds
  FROM ordered
  WHERE previous_at IS NOT NULL
)
SELECT COALESCE(SUM(upload_delta), 0), COALESCE(SUM(download_delta), 0),
  COALESCE(SUM(upload_delta + download_delta), 0), COALESCE(SUM(interval_seconds), 0), COUNT(*)
FROM deltas
WHERE collected_at >= ? AND collected_at <= ? AND interval_seconds > 0`,
		from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano), from.Format(time.RFC3339Nano),
		from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano)).Scan(&uploadBytes, &downloadBytes, &totalBytes, &observedSeconds, &sampleCount)
	if err != nil {
		return result, err
	}
	result.trend.Summary = dashboardTrafficSummary{UploadBytes: uploadBytes, DownloadBytes: downloadBytes, TotalBytes: totalBytes, SampleCount: int(sampleCount)}
	if result.eligibleCount > 0 && spec.duration > 0 {
		result.trend.Summary.Coverage = observedSeconds / (spec.duration.Seconds() * float64(result.eligibleCount))
		if result.trend.Summary.Coverage > 1 {
			result.trend.Summary.Coverage = 1
		}
	}
	return result, nil
}

func dashboardNullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (s *Server) dashboardNodeRanking(traffic dashboardTrafficAggregate) []dashboardNodeTrafficItem {
	items := make([]dashboardNodeTrafficItem, 0, len(traffic.byNode))
	for _, node := range traffic.byNode {
		share := 0.0
		if traffic.trend.Summary.TotalBytes > 0 {
			share = float64(node.totalBytes) / float64(traffic.trend.Summary.TotalBytes)
		}
		items = append(items, dashboardNodeTrafficItem{NodeID: node.nodeID, NodeName: node.nodeName, NodeType: node.nodeType, Status: node.status, UploadBytes: node.uploadBytes, DownloadBytes: node.downloadBytes, TotalBytes: node.totalBytes, Share: share, DataAt: traffic.trend.DataAt})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].TotalBytes > items[j].TotalBytes })
	if len(items) > 10 {
		items = items[:10]
	}
	return items
}

func (s *Server) dashboardUserRanking(traffic dashboardTrafficAggregate, inbounds map[string]dashboardInbound) ([]dashboardUserTrafficItem, error) {
	items := make([]dashboardUserTrafficItem, 0, len(traffic.byInbound))
	for inboundID, trafficItem := range traffic.byInbound {
		inbound, ok := inbounds[inboundID]
		if !ok {
			continue
		}
		items = append(items, dashboardUserTrafficItem{UserID: inbound.userID, UserName: inbound.userName, NodeID: inbound.nodeID, NodeName: inbound.nodeName, InboundID: inbound.id, InboundTag: inbound.inboundTag, ClientCount: inbound.clientCount, Status: inbound.status, ExpiresAt: inbound.expiresAt, UploadBytes: trafficItem.uploadBytes, DownloadBytes: trafficItem.downloadBytes, TotalBytes: trafficItem.totalBytes, LastActivityAt: trafficItem.lastActivity, DataAt: traffic.trend.DataAt})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].TotalBytes > items[j].TotalBytes })
	if len(items) > 10 {
		items = items[:10]
	}
	return items, nil
}
