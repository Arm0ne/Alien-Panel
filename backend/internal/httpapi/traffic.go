package httpapi

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"
)

type trafficRangeSpec struct {
	name     string
	duration time.Duration
	bucket   time.Duration
}

var trafficRanges = map[string]trafficRangeSpec{
	"1h": {name: "1h", duration: time.Hour, bucket: time.Minute},
	"6h": {name: "6h", duration: 6 * time.Hour, bucket: 5 * time.Minute},
	"1d": {name: "1d", duration: 24 * time.Hour, bucket: 15 * time.Minute},
	"7d": {name: "7d", duration: 7 * 24 * time.Hour, bucket: time.Hour},
}

func trafficBucketLabel(bucket time.Duration) string {
	switch bucket {
	case time.Minute:
		return "1m"
	case 5 * time.Minute:
		return "5m"
	case 15 * time.Minute:
		return "15m"
	case time.Hour:
		return "1h"
	default:
		return bucket.String()
	}
}

type trafficTrendPoint struct {
	Time          string  `json:"time"`
	UploadBytes   int64   `json:"uploadBytes"`
	DownloadBytes int64   `json:"downloadBytes"`
	UploadRate    float64 `json:"uploadRate"`
	DownloadRate  float64 `json:"downloadRate"`
	SampleCount   int     `json:"sampleCount"`
	ResetDetected bool    `json:"resetDetected"`
	HasGap        bool    `json:"hasGap"`
}

type trafficTrendSummary struct {
	UploadBytes         int64   `json:"uploadBytes"`
	DownloadBytes       int64   `json:"downloadBytes"`
	TotalBytes          int64   `json:"totalBytes"`
	AverageUploadRate   float64 `json:"averageUploadRate"`
	AverageDownloadRate float64 `json:"averageDownloadRate"`
	PeakUploadRate      float64 `json:"peakUploadRate"`
	PeakDownloadRate    float64 `json:"peakDownloadRate"`
	SampleCount         int     `json:"sampleCount"`
	Coverage            float64 `json:"coverage"`
}

type trafficTrendResponse struct {
	Range   string              `json:"range"`
	From    string              `json:"from"`
	To      string              `json:"to"`
	Bucket  string              `json:"bucket"`
	DataAt  *string             `json:"dataAt"`
	Points  []trafficTrendPoint `json:"points"`
	Summary trafficTrendSummary `json:"summary"`
}

type trafficSnapshotRow struct {
	collectedAt time.Time
	up          int64
	down        int64
	reset       bool
}

type trafficBucket struct {
	start         time.Time
	uploadBytes   int64
	downloadBytes int64
	observedSecs  float64
	sampleCount   int
	resetDetected bool
	hasGap        bool
}

func (s *Server) userTraffic(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}

	rangeName := strings.TrimSpace(r.URL.Query().Get("range"))
	if rangeName == "" {
		rangeName = "1h"
	}
	spec, ok := trafficRanges[rangeName]
	if !ok {
		writeFailure(w, http.StatusBadRequest, validationCode, "range must be one of 1h, 6h, 1d, 7d")
		return
	}

	var inboundID string
	err := s.db.QueryRow(`SELECT COALESCE(i.id, '')
FROM users u
LEFT JOIN user_inbounds ui ON ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL
LEFT JOIN inbounds i ON i.id = ui.inbound_id AND i.kind = 'user' AND i.deleted_at IS NULL
 AND i.node_id IN (SELECT id FROM nodes WHERE type = 'relay' AND deleted_at IS NULL)
WHERE u.id = ?`, id).Scan(&inboundID)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
		return
	}
	if err != nil {
		s.logger.Error("read user traffic inbound", "user_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user traffic")
		return
	}

	to := time.Now().UTC()
	from := to.Add(-spec.duration)
	response, err := s.buildTrafficTrend(inboundID, from, to, spec)
	if err != nil {
		s.logger.Error("read user traffic trend", "user_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user traffic")
		return
	}
	writeSuccess(w, response)
}

