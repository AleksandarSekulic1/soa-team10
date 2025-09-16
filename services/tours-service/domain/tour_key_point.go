// domain/tour_key_point.go

package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TourKeyPoint definiše model ključne tačke unutar ture
type TourKeyPoint struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TourId      string             `bson:"tourId" json:"tourId"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Description string             `bson:"description" json:"description"`
	Latitude    float64            `bson:"latitude" json:"latitude" binding:"required"`
	Longitude   float64            `bson:"longitude" json:"longitude" binding:"required"`
	ImageUrl    string             `bson:"imageUrl" json:"imageUrl"`
}