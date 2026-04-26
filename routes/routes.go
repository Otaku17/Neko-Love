package routes

import (
	"neko-love/config"
	"neko-love/middlewares"
	"neko-love/services/cache"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupRoutes configures the main application routes for the Fiber app.
// It applies middleware to disable caching, sets up API version 4 routes for images and filters,
// and registers debug routes under the "/debug" path.
//
// Parameters:
//   - app: The Fiber application instance to which the routes will be attached.
func SetupRoutes(app *fiber.App, cfg config.Config, imageCache *cache.ImageCache) {
	if cfg.CORSAllowOrigins != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: cfg.CORSAllowOrigins,
			AllowMethods: "GET,HEAD,OPTIONS",
		}))
	}

	app.Use(middlewares.NoCache())
	RegisterOpenAPIRoutes(app)
	RegisterDocsRoutes(app)

	api := app.Group("/api/v4")
	RegisterImageRoutes(api, imageCache)
	RegisterFilterRoutes(api, NewFilterHandler(cfg.FilterTimeout, cfg.FilterMaxBytes))

	if cfg.EnableDebugRoutes {
		debug := app.Group("/debug")
		RegisterDebugRoutes(debug, imageCache)
	}
}
