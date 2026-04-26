package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultPort              = "3030"
	defaultAssetsDir         = "./assets"
	defaultFilterTimeout     = 10 * time.Second
	defaultFilterMaxBytes    = 10 << 20 // 10 MiB
	defaultEnableDebugRoutes = false
	defaultCORSAllowOrigins  = ""
)

type Config struct {
	Port              string
	AssetsDir         string
	EnableDebugRoutes bool
	FilterTimeout     time.Duration
	FilterMaxBytes    int64
	CORSAllowOrigins  string
}

func Load() Config {
	return Config{
		Port:              getenv("PORT", defaultPort),
		AssetsDir:         getenv("ASSETS_DIR", defaultAssetsDir),
		EnableDebugRoutes: getBool("ENABLE_DEBUG_ROUTES", defaultEnableDebugRoutes),
		FilterTimeout:     getDuration("FILTER_TIMEOUT", defaultFilterTimeout),
		FilterMaxBytes:    getInt64("FILTER_MAX_BYTES", defaultFilterMaxBytes),
		CORSAllowOrigins:  getenv("CORS_ALLOW_ORIGINS", defaultCORSAllowOrigins),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func getInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
