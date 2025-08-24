package service

import (
	"time"
	"tours-service/domain"
	"tours-service/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourService interface {
	Create(tour *domain.Tour) error
	GetByAuthorId(authorId string) ([]*domain.Tour, error)
	GetAll() ([]*domain.Tour, error)                          // <-- DODATO
	AddReview(tourId string, review *domain.TourReview) error // <-- DODATO
}

type tourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) TourService {
	return &tourService{repo: repo}
}

func (s *tourService) Create(tour *domain.Tour) error {
	tour.Status = "draft"
	tour.Price = 0.0
	return s.repo.Create(tour)
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
