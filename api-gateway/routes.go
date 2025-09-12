package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

// setupStakeholdersRoutes sets up routes for stakeholders service
func setupStakeholdersRoutes(router *mux.Router) {
	stakeholdersURL, _ := url.Parse("http://stakeholders-service:8081")
	stakeholdersProxy := httputil.NewSingleHostReverseProxy(stakeholdersURL)

	// Custom response modifier to ensure CORS headers
	stakeholdersProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Max-Age", "86400")
		return nil
	}

	// Handle OPTIONS preflight requests for stakeholders
	router.PathPrefix("/api/stakeholders/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
	})

	// Stakeholders service routes - keep full path
	router.PathPrefix("/api/stakeholders/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stakeholdersProxy.ServeHTTP(w, r)
	})
}

// setupPaymentsRoutes sets up routes for payments service
func setupPaymentsRoutes(router *mux.Router) {
	paymentsURL, _ := url.Parse("http://payments-service:8080")
	paymentsProxy := httputil.NewSingleHostReverseProxy(paymentsURL)

	// Custom response modifier to ensure CORS headers
	paymentsProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Max-Age", "86400")
		return nil
	}

	// Payments service routes - map /api/payments to /api/shopping-cart
	router.PathPrefix("/api/payments/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Replace /api/payments with /api/shopping-cart for the backend service
		r.URL.Path = "/api/shopping-cart" + r.URL.Path[len("/api/payments"):]
		paymentsProxy.ServeHTTP(w, r)
	})
}

// setupBlogRoutes sets up routes for blog service
func setupBlogRoutes(router *mux.Router) {
	blogURL, _ := url.Parse("http://blog-service:8082")
	blogProxy := httputil.NewSingleHostReverseProxy(blogURL)

	// Custom response modifier to ensure CORS headers
	blogProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Max-Age", "86400")
		return nil
	}

	// Handle /api/blog without trailing slash FIRST (more specific route)
	router.HandleFunc("/api/blog", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/api/blogs"
		blogProxy.ServeHTTP(w, r)
	})
	
	// Blog service routes - map /api/blog to /api/blogs
	router.PathPrefix("/api/blog/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Replace /api/blog with /api/blogs for the backend service
		r.URL.Path = "/api/blogs" + r.URL.Path[len("/api/blog"):]
		blogProxy.ServeHTTP(w, r)
	})
}

// setupToursRoutes sets up routes for tours service
func setupToursRoutes(router *mux.Router) {
	toursURL, _ := url.Parse("http://tours-service:8083")
	toursProxy := httputil.NewSingleHostReverseProxy(toursURL)

	// Custom response modifier to ensure CORS headers
	toursProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Max-Age", "86400")
		return nil
	}

	// Tours service routes - both with and without trailing slash
	router.PathPrefix("/api/tours/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		toursProxy.ServeHTTP(w, r)
	})
	router.HandleFunc("/api/tours", func(w http.ResponseWriter, r *http.Request) {
		toursProxy.ServeHTTP(w, r)
	})
}

// setupEncountersRoutes sets up routes for encounters service
func setupEncountersRoutes(router *mux.Router) {
	encountersURL, _ := url.Parse("http://encounters-service:8084")
	encountersProxy := httputil.NewSingleHostReverseProxy(encountersURL)

	// Custom response modifier to ensure CORS headers
	encountersProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Max-Age", "86400")
		return nil
	}

	// Encounters service routes for tourist position
	router.PathPrefix("/api/tourist-position/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encountersProxy.ServeHTTP(w, r)
	})
	
	// Encounters service routes for tour executions
	router.PathPrefix("/api/tour-executions/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encountersProxy.ServeHTTP(w, r)
	})
	
	// Keep the original encounters prefix for compatibility
	router.PathPrefix("/api/encounters/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encountersProxy.ServeHTTP(w, r)
	})
}
