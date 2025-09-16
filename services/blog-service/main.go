package main

import (
	"blog-service/startup"
)

// @title Blog Service API
// @version 1.0
// @description API za upravljanje blogovima u turističkoj aplikaciji.
// @host localhost:8082
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	server := startup.NewServer()
	server.Start()
}