package routes

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"neko-love/services/cache"

	"github.com/gofiber/fiber/v2"
)

func BenchmarkGetRandomImageMetaRoute(b *testing.B) {
	root := b.TempDir()
	categoryDir := filepath.Join(root, "neko")
	if err := os.Mkdir(categoryDir, 0o755); err != nil {
		b.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(categoryDir, "cat.png"), pngFixture(), 0o644); err != nil {
		b.Fatalf("write file failed: %v", err)
	}

	imageCache, err := cache.New(root)
	if err != nil {
		b.Fatalf("cache init failed: %v", err)
	}
	defer imageCache.Close()

	app := fiber.New()
	RegisterImageRoutes(app, imageCache)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/neko", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			b.Fatalf("request failed: %v", err)
		}
		_ = resp.Body.Close()
	}
}
