package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Setup CORS middleware
	router.Use(corsMiddleware)
	
	// Setup logging middleware
	router.Use(loggingMiddleware)

	// Health check endpoint
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Setup routes for different services
	setupStakeholdersRoutes(router)
	setupPaymentsRoutes(router)
	setupBlogRoutes(router)
	setupToursRoutes(router)
	setupEncountersRoutes(router)
	setupFollowerRoutes(router)

	log.Println("API Gateway starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is healthy"))
}
