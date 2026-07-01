package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"public-apis/internal/config"
	"public-apis/internal/lib/browser"
	libhtml2md "public-apis/internal/lib/html2md"
	"public-apis/internal/lib/img2ascii"
	"public-apis/internal/lib/paddleocr"
	"public-apis/internal/lib/search"
	"github.com/go-chi/chi/v5"
	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/net/html"
)

func New(cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	h := &handler{
		cfg:     cfg,
		browser: browser.New(cfg),
	}

	mcpServer := server.NewMCPServer(
		"public-apis",
		"1.0.0",
	)

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

	mcpServer.AddTool(mcp.NewTool("img-to-ascii",
		mcp.WithDescription("Convert a remote image to ASCII art. Provide a URL to the image."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL of the image to convert to ASCII"),
		),
		mcp.WithNumber("maxSize",
			mcp.Description("Maximum width/height in characters (default: 80)"),
		),
	), h.handleImgToAscii)

	mcpServer.AddTool(mcp.NewTool("search",
		mcp.WithDescription("Search the web using Bing or DuckDuckGo"),
		mcp.WithString("q",
			mcp.Required(),
			mcp.Description("Search query"),
		),
		mcp.WithString("engine",
			mcp.Description("Search engine to use (default: duckduckgo)"),
		),
	), h.handleSearch)

	mcpServer.AddTool(mcp.NewTool("ocr",
		mcp.WithDescription("Extract text and layout from a document/image using OCR. Provide a URL to the document/image file."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL of the document/image to OCR"),
		),
		mcp.WithString("mode",
			mcp.Description("Output mode: 'layout' includes position/label annotations, 'markdown' returns plain markdown without layout info (default: 'markdown')"),
		),
		mcp.WithNumber("offset",
			mcp.Description("Character index to start output from (default: 0)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum character length of the output (default: 30000)"),
		),
	), h.handleOCR)

	sseServer := server.NewSSEServer(mcpServer)

	r.Handle("/sse", sseServer)
	r.Handle("/message", sseServer)

	return r
}

type handler struct {
	cfg     *config.Config
	browser *browser.Manager
}

func (h *handler) handleHTMLToMarkdown(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := getStringArg(request, "url")
	if url == "" {
		return mcp.NewToolResultError("Missing url parameter"), nil
	}
	offset := int(getNumberArg(request, "offset"))
	limit := int(getNumberArg(request, "limit"))
	if limit == 0 {
		limit = 30000
	}
	mode := getStringArg(request, "mode")
	if mode == "" {
		mode = "readable"
	}

	full, err := h.htmlToMarkdown(ctx, url, mode)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := truncate(full, offset, limit)
	return mcp.NewToolResultText(result), nil
}

func (h *handler) handleImgToAscii(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := getStringArg(request, "url")
	if url == "" {
		return mcp.NewToolResultError("Missing url parameter"), nil
	}
	maxSize := int(getNumberArg(request, "maxSize"))
	if maxSize == 0 {
		maxSize = 80
	}

	resp, err := http.Get(url)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch image: %v", err)), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch image: %d %s", resp.StatusCode, resp.Status)), nil
	}

	data, _ := io.ReadAll(resp.Body)
	opts := &img2ascii.Options{MaxWidth: maxSize, MaxHeight: maxSize}
	ascii, err := img2ascii.Convert(bytes.NewReader(data), opts)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(ascii), nil
}

func (h *handler) handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	q := getStringArg(request, "q")
	if q == "" {
		return mcp.NewToolResultError("Missing q parameter"), nil
	}
	engine := getStringArg(request, "engine")
	if engine == "" {
		engine = "duckduckgo"
	}

	var resp *search.Response
	var err error

	if engine == "bing" {
		resp, err = search.Bing(ctx, h.browser, q)
	} else {
		resp, err = search.DuckDuckGo(q)
	}

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(search.ToMarkdown(resp)), nil
}

func (h *handler) handleOCR(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := getStringArg(request, "url")
	if url == "" {
		return mcp.NewToolResultError("Missing url parameter"), nil
	}
	mode := getStringArg(request, "mode")
	if mode == "" {
		mode = "markdown"
	}
	offset := int(getNumberArg(request, "offset"))
	limit := int(getNumberArg(request, "limit"))
	if limit == 0 {
		limit = 30000
	}

	token := h.cfg.OCRToken
	if token == "" {
		return mcp.NewToolResultError("Missing x-ocr-token header"), nil
	}

	results, err := paddleocr.Run(url, paddleocr.Options{
		Token:   token,
		Model:   "PaddleOCR-VL-1.6",
		Timeout: 0,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var pages []string
	for _, item := range results {
		for idx, page := range item.Result.LayoutParsingResults {
			if mode == "markdown" {
				pages = append(pages, page.Markdown.Text)
				_ = idx
			} else {
				w := page.PrunedResult.Width
				height := page.PrunedResult.Height
				header := fmt.Sprintf("## Page %d (%d×%d)\n\n", idx+1, w, height)

				var blocks []string
				for _, block := range page.PrunedResult.ParsingResList {
					pos := formatBbox(block.BlockBbox)
					blocks = append(blocks, fmt.Sprintf("<!-- [%s] %s -->\n%s", block.BlockLabel, pos, block.BlockContent))
				}
				pages = append(pages, header+strings.Join(blocks, "\n\n"))
			}
		}
	}
	full := strings.Join(pages, "\n\n")

	result := truncate(full, offset, limit)
	return mcp.NewToolResultText(result), nil
}

func (h *handler) htmlToMarkdown(ctx context.Context, url string, mode string) (string, error) {
	resp, err := http.Get(url)
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
				opts := &libhtml2md.Options{Mode: mode}
				return libhtml2md.Convert(htmlStr, opts)
			}

			browserHTML, browserErr := h.browser.FetchHTML(ctx, url)
			if browserErr == nil && browserHTML != "" {
				opts := &libhtml2md.Options{Mode: mode}
				return libhtml2md.Convert(browserHTML, opts)
			}

			opts := &libhtml2md.Options{Mode: mode}
			return libhtml2md.Convert(htmlStr, opts)
		}
	}

	browserHTML, err := h.browser.FetchHTML(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch url: %w", err)
	}

	opts := &libhtml2md.Options{Mode: mode}
	return libhtml2md.Convert(browserHTML, opts)
}

func getArgs(request mcp.CallToolRequest) map[string]interface{} {
	if m, ok := request.Params.Arguments.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func getStringArg(request mcp.CallToolRequest, key string) string {
	args := getArgs(request)
	if args == nil {
		return ""
	}
	if v, ok := args[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		}
	}
	return ""
}

func getNumberArg(request mcp.CallToolRequest, key string) float64 {
	args := getArgs(request)
	if args == nil {
		return 0
	}
	if v, ok := args[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case json.Number:
			f, _ := val.Float64()
			return f
		}
	}
	return 0
}

func truncate(s string, offset, limit int) string {
	if offset > len(s) {
		return ""
	}
	result := s[offset:]
	if limit > 0 && len(result) > limit {
		result = result[:limit] + "\n\n... (truncated)"
	}
	return result
}

func formatBbox(bbox []float64) string {
	if len(bbox) < 4 {
		return ""
	}
	x1 := int(bbox[0] + 0.5)
	y1 := int(bbox[1] + 0.5)
	x2 := int(bbox[2] + 0.5)
	y2 := int(bbox[3] + 0.5)
	return fmt.Sprintf("(%d,%d)-(%d,%d)", x1, y1, x2, y2)
}
