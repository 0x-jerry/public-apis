package qr

import (
	"fmt"
	"net/http"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/go-chi/chi/v5"
)

func New() http.Handler {
	r := chi.NewRouter()
	r.Get("/generate", generate)
	r.Post("/scan", scan)
	return r
}

func generate(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("c")
	if content == "" {
		content = r.URL.Query().Get("content")
	}
	if content == "" {
		http.Error(w, "Missing c parameter", http.StatusBadRequest)
		return
	}

	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	svg := renderSVG(qr)
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}

func renderSVG(qr *qrcode.QRCode) string {
	bitmap := qr.Bitmap()
	size := len(bitmap)
	scale := 10

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`,
		size*scale, size*scale, size*scale, size*scale))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white"/>`, size*scale, size*scale))

	for y := range size {
		for x := range size {
			if bitmap[y][x] {
				sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="black"/>`,
					x*scale, y*scale, scale, scale))
			}
		}
	}

	sb.WriteString("</svg>")
	return sb.String()
}
