// repository/tourist_position_repository.go

package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"encounters-service/domain"
)

type TouristPositionRepository interface {
	GetByUserId(userId string) (*domain.TouristPosition, error) // IZMENA: userId je sada string
	Update(position *domain.TouristPosition) (*domain.TouristPosition, error)
}

type touristPositionRepository struct {
	positions *mongo.Collection
}

func NewTouristPositionRepository(client *mongo.Client) TouristPositionRepository {
	positions := client.Database("tours-db").Collection("tourist-positions")
	return &touristPositionRepository{positions: positions}
}

func (r *touristPositionRepository) GetByUserId(userId string) (*domain.TouristPosition, error) { // IZMENA: userId je sada string
	var position domain.TouristPosition
	filter := bson.M{"userId": userId}
	err := r.positions.FindOne(context.TODO(), filter).Decode(&position)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// Update koristi "FindOneAndUpdate" da bi odmah vratio ažurirani/kreirani dokument.
func (r *touristPositionRepository) Update(position *domain.TouristPosition) (*domain.TouristPosition, error) {
	filter := bson.M{"userId": position.UserId}
	update := bson.M{"$set": position}
	
	// Opcije:
	// - SetUpsert(true): Ako dokument ne postoji, kreiraj ga.
	// - SetReturnDocument(options.After): Vrati dokument NAKON što je ažuriran.
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedPosition domain.TouristPosition
	err := r.positions.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedPosition)
	if err != nil {
		return nil, err
	}
	
	return &updatedPosition, nil
}