package main

import (
	"context"
	"log"
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
	uri := "neo4j://localhost:7687"
	user := "neo4j"
	pass := "sifra1234"

	driver, err := startup.InitDB(uri, user, pass)
	if err != nil {
		log.Fatal(err)
	}
	// Moramo definisati context ovde da bi bio dostupan
	ctx := context.Background()
	defer driver.Close(ctx)

	server := startup.NewServer(driver)
	server.Start()
}
