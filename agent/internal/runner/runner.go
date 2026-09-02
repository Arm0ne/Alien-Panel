package runner

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"xpanel-central/agent/internal/central"
	"xpanel-central/agent/internal/collector"
)

type Runner struct {
	collector *collector.Collector
	central   *central.Client
	interval  time.Duration
	logger    *slog.Logger
}

func New(snapshotCollector *collector.Collector, centralClient *central.Client, interval time.Duration, logger *slog.Logger) (*Runner, error) {
	if snapshotCollector == nil {
		return nil, errors.New("collector is required")
	}
	if centralClient == nil {
		return nil, errors.New("central client is required")
	}
	if interval <= 0 {
		return nil, errors.New("sync interval must be positive")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{collector: snapshotCollector, central: centralClient, interval: interval, logger: logger}, nil
}

func (runner *Runner) Run(ctx context.Context) error {
	runner.syncOnce(ctx)
	ticker := time.NewTicker(runner.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			runner.syncOnce(ctx)
		}
	}
}

func (runner *Runner) syncOnce(ctx context.Context) {
	snapshot, err := runner.collector.Collect(ctx)
	if err != nil {
		runner.logger.Warn("collect xpanel snapshot", "error", err)
		return
	}
	if err := runner.central.Heartbeat(ctx, snapshot.Heartbeat()); err != nil {
		runner.logger.Warn("send node heartbeat", "sync_id", snapshot.SyncID, "error", err)
	}
	var result struct {
		SyncID string `json:"sync_id"`
		Status string `json:"status"`
	}
	if err := runner.central.Sync(ctx, snapshot, &result); err != nil {
		runner.logger.Warn("send xpanel snapshot", "sync_id", snapshot.SyncID, "error", err)
		return
	}
	runner.logger.Info("xpanel snapshot synchronized", "sync_id", snapshot.SyncID, "inbound_count", len(snapshot.Inbounds), "status", result.Status)
}
