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

Fetches a URL, strips CSS/JS, and converts the HTML content to Markdown with a YAML frontmatter block containing page metadata.

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

The MCP endpoint is available at `https://public-apis.0x-jerry.deno.net/mcp`.

#### Usage with Claude Desktop / MCP Clients

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "public-apis": {
      "type": "streamableHttp",
      "url": "https://public-apis.0x-jerry.deno.net/mcp"
    }
  }
}
```

#### Tools

##### `html-to-markdown`

Fetches a URL, extracts page metadata, and converts the HTML content to Markdown with YAML frontmatter. Behaves identically to the [GET /html/to-markdown](#html-to-markdown) endpoint.

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

### RSS Feeds

#### Machine Intelligence RSS
```
GET /rss/jiqizhixin.com
```

Proxies and reformats the RSS feed from [jiqizhixin.com](https://www.jiqizhixin.com/rss) into a clean RSS 2.0 feed.

**Response:** `application/rss+xml`

---

#### Ollama Blog RSS
```
GET /rss/ollama.com/blog
```

Scrapes the [Ollama blog](https://ollama.com/blog) and generates an RSS 2.0 feed from the page content.

**Response:** `application/rss+xml`

---

## Project Structure

```
src/
├── _app.ts              # Main Hono app instance with error handling
├── _index.tsx            # Homepage route (renders readme.md)
├── index.ts              # Root router aggregator
├── _libs/
│   └── img-to-ascii/
│       └── index.ts       # Shared imgToAscii() utility
  ├── html/
  │   ├── _app.ts           # Sub-app for /html routes
  │   ├── convert.ts         # Shared HTML-to-Markdown conversion logic
  │   ├── to-markdown.ts    # GET /html/to-markdown
  │   └── index.ts          # Router aggregator
  ├── img-to-ascii/
  │   ├── _app.ts           # Sub-app for /img-to-ascii routes
  │   ├── convert.ts        # GET /img-to-ascii/convert
  │   └── index.ts          # Router aggregator
  ├── link/
  │   ├── _app.ts           # Sub-app for /link routes
  │   ├── metadata.ts       # GET /link/metadata
  │   └── index.ts          # Router aggregator
  ├── mcp/
  │   ├── _app.ts           # MCP Hono app + server + transport
  │   ├── html-to-markdown.ts # MCP tool: html-to-markdown
  │   ├── img-to-ascii.ts   # MCP tool: img-to-ascii
  │   ├── ocr.ts            # MCP tool: ocr
  │   └── index.ts          # Router aggregator
├── qr/
│   ├── _app.ts           # Sub-app for /qr routes
│   ├── generate.ts       # GET /qr/generate
│   ├── scan.ts           # POST /qr/scan
│   └── index.ts          # Router aggregator
├── rss/
│   ├── _app.ts           # Sub-app for /rss routes
│   ├── _utils.ts         # HTTP fetch utility
│   ├── jiqizhixin.ts     # GET /rss/jiqizhixin.com
│   ├── ollama.ts         # GET /rss/ollama.com/blog
│   └── index.ts          # Router aggregator
└── svg/
    ├── _app.ts           # Sub-app for /svg routes
    ├── png.ts            # POST /svg/png
    └── index.ts          # Router aggregator
```

## Tech Stack

- **Runtime:** Node.js via [tsx](https://github.com/privatenumber/tsx)
- **Framework:** [Hono](https://hono.dev)
- **Server:** [@hono/node-server](https://github.com/honojs/node-server)
- **Package Manager:** Bun
- **Language:** TypeScript
