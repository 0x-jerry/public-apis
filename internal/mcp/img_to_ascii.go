package mcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"public-apis/internal/lib/img2ascii"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterImgToAscii(mcpServer *server.MCPServer, h *Handler) {
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
}

func (h *Handler) handleImgToAscii(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := GetStringArg(request, "url")
	if url == "" {
		return mcp.NewToolResultError("Missing url parameter"), nil
	}
	maxSize := int(GetNumberArg(request, "maxSize"))
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
