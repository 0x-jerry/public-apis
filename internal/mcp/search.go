package mcp

import (
	"context"

	"public-apis/internal/lib/search"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterSearch(mcpServer *server.MCPServer, h *Handler) {
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
}

func (h *Handler) handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	q := request.GetString("q", "")
	if q == "" {
		return mcp.NewToolResultError("Missing q parameter"), nil
	}
	engine := request.GetString("engine", "")
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
