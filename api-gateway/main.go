package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Prometheus client for host metrics
	initPrometheusClient()
	
	router := mux.NewRouter()

	// Setup CORS middleware
	router.Use(corsMiddleware)
	
	// Setup logging middleware
	router.Use(loggingMiddleware)

	// Health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")
	
	// Monitoring endpoints
	router.HandleFunc("/metrics", metricsHandler).Methods("GET")
	router.HandleFunc("/monitoring", monitoringHandler).Methods("GET")
	
	// Host metrics endpoints
	router.HandleFunc("/api/monitoring/host/cpu", hostCPUHandler).Methods("GET")
	router.HandleFunc("/api/monitoring/host/memory", hostMemoryHandler).Methods("GET")
	router.HandleFunc("/api/monitoring/host/disk", hostDiskHandler).Methods("GET")
	router.HandleFunc("/api/monitoring/host/network", hostNetworkHandler).Methods("GET")
	router.HandleFunc("/api/monitoring/host/all", hostAllMetricsHandler).Methods("GET")
	
	// CORS preflight handlers for host metrics endpoints
	router.HandleFunc("/api/monitoring/host/cpu", corsPreflightHandler).Methods("OPTIONS")
	router.HandleFunc("/api/monitoring/host/memory", corsPreflightHandler).Methods("OPTIONS")
	router.HandleFunc("/api/monitoring/host/disk", corsPreflightHandler).Methods("OPTIONS")
	router.HandleFunc("/api/monitoring/host/network", corsPreflightHandler).Methods("OPTIONS")
	router.HandleFunc("/api/monitoring/host/all", corsPreflightHandler).Methods("OPTIONS")
	
	// API Monitoring endpoint with CORS support
	router.HandleFunc("/api/monitoring", monitoringHandler).Methods("GET")
	router.HandleFunc("/api/monitoring", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	// Setup routes for different services
	setupStakeholdersRoutes(router)
	setupPaymentsRoutes(router)
	setupBlogRoutes(router)
	setupToursRoutes(router)
	setupEncountersRoutes(router)
	setupFollowerRoutes(router)

	log.Println("🚀 API Gateway starting on port 8000...")
	log.Println("📊 Metrics available at http://localhost:8000/metrics")
	log.Println("📈 Monitoring dashboard at http://localhost:8000/monitoring")
	log.Println("🔗 Prometheus client initialized for host metrics")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is healthy"))
}

func corsPreflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusOK)
}
