package html2md

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"public-apis/internal/config"
	"public-apis/internal/lib/browser"
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
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/markdown")

	resp, clientErr := http.DefaultClient.Do(req)
	if clientErr == nil {
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
				opts := &libhtml2md.Options{Mode: "readable", URL: url}
				return libhtml2md.Convert(htmlStr, opts)
			}

			browserHTML, browserErr := h.browser.FetchHTML(ctx, url)
			if browserErr == nil && browserHTML != "" {
				opts := &libhtml2md.Options{Mode: "readable", URL: url}
				return libhtml2md.Convert(browserHTML, opts)
			}

			opts := &libhtml2md.Options{Mode: "readable", URL: url}
			return libhtml2md.Convert(htmlStr, opts)
		}
	}

	browserHTML, err := h.browser.FetchHTML(ctx, url)
	if err != nil {
		req2, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		req2.Header.Set("Accept", "text/html")
		resp2, err2 := http.DefaultClient.Do(req2)
		if err2 != nil {
			return "", fmt.Errorf("failed to fetch url: %w", err2)
		}
		defer resp2.Body.Close()
		data, _ := io.ReadAll(resp2.Body)
		opts := &libhtml2md.Options{Mode: "readable", URL: url}
		return libhtml2md.Convert(string(data), opts)
	}

	opts := &libhtml2md.Options{Mode: "readable", URL: url}
	return libhtml2md.Convert(browserHTML, opts)
}
