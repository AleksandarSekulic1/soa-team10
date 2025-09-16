package startup

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// Čitamo URI ISKLJUČIVO iz environment varijable
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment varijabla nije postavljena!")
	}

	log.Printf("Povezivanje na MongoDB na adresi: %s", uri)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("Neuspešno povezivanje na MongoDB: %v", err)
	}

	// Proveravamo konekciju
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Neuspešan ping na MongoDB: %v", err)
	}

	log.Println("Uspešno povezan na MongoDB.")
	return client
}
