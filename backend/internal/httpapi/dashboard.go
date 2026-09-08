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
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		from = now.Add(-spec.duration)
	}
	return spec, from, now, nil
}

func parseDashboardTime(value string, endOfDay bool) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.UTC(), nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	if endOfDay {
		return parsed.Add(24 * time.Hour), nil
	}
	return parsed, nil
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

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	s.refreshOperationalStatuses(now)
	spec, from, to, err := parseDashboardRange(r, now)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	traffic, err := s.dashboardTraffic(from, to, spec)
	if err != nil {
		s.logger.Error("dashboard traffic query", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read dashboard traffic")
		return
	}
	todaySpec := dashboardRanges["today"]
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayTraffic, err := s.dashboardTraffic(todayStart, now, todaySpec)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read today traffic")
		return
	}
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthSpec := dashboardRangeSpec{name: "month", duration: now.Sub(monthStart), bucket: 24 * time.Hour}
	monthTraffic, err := s.dashboardTraffic(monthStart, now, monthSpec)
	if err != nil {
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
	userRanking, err := s.dashboardUserRanking(traffic)
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
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events WHERE visibility = 'public' AND requires_action = 1 AND event_status NOT IN ('resolved', 'dismissed')`).Scan(&result.PendingCount); err != nil {
		return result, err
	}
	rows, err := s.db.Query(`SELECT e.id, e.event_type, e.event_category, e.severity, COALESCE(e.title, ''),
COALESCE(e.node_id, ''), COALESCE(n.name, ''), e.message, e.created_at, e.acknowledged, e.requires_action, e.event_status
FROM node_events e LEFT JOIN nodes n ON n.id = e.node_id
WHERE e.visibility = 'public' AND e.requires_action = 1 AND e.event_status NOT IN ('resolved', 'dismissed')
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
COALESCE((SELECT MAX(c.last_seen_at) FROM clients c WHERE c.inbound_id = i.id), COALESCE(i.last_seen_at, ''))
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

func (s *Server) dashboardTraffic(from, to time.Time, spec dashboardRangeSpec) (dashboardTrafficAggregate, error) {
	result := dashboardTrafficAggregate{
		trend:     dashboardTrafficTrend{Range: spec.name, From: from.Format(time.RFC3339Nano), To: to.Format(time.RFC3339Nano), Bucket: dashboardBucketLabel(spec.bucket), Points: make([]dashboardTrafficPoint, 0)},
		byInbound: make(map[string]*dashboardInboundTraffic), byNode: make(map[string]*dashboardNodeTraffic),
	}
	inbounds, err := s.dashboardInbounds()
	if err != nil {
		return result, err
	}
	result.eligibleCount = len(inbounds)
	for _, inbound := range inbounds {
		result.byInbound[inbound.id] = &dashboardInboundTraffic{lastActivity: inbound.lastActivity}
		result.byNode[inbound.nodeID] = &dashboardNodeTraffic{nodeID: inbound.nodeID, nodeName: inbound.nodeName, nodeType: "relay", status: inbound.nodeStatus}
	}
	if len(inbounds) == 0 {
		return result, nil
	}
	rows, err := s.db.Query(`SELECT t.inbound_id, t.collected_at, t.up, t.down, t.reset_detected
FROM traffic_snapshots t JOIN inbounds i ON i.id = t.inbound_id JOIN nodes n ON n.id = i.node_id
WHERE i.kind = 'user' AND i.deleted_at IS NULL AND n.type = 'relay' AND n.deleted_at IS NULL
  AND EXISTS (SELECT 1 FROM user_inbounds ui WHERE ui.user_id = i.user_id AND ui.inbound_id = i.id AND ui.is_primary = 1 AND ui.active_to IS NULL)
  AND ((t.collected_at >= ? AND t.collected_at <= ?) OR t.collected_at = (
    SELECT MAX(p.collected_at) FROM traffic_snapshots p WHERE p.inbound_id = t.inbound_id AND p.collected_at < ?
  ))
ORDER BY t.inbound_id, t.collected_at`, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano), from.Format(time.RFC3339Nano))
	if err != nil {
		return result, err
	}
	defer rows.Close()
	buckets := make(map[int64]*dashboardTrafficPoint)
	var currentInbound string
	var previous dashboardSnapshot
	var hasPrevious bool
	var uploadTotal, downloadTotal, observedSeconds float64
	var sampleCount int
	latestAt := time.Time{}
	for rows.Next() {
		var snapshot dashboardSnapshot
		var collected string
		var reset int
		if err := rows.Scan(&snapshot.inboundID, &collected, &snapshot.up, &snapshot.down, &reset); err != nil {
			return result, err
		}
		parsed, err := time.Parse(time.RFC3339Nano, collected)
		if err != nil {
			continue
		}
		snapshot.at = parsed.UTC()
		snapshot.reset = reset == 1
		if snapshot.at.After(to) {
			continue
		}
		if snapshot.inboundID != currentInbound {
			currentInbound = snapshot.inboundID
			hasPrevious = false
		}
		if snapshot.at.After(latestAt) {
			latestAt = snapshot.at
		}
		if !hasPrevious {
			previous = snapshot
			hasPrevious = true
			continue
		}
		if snapshot.at.Before(from) {
			previous = snapshot
			continue
		}
		interval := snapshot.at.Sub(previous.at)
		if interval <= 0 {
			previous = snapshot
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
		inbound := inbounds[snapshot.inboundID]
		inboundTraffic := result.byInbound[snapshot.inboundID]
		inboundTraffic.uploadBytes += uploadDelta
		inboundTraffic.downloadBytes += downDelta
		inboundTraffic.totalBytes += uploadDelta + downDelta
		activity := snapshot.at.Format(time.RFC3339Nano)
		inboundTraffic.lastActivity = &activity
		nodeTraffic := result.byNode[inbound.nodeID]
		nodeTraffic.uploadBytes += uploadDelta
		nodeTraffic.downloadBytes += downDelta
		nodeTraffic.totalBytes += uploadDelta + downDelta
		bucketStart := time.Unix((snapshot.at.Unix()/int64(spec.bucket.Seconds()))*int64(spec.bucket.Seconds()), 0).UTC()
		bucketKey := bucketStart.Unix()
		point := buckets[bucketKey]
		if point == nil {
			point = &dashboardTrafficPoint{Time: bucketStart.Format(time.RFC3339Nano)}
			buckets[bucketKey] = point
		}
		point.UploadBytes += uploadDelta
		point.DownloadBytes += downDelta
		point.TotalBytes += uploadDelta + downDelta
		point.SampleCount++
		point.ResetDetected = point.ResetDetected || resetDetected
		point.HasGap = point.HasGap || interval > spec.bucket*3
		uploadTotal += float64(uploadDelta)
		downloadTotal += float64(downDelta)
		observedSeconds += interval.Seconds()
		sampleCount++
		previous = snapshot
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	if !latestAt.IsZero() {
		value := latestAt.Format(time.RFC3339Nano)
		result.trend.DataAt = &value
	}
	keys := make([]int64, 0, len(buckets))
	for key := range buckets {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	var cumulativeUpload, cumulativeDownload int64
	for _, key := range keys {
		point := *buckets[key]
		cumulativeUpload += point.UploadBytes
		cumulativeDownload += point.DownloadBytes
		point.UploadBytes = cumulativeUpload
		point.DownloadBytes = cumulativeDownload
		point.TotalBytes = cumulativeUpload + cumulativeDownload
		result.trend.Points = append(result.trend.Points, point)
	}
	coverage := 0.0
	if result.eligibleCount > 0 && spec.duration > 0 {
		coverage = observedSeconds / (spec.duration.Seconds() * float64(result.eligibleCount))
		if coverage > 1 {
			coverage = 1
		}
	}
	result.trend.Summary = dashboardTrafficSummary{UploadBytes: int64(uploadTotal), DownloadBytes: int64(downloadTotal), TotalBytes: int64(uploadTotal + downloadTotal), SampleCount: sampleCount, Coverage: coverage}
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

func (s *Server) dashboardUserRanking(traffic dashboardTrafficAggregate) ([]dashboardUserTrafficItem, error) {
	items := make([]dashboardUserTrafficItem, 0, len(traffic.byInbound))
	inbounds, err := s.dashboardInbounds()
	if err != nil {
		return items, err
	}
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
