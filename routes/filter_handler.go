package routes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"neko-love/services"

	"github.com/gofiber/fiber/v2"
)

var errImageTooLarge = errors.New("image exceeds maximum allowed size")

type FilterHandler struct {
	client   *http.Client
	maxBytes int64
}

func NewFilterHandler(timeout time.Duration, maxBytes int64) *FilterHandler {
	dialer := &net.Dialer{Timeout: timeout}

	transport := &http.Transport{
		Proxy: nil,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				host = address
			}

			if err := validateHost(host); err != nil {
				return nil, err
			}

			return dialer.DialContext(ctx, network, address)
		},
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return errors.New("too many redirects")
			}
			return validateSourceURL(req.URL.String())
		},
	}

	return &FilterHandler{
		client:   client,
		maxBytes: maxBytes,
	}
}

func NewFilterHandlerWithClient(client *http.Client, maxBytes int64) *FilterHandler {
	if client == nil {
		client = http.DefaultClient
	}

	return &FilterHandler{
		client:   client,
		maxBytes: maxBytes,
	}
}

func (h *FilterHandler) Handle(c *fiber.Ctx) error {
	filter := c.Params("filter")
	if filter == "" {
		return fiber.ErrNotFound
	}

	options, err := parseFilterOptions(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	imageURL := c.Query("image")
	if imageURL == "" {
		return fiber.NewError(fiber.StatusBadRequest, "image URL is required")
	}

	data, err := h.fetchImage(c.UserContext(), imageURL)
	if err != nil {
		switch {
		case errors.Is(err, errImageTooLarge):
			return fiber.NewError(fiber.StatusRequestEntityTooLarge, "image exceeds maximum allowed size")
		default:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	format := http.DetectContentType(data)
	if !strings.HasPrefix(format, "image/") {
		return fiber.NewError(fiber.StatusBadRequest, "remote file is not a supported image")
	}

	c.Locals("noCache", true)

	if strings.HasPrefix(format, "image/gif") {
		return handleGIF(c, filter, data, options)
	}

	srcImg, formatStr, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to decode image")
	}

	result := services.ApplyFilterWithOptions(filter, srcImg, options)
	return services.EncodeAndSetContentType(c, result, formatStr)
}

func RegisterFilterRoutes(router fiber.Router, handler *FilterHandler) {
	router.Get("/filters/:filter", handler.Handle)
}

func (h *FilterHandler) fetchImage(ctx context.Context, rawURL string) ([]byte, error) {
	if err := validateSourceURL(rawURL); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, errors.New("invalid image URL")
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
	}

	limited := io.LimitReader(resp.Body, h.maxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	if int64(len(data)) > h.maxBytes {
		return nil, errImageTooLarge
	}

	return data, nil
}

func validateSourceURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return errors.New("invalid image URL")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("only http and https image URLs are allowed")
	}

	if parsed.Hostname() == "" {
		return errors.New("image URL must include a hostname")
	}

	if parsed.User != nil {
		return errors.New("image URL must not include credentials")
	}

	if err := validateHost(parsed.Hostname()); err != nil {
		return err
	}

	return nil
}

func validateHost(host string) error {
	if strings.EqualFold(host, "localhost") {
		return errors.New("localhost image URLs are not allowed")
	}

	if ip := net.ParseIP(host); ip != nil {
		if isPrivateIP(ip) {
			return errors.New("private or local image URLs are not allowed")
		}
		return nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return errors.New("failed to resolve image host")
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return errors.New("private or local image URLs are not allowed")
		}
	}

	return nil
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",
		"169.254.0.0/16",
		"fc00::/7",
		"fe80::/10",
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil && network.Contains(ip) {
			return true
		}
	}

	return false
}

func parseFilterOptions(c *fiber.Ctx) (services.FilterOptions, error) {
	options := services.DefaultFilterOptions()

	if value := c.Query("pixel_size"); value != "" {
		pixelSize, err := strconv.Atoi(value)
		if err != nil {
			return options, errors.New("pixel_size must be an integer")
		}
		if pixelSize < 2 || pixelSize > 64 {
			return options, errors.New("pixel_size must be between 2 and 64")
		}
		options.PixelSize = pixelSize
	}

	return options, nil
}

func handleGIF(c *fiber.Ctx, filter string, data []byte, options services.FilterOptions) error {
	gifReader := bytes.NewReader(data)
	gifData, err := gif.DecodeAll(gifReader)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to decode GIF")
	}
	filteredGIF, err := services.ProcessGIF(filter, gifData, options)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to process GIF")
	}
	c.Set("Content-Type", "image/gif")
	return gif.EncodeAll(c.Context().Response.BodyWriter(), filteredGIF)
}
