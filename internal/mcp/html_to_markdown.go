package mcp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	libhtml2md "public-apis/internal/lib/html2md"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/net/html"
)

func RegisterHTMLToMarkdown(mcpServer *server.MCPServer, h *Handler) {
	mcpServer.AddTool(mcp.NewTool("html-to-markdown",
		mcp.WithDescription("Fetch a URL, extract metadata and convert HTML to Markdown"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL of the webpage to convert to Markdown"),
		),
		mcp.WithNumber("offset",
			mcp.Description("Character index to start output from (default: 0)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum character length of the output (default: 30000)"),
		),
		mcp.WithString("mode",
			mcp.Description(`Conversion mode: "readable" extracts article content, "full" keeps entire body`),
		),
	), h.handleHTMLToMarkdown)
}

func (h *Handler) handleHTMLToMarkdown(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := request.GetString("url", "")
	if url == "" {
		return mcp.NewToolResultError("Missing url parameter"), nil
	}
	offset := request.GetInt("offset", 0)
	limit := request.GetInt("limit", 0)
	if limit == 0 {
		limit = 30000
	}
	mode := request.GetString("mode", "")
	if mode == "" {
		mode = "readable"
	}

	full, err := h.htmlToMarkdown(ctx, url, mode)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := Truncate(full, offset, limit)
	return mcp.NewToolResultText(result), nil
}

func (h *Handler) htmlToMarkdown(ctx context.Context, url string, mode string) (string, error) {
	resp, err := http.Get(url)
	
	opts := &libhtml2md.Options{Mode: mode, URL: url}
	
	if err == nil {
		defer resp.Body.Close()
		contentType := resp.Header.Get("Content-Type")

		if resp.StatusCode < 400 && strings.Contains(contentType, "text/markdown") {
			data, _ := io.ReadAll(resp.Body)
			return string(data), nil
		}
		if resp.StatusCode < 400 && strings.Contains(contentType, "text/html") {
			data, _ := io.ReadAll(resp.Body)
			htmlStr := string(data)

			doc, parseErr := html.Parse(strings.NewReader(htmlStr))
			if parseErr != nil {
				return "", parseErr
			}

			if readability.CheckDocument(doc) {
				return libhtml2md.Convert(htmlStr, opts)
			}

			browserHTML, browserErr := h.browser.FetchHTML(ctx, url)
			if browserErr == nil && browserHTML != "" {
				return libhtml2md.Convert(browserHTML, opts)
			}

			return libhtml2md.Convert(htmlStr, opts)
		}
	}

	browserHTML, err := h.browser.FetchHTML(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch url: %w", err)
	}

	return libhtml2md.Convert(browserHTML, opts)
}
