// service/user_service.go

package service

import (
	"stakeholders-service/domain"
	"stakeholders-service/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *domain.User) error {
	// 1. Hešovanje lozinke (poslovna logika)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 2. Pozivanje repozitorijuma da sačuva korisnika
	return s.repo.Create(user)
}
