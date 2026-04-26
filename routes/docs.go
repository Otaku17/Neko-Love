package routes

import (
	"html/template"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type docsPageData struct {
	OpenAPIPath string
	RedocPath   string
}

// RegisterDocsRoutes exposes a polished documentation experience backed by the local OpenAPI document.
func RegisterDocsRoutes(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		return renderDocsPage(c, "scalar-docs", scalarDocsHTML, docsPageData{
			OpenAPIPath: "/openapi.json",
			RedocPath:   "/docs/redoc",
		})
	})

	app.Get("/docs/redoc", func(c *fiber.Ctx) error {
		return renderDocsPage(c, "redoc-docs", redocDocsHTML, docsPageData{
			OpenAPIPath: "/openapi.json",
			RedocPath:   "/docs/redoc",
		})
	})
}

func renderDocsPage(c *fiber.Ctx, name, src string, data docsPageData) error {
	tmpl, err := template.New(name).Parse(src)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to render docs page")
	}

	var out strings.Builder
	if err := tmpl.Execute(&out, data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to render docs page")
	}

	c.Type("html", "utf-8")
	return c.SendString(out.String())
}

const scalarDocsHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Neko-Love API Docs</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;700&family=IBM+Plex+Mono:wght@400;500&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg: #f6ede1;
      --bg-deep: #ecd4b5;
      --paper: rgba(255, 251, 246, 0.9);
      --paper-strong: rgba(255, 248, 239, 0.98);
      --ink: #1d1916;
      --muted: #6f665f;
      --line: rgba(32, 22, 14, 0.12);
      --line-strong: rgba(32, 22, 14, 0.18);
      --accent: #de6d3f;
      --accent-deep: #b84b22;
      --teal: #1e5756;
      --teal-soft: rgba(30, 87, 86, 0.12);
      --shadow: 0 22px 60px rgba(68, 47, 27, 0.16);
      --shadow-soft: 0 14px 36px rgba(68, 47, 27, 0.09);
      --sans: "Space Grotesk", "Segoe UI", sans-serif;
      --mono: "IBM Plex Mono", Consolas, monospace;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: var(--sans);
      color: var(--ink);
      background:
        radial-gradient(circle at top left, rgba(222,109,63,0.2), transparent 26%),
        radial-gradient(circle at 92% 8%, rgba(30,87,86,0.16), transparent 22%),
        linear-gradient(180deg, var(--bg) 0%, var(--bg-deep) 100%);
    }
    a { color: inherit; }
    .shell {
      max-width: 1460px;
      margin: 0 auto;
      padding: 26px 18px 56px;
    }
    .hero {
      position: relative;
      overflow: hidden;
      padding: 34px;
      border-radius: 34px;
      border: 1px solid var(--line);
      background: linear-gradient(145deg, rgba(255,251,246,0.98), rgba(248,239,227,0.92));
      box-shadow: var(--shadow);
    }
    .hero::after {
      content: "";
      position: absolute;
      right: -60px;
      bottom: -60px;
      width: 280px;
      height: 280px;
      border-radius: 50%;
      background:
        radial-gradient(circle at 32% 32%, rgba(222,109,63,0.22), transparent 48%),
        radial-gradient(circle at 68% 68%, rgba(30,87,86,0.18), transparent 42%);
      pointer-events: none;
    }
    .eyebrow {
      margin: 0 0 10px;
      color: var(--teal);
      text-transform: uppercase;
      letter-spacing: 0.18em;
      font-size: 12px;
      font-weight: 800;
    }
    h1 {
      margin: 0;
      font-size: clamp(2.8rem, 6vw, 5.4rem);
      line-height: 0.9;
      letter-spacing: -0.06em;
      max-width: 8ch;
    }
    .lead {
      margin: 16px 0 0;
      max-width: 62ch;
      color: var(--muted);
      font-size: 1.05rem;
      line-height: 1.6;
    }
    .hero-grid {
      display: grid;
      gap: 16px;
      margin-top: 28px;
    }
    @media (min-width: 1040px) {
      .hero-grid {
        grid-template-columns: minmax(0, 1.2fr) minmax(360px, 0.8fr);
        align-items: start;
      }
    }
    .hero-panel,
    .mini-panel,
    .reference-shell {
      border-radius: 28px;
      border: 1px solid var(--line);
      background: var(--paper);
      box-shadow: var(--shadow-soft);
      backdrop-filter: blur(10px);
    }
    .hero-panel {
      padding: 20px;
    }
    .hero-panel h2,
    .mini-panel h2 {
      margin: 0 0 10px;
      font-size: 1rem;
    }
    .hero-panel p,
    .mini-panel p {
      margin: 0;
      color: var(--muted);
      line-height: 1.6;
    }
    .actions {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      margin-top: 20px;
    }
    .button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 46px;
      padding: 0 18px;
      border-radius: 999px;
      border: 1px solid transparent;
      text-decoration: none;
      font-weight: 700;
      transition: transform 140ms ease, box-shadow 140ms ease, background 140ms ease;
    }
    .button:hover {
      transform: translateY(-1px);
      box-shadow: 0 8px 20px rgba(72, 52, 36, 0.12);
    }
    .button.primary {
      background: linear-gradient(135deg, var(--accent), var(--accent-deep));
      color: #fff;
    }
    .button.secondary {
      background: rgba(255,255,255,0.72);
      color: var(--ink);
      border-color: var(--line);
    }
    .mini-stack {
      display: grid;
      gap: 16px;
    }
    .mini-panel {
      padding: 18px 20px;
    }
    .stats {
      display: grid;
      gap: 14px;
      margin-top: 18px;
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .stat {
      padding: 14px;
      border-radius: 18px;
      background: rgba(255,255,255,0.66);
      border: 1px solid var(--line);
    }
    .stat strong {
      display: block;
      font-size: 1.35rem;
      margin-bottom: 6px;
    }
    .stat span {
      color: var(--muted);
      font-size: 0.92rem;
    }
    .section {
      margin-top: 22px;
      padding: 22px;
      border-radius: 30px;
      border: 1px solid var(--line);
      background: var(--paper-strong);
      box-shadow: var(--shadow);
    }
    .section-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      margin-bottom: 18px;
      flex-wrap: wrap;
    }
    .section-head h2 {
      margin: 0;
      font-size: 1.3rem;
    }
    .section-head p {
      margin: 0;
      color: var(--muted);
    }
    .examples {
      display: grid;
      gap: 16px;
    }
    @media (min-width: 980px) {
      .examples {
        grid-template-columns: repeat(3, minmax(0, 1fr));
      }
    }
    .example-card {
      overflow: hidden;
      border-radius: 22px;
      border: 1px solid var(--line);
      background: linear-gradient(180deg, rgba(255,255,255,0.98), rgba(246,239,231,0.92));
      box-shadow: var(--shadow-soft);
    }
    .example-card h3 {
      margin: 0;
      padding: 14px 16px;
      font-size: 0.96rem;
      border-bottom: 1px solid var(--line);
      background: linear-gradient(90deg, rgba(222,109,63,0.1), rgba(30,87,86,0.08));
    }
    pre {
      margin: 0;
      padding: 16px;
      overflow: auto;
      font-family: var(--mono);
      font-size: 0.85rem;
      line-height: 1.55;
      color: #1f2933;
      background: rgba(255,255,255,0.55);
    }
    #api-reference {
      min-height: 960px;
      border-radius: 24px;
      overflow: hidden;
      border: 1px solid var(--line);
      background: #fff;
    }
    .footer-note {
      margin-top: 12px;
      color: var(--muted);
      font-size: 0.92rem;
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <p class="eyebrow">Neko-Love Is Back</p>
      <h1>Docs that look like a product.</h1>
      <p class="lead">This reference is powered by Scalar on top of your local OpenAPI spec, so you keep clean endpoint docs, examples, schemas, and browser-side request testing without the rough default Swagger look.</p>
      <div class="hero-grid">
        <div class="hero-panel">
          <h2>What you get here</h2>
          <p>Browse every endpoint, inspect request and response shapes, test filters like <code>pixelate</code> live, and share a doc URL that feels closer to a real product launch than a raw dev tool.</p>
          <div class="actions">
            <a class="button primary" href="#api-reference">Open API Reference</a>
            <a class="button secondary" href="{{ .RedocPath }}">Open Redoc View</a>
            <a class="button secondary" href="{{ .OpenAPIPath }}" target="_blank" rel="noreferrer">Download OpenAPI JSON</a>
          </div>
          <div class="stats">
            <div class="stat">
              <strong>3</strong>
              <span>Main endpoint groups</span>
            </div>
            <div class="stat">
              <strong>18</strong>
              <span>Available filters documented</span>
            </div>
            <div class="stat">
              <strong>OpenAPI 3.1</strong>
              <span>Ready for tooling and SDKs</span>
            </div>
          </div>
        </div>
        <div class="mini-stack">
          <div class="mini-panel">
            <h2>Usage ideas</h2>
            <p>Perfect for wrapper libraries, bot integrations, and image generation utilities that need stable metadata and direct asset routes.</p>
          </div>
          <div class="mini-panel">
            <h2>For GitHub Pages too</h2>
            <p>The same OpenAPI structure can power a fully static doc build, so your landing page can live on GitHub Pages while the Go API stays on its own server.</p>
          </div>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <div>
          <h2>Quick examples</h2>
          <p>Copy the snippets, then tweak them directly in the interactive reference below.</p>
        </div>
      </div>
      <div class="examples">
        <article class="example-card">
          <h3>cURL</h3>
          <pre><code>curl "http://localhost:3030/api/v4/neko"

curl "http://localhost:3030/api/v4/filters/pixelate?image=https://example.com/image.png&pixel_size=12"</code></pre>
        </article>
        <article class="example-card">
          <h3>JavaScript</h3>
          <pre><code>const meta = await fetch("http://localhost:3030/api/v4/neko")
  .then((res) => res.json())

const url = new URL("http://localhost:3030/api/v4/filters/pixelate")
url.searchParams.set("image", "https://example.com/image.png")
url.searchParams.set("pixel_size", "12")</code></pre>
        </article>
        <article class="example-card">
          <h3>Python</h3>
          <pre><code>import requests

meta = requests.get("http://localhost:3030/api/v4/neko").json()

img = requests.get(
    "http://localhost:3030/api/v4/filters/pixelate",
    params={"image": "https://example.com/image.png", "pixel_size": 12},
)</code></pre>
        </article>
      </div>
    </section>

    <section class="section">
      <div class="section-head">
        <div>
          <h2>Scalar Reference</h2>
          <p>Interactive docs rendered from the same OpenAPI document your app exposes at runtime.</p>
        </div>
      </div>
      <div id="api-reference"></div>
      <p class="footer-note">Prefer a more classic documentation layout? Try the <a href="{{ .RedocPath }}">Redoc version</a>.</p>
    </section>
  </div>

  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  <script>
    Scalar.createApiReference("#api-reference", {
      url: "{{ .OpenAPIPath }}",
      darkMode: false,
      hideModels: false,
      searchHotKey: "k",
      layout: "modern",
      defaultHttpClient: {
        targetKey: "js",
        clientKey: "fetch"
      }
    })
  </script>
</body>
</html>`

const redocDocsHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Neko-Love API Docs | Redoc</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;700&family=IBM+Plex+Mono:wght@400;500&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg: #f6ede1;
      --paper: rgba(255,251,246,0.94);
      --ink: #1d1916;
      --muted: #6f665f;
      --line: rgba(32,22,14,0.12);
      --accent: #de6d3f;
      --teal: #1e5756;
      --shadow: 0 20px 56px rgba(68, 47, 27, 0.12);
      --sans: "Space Grotesk", "Segoe UI", sans-serif;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: var(--sans);
      color: var(--ink);
      background:
        radial-gradient(circle at top left, rgba(222,109,63,0.18), transparent 24%),
        linear-gradient(180deg, #f7efe3 0%, #ecddc8 100%);
    }
    .shell {
      max-width: 1420px;
      margin: 0 auto;
      padding: 24px 18px 44px;
    }
    .topbar {
      margin-bottom: 18px;
      padding: 22px 24px;
      border-radius: 28px;
      border: 1px solid var(--line);
      background: var(--paper);
      box-shadow: var(--shadow);
    }
    .topbar p {
      margin: 0 0 8px;
      color: var(--teal);
      text-transform: uppercase;
      letter-spacing: 0.16em;
      font-size: 12px;
      font-weight: 800;
    }
    .topbar h1 {
      margin: 0;
      font-size: clamp(2rem, 5vw, 3.6rem);
      line-height: 1;
      letter-spacing: -0.05em;
    }
    .topbar .sub {
      margin-top: 12px;
      color: var(--muted);
      max-width: 70ch;
      line-height: 1.6;
    }
    .links {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      margin-top: 18px;
    }
    .links a {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      padding: 0 16px;
      border-radius: 999px;
      text-decoration: none;
      font-weight: 700;
      border: 1px solid var(--line);
      background: rgba(255,255,255,0.78);
      color: var(--ink);
    }
    .links a.primary {
      color: #fff;
      border-color: transparent;
      background: linear-gradient(135deg, var(--accent), #b84b22);
    }
    #redoc-wrap {
      overflow: hidden;
      border-radius: 28px;
      border: 1px solid var(--line);
      box-shadow: var(--shadow);
      background: #fff;
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="topbar">
      <p>Alternative Layout</p>
      <h1>Neko-Love API with Redoc</h1>
      <div class="sub">Redoc gives you the more classic three-panel API reference feel. Keep this route if you want a familiar, polished spec browser alongside the more product-feeling Scalar page.</div>
      <div class="links">
        <a class="primary" href="/docs">Open Scalar View</a>
        <a href="{{ .OpenAPIPath }}" target="_blank" rel="noreferrer">OpenAPI JSON</a>
      </div>
    </section>

    <section id="redoc-wrap">
      <redoc spec-url="{{ .OpenAPIPath }}"></redoc>
    </section>
  </div>

  <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>`
