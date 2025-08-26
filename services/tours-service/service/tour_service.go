package service

import (
	"time"
	"tours-service/domain"
	"tours-service/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourService interface {
	Create(tour *domain.Tour) error
	GetById(tourId string) (*domain.Tour, error) // <-- DODATI
	GetByAuthorId(authorId string) ([]*domain.Tour, error)
	GetAll() ([]*domain.Tour, error)                          // <-- DODATO
	AddReview(tourId string, review *domain.TourReview) error // <-- DODATO
	// NOVE METODE ZA KEY POINTS
	AddKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error
	UpdateKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error
	DeleteKeyPoint(tourId, keyPointId string) error
}

type tourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) TourService {
	return &tourService{repo: repo}
}

// service/tour_service.go

func (s *tourService) Create(tour *domain.Tour) error {
	tour.Status = "draft"
	tour.Price = 0.0
	// --- DODATI OVE DVE LINIJE ---
	tour.KeyPoints = []domain.TourKeyPoint{}
	//tour.Reviews = []domain.TourReview{}
	// ---------------------------
	return s.repo.Create(tour)
}

func (s *tourService) GetById(tourId string) (*domain.Tour, error) {
	return s.repo.GetById(tourId)
}

func (s *tourService) GetByAuthorId(authorId string) ([]*domain.Tour, error) {
	return s.repo.GetByAuthorId(authorId)
}

func (s *tourService) GetAll() ([]*domain.Tour, error) {
	return s.repo.GetAll()
}

func (s *tourService) AddReview(tourId string, review *domain.TourReview) error {
	// Postavljamo serverske vrednosti
	review.ID = primitive.NewObjectID()
	review.CommentDate = time.Now()
	return s.repo.AddReview(tourId, review)
}

func (s *tourService) AddKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error {
	// Servis je odgovoran za generisanje ID-ja
	keyPoint.ID = primitive.NewObjectID()
	keyPoint.TourId = tourId
	return s.repo.AddKeyPoint(tourId, keyPoint)
}

func (s *tourService) UpdateKeyPoint(tourId string, keyPoint *domain.TourKeyPoint) error {
	// Logika za validaciju ili dodatne provere bi išla ovde
	return s.repo.UpdateKeyPoint(tourId, keyPoint)
}

func (s *tourService) DeleteKeyPoint(tourId, keyPointId string) error {
	return s.repo.DeleteKeyPoint(tourId, keyPointId)
}