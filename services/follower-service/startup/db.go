package startup

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetNeo4jDriver() (neo4j.DriverWithContext, error) {
	uri := os.Getenv("NEO4J_URI")
	if uri == "" {
		uri = "bolt://neo4j-db:7687" // Default for Docker
	}

	username := os.Getenv("NEO4J_USERNAME")
	if username == "" {
		username = "neo4j"
	}

	password := os.Getenv("NEO4J_PASSWORD")
	if password == "" {
		password = "password123"
	}

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	// Test the connection
	ctx := context.Background()
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to Neo4j at %s", uri)
	return driver, nil
}

func CreateConstraints(driver neo4j.DriverWithContext) error {
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// Create unique constraint on User ID
	_, err := session.Run(ctx, "CREATE CONSTRAINT user_id_unique IF NOT EXISTS FOR (u:User) REQUIRE u.id IS UNIQUE", nil)
	if err != nil {
		log.Printf("Error creating user ID constraint: %v", err)
		return err
	}

	log.Println("Neo4j constraints created successfully")
	return nil
}
