package routes

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"neko-love/services/cache"

	"github.com/gofiber/fiber/v2"
)

func TestSetImageHeadersSetsCacheMetadata(t *testing.T) {
	t.Parallel()

	imageCache := newTestImageCache(t)
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		if setImageHeaders(c, "neko", "cat.png", imageCache) {
			return c.SendStatus(fiber.StatusNotModified)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Cache-Control"); got == "" {
		t.Fatal("expected Cache-Control header")
	}
	if got := resp.Header.Get("ETag"); got == "" {
		t.Fatal("expected ETag header")
	}
	if got := resp.Header.Get("Last-Modified"); got == "" {
		t.Fatal("expected Last-Modified header")
	}
}

func TestSetImageHeadersReturnsNotModifiedOnMatchingETag(t *testing.T) {
	t.Parallel()

	imageCache := newTestImageCache(t)
	app := fiber.New()

	meta, ok := imageCache.GetImageMeta("neko", "cat.png")
	if !ok {
		t.Fatal("expected metadata")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		if setImageHeaders(c, "neko", "cat.png", imageCache) {
			return c.SendStatus(fiber.StatusNotModified)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("If-None-Match", buildImageETag(meta))

	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotModified {
		t.Fatalf("expected status 304, got %d", resp.StatusCode)
	}
}

func newTestImageCache(t *testing.T) *cache.ImageCache {
	t.Helper()

	root := t.TempDir()
	categoryDir := filepath.Join(root, "neko")
	if err := os.Mkdir(categoryDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(categoryDir, "cat.png"), pngFixture(), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	imageCache, err := cache.New(root)
	if err != nil {
		t.Fatalf("cache init failed: %v", err)
	}
	t.Cleanup(func() {
		_ = imageCache.Close()
	})

	return imageCache
}

func pngFixture() []byte {
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0xf8, 0xcf, 0xc0, 0xf0,
		0x1f, 0x00, 0x05, 0x00, 0x01, 0xff, 0x89, 0x99,
		0x3d, 0x1d, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
		0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
}
