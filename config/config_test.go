package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaultsForInvalidValues(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ASSETS_DIR", "")
	t.Setenv("ENABLE_DEBUG_ROUTES", "not-a-bool")
	t.Setenv("FILTER_TIMEOUT", "oops")
	t.Setenv("FILTER_MAX_BYTES", "-2")

	cfg := Load()

	if cfg.Port != defaultPort {
		t.Fatalf("expected default port, got %q", cfg.Port)
	}
	if cfg.AssetsDir != defaultAssetsDir {
		t.Fatalf("expected default assets dir, got %q", cfg.AssetsDir)
	}
	if cfg.EnableDebugRoutes != defaultEnableDebugRoutes {
		t.Fatalf("expected default debug flag, got %v", cfg.EnableDebugRoutes)
	}
	if cfg.FilterTimeout != defaultFilterTimeout {
		t.Fatalf("expected default timeout, got %v", cfg.FilterTimeout)
	}
	if cfg.FilterMaxBytes != defaultFilterMaxBytes {
		t.Fatalf("expected default max bytes, got %d", cfg.FilterMaxBytes)
	}
	if cfg.CORSAllowOrigins != defaultCORSAllowOrigins {
		t.Fatalf("expected default CORS origins, got %q", cfg.CORSAllowOrigins)
	}
}

func TestLoadReadsEnvironmentOverrides(t *testing.T) {
	t.Setenv("PORT", "4040")
	t.Setenv("ASSETS_DIR", "./custom-assets")
	t.Setenv("ENABLE_DEBUG_ROUTES", "true")
	t.Setenv("FILTER_TIMEOUT", "3s")
	t.Setenv("FILTER_MAX_BYTES", "2048")
	t.Setenv("CORS_ALLOW_ORIGINS", "https://docs.neko-love.com")

	cfg := Load()

	if cfg.Port != "4040" {
		t.Fatalf("expected port override, got %q", cfg.Port)
	}
	if cfg.AssetsDir != "./custom-assets" {
		t.Fatalf("expected assets dir override, got %q", cfg.AssetsDir)
	}
	if !cfg.EnableDebugRoutes {
		t.Fatal("expected debug routes to be enabled")
	}
	if cfg.FilterTimeout != 3*time.Second {
		t.Fatalf("expected timeout override, got %v", cfg.FilterTimeout)
	}
	if cfg.FilterMaxBytes != 2048 {
		t.Fatalf("expected max bytes override, got %d", cfg.FilterMaxBytes)
	}
	if cfg.CORSAllowOrigins != "https://docs.neko-love.com" {
		t.Fatalf("expected CORS origins override, got %q", cfg.CORSAllowOrigins)
	}
}
