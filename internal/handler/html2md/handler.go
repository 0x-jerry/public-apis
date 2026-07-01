package html2md

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"public-apis/internal/config"
	"public-apis/internal/lib/browser"
	"public-apis/internal/lib/fetch"
	libhtml2md "public-apis/internal/lib/html2md"
	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/html"
)

type Handler struct {
	cfg     *config.Config
	browser *browser.Manager
}

func New(cfg *config.Config) http.Handler {
	h := &Handler{
		cfg:     cfg,
		browser: browser.New(cfg),
	}

	r := chi.NewRouter()
	r.Get("/to-markdown", h.toMarkdown)
	return r
}

func (h *Handler) toMarkdown(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	result, err := h.convert(r.Context(), url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}

func (h *Handler) convert(ctx context.Context, url string) (string, error) {
	fr, fetchErr := fetch.Fetch(url, "text/markdown")
	opts := &libhtml2md.Options{Mode: "readable", URL: url}

	if fetchErr == nil {
		if fr.StatusCode < 400 && strings.Contains(fr.ContentType, "text/markdown") {
			return fr.Body, nil
		}

		if fr.StatusCode < 400 && strings.Contains(fr.ContentType, "text/html") {
			doc, parseErr := html.Parse(strings.NewReader(fr.Body))
			if parseErr != nil {
				return "", parseErr
			}

			if readability.CheckDocument(doc) {
				return libhtml2md.Convert(fr.Body, opts)
			}

			browserHTML, browserErr := h.browser.FetchHTML(ctx, url)
			if browserErr == nil && browserHTML != "" {
				return libhtml2md.Convert(browserHTML, opts)
			}

			return libhtml2md.Convert(fr.Body, opts)
		}
	}

	browserHTML, err := h.browser.FetchHTML(ctx, url)
	if err != nil {
		fr2, fetchErr2 := fetch.Fetch(url, "text/html")
		if fetchErr2 != nil {
			return "", fmt.Errorf("failed to fetch url: %w", fetchErr2)
		}
		return libhtml2md.Convert(fr2.Body, opts)
	}

	return libhtml2md.Convert(browserHTML, opts)
}
