// services/tours-service/main.go
package main

import (
	"log"
	"net"
	"tours-service/api"
	"tours-service/repository"
	"tours-service/service"
	"tours-service/startup"

	"google.golang.org/grpc"
	"tours-service/proto/tours" // Uvozimo generisani proto kod
)

// @title           Tours Service API
// @version         1.0
// @description     API for managing tours.
// @host            localhost:8083
// @BasePath        /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Pokrećemo gRPC server u pozadini (u posebnoj "gorutini")
	// da ne bi blokirao pokretanje HTTP servera.
	go startGrpcServer()

	// Pokrećemo HTTP (REST) server kao i do sada
	server := startup.NewServer()
	server.Start()
}

// Nova funkcija za pokretanje gRPC servera
func startGrpcServer() {
	// Definišemo port za gRPC komunikaciju (npr. 8086, važno je da bude drugačiji od REST porta)
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	// Inicijalizujemo slojeve, slično kao za REST
	mongoClient := startup.ConnectDB()
	tourRepo := repository.NewTourRepository(mongoClient)
	tourService := service.NewTourService(tourRepo)
	
	// Kreiramo instancu našeg novog gRPC hendlera
	toursHandler := api.NewToursGrpcHandler(tourService)

	// Kreiramo novi gRPC server
	grpcServer := grpc.NewServer()
	
	// Registrujemo naš hendler na server
	tours.RegisterToursServiceServer(grpcServer, toursHandler)

	log.Println("gRPC server is listening on port 8086")
	// Pokrećemo server da sluša na definisanom portu
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}