package proxy

import (
	"net"
	"net/http"
	"net/url"

	"public-apis/internal/config"
	"public-apis/internal/lib/browser"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	browser *browser.Manager
}

func New(cfg *config.Config) http.Handler {
	h := &Handler{browser: browser.New(cfg)}

	r := chi.NewRouter()
	r.Get("/", h.proxy)
	return r
}

func (h *Handler) proxy(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("url")
	if raw == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	target, err := url.Parse(raw)
	if err != nil || target.Host == "" ||
		(target.Scheme != "http" && target.Scheme != "https") {
		http.Error(w, "Invalid url parameter", http.StatusBadRequest)
		return
	}

	// Resolve the target and reject private, loopback, link-local and
	// multicast addresses before handing it to the browser.
	blocked, resolveErr := blockedTarget(target)
	if resolveErr != nil {
		http.Error(w, "Failed to resolve host: "+resolveErr.Error(), http.StatusBadGateway)
		return
	}
	if blocked != nil {
		http.Error(w, "blocked address: "+blocked.String(), http.StatusBadGateway)
		return
	}

	result, err := h.browser.Fetch(r.Context(), target.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if result.ContentType != "" {
		w.Header().Set("Content-Type", result.ContentType)
	}
	if result.StatusCode != 0 {
		w.WriteHeader(result.StatusCode)
	}
	w.Write(result.Body)
}

// blockedTarget resolves the URL host and returns the first IP that maps to a
// private, loopback, link-local or multicast address. It returns a nil IP when
// the host is allowed, and an error when the host cannot be resolved.
func blockedTarget(target *url.URL) (net.IP, error) {
	ips, err := net.LookupIP(target.Hostname())
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if isBlockedIP(ip) {
			return ip, nil
		}
	}

	return nil, nil
}

func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() ||
		ip.IsUnspecified() ||
		ip.IsMulticast()
}
