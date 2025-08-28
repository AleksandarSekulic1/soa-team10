// encounters-service/main.go

package main

import "encounters-service/startup"

// @title           Encounters Service API
// @version         1.0
// @description     This is a server for managing encounters and tourist positions.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8085
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	server := startup.NewServer()
	server.Start()
}