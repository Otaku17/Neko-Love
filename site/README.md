# Static Site

This folder is ready to publish on GitHub Pages.

Files:

- `index.html`: marketing-style landing page
- `docs/index.html`: Scalar API reference
- `docs/redoc.html`: Redoc API reference
- `openapi.json`: static OpenAPI document used by the Pages docs

Before publishing:

1. Update `site/openapi.json` and replace `https://api.neko-love.com` with your real API domain.
2. Configure the Go API with `CORS_ALLOW_ORIGINS=https://<your-pages-domain>` if you want the static docs to use live browser requests.
3. Publish the `site/` folder with the GitHub Pages workflow included in this repository.

