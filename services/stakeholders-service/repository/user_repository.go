// repository/user_repository.go

package repository

import (
	"context"
	"stakeholders-service/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository interface {
	Create(user *domain.User) error
}

type userRepository struct {
	driver neo4j.DriverWithContext
}

// "Konstruktor" sada prihvata drajver kao zavisnost
func NewUserRepository(driver neo4j.DriverWithContext) UserRepository {
	return &userRepository{driver: driver}
}

func (r *userRepository) Create(user *domain.User) error {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			// Cypher upit za kreiranje čvora (node) sa labelom 'User'
			query := "CREATE (u:User {username: $username, password: $password, email: $email, role: $role})"

			// Parametri za upit, mapirani iz našeg User objekta
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
