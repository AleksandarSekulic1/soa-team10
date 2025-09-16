// domain/tour.go

package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// --- NOVO: Definišemo tip za status ture ---
type TourStatus string

const (
	TourStatusDraft     TourStatus = "draft"
	TourStatusPublished TourStatus = "published"
	TourStatusArchived  TourStatus = "archived"
)

// --- NOVO: Definišemo tip za prevoz ---
type TransportType string

const (
	TransportTypeWalking TransportType = "walking"
	TransportTypeBicycle TransportType = "bicycle"
	TransportTypeCar     TransportType = "car"
)

// --- NOVO: Struktura za čuvanje vremena putovanja ---
type TourTransport struct {
	Type          TransportType `bson:"type" json:"type"`
	TimeInMinutes int           `bson:"timeInMinutes" json:"timeInMinutes"`
}

// Tour definiše model ture u sistemu
type Tour struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AuthorId      string             `bson:"authorId" json:"authorId"`
	Name          string             `bson:"name" json:"name" binding:"required"`
	Description   string             `bson:"description" json:"description" binding:"required"`
	Difficulty    int                `bson:"difficulty" json:"difficulty" binding:"required"`
	Tags          []string           `bson:"tags" json:"tags"`
	Status        TourStatus         `bson:"status" json:"status"` // IZMENA: Koristimo novi tip
	Price         float64            `bson:"price" json:"price"`
	Reviews       []TourReview       `bson:"reviews" json:"reviews"`
	KeyPoints     []TourKeyPoint     `bson:"keyPoints" json:"keyPoints"`
	
	// --- NOVA POLJA ---
	Distance      float64         `bson:"distance" json:"distance"` // Distanca u kilometrima
	TransportInfo []TourTransport `bson:"transportInfo" json:"transportInfo"` // Lista vremena putovanja
	PublishedAt   *time.Time      `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"` // Vreme objave
	ArchivedAt    *time.Time      `bson:"archivedAt,omitempty" json:"archivedAt,omitempty"` // Vreme arhiviranja
}