package upload

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"public-apis/internal/config"
	"public-apis/internal/lib/storage"
	"github.com/go-chi/chi/v5"
)

func New(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Get("/", uploadPage)

	h := &handler{cfg: cfg}
	r.Get("/file/{name}", h.download)
	r.Post("/file", h.upload)

	return r
}

type handler struct {
	cfg *config.Config
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.New("upload").Parse(uploadHTML))
	tmpl.Execute(w, nil)
}

func (h *handler) upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Not found file"})
		return
	}
	defer file.Close()

	dir := h.cfg.AbsStoragePath()
	name, size, err := storage.Save(dir, file, header.Filename)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file too large, max 100MB" {
			status = http.StatusRequestEntityTooLarge
		}
		writeJSON(w, status, map[string]string{"message": err.Error()})
		return
	}

	storage.Evict(dir)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name": name,
		"size": size,
		"type": header.Header.Get("Content-Type"),
	})
}

func (h *handler) download(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	dir := h.cfg.AbsStoragePath()

	filePath := filepath.Join(dir, name)
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(dir)) {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
