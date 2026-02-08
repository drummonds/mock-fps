package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                 string
	LifecycleStepDelayMs int
	WebhookWorkers       int
	WebhookBufferSize    int
}

func Load() Config {
	return Config{
		Port:                 envOrDefault("PORT", "8080"),
		LifecycleStepDelayMs: envIntOrDefault("LIFECYCLE_STEP_DELAY_MS", 500),
		WebhookWorkers:       envIntOrDefault("WEBHOOK_WORKERS", 4),
		WebhookBufferSize:    envIntOrDefault("WEBHOOK_BUFFER_SIZE", 1000),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
