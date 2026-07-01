package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"public-apis/internal/config"
)

type contextKey string

const TokenKey contextKey = "auth_token"

func Required(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.IsAuthConfigured() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(map[string]string{"message": "Auth is not enabled"})
				return
			}
			if !verifyToken(w, r, cfg.AuthToken) {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Optional(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.IsAuthConfigured() {
				if !verifyToken(w, r, cfg.AuthToken) {
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func verifyToken(w http.ResponseWriter, r *http.Request, expected string) bool {
	header := r.Header.Get("Authorization")
	token := strings.TrimPrefix(header, "Bearer ")
	if header == token {
		token = ""
	}

	if token != expected {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return false
	}
	return true
}

func GetToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	return strings.TrimPrefix(h, "Bearer ")
}
