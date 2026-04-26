package routes

import "github.com/gofiber/fiber/v2"

// RegisterOpenAPIRoutes exposes the generated OpenAPI description for the API.
func RegisterOpenAPIRoutes(app *fiber.App) {
	app.Get("/openapi.json", func(c *fiber.Ctx) error {
		c.Type("json", "utf-8")
		return c.SendString(openAPISpec)
	})
}

const openAPISpec = `{
  "openapi": "3.1.0",
  "info": {
    "title": "Neko-Love API",
    "version": "v4",
    "description": "Random image API with local asset categories and image filter endpoints."
  },
  "servers": [
    {
      "url": "/",
      "description": "Current server"
    }
  ],
  "paths": {
    "/api/v4/{category}": {
      "get": {
        "summary": "Get random image metadata for a category",
        "parameters": [
          {
            "name": "category",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Random image metadata",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ImageMeta"
                },
                "examples": {
                  "neko": {
                    "summary": "Example metadata response",
                    "value": {
                      "name": "04.webp",
                      "category": "neko",
                      "path": "/api/v4/images/neko/04.webp",
                      "size": "438.20 KB",
                      "size_bytes": 448717,
                      "modified_at": 1720018294,
                      "mime_type": "image/webp"
                    }
                  }
                }
              }
            }
          },
          "404": {
            "description": "Category not found or empty"
          }
        }
      }
    },
    "/api/v4/images/{category}/{name}": {
      "get": {
        "summary": "Serve an image file directly",
        "parameters": [
          {
            "name": "category",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "name",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Image binary response",
            "content": {
              "image/png": {},
              "image/jpeg": {},
              "image/webp": {}
            }
          },
          "304": {
            "description": "Cached client copy is still valid"
          },
          "404": {
            "description": "Image not found"
          }
        }
      }
    },
    "/api/v4/filters/{filter}": {
      "get": {
        "summary": "Apply a filter to a remote public image",
        "parameters": [
          {
            "name": "filter",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "enum": [
                "amber",
                "anime_outline",
                "aqua",
                "blurple",
                "bubblegum",
                "crimson",
                "deepfry",
                "fuchsia",
                "glitch",
                "greyscale",
                "holographic",
                "mint",
                "negative",
                "pixelate",
                "poppink",
                "posterize",
                "sunset",
                "vaporwave"
              ]
            }
          },
          {
            "name": "image",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "format": "uri"
            },
            "description": "Public image URL. Private and localhost addresses are rejected."
          },
          {
            "name": "pixel_size",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "minimum": 2,
              "maximum": 64,
              "default": 6
            },
            "description": "Optional block size for the pixelate filter only. Larger values produce stronger pixelation."
          }
        ],
        "responses": {
          "200": {
            "description": "Filtered image response",
            "content": {
              "image/png": {},
              "image/jpeg": {},
              "image/webp": {},
              "image/gif": {}
            }
          },
          "400": {
            "description": "Invalid input or unsupported remote image"
          },
          "413": {
            "description": "Remote image exceeds configured size limit"
          }
        }
      }
    },
    "/openapi.json": {
      "get": {
        "summary": "OpenAPI document",
        "responses": {
          "200": {
            "description": "OpenAPI JSON document",
            "content": {
              "application/json": {}
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "ImageMeta": {
        "type": "object",
        "required": [
          "name",
          "category",
          "path",
          "size",
          "size_bytes",
          "modified_at",
          "mime_type"
        ],
        "properties": {
          "name": {
            "type": "string"
          },
          "category": {
            "type": "string"
          },
          "path": {
            "type": "string"
          },
          "size": {
            "type": "string"
          },
          "size_bytes": {
            "type": "integer"
          },
          "modified_at": {
            "type": "integer"
          },
          "mime_type": {
            "type": "string"
          }
        }
      }
    }
  }
}`
