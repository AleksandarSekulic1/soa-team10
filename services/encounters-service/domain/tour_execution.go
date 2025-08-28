// encounters-service/domain/tour_execution.go
package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Definišemo tip za status izvođenja ture
type TourExecutionStatus string

const (
	ExecutionStatusActive    TourExecutionStatus = "Active"
	ExecutionStatusCompleted TourExecutionStatus = "Completed"
	ExecutionStatusAbandoned TourExecutionStatus = "Abandoned"
)

// Struktura koja beleži kada je koja ključna tačka kompletirana
type CompletedKeyPoint struct {
	KeyPointId     primitive.ObjectID `bson:"keyPointId"`
	CompletionTime time.Time          `bson:"completionTime"`
}

// Glavni model za praćenje sesije ture
type TourExecution struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty"`
	TourId             string              `bson:"tourId"` // <-- ISPRAVKA: Promenjeno sa int64 na string
	UserId             string              `bson:"userId"` // ID turiste
	Status             TourExecutionStatus `bson:"status"`
	CompletedKeyPoints []CompletedKeyPoint `bson:"completedKeyPoints,omitempty"`
	LastActivity       time.Time           `bson:"lastActivity"`
	StartTime          time.Time           `bson:"startTime"`
	EndTime            *time.Time          `bson:"endTime,omitempty"` // Pointer da može biti null
}