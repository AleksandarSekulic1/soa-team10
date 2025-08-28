// encounters-service/repository/tour_execution_repository.go
package repository

import (
	"context"
	"encounters-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourExecutionRepository interface {
	GetById(id string) (*domain.TourExecution, error)
	GetActiveByUser(userId string) (*domain.TourExecution, error)
	Create(execution *domain.TourExecution) error
	Update(execution *domain.TourExecution) error
}

type tourExecutionRepository struct {
	executions *mongo.Collection
}

func NewTourExecutionRepository(client *mongo.Client) TourExecutionRepository {
	executions := client.Database("encounters-db").Collection("tour-executions")
	return &tourExecutionRepository{executions: executions}
}

func (r *tourExecutionRepository) GetById(id string) (*domain.TourExecution, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var execution domain.TourExecution
	filter := bson.M{"_id": objID}
	err = r.executions.FindOne(context.TODO(), filter).Decode(&execution)
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *tourExecutionRepository) GetActiveByUser(userId string) (*domain.TourExecution, error) {
	var execution domain.TourExecution
	// Pronalazi sesiju za datog korisnika koja je trenutno aktivna
	filter := bson.M{"userId": userId, "status": domain.ExecutionStatusActive}
	err := r.executions.FindOne(context.TODO(), filter).Decode(&execution)
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *tourExecutionRepository) Create(execution *domain.TourExecution) error {
	result, err := r.executions.InsertOne(context.TODO(), execution)
	if err != nil {
		return err
	}
	execution.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *tourExecutionRepository) Update(execution *domain.TourExecution) error {
	filter := bson.M{"_id": execution.ID}
	update := bson.M{"$set": execution}
	_, err := r.executions.UpdateOne(context.TODO(), filter, update)
	return err
}