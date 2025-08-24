package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Tour definiše model ture u sistemu
type Tour struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AuthorId    string             `bson:"authorId" json:"authorId"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Description string             `bson:"description" json:"description" binding:"required"`
	Difficulty  int                `bson:"difficulty" json:"difficulty" binding:"required"`
	Tags        []string           `bson:"tags" json:"tags"`
	Status      string             `bson:"status" json:"status"`
	Price       float64            `bson:"price" json:"price"`
}
