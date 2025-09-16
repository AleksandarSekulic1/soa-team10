// encounters-service/startup/db.go
package startup

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// Podesite vaš connection string za MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://mongo-db:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Provera konekcije
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}