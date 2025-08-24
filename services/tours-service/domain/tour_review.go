package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TourReview predstavlja model recenzije za turu
type TourReview struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Rating      int                `bson:"rating" json:"rating" binding:"required,min=1,max=5"`
	Comment     string             `bson:"comment" json:"comment"`
	TouristId   string             `bson:"touristId" json:"touristId"`
	VisitDate   time.Time          `bson:"visitDate" json:"visitDate" binding:"required"`
	CommentDate time.Time          `bson:"commentDate" json:"commentDate"`
	ImageUrls   []string           `bson:"imageUrls" json:"imageUrls"`
}
