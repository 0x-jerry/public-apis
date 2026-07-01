package mcp

import (
	"fmt"

	"public-apis/internal/config"
	"public-apis/internal/lib/browser"

	"github.com/mark3labs/mcp-go/server"
)

type Handler struct {
	cfg     *config.Config
	browser *browser.Manager
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg:     cfg,
		browser: browser.New(cfg),
	}
}

func NewServer(cfg *config.Config) *server.MCPServer {
	h := NewHandler(cfg)

	s := server.NewMCPServer(
		"public-apis",
		"1.0.0",
	)

	RegisterHTMLToMarkdown(s, h)
	RegisterImgToAscii(s, h)
	RegisterSearch(s, h)
	RegisterOCR(s, h)

	return s
}

func Truncate(s string, offset, limit int) string {
	if offset > len(s) {
		return ""
	}
	result := s[offset:]
	if limit > 0 && len(result) > limit {
		result = result[:limit] + "\n\n... (truncated)"
	}
	return result
}

func FormatBbox(bbox []float64) string {
	if len(bbox) < 4 {
		return ""
	}
	x1 := int(bbox[0] + 0.5)
	y1 := int(bbox[1] + 0.5)
	x2 := int(bbox[2] + 0.5)
	y2 := int(bbox[3] + 0.5)
	return fmt.Sprintf("(%d,%d)-(%d,%d)", x1, y1, x2, y2)
}
