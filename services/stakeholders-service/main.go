// main.go

package main

import (
	"context"
	"log"
	"stakeholders-service/startup"
)

// @title Stakeholders Service API ...
func main() {
	// Kredencijali iz docker-compose ili docker run komande
	uri := "neo4j://localhost:7687"
	user := "neo4j"
	pass := "sifra1234"

	driver, err := startup.InitDB(uri, user, pass)
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close(context.Background())

	server := startup.NewServer(driver) // Prosleđujemo drajver
	server.Start()
}
