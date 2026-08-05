package main

import (
	"net/http"
	"strings"
)

var corsAllowedHeaders = strings.Join([]string{
	"Content-Type",
	"Authorization",
	"Connect-Protocol-Version",
	"Connect-Timeout-Ms",
	"Grpc-Timeout",
	"X-Grpc-Web",
	"X-User-Agent",
	"X-User-Id",
}, ", ")

var corsExposedHeaders = strings.Join([]string{
	"Grpc-Status",
	"Grpc-Message",
	"Grpc-Status-Details-Bin",
}, ", ")

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if app.originAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", corsAllowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", corsExposedHeaders)
			w.Header().Set("Access-Control-Max-Age", "7200")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) originAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	if app.config.env == "develop" {
		switch origin {
		case "http://localhost:3000",
			"https://localhost:3000",
			"http://127.0.0.1:3000",
			"https://127.0.0.1:3000":
			return true
		}
		return false
	}
	return origin == "https://"+app.config.ui_address+":"+app.config.ui_port
}
