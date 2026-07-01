package img2ascii

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"public-apis/internal/lib/img2ascii"
	"github.com/go-chi/chi/v5"
)

func New() http.Handler {
	r := chi.NewRouter()
	r.Get("/convert", convert)
	return r
}

func convert(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	maxSize := 80
	if v := r.URL.Query().Get("maxSize"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxSize = n
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch image: %v", err), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		http.Error(w, fmt.Sprintf("Failed to fetch image: %d %s", resp.StatusCode, resp.Status), http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}

	opts := &img2ascii.Options{MaxWidth: maxSize, MaxHeight: maxSize}
	ascii, err := img2ascii.Convert(bytes.NewReader(data), opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(ascii))
}
