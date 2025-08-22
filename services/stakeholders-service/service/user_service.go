package service

import (
	"errors"
	"time"

	"stakeholders-service/domain"
	"stakeholders-service/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Definišemo tajni ključ. Sada počinje VELIKIM SLOVOM da bi bio vidljiv van paketa.
var JwtKey = []byte("super_secret_key")

type UserService interface {
	Register(user *domain.User) error
	GetAll() ([]*domain.User, error)
	Login(username, password string) (string, error)
	GetProfile(username string) (*domain.User, error)
	UpdateProfile(user *domain.User) (*domain.User, error)
	SetBlockedStatus(username string, isBlocked bool) error
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.Create(user)
}

func (s *userService) GetAll() ([]*domain.User, error) {
	return s.repo.GetAll()
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if user.IsBlocked {
		return "", errors.New("user is blocked")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // Koristimo novo, veliko slovo 'J'
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *userService) GetProfile(username string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateProfile(user *domain.User) (*domain.User, error) {
	err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return s.GetProfile(user.Username)
}

func (s *userService) SetBlockedStatus(username string, isBlocked bool) error {
	return s.repo.UpdateBlockedStatus(username, isBlocked)
}