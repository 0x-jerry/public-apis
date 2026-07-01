package index

import (
	"bytes"
	"net/http"

	publicapis "public-apis"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(extension.Table))
	if err := md.Convert(publicapis.Readme, &buf); err != nil {
		http.Error(w, "failed to render", http.StatusInternalServerError)
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Some API</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-markdown-css/github-markdown.css">
</head>
<body>
<div data-color-mode="auto" data-dark-theme="auto" class="markdown-body">
<div style="width: 768px; margin: auto; padding: 32px 0;">` + buf.String() + `</div>
</div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
