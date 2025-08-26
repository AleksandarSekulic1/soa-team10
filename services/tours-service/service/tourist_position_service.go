// service/tourist_position_service.go

package service

import (
	"tours-service/domain"
	"tours-service/repository"
	"time"
)

type TouristPositionService interface {
	GetByUserId(userId string) (*domain.TouristPosition, error) // IZMENA: userId je sada string
	Update(position *domain.TouristPosition) (*domain.TouristPosition, error)
}

type touristPositionService struct {
	repo repository.TouristPositionRepository
}

func NewTouristPositionService(repo repository.TouristPositionRepository) TouristPositionService {
	return &touristPositionService{repo: repo}
}

func (s *touristPositionService) GetByUserId(userId string) (*domain.TouristPosition, error) { // IZMENA: userId je sada string
	return s.repo.GetByUserId(userId)
}

func (s *touristPositionService) Update(position *domain.TouristPosition) (*domain.TouristPosition, error) {
	position.UpdatedAt = time.Now()
	return s.repo.Update(position)
}