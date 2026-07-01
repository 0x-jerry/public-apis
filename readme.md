# Public APIs

A collection of useful public utility APIs built with Go and [chi](https://github.com/go-chi/chi).

## Getting Started

```bash
# Install dependencies
go mod download

# Start development server
go run ./cmd/server

# Build
go build -o server ./cmd/server
```

The server starts at `http://localhost:3000`.

## Environment Variables

| Variable            | Default                 | Description                                       |
|--------------------|-------------------------|---------------------------------------------------|
| `AUTH_ENABLED`       | —                       | Set to `true` or `1` to require authentication     |
| `AUTH_TOKEN`         | —                       | Bearer token required when auth is enabled          |
| `BROWSER_WS`         | `ws://127.0.0.1:9222`   | Chrome DevTools WebSocket endpoint for headless browser   |
| `BROWSER_WS_ENABLED` | —                       | Set to `true` or `1` to use headless browser; falls back to plain fetch otherwise |
| `STORAGE_PATH`       | `files`                 | Directory for uploaded files                        |
| `OCR_TOKEN`          | —                       | PaddleOCR API authentication token                  |

## Docker

```bash
# Start Lightpanda headless browser
docker compose -f compose.dev.yaml up -d

# Build and run
docker build -t public-apis .
docker run -p 3000:3000 --env-file .env public-apis
```

## API Reference

### QR Code

#### Generate QR Code
```
GET /qr/generate?c=<content>
```

Generates a QR code as an SVG image.

**Query Parameters:**

| Param | Type   | Description          |
|-------|--------|----------------------|
| `c`   | string | Content to encode    |

**Response:** `image/svg+xml`

---

#### Scan QR Code
```
POST /qr/scan
```

Decodes a QR code from an uploaded image file.

**Request:** `multipart/form-data`

**Response:** `application/json`

---

### HTML to Markdown

#### Convert HTML to Markdown
```
GET /html/to-markdown?url=<url>
```

Fetches a URL using a headless browser (via rod), strips CSS/JS, and converts the rendered HTML content to Markdown with a YAML frontmatter block containing page metadata.

---

### Image to ASCII

#### Convert Image to ASCII Art
```
GET /img-to-ascii/convert?url=<url>&maxSize=<maxSize>
```

Fetches a remote image and converts it to ASCII art.

---

### MCP Server

This project exposes an [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server using the SSE transport.

The MCP endpoint is available at `http://localhost:3000/mcp/sse`.

#### Tools

- **html-to-markdown** — Fetch URL and convert to Markdown
- **img-to-ascii** — Convert remote image to ASCII art
- **search** — Search the web using Bing or DuckDuckGo
- **ocr** — Extract text from documents using PaddleOCR

---

### File Upload

#### Upload Page
```
GET /upload
```

A web interface for uploading files with drag-and-drop support.

#### Upload File
```
POST /upload/file
```

Upload a file. Requires authentication.

#### Download File
```
GET /upload/file/{name}
```

Download a previously uploaded file by name.

---

## Tech Stack

- **Runtime:** Go
- **Framework:** [chi](https://github.com/go-chi/chi)
- **Browser:** [rod](https://github.com/go-rod/rod)
- **HTML→Markdown:** [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown)
- **Readability:** [go-readability](https://codeberg.org/readeck/go-readability)
- **MCP:** [mcp-go](https://github.com/mark3labs/mcp-go)
