package main

import (
	"fmt"
	"net/http"
	"snippetbox/internal/config"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "GO")

		next.ServeHTTP(w, r)
	})
}

func logRequest(app *config.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			url    = r.URL.RequestURI()
		)

		app.Logger.Info("Received request", "ip", ip, "proto", proto, "method", method, "url", url)

		next.ServeHTTP(w, r)
	})
}

func recoverPanic(app *config.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				serverError(app, r, fmt.Errorf("%s", err))
				clientError(http.StatusInternalServerError, w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
