package service

import (
	"tours-service/domain"
	"tours-service/repository"
)

type TourService interface {
	Create(tour *domain.Tour) error
}

type tourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) TourService {
	return &tourService{repo: repo}
}

func (s *tourService) Create(tour *domain.Tour) error {
	// Pri početnom kreiranju ture, tura treba da ima status draft i da joj je cena postavljena na 0. [cite: 79]
	tour.Status = "draft"
	tour.Price = 0.0

	return s.repo.Create(tour)
}
