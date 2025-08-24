package main

import (
	"tours-service/startup"
)

// @title           Tours Service API
// @version         1.0
// @description     API za upravljanje turama.
// @host            localhost:8083
// @BasePath        /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	server := startup.NewServer()
	server.Start()
}
