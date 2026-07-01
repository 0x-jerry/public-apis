package mcp

import (
	"encoding/json"
	"fmt"

	"public-apis/internal/config"
	"public-apis/internal/lib/browser"

	"github.com/mark3labs/mcp-go/mcp"
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

func getArgs(request mcp.CallToolRequest) map[string]interface{} {
	if m, ok := request.Params.Arguments.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func GetStringArg(request mcp.CallToolRequest, key string) string {
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

func GetNumberArg(request mcp.CallToolRequest, key string) float64 {
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
