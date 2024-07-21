package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_time := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// call the actual handler
		nextHandler.ServeHTTP(wrapped, r)

		// print log
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, r.RemoteAddr, time.Since(req_time))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// don't redirect the login path itself
		if strings.HasPrefix(r.URL.Path, "/auth/") ||
			strings.HasPrefix(r.URL.Path, "/static/css/") ||
			r.URL.Path == "/favicon.ico" {
			next.ServeHTTP(w, r)
			return
		}

		_, ok := checkSession(r)
		if !ok {
			http.Redirect(w, r, "/auth/login/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
