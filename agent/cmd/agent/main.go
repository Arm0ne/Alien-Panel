package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"xpanel-central/agent/internal/central"
	"xpanel-central/agent/internal/config"
	"xpanel-central/agent/internal/xpanel"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	defaultPath := os.Getenv("XPANEL_AGENT_CONFIG")
	if defaultPath == "" {
		defaultPath = "./agent.yaml"
	}
	configPath := flag.String("config", defaultPath, "path to the agent YAML configuration")
	flag.Parse()

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
	_ = xpanelClient
	_ = centralClient
	logger.Info("agent communication layer ready", "node_key", cfg.NodeKey, "node_type", cfg.NodeType)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)
	<-stop
	logger.Info("agent stopped")
}
