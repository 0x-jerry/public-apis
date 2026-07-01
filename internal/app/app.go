package app

import (
	"encoding/json"
	"log"
	"net/http"

	"public-apis/internal/auth"
	"public-apis/internal/config"
	"public-apis/internal/handler/html2md"
	"public-apis/internal/handler/img2ascii"
	"public-apis/internal/handler/index"
	"public-apis/internal/handler/mcp"
	"public-apis/internal/handler/qr"
	"public-apis/internal/handler/upload"
	"public-apis/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func New(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.TrimTrailingSlash)
	r.Use(middleware.Recover)
	r.Use(errorHandler)

	r.Mount("/qr", qr.New())
	r.Mount("/html", html2md.New(cfg))
	r.Mount("/img-to-ascii", img2ascii.New())
	r.Mount("/mcp", mcp.New(cfg))

	r.With(auth.Required(cfg)).Mount("/upload", upload.New(cfg))

	r.Get("/", index.Handler)

	return r
}

func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				msg := "Bad request"
				if e, ok := err.(error); ok {
					msg = e.Error()
				} else if s, ok := err.(string); ok {
					msg = s
				}
				json.NewEncoder(w).Encode(map[string]string{"message": msg})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
