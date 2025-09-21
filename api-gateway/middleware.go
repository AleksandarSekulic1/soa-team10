package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
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

// loggingMiddleware logs incoming requests with tracing
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Generate simple trace ID
		traceID := generateTraceID()
		
		// Add trace ID to response header
		w.Header().Set("X-Trace-ID", traceID)
		
		// Increment metrics
		atomic.AddInt64(&activeRequests, 1)
		atomic.AddInt64(&requestCounter, 1)
		
		log.Printf("🔀 [GATEWAY] [TRACE:%s] %s %s from %s", traceID, r.Method, r.URL.Path, r.RemoteAddr)
		log.Printf("📨 [GATEWAY] [TRACE:%s] Headers: User-Agent=%s, Origin=%s", traceID, r.Header.Get("User-Agent"), r.Header.Get("Origin"))
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		atomic.AddInt64(&activeRequests, -1)
		
		// Store duration for metrics
		requestDurations = append(requestDurations, float64(duration.Nanoseconds())/1e6)
		if len(requestDurations) > 100 { // Keep only last 100 measurements
			requestDurations = requestDurations[1:]
		}
		
		log.Printf("✅ [GATEWAY] [TRACE:%s] Response sent for %s %s (Duration: %v)", traceID, r.Method, r.URL.Path, duration)
	})
}

// Generate simple trace ID
func generateTraceID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
