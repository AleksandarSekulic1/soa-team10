package startup

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"follower-service/api"
	"follower-service/repository"
	"follower-service/service"

	"github.com/gorilla/mux"
)

func StartServer() {
	// Initialize Neo4j connection
	driver, err := GetNeo4jDriver()
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer driver.Close(context.Background())

	// Create constraints
	if err := CreateConstraints(driver); err != nil {
		log.Printf("Warning: Failed to create constraints: %v", err)
	}

	// Initialize layers
	followerRepo := repository.NewFollowerRepository(driver)
	followerService := service.NewFollowerService(followerRepo)
	followerHandler := api.NewFollowerHandler(followerService)

	// Setup routes
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Follower Service is healthy"))
	}).Methods("GET")

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()

	// User management
	apiRouter.HandleFunc("/users", followerHandler.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users", followerHandler.GetAllUsers).Methods("GET")

	// Follow operations
	apiRouter.HandleFunc("/follow", followerHandler.FollowUser).Methods("POST")
	apiRouter.HandleFunc("/unfollow", followerHandler.UnfollowUser).Methods("POST")

	// Get relationships
	apiRouter.HandleFunc("/users/{userId}/followers", followerHandler.GetFollowers).Methods("GET")
	apiRouter.HandleFunc("/users/{userId}/following", followerHandler.GetFollowing).Methods("GET")
	apiRouter.HandleFunc("/users/{userId}/is-following", followerHandler.IsFollowing).Methods("GET")

	// Recommendations
	apiRouter.HandleFunc("/recommendations", followerHandler.GetRecommendations).Methods("GET")

	// For blog service integration
	apiRouter.HandleFunc("/followed-users", followerHandler.GetFollowedUserIds).Methods("GET")

	// Logging middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Follower Service starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
