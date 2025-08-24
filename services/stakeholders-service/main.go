package main

import (
	"stakeholders-service/startup"
)

// @title Stakeholders Service API
// @version 1.0
// @description API za upravljanje korisnicima (Stakeholders) u turističkoj aplikaciji.
// @host localhost:8081
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 1. Inicijalizujemo konekciju sa bazom
	driver := startup.NewDriver()

	// 2. Kreiramo novi server i prosleđujemo mu drajver za bazu
	server := startup.NewServer(driver)

	// 3. Pokrećemo server
	server.Start()
}