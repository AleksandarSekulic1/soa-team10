// repository/user_repository.go

package repository

import (
	"context"
	"stakeholders-service/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetAll() ([]*domain.User, error)
	GetByUsername(username string) (*domain.User, error) // <-- NOVA METODA
}

type userRepository struct {
	driver neo4j.DriverWithContext
}

func NewUserRepository(driver neo4j.DriverWithContext) UserRepository {
	return &userRepository{driver: driver}
}

func (r *userRepository) Create(user *domain.User) error {
	// ... postojeća Create metoda ostaje ista ...
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			query := "CREATE (u:User {username: $username, password: $password, email: $email, role: $role})"
			parameters := map[string]any{
				"username": user.Username,
				"password": user.Password,
				"email":    user.Email,
				"role":     user.Role,
			}
			_, err := tx.Run(ctx, query, parameters)
			return nil, err
		})
	return err
}

func (r *userRepository) GetAll() ([]*domain.User, error) {
	// ... postojeća GetAll metoda ostaje ista ...
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	users := []*domain.User{}
	_, err := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			query := "MATCH (u:User) RETURN u.username AS username, u.email AS email, u.role AS role"
			records, err := tx.Run(ctx, query, nil)
			if err != nil {
				return nil, err
			}
			for records.Next(ctx) {
				record := records.Record()
				username, _ := record.Get("username")
				email, _ := record.Get("email")
				role, _ := record.Get("role")
				users = append(users, &domain.User{
					Username: username.(string),
					Email:    email.(string),
					Role:     role.(string),
				})
			}
			return nil, records.Err()
		})
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Implementacija nove metode
func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	user, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (u:User {username: $username}) RETURN u.username AS username, u.password AS password, u.email AS email, u.role AS role"
		parameters := map[string]any{"username": username}

		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}

		// Očekujemo samo jedan rezultat
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err // Vraća grešku ako korisnik nije pronađen
		}

		// Parsiramo podatke u User objekat
		return &domain.User{
			Username: record.Values[0].(string),
			Password: record.Values[1].(string),
			Email:    record.Values[2].(string),
			Role:     record.Values[3].(string),
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return user.(*domain.User), nil
}
