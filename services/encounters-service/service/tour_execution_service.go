// encounters-service/service/tour_execution_service.go
package service

import (
	"encounters-service/domain"
	"encounters-service/repository"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive" // <-- 1. DODAT IMPORT
	"math"
	"time"
)

// 2. Definisana je lokalna `KeyPoint` struktura
// Ona sadrži samo polja koja su nam potrebna u ovom servisu
type KeyPoint struct {
	ID        primitive.ObjectID
	Latitude  float64
	Longitude float64
}

const KeyPointCompletionThreshold = 50.0

type TourExecutionService interface {
	StartTour(execution *domain.TourExecution) (*domain.TourExecution, error)
	CheckPosition(userId string, currentLatitude, currentLongitude float64) (*domain.TourExecution, error)
	CompleteTour(executionId string) (*domain.TourExecution, error)
	AbandonTour(executionId string) (*domain.TourExecution, error)
}

type tourExecutionService struct {
	repo repository.TourExecutionRepository
}

func NewTourExecutionService(repo repository.TourExecutionRepository) TourExecutionService {
	return &tourExecutionService{repo: repo}
}

func (s *tourExecutionService) StartTour(execution *domain.TourExecution) (*domain.TourExecution, error) {
	activeExecution, _ := s.repo.GetActiveByUser(execution.UserId)
	if activeExecution != nil {
		return nil, errors.New("user already has an active tour")
	}

	execution.Status = domain.ExecutionStatusActive
	execution.StartTime = time.Now()
	execution.LastActivity = time.Now()
	execution.CompletedKeyPoints = []domain.CompletedKeyPoint{}

	err := s.repo.Create(execution)
	if err != nil {
		return nil, err
	}
	return execution, nil
}

func (s *tourExecutionService) CheckPosition(userId string, currentLatitude, currentLongitude float64) (*domain.TourExecution, error) {
	activeExecution, err := s.repo.GetActiveByUser(userId)
	if err != nil {
		return nil, errors.New("no active tour found for this user")
	}

	// TODO: gRPC poziv ka tours-service da se dobiju ključne tačke za activeExecution.TourId
	mockKeyPoints := getMockKeyPoints() // 3. Ažurirana `getMockKeyPoints` funkcija

	// 4. Ažurirana `CheckPosition` metoda da koristi lokalnu KeyPoint strukturu
	for _, keyPoint := range mockKeyPoints {
		isAlreadyCompleted := false
		for _, completedKp := range activeExecution.CompletedKeyPoints {
			if completedKp.KeyPointId == keyPoint.ID {
				isAlreadyCompleted = true
				break
			}
		}

		if !isAlreadyCompleted {
			distance := calculateDistance(currentLatitude, currentLongitude, keyPoint.Latitude, keyPoint.Longitude)

			if distance <= KeyPointCompletionThreshold {
				completedPoint := domain.CompletedKeyPoint{
					KeyPointId:     keyPoint.ID,
					CompletionTime: time.Now(),
				}
				activeExecution.CompletedKeyPoints = append(activeExecution.CompletedKeyPoints, completedPoint)
			}
		}
	}

	activeExecution.LastActivity = time.Now()
	err = s.repo.Update(activeExecution)
	if err != nil {
		return nil, err
	}

	return activeExecution, nil
}

func (s *tourExecutionService) CompleteTour(executionId string) (*domain.TourExecution, error) {
	execution, err := s.repo.GetById(executionId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	execution.Status = domain.ExecutionStatusCompleted
	execution.EndTime = &now
	err = s.repo.Update(execution)
	return execution, err
}

func (s *tourExecutionService) AbandonTour(executionId string) (*domain.TourExecution, error) {
	execution, err := s.repo.GetById(executionId)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	execution.Status = domain.ExecutionStatusAbandoned
	execution.EndTime = &now
	err = s.repo.Update(execution)
	return execution, err
}


// Pomoćna funkcija za računanje distance (Haversine formula)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	radLat1 := lat1 * math.Pi / 180
	radLat2 := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(radLat1)*math.Cos(radLat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// Pomoćna funkcija koja simulira dobavljanje ključnih tačaka i sada koristi lokalnu strukturu
func getMockKeyPoints() []KeyPoint {
	kp1ID, _ := primitive.ObjectIDFromHex("68af26e5d4000566777bb098")
	kp2ID, _ := primitive.ObjectIDFromHex("68af26fdd4000566777bb099")
	return []KeyPoint{
		{ID: kp1ID, Latitude: 44.7951, Longitude: 20.4568},
		{ID: kp2ID, Latitude: 44.7828, Longitude: 20.4810},
	}
}