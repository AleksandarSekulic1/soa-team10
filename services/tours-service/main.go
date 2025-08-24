package main

import (
	"log"
	"os"
	"tours-service/startup"
)

// @title Tours Service API
// @version 1.0
// @description API za upravljanje turama.
// @host localhost:8082
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := startup.InitDB(uri)
	if err != nil {
		log.Fatal(err)
	}

	server := startup.NewServer(client)
	server.Start()
}
