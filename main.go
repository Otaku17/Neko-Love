package main

import (
	"log"
	"neko-love/config"
	"neko-love/routes"
	"neko-love/services/cache"

	"github.com/gofiber/fiber/v2"
)

// main is the entry point of the application. It initializes a new Fiber web server,
// starts watching for asset changes, sets up the application routes, and begins
// listening for incoming HTTP requests on port 3030.
func main() {
	cfg := config.Load()
	app := fiber.New()

	cacheAssets, err := cache.New(cfg.AssetsDir)
	if err != nil {
		log.Fatalf("failed to initialize image cache: %v", err)
	}
	defer cacheAssets.Close()

	routes.SetupRoutes(app, cfg, cacheAssets)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
