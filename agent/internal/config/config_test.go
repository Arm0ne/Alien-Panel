package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadYAMLAndEnvironmentOverrides(t *testing.T) {
	path := filepath.Join(t.TempDir(), "agent.yaml")
	content := []byte(`node_key: relay-yaml
node_name: Relay YAML
node_type: relay
central_url: https://central.example.test
central_token: yaml-token
xpanel_url: http://127.0.0.1:2053
xpanel_base_path: /panel/
xpanel_username: yaml-user
xpanel_password: yaml-password
sync_interval: 2m
http_timeout: 4s
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XPANEL_NODE_NAME", "Relay Env")
	t.Setenv("XPANEL_SYNC_INTERVAL", "7s")
	t.Setenv("XPANEL_HTTP_TIMEOUT", "invalid")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.NodeName != "Relay Env" {
		t.Fatalf("NodeName = %q, want environment override", cfg.NodeName)
	}
	if cfg.SyncInterval != 7*time.Second {
		t.Fatalf("SyncInterval = %s, want 7s", cfg.SyncInterval)
	}
	if cfg.HTTPTimeout != 4*time.Second {
		t.Fatalf("HTTPTimeout = %s, want YAML value after invalid override", cfg.HTTPTimeout)
	}
	if cfg.XPanelBasePath != "/panel/" {
		t.Fatalf("XPanelBasePath = %q", cfg.XPanelBasePath)
	}
}

func TestLoadMissingFileStillValidates(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("Load() error = nil, want required configuration error")
	}
}

func TestValidateRejectsNonPositiveDurations(t *testing.T) {
	cfg := Config{
		NodeKey:        "node-1",
		CentralURL:     "https://central.example.test",
		XPanelURL:      "http://127.0.0.1:2053",
		XPanelUsername: "admin",
		XPanelPassword: "secret",
		SyncInterval:   0,
		HTTPTimeout:    time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want duration validation error")
	}
}

func TestParseBool(t *testing.T) {
	if !ParseBool("true", false) {
		t.Fatal("ParseBool(true) = false")
	}
	if ParseBool("not-bool", true) != true {
		t.Fatal("ParseBool(invalid) did not use fallback")
	}
}
