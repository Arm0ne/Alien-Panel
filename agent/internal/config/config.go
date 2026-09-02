package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	NodeKey        string        `yaml:"node_key"`
	NodeName       string        `yaml:"node_name"`
	NodeType       string        `yaml:"node_type"`
	CentralURL     string        `yaml:"central_url"`
	CentralToken   string        `yaml:"central_token"`
	XPanelURL      string        `yaml:"xpanel_url"`
	XPanelBasePath string        `yaml:"xpanel_base_path"`
	XPanelUsername string        `yaml:"xpanel_username"`
	XPanelPassword string        `yaml:"xpanel_password"`
	SyncInterval   time.Duration `yaml:"sync_interval"`
	HTTPTimeout    time.Duration `yaml:"http_timeout"`
}

func Load(path string) (Config, error) {
	cfg := Config{SyncInterval: 60 * time.Second, HTTPTimeout: 15 * time.Second, XPanelBasePath: "/"}
	if path != "" {
		content, err := os.ReadFile(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
		if len(content) > 0 {
			if err := yaml.Unmarshal(content, &cfg); err != nil {
				return Config{}, fmt.Errorf("parse config: %w", err)
			}
		}
	}
	applyEnv(&cfg)
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg Config) Validate() error {
	missing := make([]string, 0, 5)
	if strings.TrimSpace(cfg.NodeKey) == "" {
		missing = append(missing, "node_key")
	}
	if strings.TrimSpace(cfg.CentralURL) == "" {
		missing = append(missing, "central_url")
	}
	if strings.TrimSpace(cfg.XPanelURL) == "" {
		missing = append(missing, "xpanel_url")
	}
	if strings.TrimSpace(cfg.XPanelUsername) == "" {
		missing = append(missing, "xpanel_username")
	}
	if cfg.XPanelPassword == "" {
		missing = append(missing, "xpanel_password")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required agent configuration: %s", strings.Join(missing, ", "))
	}
	if cfg.SyncInterval <= 0 || cfg.HTTPTimeout <= 0 {
		return errors.New("sync_interval and http_timeout must be positive")
	}
	return nil
}

func applyEnv(cfg *Config) {
	setString(&cfg.NodeKey, "XPANEL_NODE_KEY")
	setString(&cfg.NodeName, "XPANEL_NODE_NAME")
	setString(&cfg.NodeType, "XPANEL_NODE_TYPE")
	setString(&cfg.CentralURL, "XPANEL_CENTRAL_URL")
	setString(&cfg.CentralToken, "XPANEL_CENTRAL_TOKEN")
	setString(&cfg.XPanelURL, "XPANEL_XPANEL_URL")
	setString(&cfg.XPanelBasePath, "XPANEL_XPANEL_BASE_PATH")
	setString(&cfg.XPanelUsername, "XPANEL_XPANEL_USERNAME")
	setString(&cfg.XPanelPassword, "XPANEL_XPANEL_PASSWORD")
	if value := strings.TrimSpace(os.Getenv("XPANEL_SYNC_INTERVAL")); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			cfg.SyncInterval = parsed
		}
	}
	if value := strings.TrimSpace(os.Getenv("XPANEL_HTTP_TIMEOUT")); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			cfg.HTTPTimeout = parsed
		}
	}
}

func setString(target *string, key string) {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		*target = value
	}
}

func ParseBool(value string, fallback bool) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
