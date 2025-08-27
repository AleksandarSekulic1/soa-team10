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
	GetById(tourId string) (*domain.Tour, error) // <-- DODATI
	GetByAuthorId(authorId string) ([]*domain.Tour, error)
	GetAll() ([]*domain.Tour, error)                          // <-- DODATO
	AddReview(tourId string, review *domain.TourReview) error // <-- DODATO
	// NOVE METODE ZA KEY POINTS
	AddKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error
	UpdateKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error
	DeleteKeyPoint(tourId, keyPointId string) error
	Update(tour *domain.Tour) error // <-- DODATI NOVU METODU
	GetPublished() ([]*domain.Tour, error) // <-- DODATA NOVA METODA
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

func (r *tourRepository) GetById(tourId string) (*domain.Tour, error) {
	var tour domain.Tour
	objID, err := primitive.ObjectIDFromHex(tourId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	err = r.tours.FindOne(context.TODO(), filter).Decode(&tour)
	if err != nil {
		return nil, err
	}
	return &tour, nil
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

func (r *tourRepository) AddKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error {
	tourObjID, err := primitive.ObjectIDFromHex(tourId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": tourObjID}
	update := bson.M{"$push": bson.M{"keyPoints": keyPoint}}
	_, err = r.tours.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *tourRepository) UpdateKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error {
	tourObjID, err := primitive.ObjectIDFromHex(tourId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": tourObjID, "keyPoints._id": keyPoint.ID}
	update := bson.M{"$set": bson.M{"keyPoints.$": keyPoint}}
	_, err = r.tours.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *tourRepository) DeleteKeyPoint(tourId, keyPointId string) error {
	tourObjID, err := primitive.ObjectIDFromHex(tourId)
	if err != nil {
		return err
	}
	keyPointObjID, err := primitive.ObjectIDFromHex(keyPointId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": tourObjID}
	update := bson.M{"$pull": bson.M{"keyPoints": bson.M{"_id": keyPointObjID}}}
	_, err = r.tours.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *tourRepository) Update(tour *domain.Tour) error {
	filter := bson.M{"_id": tour.ID}
	update := bson.M{"$set": tour}
	_, err := r.tours.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *tourRepository) GetPublished() ([]*domain.Tour, error) {
	var tours []*domain.Tour
	// Filter koji vraća samo dokumente gde je status "published"
	filter := bson.M{"status": domain.TourStatusPublished}
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