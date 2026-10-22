package httpapi

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// eventsStats aggregates the operator-facing event stream so the UI does not
// have to fetch a large page and compute counters client-side.
func (s *Server) eventsStats(w http.ResponseWriter, r *http.Request) {
	days := parseWindowDays(r.URL.Query().Get("window"), 7)
	now := time.Now().UTC()
	since := now.AddDate(0, 0, -days).Format(time.RFC3339Nano)
	todayStart := now.Format("2006-01-02") + "T00:00:00Z"

	var pending, unread, attention, processedToday, resolvedWindow, totalCount, autoRecoveredToday int
	var avgHandle float64

	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events WHERE visibility = 'public'`).Scan(&totalCount); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count events")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE visibility = 'public' AND event_status = 'resolved' AND resolved_at IS NOT NULL AND datetime(resolved_at) >= datetime(?)
  AND event_type IN ('xpanel_sync_recovered', 'node_recovered')`, todayStart).Scan(&autoRecoveredToday); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count auto-recovered events")
		return
	}

	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events e WHERE e.visibility = 'public' AND ` + pendingEventFilter("e")).Scan(&pending); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count pending events")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events e WHERE e.visibility = 'public' AND e.acknowledged = 0 AND e.event_status NOT IN ('resolved', 'dismissed')`).Scan(&unread); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count unread events")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events e WHERE e.visibility = 'public' AND ` + pendingEventFilter("e") + ` AND e.severity IN ('warning', 'error')`).Scan(&attention); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count attention events")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events WHERE visibility = 'public' AND event_status = 'resolved' AND resolved_at IS NOT NULL AND datetime(resolved_at) >= datetime(?)`, todayStart).Scan(&processedToday); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count resolved events")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(AVG((julianday(resolved_at) - julianday(created_at)) * 1440), 0)
FROM node_events
WHERE visibility = 'public' AND event_status = 'resolved' AND resolved_at IS NOT NULL AND datetime(created_at) >= datetime(?)`, since).Scan(&resolvedWindow, &avgHandle); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read event handling time")
		return
	}

	medianHandle := 0.0
	durationRows, err := s.db.Query(`SELECT (julianday(resolved_at) - julianday(created_at)) * 1440
FROM node_events
WHERE visibility = 'public' AND event_status = 'resolved' AND resolved_at IS NOT NULL AND datetime(created_at) >= datetime(?)
ORDER BY resolved_at`, since)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read event durations")
		return
	}
	durations := make([]float64, 0)
	for durationRows.Next() {
		var minutes float64
		if err := durationRows.Scan(&minutes); err != nil {
			_ = durationRows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode event durations")
			return
		}
		if minutes >= 0 {
			durations = append(durations, minutes)
		}
	}
	if err := durationRows.Close(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close event durations")
		return
	}
	if len(durations) > 0 {
		middle := len(durations) / 2
		if len(durations)%2 == 1 {
			medianHandle = durations[middle]
		} else {
			medianHandle = (durations[middle-1] + durations[middle]) / 2
		}
	}

	byCategory := make([]map[string]any, 0)
	rows, err := s.db.Query(`SELECT event_category, COUNT(*) FROM node_events
WHERE visibility = 'public' AND datetime(created_at) >= datetime(?)
GROUP BY event_category ORDER BY COUNT(*) DESC`, since)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not aggregate event categories")
		return
	}
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			_ = rows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode event categories")
			return
		}
		byCategory = append(byCategory, map[string]any{"category": category, "count": count})
	}
	if err := rows.Close(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close event categories")
		return
	}

	bySeverity := make([]map[string]any, 0)
	rows, err = s.db.Query(`SELECT severity, COUNT(*) FROM node_events
WHERE visibility = 'public' AND datetime(created_at) >= datetime(?)
GROUP BY severity ORDER BY COUNT(*) DESC`, since)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not aggregate event severities")
		return
	}
	for rows.Next() {
		var severity string
		var count int
		if err := rows.Scan(&severity, &count); err != nil {
			_ = rows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode event severities")
			return
		}
		bySeverity = append(bySeverity, map[string]any{"severity": severity, "count": count})
	}
	if err := rows.Close(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close event severities")
		return
	}

	byType := make([]map[string]any, 0)
	rows, err = s.db.Query(`SELECT event_type, COUNT(*) FROM node_events
WHERE visibility = 'public' AND datetime(created_at) >= datetime(?)
GROUP BY event_type ORDER BY COUNT(*) DESC`, since)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not aggregate event types")
		return
	}
	for rows.Next() {
		var eventType string
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			_ = rows.Close()
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode event types")
			return
		}
		byType = append(byType, map[string]any{"type": eventType, "count": count})
	}
	if err := rows.Close(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close event types")
		return
	}

	writeSuccess(w, map[string]any{
		"windowDays":      days,
		"pendingCount":    pending,
		"unreadCount":     unread,
		"attentionCount":  attention,
		"processedToday":  processedToday,
		"resolvedInWindow": resolvedWindow,
		"avgHandleMinutes": avgHandle,
		"byCategory":      byCategory,
		"bySeverity":      bySeverity,
		"byType":          byType,
		"totalCount":      totalCount,
		"autoRecoveredToday": autoRecoveredToday,
		"medianHandleMinutes": medianHandle,
	})
}

