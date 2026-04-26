package routes

import (
	"neko-love/services/cache"

	"github.com/gofiber/fiber/v2"
)

// RegisterDebugRoutes registers debug-related routes to the provided Fiber router.
// Specifically, it adds a GET endpoint at "/cache/:category" that returns a JSON
// response containing the list of cached files for the specified category.
func RegisterDebugRoutes(router fiber.Router, imageCache *cache.ImageCache) {
	router.Get("/cache/:category", func(c *fiber.Ctx) error {
		category := c.Params("category")
		files := imageCache.GetFiles(category)
		return c.JSON(fiber.Map{
			"category": category,
			"count":    len(files),
			"files":    files,
		})
	})
}
