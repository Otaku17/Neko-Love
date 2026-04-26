package routes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestValidateSourceURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{name: "rejects empty host", rawURL: "https:///image.png", wantErr: true},
		{name: "rejects unsupported scheme", rawURL: "file:///tmp/image.png", wantErr: true},
		{name: "rejects localhost", rawURL: "http://localhost/image.png", wantErr: true},
		{name: "rejects private ip", rawURL: "http://192.168.1.10/image.png", wantErr: true},
		{name: "accepts public ip", rawURL: "https://1.1.1.1/image.png", wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateSourceURL(tt.rawURL)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for %q", tt.rawURL)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.rawURL, err)
			}
		})
	}
}

func TestFetchImageRejectsOversizedResponse(t *testing.T) {
	t.Parallel()

	handler := NewFilterHandlerWithClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(bytes.Repeat([]byte("a"), 6))),
				Header:     make(http.Header),
			}, nil
		}),
		Timeout: time.Second,
	}, 5)

	_, err := handler.fetchImage(context.Background(), "https://1.1.1.1/image.png")
	if !errors.Is(err, errImageTooLarge) {
		t.Fatalf("expected errImageTooLarge, got %v", err)
	}
}

func TestFilterRouteRejectsMissingImageQuery(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	RegisterFilterRoutes(app, NewFilterHandlerWithClient(http.DefaultClient, 1024))

	req := httptest.NewRequest(http.MethodGet, "/filters/blurple", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestFilterRouteProcessesPNG(t *testing.T) {
	t.Parallel()

	pngData := []byte{
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

	handler := NewFilterHandlerWithClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(pngData)),
				Header:     http.Header{"Content-Type": []string{"image/png"}},
			}, nil
		}),
		Timeout: time.Second,
	}, int64(len(pngData)+10))

	app := fiber.New()
	RegisterFilterRoutes(app, handler)

	req := httptest.NewRequest(http.MethodGet, "/filters/negative?image=https://1.1.1.1/test.png", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if got := resp.Header.Get("Content-Type"); got != "image/png" {
		t.Fatalf("expected image/png content type, got %q", got)
	}
}

func TestParseFilterOptionsRejectsInvalidPixelSize(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		_, err := parseFilterOptions(c)
		if err == nil {
			return c.SendStatus(fiber.StatusOK)
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	})

	req := httptest.NewRequest(http.MethodGet, "/?pixel_size=abc", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestParseFilterOptionsAcceptsPixelSize(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		options, err := parseFilterOptions(c)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.JSON(fiber.Map{"pixel_size": options.PixelSize})
	})

	req := httptest.NewRequest(http.MethodGet, "/?pixel_size=12", nil)
	resp, err := app.Test(req, int(time.Second.Milliseconds()))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}
