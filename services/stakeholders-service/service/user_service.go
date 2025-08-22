// service/user_service.go

package service

import (
	"errors" // Import za greške
	"time"   // Import za vreme (za trajanje tokena)

	"stakeholders-service/domain"
	"stakeholders-service/repository"

	"github.com/golang-jwt/jwt/v5" // Import za JWT
	"golang.org/x/crypto/bcrypt"
)

// Definišemo tajni ključ za potpisivanje tokena.
// U pravoj aplikaciji, ovo bi trebalo da bude u environment varijabli.
var jwtKey = []byte("super_secret_key")

type UserService interface {
	Register(user *domain.User) error
	GetAll() ([]*domain.User, error)
	Login(username, password string) (string, error) // <-- NOVA METODA
}

// Claims struktura za JWT token
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
	// ... postojeća Register metoda ostaje ista ...
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.Create(user)
}

func (s *userService) GetAll() ([]*domain.User, error) {
	// ... postojeća GetAll metoda ostaje ista ...
	return s.repo.GetAll()
}

// Implementacija nove metode
func (s *userService) Login(username, password string) (string, error) {
	// 1. Pronađi korisnika u bazi
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. Uporedi heširanu lozinku iz baze sa lozinkom koju je korisnik poslao
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Kreiraj JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // Token traje 24 sata
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
