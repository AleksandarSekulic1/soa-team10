package repository

import (
	"context"
	"tours-service/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourRepository interface {
	Create(tour *domain.Tour) error
	GetByAuthorId(authorId string) ([]*domain.Tour, error)
	GetAll() ([]*domain.Tour, error)                          // <-- DODATO
	AddReview(tourId string, review *domain.TourReview) error // <-- DODATO
}

type tourRepository struct {
	tours *mongo.Collection
}

func NewTourRepository(client *mongo.Client) TourRepository {
	tours := client.Database("tours-db").Collection("tours")
	return &tourRepository{tours: tours}
}

func (r *tourRepository) Create(tour *domain.Tour) error {
	result, err := r.tours.InsertOne(context.TODO(), tour)
	if err != nil {
		return err
	}
	tour.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *tourRepository) GetByAuthorId(authorId string) ([]*domain.Tour, error) {
	var tours []*domain.Tour
	filter := bson.M{"authorId": authorId}
	cursor, err := r.tours.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	if err = cursor.All(context.TODO(), &tours); err != nil {
		return nil, err
	}
	return tours, nil
}

func (r *tourRepository) GetAll() ([]*domain.Tour, error) {
	var tours []*domain.Tour
	cursor, err := r.tours.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tours); err != nil {
		return nil, err
	}
	return tours, nil
}

func (r *tourRepository) AddReview(tourId string, review *domain.TourReview) error {
	objID, err := primitive.ObjectIDFromHex(tourId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$push": bson.M{"reviews": review}}
	_, err = r.tours.UpdateOne(context.TODO(), filter, update)
	return err
}
