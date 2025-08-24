// u services/blog-service/startup/db.go

package startup

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// ISPRAVKA: Umesto imena servisa 'mongo-db', koristimo 'localhost'
	// jer se program izvršava van Docker mreže.
	clientOptions := options.Client().ApplyURI("mongodb://mongo-db:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB: ", err)
	}

	log.Println("Successfully connected to MongoDB!")
	return client
}