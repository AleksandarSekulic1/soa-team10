package main

import (
	"log"
	"net/http"
	"strings"
)

// corsMiddleware handles CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only set CORS headers for non-API routes (like /health)
		// API routes will get CORS headers from ModifyResponse in reverse proxy
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")
			
			// Handle preflight OPTIONS requests for non-API routes
			if r.Method == "OPTIONS" {
				log.Printf("🔄 [CORS] OPTIONS preflight request for %s from %s", r.URL.Path, r.Header.Get("Origin"))
				w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			// For API routes, just log OPTIONS requests but don't handle them here
			if r.Method == "OPTIONS" {
				log.Printf("🔄 [CORS] OPTIONS preflight request for API %s from %s", r.URL.Path, r.Header.Get("Origin"))
			}
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🔀 [GATEWAY] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		log.Printf("📨 [GATEWAY] Headers: User-Agent=%s, Origin=%s", r.Header.Get("User-Agent"), r.Header.Get("Origin"))
		next.ServeHTTP(w, r)
		log.Printf("✅ [GATEWAY] Response sent for %s %s", r.Method, r.URL.Path)
	})
}