func (s *Server) buildTrafficTrend(inboundID string, from, to time.Time, spec trafficRangeSpec) (trafficTrendResponse, error) {
	response := trafficTrendResponse{
		Range:  spec.name,
		From:   from.Format(time.RFC3339Nano),
		To:     to.Format(time.RFC3339Nano),
		Bucket: trafficBucketLabel(spec.bucket),
		Points: make([]trafficTrendPoint, 0),
	}
	if inboundID == "" {
		return response, nil
	}

	var latestAt string
	if err := s.db.QueryRow(`SELECT collected_at FROM traffic_snapshots WHERE inbound_id = ? AND collected_at <= ? ORDER BY collected_at DESC LIMIT 1`, inboundID, to.Format(time.RFC3339Nano)).Scan(&latestAt); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return response, err
	} else if latestAt != "" {
		parsed, parseErr := time.Parse(time.RFC3339Nano, latestAt)
		if parseErr == nil {
			value := parsed.UTC().Format(time.RFC3339Nano)
			response.DataAt = &value
		}
	}

	rows, err := s.db.Query(`SELECT collected_at, up, down, reset_detected
FROM traffic_snapshots
WHERE inbound_id = ? AND collected_at >= ? AND collected_at <= ?
ORDER BY collected_at ASC`, inboundID, from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano))
	if err != nil {
		return response, err
	}
	defer rows.Close()

	snapshots := make([]trafficSnapshotRow, 0)
	for rows.Next() {
		var collectedAt string
		var up, down int64
		var resetDetected int
		if err := rows.Scan(&collectedAt, &up, &down, &resetDetected); err != nil {
			return response, err
		}
		parsed, err := time.Parse(time.RFC3339Nano, collectedAt)
		if err != nil {
			continue
		}
		snapshots = append(snapshots, trafficSnapshotRow{collectedAt: parsed.UTC(), up: up, down: down, reset: resetDetected == 1})
	}
	if err := rows.Err(); err != nil {
		return response, err
	}
	if len(snapshots) < 2 {
		return response, nil
	}

	buckets := make(map[int64]*trafficBucket)
	var previous trafficSnapshotRow
	var hasPrevious bool
	var uploadTotal, downloadTotal, observedSeconds float64
	var sampleCount int
	var peakUpload, peakDownload float64
	maxGap := spec.bucket * 3
	for _, current := range snapshots {
		if !hasPrevious {
			// The first sample in the selected window is a baseline. We do not
			// attribute its cumulative value to the current period.
			previous = current
			hasPrevious = true
			continue
		}
		interval := current.collectedAt.Sub(previous.collectedAt)
		if interval <= 0 {
			previous = current
			continue
		}
		resetDetected := current.reset || current.up < previous.up || current.down < previous.down
		uploadDelta := current.up - previous.up
		downloadDelta := current.down - previous.down
		if resetDetected {
			uploadDelta = current.up
			downloadDelta = current.down
		}
		if uploadDelta < 0 {
			uploadDelta = 0
		}
		if downloadDelta < 0 {
			downloadDelta = 0
		}
		intervalSeconds := interval.Seconds()
		uploadRate := float64(uploadDelta) / intervalSeconds
		downloadRate := float64(downloadDelta) / intervalSeconds
		hasGap := interval > maxGap
		bucketStart := time.Unix((current.collectedAt.Unix()/int64(spec.bucket.Seconds()))*int64(spec.bucket.Seconds()), 0).UTC()
		bucketKey := bucketStart.Unix()
		bucket := buckets[bucketKey]
		if bucket == nil {
			bucket = &trafficBucket{start: bucketStart}
			buckets[bucketKey] = bucket
		}
		bucket.uploadBytes += uploadDelta
		bucket.downloadBytes += downloadDelta
		bucket.observedSecs += intervalSeconds
		bucket.sampleCount++
		bucket.resetDetected = bucket.resetDetected || resetDetected
		bucket.hasGap = bucket.hasGap || hasGap
		uploadTotal += float64(uploadDelta)
		downloadTotal += float64(downloadDelta)
		observedSeconds += intervalSeconds
		sampleCount++
		if uploadRate > peakUpload {
			peakUpload = uploadRate
		}
		if downloadRate > peakDownload {
			peakDownload = downloadRate
		}
		previous = current
	}

	keys := make([]int64, 0, len(buckets))
	for key := range buckets {
		keys = append(keys, key)
	}
	for index := 0; index < len(keys); index++ {
		for next := index + 1; next < len(keys); next++ {
			if keys[next] < keys[index] {
				keys[index], keys[next] = keys[next], keys[index]
			}
		}
	}
	for _, key := range keys {
		bucket := buckets[key]
		seconds := bucket.observedSecs
		if seconds <= 0 {
			seconds = spec.bucket.Seconds()
		}
		response.Points = append(response.Points, trafficTrendPoint{
			Time:          bucket.start.Format(time.RFC3339Nano),
			UploadBytes:   bucket.uploadBytes,
			DownloadBytes: bucket.downloadBytes,
			UploadRate:    float64(bucket.uploadBytes) / seconds,
			DownloadRate:  float64(bucket.downloadBytes) / seconds,
			SampleCount:   bucket.sampleCount,
			ResetDetected: bucket.resetDetected,
			HasGap:        bucket.hasGap,
		})
	}

	coverage := observedSeconds / spec.duration.Seconds()
	if coverage > 1 {
		coverage = 1
	}
	averageUploadRate := 0.0
	averageDownloadRate := 0.0
	if observedSeconds > 0 {
		averageUploadRate = uploadTotal / observedSeconds
		averageDownloadRate = downloadTotal / observedSeconds
	}
	response.Summary = trafficTrendSummary{
		UploadBytes:         int64(uploadTotal),
		DownloadBytes:       int64(downloadTotal),
		TotalBytes:          int64(uploadTotal + downloadTotal),
		AverageUploadRate:   averageUploadRate,
		AverageDownloadRate: averageDownloadRate,
		PeakUploadRate:      peakUpload,
		PeakDownloadRate:    peakDownload,
		SampleCount:         sampleCount,
		Coverage:            coverage,
	}
	return response, nil
}
