package mcp

import (
	"net/http"

	"public-apis/internal/config"
	mcplib "public-apis/internal/mcp"

	"github.com/go-chi/chi/v5"
	"github.com/mark3labs/mcp-go/server"
)

func New(cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	mcpServer := mcplib.NewServer(cfg)

	streamableServer := server.NewStreamableHTTPServer(mcpServer, server.WithEndpointPath("/"))
	r.Mount("/", streamableServer)

	return r
}
