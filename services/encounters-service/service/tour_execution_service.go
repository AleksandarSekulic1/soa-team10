// encounters-service/service/tour_execution_service.go
package service

import (
	"context"
	"encounters-service/domain"
	"encounters-service/repository"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math"
	"time"
	
	"encounters-service/proto/tours" // Uvozimo proto kod
)

// Definišemo lokalnu KeyPoint strukturu
type KeyPoint struct {
	ID        primitive.ObjectID
	Latitude  float64
	Longitude float64
	Name      string // Dodajemo ime radi logovanja i budućih potreba
}

const KeyPointCompletionThreshold = 500.0 // Povećano na 500 metara radi lakšeg testiranja

type TourExecutionService interface {
	StartTour(execution *domain.TourExecution) (*domain.TourExecution, error)
	CheckPosition(userId string, currentLatitude, currentLongitude float64) (*domain.TourExecution, error)
	CompleteTour(executionId string) (*domain.TourExecution, error)
	AbandonTour(executionId string) (*domain.TourExecution, error)
	GetActiveByUser(userId string) (*domain.TourExecution, error) // <-- DODAJTE OVU METODU
}

type tourExecutionService struct {
	repo        repository.TourExecutionRepository
	toursClient tours.ToursServiceClient // <-- NOVO: gRPC klijent
}

// Ažuriramo konstruktor da kreira gRPC klijenta
func NewTourExecutionService(repo repository.TourExecutionRepository) TourExecutionService {
	// Uspostavljamo gRPC konekciju ka tours-service koji radi na portu 8086
	conn, err := grpc.Dial("tours-service:8086", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to tours-service: %v", err)
	}

	client := tours.NewToursServiceClient(conn)

	return &tourExecutionService{
		repo:        repo,
		toursClient: client, // <-- Čuvamo klijenta
	}
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

// --- PREPRAVLJENA CheckPosition METODA ---
// encounters-service/service/tour_execution_service.go

func (s *tourExecutionService) CheckPosition(userId string, currentLatitude, currentLongitude float64) (*domain.TourExecution, error) {
	activeExecution, err := s.repo.GetActiveByUser(userId)
	if err != nil {
		return nil, errors.New("no active tour found for this user")
	}

	log.Printf("--- Checking position for user %s ---", userId)
	log.Printf("Tourist current position: Lat=%f, Lon=%f", currentLatitude, currentLongitude)

	tourResponse, err := s.toursClient.GetTourById(context.Background(), &tours.GetTourByIdRequest{TourId: activeExecution.TourId})
	if err != nil {
		log.Printf("gRPC call to tours-service failed: %v", err)
		return nil, fmt.Errorf("could not get tour details via gRPC: %w", err)
	}

	log.Printf("Received %d key points from tours-service.", len(tourResponse.KeyPoints))
	
	// --- NOVA, PAMETNIJA LOGIKA ---
	
	// 1. Određujemo koja je sledeća ključna tačka na redu
	completedCount := len(activeExecution.CompletedKeyPoints)
	if completedCount >= len(tourResponse.KeyPoints) {
		log.Println("All key points already completed.")
		return activeExecution, nil // Sve je već završeno
	}
	
	nextKeyPointProto := tourResponse.KeyPoints[completedCount]
	nextKeyPointId, _ := primitive.ObjectIDFromHex(nextKeyPointProto.Id)

	// 2. Računamo distancu SAMO do te sledeće tačke
	distance := calculateDistance(currentLatitude, currentLongitude, nextKeyPointProto.Latitude, nextKeyPointProto.Longitude)
	log.Printf("Checking distance to NEXT key point '%s'... Distance: %.2f meters", nextKeyPointProto.Name, distance)

	// 3. Ako smo dovoljno blizu, kompletiramo je
	if distance <= KeyPointCompletionThreshold {
		completedPoint := domain.CompletedKeyPoint{
			KeyPointId:     nextKeyPointId,
			CompletionTime: time.Now(),
		}
		activeExecution.CompletedKeyPoints = append(activeExecution.CompletedKeyPoints, completedPoint)
		log.Printf("!!! SUCCESS: User '%s' completed key point '%s'", userId, nextKeyPointProto.Name)
	}

	// Uvek ažuriramo vreme poslednje aktivnosti
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
/*func getMockKeyPoints() []KeyPoint {
	kp1ID, _ := primitive.ObjectIDFromHex("68af26e5d4000566777bb098")
	kp2ID, _ := primitive.ObjectIDFromHex("68af26fdd4000566777bb099")
	return []KeyPoint{
		{ID: kp1ID, Latitude: 44.7951, Longitude: 20.4568},
		{ID: kp2ID, Latitude: 44.7828, Longitude: 20.4810},
	}
}*/

func (s *tourExecutionService) GetActiveByUser(userId string) (*domain.TourExecution, error) {
	return s.repo.GetActiveByUser(userId)
}