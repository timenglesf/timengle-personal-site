package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
)

func (app *application) commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			fmt.Sprintf("default-src 'self'; "+
				"style-src 'self' 'unsafe-inline' fonts.googleapis.com %s; "+
				"font-src fonts.gstatic.com; "+
				"img-src 'self' data: https: %s; "+
				"media-src 'self' data: https: %s; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' www.youtube.com s.ytimg.com %s; "+
				"frame-src 'self' www.youtube.com",
				app.cfg.objectStorage.objectStorageURL,
				app.cfg.objectStorage.objectStorageURL,
				app.cfg.objectStorage.objectStorageURL,
				app.cfg.objectStorage.objectStorageURL))

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.ExemptPath("/admin/logout")
	csrfHandler.ExemptGlob("/posts/update/*")
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   app.cfg.secureCookies,
		Path:     "/",
	})
	return csrfHandler
}

func (app *application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAdmin(r) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) staticCacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Add("Cache-Control", "public, max-age=1800")
		}
		next.ServeHTTP(w, r)
	})
}
