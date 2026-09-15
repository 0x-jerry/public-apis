package middleware

import (
	"bytes"
	"log"
	"net/http"
	"strings"
)

// maxErrorDetail bounds how much of a failed response body is logged.
const maxErrorDetail = 1024 * 2

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &errorWriter{ResponseWriter: w}
		next.ServeHTTP(ww, r)
		if detail := strings.TrimSpace(ww.detail.String()); detail != "" {
			log.Printf("%s %s %d %s", r.Method, r.URL.Path, ww.status, detail)
			return
		}
		log.Printf("%s %s %d", r.Method, r.URL.Path, ww.status)
	})
}

// errorWriter records the status code and keeps the beginning of the body of
// error responses, which is where http.Error and the app's error handler put
// the message.
type errorWriter struct {
	http.ResponseWriter
	status int
	detail bytes.Buffer
}

func (w *errorWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *errorWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	if w.status >= http.StatusBadRequest && w.detail.Len() < maxErrorDetail {
		w.detail.Write(b[:min(len(b), maxErrorDetail-w.detail.Len())])
	}
	return w.ResponseWriter.Write(b)
}

func (w *errorWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

func TrimTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
