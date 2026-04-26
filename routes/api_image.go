package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"neko-love/services/cache"

	"github.com/gofiber/fiber/v2"
)

const httpTimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

type ImageHandler struct {
	cache *cache.ImageCache
}

// NewImageHandler creates and returns a new ImageHandler instance with the provided ImageCache.
// The cache parameter is used to store and retrieve image data efficiently.
func NewImageHandler(c *cache.ImageCache) *ImageHandler {
	return &ImageHandler{cache: c}
}

// GetRandomImage handles HTTP requests to retrieve a random image from a specified category.
// It extracts the category from the route parameters, fetches a random image name from the cache,
// and retrieves the corresponding image path. If successful, it sets appropriate response headers
// and serves the image file. Returns a 404 error if the category or image is not found.
func (h *ImageHandler) GetRandomImage(c *fiber.Ctx) error {
	category := c.Params("category")

	name, err := h.cache.GetRandom(category)
	if err != nil {
		return fiber.ErrNotFound
	}

	imagePath, ok := h.cache.GetImagePath(category, name)
	if !ok {
		return fiber.ErrNotFound
	}

	c.Locals("noCache", true)
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, name))
	_ = setImageHeaders(c, category, name, h.cache)
	return c.SendFile(imagePath)
}

// GetRandomImageMeta handles the HTTP request to retrieve metadata for a random image
// within a specified category.
func (h *ImageHandler) GetRandomImageMeta(c *fiber.Ctx) error {
	category := c.Params("category")

	name, err := h.cache.GetRandom(category)
	if err != nil {
		return fiber.ErrNotFound
	}

	meta, ok := h.cache.GetImageMeta(category, name)
	if !ok {
		return fiber.ErrNotFound
	}

	return c.JSON(fiber.Map{
		"name":        name,
		"category":    category,
		"path":        fmt.Sprintf("/api/v4/images/%s/%s", category, name),
		"size":        meta.Readable,
		"size_bytes":  meta.Size,
		"modified_at": meta.ModifiedAt,
		"mime_type":   meta.MimeType,
	})
}

// ServeImage handles HTTP requests to serve an image file based on the provided
// category and name parameters in the URL. It retrieves the image path from the
// cache and sends the file as a response. If the image is not found in the cache,
// it returns a 404 Not Found error.
func (h *ImageHandler) ServeImage(c *fiber.Ctx) error {
	category := c.Params("category")
	name := c.Params("name")

	path, ok := h.cache.GetImagePath(category, name)
	if !ok {
		return fiber.ErrNotFound
	}

	if setImageHeaders(c, category, name, h.cache) {
		return c.SendStatus(fiber.StatusNotModified)
	}

	return c.SendFile(path)
}

// RegisterImageRoutes registers image-related API routes to the provided Fiber router.
func RegisterImageRoutes(router fiber.Router, imageCache *cache.ImageCache) {
	handler := NewImageHandler(imageCache)

	router.Get("/:category", handler.GetRandomImageMeta)
	router.Get("/images/:category/:name", handler.ServeImage)
}

func setImageHeaders(c *fiber.Ctx, category, name string, imageCache *cache.ImageCache) bool {
	meta, ok := imageCache.GetImageMeta(category, name)
	if !ok {
		return false
	}

	etag := buildImageETag(meta)
	lastModified := time.Unix(meta.ModifiedAt, 0).UTC().Format(httpTimeFormat)

	c.Set(fiber.HeaderCacheControl, "public, max-age=31536000, immutable")
	c.Set(fiber.HeaderETag, etag)
	c.Set(fiber.HeaderLastModified, lastModified)
	if meta.MimeType != "" {
		c.Set(fiber.HeaderContentType, meta.MimeType)
	}

	if match := c.Get(fiber.HeaderIfNoneMatch); match != "" && match == etag {
		return true
	}

	if modifiedSince := c.Get(fiber.HeaderIfModifiedSince); modifiedSince != "" {
		if t, err := time.Parse(http.TimeFormat, modifiedSince); err == nil && !t.Before(time.Unix(meta.ModifiedAt, 0).UTC()) {
			return true
		}
	}

	return false
}

func buildImageETag(meta cache.FileMeta) string {
	return `W/"` + strconv.FormatInt(meta.ModifiedAt, 10) + `-` + strconv.FormatInt(meta.Size, 10) + `"`
}
