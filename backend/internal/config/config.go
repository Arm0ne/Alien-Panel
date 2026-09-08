package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListenAddress          string
	DatabasePath           string
	AdminUsername          string
	AdminPassword          string
	AgentRegistrationToken string
	SessionTTL             time.Duration
	MaintenanceInterval    time.Duration
	CorsOrigins            []string
	AllowInsecureCookie    bool
}

func Load() (Config, error) {
	cfg := Config{
		ListenAddress:          envOrDefault("XPANEL_LISTEN", ":8090"),
		DatabasePath:           envOrDefault("XPANEL_DATABASE", "./data/panel.db"),
		AdminUsername:          envOrDefault("XPANEL_ADMIN_USER", "admin"),
		AdminPassword:          os.Getenv("XPANEL_ADMIN_PASSWORD"),
		AgentRegistrationToken: os.Getenv("XPANEL_AGENT_REGISTRATION_TOKEN"),
		SessionTTL:             durationOrDefault("XPANEL_SESSION_TTL", 24*time.Hour),
		MaintenanceInterval:    durationOrDefault("XPANEL_MAINTENANCE_INTERVAL", time.Minute),
		CorsOrigins:            splitCSV(envOrDefault("XPANEL_CORS_ORIGINS", "http://localhost:9527,http://127.0.0.1:9527")),
		AllowInsecureCookie:    boolOrDefault("XPANEL_ALLOW_INSECURE_COOKIE", true),
	}

	if cfg.AdminPassword == "" {
		return Config{}, errors.New("XPANEL_ADMIN_PASSWORD must be set before starting the server")
	}
	if strings.TrimSpace(cfg.AdminUsername) == "" {
		return Config{}, errors.New("XPANEL_ADMIN_USER must not be empty")
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func durationOrDefault(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return fallback
	}
	return duration
}

func boolOrDefault(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
