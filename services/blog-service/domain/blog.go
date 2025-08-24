// u services/blog-service/domain/blog.go

package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	AuthorID  string             `bson:"authorId"` // Koristićemo username iz JWT-a
	CreatedAt time.Time          `bson:"createdAt"`
	Images    []string           `bson:"images,omitempty"`
	Comments  []Comment          `bson:"comments"`
	Likes     []Like             `bson:"likes"`
}

type Comment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"` // <-- DODATO POLJE
	AuthorID    string    `bson:"authorId"`
	Text        string    `bson:"text"`
	CreatedAt   time.Time `bson:"createdAt"`
	LastUpdatedAt time.Time `bson:"lastUpdatedAt"`
}

type Like struct {
	UserID    string    `bson:"userId"`
	CreatedAt time.Time `bson:"createdAt"`
}