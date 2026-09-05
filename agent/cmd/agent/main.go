package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"xpanel-central/agent/internal/buildinfo"
	"xpanel-central/agent/internal/central"
	"xpanel-central/agent/internal/collector"
	"xpanel-central/agent/internal/config"
	"xpanel-central/agent/internal/runner"
	"xpanel-central/agent/internal/xpanel"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	defaultPath := os.Getenv("XPANEL_AGENT_CONFIG")
	if defaultPath == "" {
		defaultPath = "./agent.yaml"
	}
	configPath := flag.String("config", defaultPath, "path to the agent YAML configuration")
	showVersion := flag.Bool("version", false, "print the embedded Agent version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println(buildinfo.DisplayVersion())
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("load agent configuration", "error", err)
		os.Exit(1)
	}
	xpanelClient, err := xpanel.NewClient(cfg.XPanelURL, cfg.XPanelBasePath, cfg.XPanelUsername, cfg.XPanelPassword, cfg.HTTPTimeout)
	if err != nil {
		logger.Error("create xpanel client", "error", err)
		os.Exit(1)
	}
	centralClient, err := central.NewClient(cfg.CentralURL, cfg.CentralToken, cfg.HTTPTimeout)
	if err != nil {
		logger.Error("create central client", "error", err)
		os.Exit(1)
	}
	snapshotCollector, err := collector.New(xpanelClient, cfg.NodeKey, buildinfo.DisplayVersion())
	if err != nil {
		logger.Error("create xpanel collector", "error", err)
		os.Exit(1)
	}
	agentRunner, err := runner.New(snapshotCollector, centralClient, cfg.SyncInterval, logger)
	if err != nil {
		logger.Error("create agent runner", "error", err)
		os.Exit(1)
	}
	logger.Info("agent started", "node_key", cfg.NodeKey, "node_type", cfg.NodeType, "agent_version", buildinfo.DisplayVersion(), "agent_build", buildinfo.BuildIdentity())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- agentRunner.Run(ctx)
	}()
	select {
	case <-stop:
		cancel()
		<-done
	case err := <-done:
		cancel()
		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("agent stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}
	logger.Info("agent stopped")
}
