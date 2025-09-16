package startup

import (
	"context"
	"log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// NewDriver inicijalizuje i vraća drajver za Neo4j bazu.
func NewDriver() neo4j.DriverWithContext {
	// Koristimo ime servisa 'neo4j-db' iz docker-compose.yml fajla.
	// Port 7687 je podrazumevani Bolt port za Neo4j.
	uri := "neo4j://neo4j-db:7687"
	
	// Pošto smo u docker-compose.yml podesili NEO4J_AUTH: none,
	// korisničko ime i lozinka su prazni stringovi.
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth("", "", ""))
	if err != nil {
		// Ako konekcija ne uspe, program se gasi uz fatalnu grešku.
		log.Fatal("Error creating Neo4j driver:", err)
	}

	// Proveravamo da li je konekcija zaista uspostavljena.
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		log.Fatal("Error verifying Neo4j connectivity:", err)
	}

	log.Println("Successfully connected to Neo4j!")
	return driver
}