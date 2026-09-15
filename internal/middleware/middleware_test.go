package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"success", http.StatusOK, "ok", "GET /x 200"},
		{"error", http.StatusBadGateway, "boom: bad thing", "GET /x 502 boom: bad thing"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			log.SetFlags(0)
			defer func() {
				log.SetOutput(log.Writer())
				log.SetFlags(log.LstdFlags)
			}()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			})
			Logger(next).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/x", nil))

			if got := strings.TrimSpace(buf.String()); got != tc.want {
				t.Fatalf("log = %q, want %q", got, tc.want)
			}
		})
	}
}
