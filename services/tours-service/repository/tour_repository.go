package repository

import (
	"context"
	"tours-service/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type TourRepository interface {
	Create(tour *domain.Tour) error
}

type tourRepository struct {
	tours *mongo.Collection
}

func NewTourRepository(client *mongo.Client) TourRepository {
	tours := client.Database("tours-db").Collection("tours")
	return &tourRepository{tours: tours}
}

func (r *tourRepository) Create(tour *domain.Tour) error {
	_, err := r.tours.InsertOne(context.TODO(), tour)
	return err
}
