package repository

import (
	"context"
	"errors"
	"stakeholders-service/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetAll() ([]*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Update(user *domain.User) error
	UpdateBlockedStatus(username string, isBlocked bool) error
}

type userRepository struct {
	driver neo4j.DriverWithContext
}

func NewUserRepository(driver neo4j.DriverWithContext) UserRepository {
	return &userRepository{driver: driver}
}

func (r *userRepository) Create(user *domain.User) error {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			query := `CREATE (u:User {
                username: $username, password: $password, email: $email, role: $role,
                isBlocked: $isBlocked, firstName: $firstName, lastName: $lastName,
                profilePicture: $profilePicture, biography: $biography, motto: $motto
            })`
			parameters := map[string]any{
				"username":       user.Username,
				"password":       user.Password,
				"email":          user.Email,
				"role":           user.Role,
				"isBlocked":      false, // Eksplicitno postavljamo početnu vrednost
				"firstName":      "",
				"lastName":       "",
				"profilePicture": "",
				"biography":      "",
				"motto":          "",
			}
			_, err := tx.Run(ctx, query, parameters)
			return nil, err
		})
	return err
}

// ### ISPRAVLJENA GETALL METODA ###
func (r *userRepository) GetAll() ([]*domain.User, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	users, err := session.ExecuteRead(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			query := "MATCH (u:User) RETURN u" // Upit koji vraća ceo User node
			result, err := tx.Run(ctx, query, nil)
			if err != nil {
				return nil, err
			}

			var users []*domain.User
			for result.Next(ctx) {
				record := result.Record()
				userNode, ok := record.Get("u")
				if !ok {
					continue // preskoči ako nešto nije u redu sa zapisom
				}
				node := userNode.(neo4j.Node)
				props := node.Props

				// Mapiramo SVA polja iz baze u našu User strukturu
				user := &domain.User{
					Username:       props["username"].(string),
					// Password namerno ne šaljemo nazad
					Email:          props["email"].(string),
					Role:           props["role"].(string),
					IsBlocked:      props["isBlocked"].(bool),
					FirstName:      props["firstName"].(string),
					LastName:       props["lastName"].(string),
					ProfilePicture: props["profilePicture"].(string),
					Biography:      props["biography"].(string),
					Motto:          props["motto"].(string),
				}
				users = append(users, user)
			}
			return users, nil
		})

	if err != nil {
		return nil, err
	}
	return users.([]*domain.User), nil
}


func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	user, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (u:User {username: $username}) RETURN u"
		parameters := map[string]any{"username": username}
		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		userNode, ok := record.Get("u")
		if !ok {
			return nil, errors.New("user not found in record")
		}
		node := userNode.(neo4j.Node)
		props := node.Props
		return &domain.User{
			Username:       props["username"].(string),
			Password:       props["password"].(string),
			Email:          props["email"].(string),
			Role:           props["role"].(string),
			IsBlocked:      props["isBlocked"].(bool),
			FirstName:      props["firstName"].(string),
			LastName:       props["lastName"].(string),
			ProfilePicture: props["profilePicture"].(string),
			Biography:      props["biography"].(string),
			Motto:          props["motto"].(string),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return user.(*domain.User), nil
}

func (r *userRepository) Update(user *domain.User) error {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
            MATCH (u:User {username: $username})
            SET u.firstName = $firstName, u.lastName = $lastName, u.profilePicture = $profilePicture,
                u.biography = $biography, u.motto = $motto
        `
		params := map[string]any{
			"username":       user.Username,
			"firstName":      user.FirstName,
			"lastName":       user.LastName,
			"profilePicture": user.ProfilePicture,
			"biography":      user.Biography,
			"motto":          user.Motto,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (r *userRepository) UpdateBlockedStatus(username string, isBlocked bool) error {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (u:User {username: $username}) SET u.isBlocked = $isBlocked"
		params := map[string]any{"username": username, "isBlocked": isBlocked}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}