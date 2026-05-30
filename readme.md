# Public APIs

A collection of useful public utility APIs built with [Hono](https://hono.dev) running on Node.js.

## Getting Started

```bash
# Install dependencies (requires Bun as package manager)
bun install

# Start development server
bun dev

# Type check
bun check
```

The server starts at `http://localhost:3000`.

## Environment Variables

| Variable            | Default                 | Description                                       |
|--------------------|-------------------------|---------------------------------------------------|
| `AUTH_ENABLED`       | —                       | Set to `true` or `1` to require authentication     |
| `AUTH_TOKEN`         | —                       | Bearer token required when auth is enabled          |
| `BROWSER_WS`         | `ws://127.0.0.1:9222`   | Chrome DevTools WebSocket endpoint for Puppeteer   |
| `BROWSER_WS_ENABLED` | —                       | Set to `true` or `1` to use headless browser; falls back to plain fetch otherwise |

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

**Example:**
```
GET /qr/generate?c=https://example.com
```

---

#### Scan QR Code
```
POST /qr/scan
```

Decodes a QR code from an uploaded image file.

**Request:** `multipart/form-data`

| Field  | Type | Description              |
|--------|------|--------------------------|
| `file` | File | Image containing QR code |

**Response:** `application/json`
```json
{
  "version": 4,
  "data": "decoded content",
  "location": { ... }
}
```

---

### Link Preview

#### Get Link Metadata
```
GET /link/metadata?url=<url>
```

Fetches metadata (title, description, image, etc.) for any URL.

**Query Parameters:**

| Param | Type   | Description     |
|-------|--------|-----------------|
| `url` | string | URL to preview  |

**Response:** `application/json`

---

### HTML to Markdown

#### Convert HTML to Markdown
```
GET /html/to-markdown?url=<url>&limit=<limit>
```

Fetches a URL using a headless Chrome browser (via Puppeteer), strips CSS/JS, and converts the rendered HTML content to Markdown with a YAML frontmatter block containing page metadata. Requires a Chrome/Chromium instance with remote debugging enabled on `BROWSER_WS`.

**Query Parameters:**

| Param   | Type   | Default | Description                                |
|---------|--------|---------|--------------------------------------------|
| `url`   | string | —       | URL to convert                             |
| `limit` | number | 50000   | Maximum character length of the output     |

**Response:** `text/plain` — Markdown with YAML frontmatter.

**Example output:**
```markdown
---
url: "https://example.com"
title: "Example Page"
description: "A description of the page"
og:title: "Open Graph Title"
language: "en"
---

# Main Heading

Page content converted to Markdown...
```

---

### Image to ASCII

#### Convert Image to ASCII Art
```
GET /img-to-ascii/convert?url=<url>&maxSize=<maxSize>
```

Fetches a remote image and converts it to ASCII art.

**Query Parameters:**

| Param     | Type   | Default | Description                                      |
|-----------|--------|---------|--------------------------------------------------|
| `url`     | string | —       | URL of the image to convert                      |
| `maxSize` | number | 80      | Maximum width/height in characters               |

**Response:** `text/plain`

**Example:**
```
GET /img-to-ascii/convert?url=https://example.com/image.png&maxSize=60
```

---

### MCP Server

This project exposes an [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server using the [Streamable HTTP](https://modelcontextprotocol.io/specification/2025-06-18/basic/transports#streamable-http) transport.

The MCP endpoint is available at `http://localhost:3000/mcp`.

#### Usage with Claude Desktop / MCP Clients

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "public-apis": {
      "type": "streamableHttp",
      "url": "http://localhost:3000/mcp"
    }
  }
}
```

#### Tools

##### `html-to-markdown`

Fetches a URL using a headless Chrome browser (via Puppeteer), extracts page metadata, and converts the HTML content to Markdown with YAML frontmatter. Behaves identically to the [GET /html/to-markdown](#html-to-markdown) endpoint.

**Input:**

| Field   | Type   | Default | Description                                    |
|---------|--------|---------|------------------------------------------------|
| `url`   | string | —       | URL of the webpage to convert                  |
| `limit` | number | 50000   | Maximum character length of the output         |

**Output:** `text` — Markdown with YAML frontmatter.

##### `img-to-ascii`

Converts a remote image to ASCII art. Behaves identically to the [GET /img-to-ascii/convert](#image-to-ascii) endpoint.

**Input:**

| Field     | Type   | Default | Description                                      |
|-----------|--------|---------|--------------------------------------------------|
| `url`     | string | —       | URL of the image to convert                      |
| `maxSize` | number | 80      | Maximum width/height in characters               |

**Output:** `text` — ASCII art string.

##### `ocr`

Extracts text and layout from a document/image using PaddleOCR. Provide a URL to the document/image file.

**Input:**

| Field | Type   | Description                          |
|-------|--------|--------------------------------------|
| `url` | string | URL of the document/image to OCR     |

**Output:** `text` — Extracted Markdown text.

**Headers:**

| Header        | Description                          |
|---------------|--------------------------------------|
| `x-ocr-token` | API authentication token (required)  |
| `x-ocr-model` | Model name (default: `PaddleOCR-VL-1.6`) |

---

### SVG to PNG

#### Convert SVG to PNG
```
POST /svg/png
```

Converts an SVG to a PNG image.

**Request:** `application/json`

| Field | Type   | Description                                |
|-------|--------|--------------------------------------------|
| `svg` | string | Raw SVG markup                             |
| `url` | string | URL to fetch SVG from (alternative to svg) |

> Provide either `svg` or `url`.

**Response:** `image/png`

---

## Tech Stack

- **Runtime:** Node.js via [tsx](https://github.com/privatenumber/tsx)
- **Framework:** [Hono](https://hono.dev)
- **Server:** [@hono/node-server](https://github.com/honojs/node-server)
- **Package Manager:** Bun
- **Language:** TypeScript
