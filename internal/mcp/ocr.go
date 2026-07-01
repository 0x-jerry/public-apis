package mcp

import (
	"context"
	"fmt"
	"strings"

	"public-apis/internal/lib/paddleocr"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterOCR(mcpServer *server.MCPServer, h *Handler) {
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
}

func (h *Handler) handleOCR(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := GetStringArg(request, "url")
	if url == "" {
		return mcp.NewToolResultError("Missing url parameter"), nil
	}
	mode := GetStringArg(request, "mode")
	if mode == "" {
		mode = "markdown"
	}
	offset := int(GetNumberArg(request, "offset"))
	limit := int(GetNumberArg(request, "limit"))
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
					pos := FormatBbox(block.BlockBbox)
					blocks = append(blocks, fmt.Sprintf("<!-- [%s] %s -->\n%s", block.BlockLabel, pos, block.BlockContent))
				}
				pages = append(pages, header+strings.Join(blocks, "\n\n"))
			}
		}
	}
	full := strings.Join(pages, "\n\n")

	result := Truncate(full, offset, limit)
	return mcp.NewToolResultText(result), nil
}
