# Neko-Love API V4

Community rewrite of the original Neko-Love API in Go with Fiber.

## Overview

- Serves random anime-style images from local folders such as `assets/neko/` or `assets/hug/`.
- Returns metadata through the API and serves the raw file through a dedicated image route.
- Includes an image filter endpoint for PNG, JPEG, WEBP, and animated GIF files.

## Requirements

- Go `1.24.4` or newer

## Run locally

```bash
git clone https://github.com/Otaku17/neko-love.git
cd neko-love
go mod tidy
go run .
```

The server starts on port `3030` by default.

Open the interactive playground in your browser:

```text
http://localhost:3030/docs
```

Open the Redoc version:

```text
http://localhost:3030/docs/redoc
```

Open the OpenAPI document:

```text
http://localhost:3030/openapi.json
```

## Environment variables

- `PORT`: HTTP port, default `3030`
- `ASSETS_DIR`: assets root directory, default `./assets`
- `ENABLE_DEBUG_ROUTES`: enables `/debug/*` routes, default `false`
- `FILTER_TIMEOUT`: timeout for remote filter fetches, default `10s`
- `FILTER_MAX_BYTES`: maximum remote image size accepted by `/filters`, default `10485760`
- `CORS_ALLOW_ORIGINS`: optional comma-separated origins for cross-origin docs usage, example `https://docs.neko-love.com`

## API

### Random image metadata

```http
GET /api/v4/:category
```

Example response:

```json
{
  "name": "04.webp",
  "category": "neko",
  "path": "/api/v4/images/neko/04.webp",
  "size": "438.20 KB",
  "size_bytes": 448717,
  "modified_at": 1720018294,
  "mime_type": "image/webp"
}
```

### Direct image

```http
GET /api/v4/images/:category/:name
```

### Filter endpoint

```http
GET /api/v4/filters/:filter?image=<url>&pixel_size=<size>
```

Supported formats:

- JPEG
- PNG
- WEBP
- GIF

Available filters:

- `amber`
- `anime_outline`
- `aqua`
- `blurple`
- `bubblegum`
- `crimson`
- `deepfry`
- `fuchsia`
- `glitch`
- `greyscale`
- `holographic`
- `mint`
- `negative`
- `pixelate`
- `poppink`
- `posterize`
- `sunset`
- `vaporwave`

Example:

```text
http://localhost:3030/api/v4/filters/deepfry?image=https://example.com/image.png
```

For the `pixelate` filter you can control the strength with `pixel_size`:

```text
http://localhost:3030/api/v4/filters/pixelate?image=https://example.com/image.png&pixel_size=12
```

- `pixel_size` is optional
- allowed range: `2` to `64`
- larger values produce stronger pixelation

## Assets

To add a new category:

1. Create a folder in `assets/<category>`.
2. Put image files in that folder.

The API discovers categories automatically.

## Debug routes

When `ENABLE_DEBUG_ROUTES=true`, the server exposes:

```http
GET /debug/cache/:category
```

This route is intended for local development only.

## Tests

- Unit tests stay next to the code in each package as `*_test.go`.
- Integration tests live in [tests/integration](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/tests/integration/api_test.go).

Run everything:

```bash
go test ./...
```

## Static site for GitHub Pages

The repository now includes a publishable static site in [site](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/site/index.html).

- [site/index.html](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/site/index.html): landing page
- [site/docs/index.html](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/site/docs/index.html): Scalar docs
- [site/docs/redoc.html](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/site/docs/redoc.html): Redoc docs
- [site/openapi.json](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/site/openapi.json): static OpenAPI document for Pages
- [.github/workflows/deploy-pages.yml](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/.github/workflows/deploy-pages.yml): GitHub Pages deployment workflow

Important:

- GitHub Pages can host the static site and docs, but not the Go API itself.
- Before publishing, replace `https://api.neko-love.com` in [site/openapi.json](c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/site/openapi.json) with your real API domain.
- If your docs are hosted on a different domain than the API and you want live browser requests from Scalar or Redoc, set `CORS_ALLOW_ORIGINS` on the API to your Pages domain.