// costsTrend returns the monthly income/cost/renewal series in one response so
// the finance page does not need one summary request per month.
func (s *Server) costsTrend(w http.ResponseWriter, r *http.Request) {
	months := positiveInt(r.URL.Query().Get("months"), 12)
	if months > 24 {
		months = 24
	}
	now := time.Now()
	items := make([]map[string]any, 0, months)
	for index := months - 1; index >= 0; index-- {
		period := businessMonthPeriod(now.AddDate(0, -index, 0))
		summary, err := s.financeSummary(period)
		if err != nil {
			s.logger.Error("compute finance trend month", "period", period, "error", err)
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not compute finance trend")
			return
		}
		renewal := summary.Renewal
		items = append(items, map[string]any{
			"period":      period,
			"income":      summary.MonthIncome,
			"cost":        summary.MonthCost,
			"cashIncome":  summary.CashIncome,
			"renewalRate": renewal["renewalRate"],
			"dueUsers":    renewal["dueUsers"],
			"renewedUsers": renewal["renewedUsers"],
		})
	}
	writeSuccess(w, map[string]any{"items": items})
}

// renewalCandidates lists pending Agent-detected renewals across all users so
// the finance page does not have to read the event stream.
func (s *Server) renewalCandidates(w http.ResponseWriter, r *http.Request) {
	period := strings.TrimSpace(r.URL.Query().Get("period"))
	where := []string{"c.status = 'pending'"}
	args := make([]any, 0, 2)
	if period != "" {
		start, end, err := monthBounds(period)
		if err != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, "period must use YYYY-MM format")
			return
		}
		where = append(where, "datetime(c.detected_at) >= datetime(?)", "datetime(c.detected_at) < datetime(?)")
		args = append(args, start, end)
	}
	rows, err := s.db.Query(`SELECT c.id, c.user_id, COALESCE(u.display_name, ''), COALESCE(c.inbound_id, ''),
COALESCE(c.old_expiry_at, ''), COALESCE(c.new_expiry_at, ''), c.detected_at, c.suggested_cycle, c.suggested_amount, c.currency, COALESCE(c.notes, '')
FROM user_renewal_candidates c LEFT JOIN users u ON u.id = c.user_id
WHERE `+strings.Join(where, " AND ")+`
ORDER BY datetime(c.detected_at) DESC LIMIT 200`, args...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read renewal candidates")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, userID, userName, inboundID, oldExpiry, newExpiry, detectedAt, cycle, currency, notes string
		var amount float64
		if err := rows.Scan(&id, &userID, &userName, &inboundID, &oldExpiry, &newExpiry, &detectedAt, &cycle, &amount, &currency, &notes); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode renewal candidates")
			return
		}
		items = append(items, map[string]any{
			"id": id, "userId": userID, "userName": userName, "inboundId": nullableString(inboundID),
			"oldExpiryAt": oldExpiry, "newExpiryAt": newExpiry, "detectedAt": detectedAt,
			"billingCycle": cycle, "amount": amount, "currency": currency, "notes": nullableString(notes),
		})
	}
	if err := rows.Err(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read renewal candidates")
		return
	}
	writeSuccess(w, map[string]any{"items": items, "total": len(items)})
}

// parseWindowDays turns a window token such as 7d / 30d / today into days.
func parseWindowDays(raw string, fallback int) int {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback
	}
	if value == "today" {
		return 1
	}
	digits := strings.TrimSuffix(value, "d")
	parsed, err := strconv.Atoi(digits)
	if err != nil || parsed <= 0 || parsed > 365 {
		return fallback
	}
	return parsed
}

