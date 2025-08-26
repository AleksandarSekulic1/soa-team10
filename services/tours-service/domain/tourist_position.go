// domain/tourist_position.go

package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TouristPosition struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    string             `bson:"userId"` // IZMENA: UserId je sada string
	Latitude  float64            `bson:"latitude"`
	Longitude float64            `bson:"longitude"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}