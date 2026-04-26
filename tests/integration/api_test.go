package integration

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"neko-love/config"
	"neko-love/middlewares"
	"neko-love/routes"
	"neko-love/services/cache"

	"github.com/gofiber/fiber/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestRandomImageMetaEndpoint(t *testing.T) {
	t.Parallel()

	app := newIntegrationApp(t, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v4/neko", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}

	expectedParts := []string{`"category":"neko"`, `"name":"cat.png"`, `"/api/v4/images/neko/cat.png"`}
	for _, part := range expectedParts {
		if !bytes.Contains(body, []byte(part)) {
			t.Fatalf("expected response body to contain %q, got %s", part, string(body))
		}
	}
}

func TestDocsRouteReturnsHTML(t *testing.T) {
	t.Parallel()

	app := newIntegrationApp(t, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}

	if !bytes.Contains(body, []byte("Scalar.createApiReference")) {
		t.Fatalf("expected docs page content, got %s", string(body))
	}
}

func TestRedocRouteReturnsHTML(t *testing.T) {
	t.Parallel()

	app := newIntegrationApp(t, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/docs/redoc", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}

	if !bytes.Contains(body, []byte("redoc.standalone.js")) {
		t.Fatalf("expected redoc page content, got %s", string(body))
	}
}

func TestOpenAPIRouteReturnsJSON(t *testing.T) {
	t.Parallel()

	app := newIntegrationApp(t, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}

	if !bytes.Contains(body, []byte(`"openapi": "3.1.0"`)) {
		t.Fatalf("expected openapi document, got %s", string(body))
	}
	if !bytes.Contains(body, []byte(`"pixel_size"`)) {
		t.Fatalf("expected pixel_size parameter in spec, got %s", string(body))
	}
}

func TestDebugRouteDisabledByDefault(t *testing.T) {
	t.Parallel()

	app := newIntegrationApp(t, false, nil)

	req := httptest.NewRequest(http.MethodGet, "/debug/cache/neko", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestDebugRouteEnabledByConfig(t *testing.T) {
	t.Parallel()

	app := newIntegrationApp(t, true, nil)

	req := httptest.NewRequest(http.MethodGet, "/debug/cache/neko", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestFilterEndpointEndToEnd(t *testing.T) {
	t.Parallel()

	pngData := pngFixture()
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(pngData)),
				Header:     http.Header{"Content-Type": []string{"image/png"}},
			}, nil
		}),
		Timeout: time.Second,
	}

	app := newIntegrationApp(t, false, client)

	req := httptest.NewRequest(http.MethodGet, "/api/v4/filters/negative?image=https://1.1.1.1/test.png", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "image/png" {
		t.Fatalf("expected image/png content type, got %q", got)
	}
}

func newIntegrationApp(t *testing.T, enableDebug bool, httpClient *http.Client) *fiber.App {
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

	cfg := config.Config{
		Port:              "3030",
		AssetsDir:         root,
		EnableDebugRoutes: enableDebug,
		FilterTimeout:     time.Second,
		FilterMaxBytes:    1 << 20,
	}

	app := fiber.New()
	if httpClient == nil {
		routes.SetupRoutes(app, cfg, imageCache)
		return app
	}

	app.Use(middlewares.NoCache())
	api := app.Group("/api/v4")
	routes.RegisterImageRoutes(api, imageCache)
	routes.RegisterFilterRoutes(api, routes.NewFilterHandlerWithClient(httpClient, cfg.FilterMaxBytes))

	if cfg.EnableDebugRoutes {
		debug := app.Group("/debug")
		routes.RegisterDebugRoutes(debug, imageCache)
	}

	return app
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
