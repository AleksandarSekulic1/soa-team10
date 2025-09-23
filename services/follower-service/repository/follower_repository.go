package repository

import (
	"context"
	"fmt"
	"log"

	"follower-service/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowerRepository struct {
	driver neo4j.DriverWithContext
}

func NewFollowerRepository(driver neo4j.DriverWithContext) *FollowerRepository {
	return &FollowerRepository{
		driver: driver,
	}
}

// CreateUser creates or updates a user node in Neo4j
func (r *FollowerRepository) CreateUser(ctx context.Context, user domain.User) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// Ako nema ID, koristi username kao ID
	userID := user.ID
	if userID == "" {
		userID = user.Username
	}

	query := `
		MERGE (u:User {id: $id})
		SET u.username = $username, u.email = $email, u.firstName = $firstName, u.lastName = $lastName
		RETURN u
	`

	_, err := session.Run(ctx, query, map[string]interface{}{
		"id":        userID,
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	})

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

// FollowUser creates a FOLLOWS relationship between two users
func (r *FollowerRepository) FollowUser(ctx context.Context, followerID, followingID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User {id: $followerId})
		MATCH (following:User {id: $followingId})
		MERGE (follower)-[f:FOLLOWS]->(following)
		SET f.createdAt = datetime()
		RETURN f
	`

	_, err := session.Run(ctx, query, map[string]interface{}{
		"followerId":  followerID,
		"followingId": followingID,
	})

	if err != nil {
		log.Printf("Error following user: %v", err)
		return err
	}

	return nil
}

// UnfollowUser removes a FOLLOWS relationship between two users
func (r *FollowerRepository) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User {id: $followerId})-[f:FOLLOWS]->(following:User {id: $followingId})
		DELETE f
	`

	_, err := session.Run(ctx, query, map[string]interface{}{
		"followerId":  followerID,
		"followingId": followingID,
	})

	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		return err
	}

	return nil
}

// GetFollowers returns all users who follow the given user
func (r *FollowerRepository) GetFollowers(ctx context.Context, userID string) ([]domain.User, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User)-[:FOLLOWS]->(user:User {id: $userId})
		RETURN follower.id as id, follower.username as username, follower.email as email
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"userId": userID,
	})

	if err != nil {
		log.Printf("Error getting followers: %v", err)
		return nil, err
	}

	var followers []domain.User
	for result.Next(ctx) {
		record := result.Record()
		
		id, _ := record.Get("id")
		username, _ := record.Get("username")
		email, _ := record.Get("email")

		user := domain.User{
			ID:       id.(string),
			Username: username.(string),
			Email:    email.(string),
		}
		followers = append(followers, user)
	}

	return followers, nil
}

// GetFollowing returns all users that the given user follows
func (r *FollowerRepository) GetFollowing(ctx context.Context, userID string) ([]domain.User, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (user:User {id: $userId})-[:FOLLOWS]->(following:User)
		RETURN following.id as id, following.username as username, following.email as email
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"userId": userID,
	})

	if err != nil {
		log.Printf("Error getting following: %v", err)
		return nil, err
	}

	var following []domain.User
	for result.Next(ctx) {
		record := result.Record()

		idVal, _ := record.Get("id")
		usernameVal, _ := record.Get("username")
		emailVal, _ := record.Get("email")

		// Skip if any required field is nil
		if idVal == nil || usernameVal == nil || emailVal == nil {
			continue
		}

		idStr, ok1 := idVal.(string)
		usernameStr, ok2 := usernameVal.(string)
		emailStr, ok3 := emailVal.(string)
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		user := domain.User{
			ID:       idStr,
			Username: usernameStr,
			Email:    emailStr,
		}
		following = append(following, user)
	}

	return following, nil
}

// IsFollowing checks if follower follows following
func (r *FollowerRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User {id: $followerId})-[f:FOLLOWS]->(following:User {id: $followingId})
		RETURN count(f) > 0 as isFollowing
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"followerId":  followerID,
		"followingId": followingID,
	})

	if err != nil {
		log.Printf("Error checking if following: %v", err)
		return false, err
	}

	if result.Next(ctx) {
		record := result.Record()
		isFollowing, _ := record.Get("isFollowing")
		return isFollowing.(bool), nil
	}

	return false, nil
}

// GetRecommendations returns user recommendations based on mutual connections
func (r *FollowerRepository) GetRecommendations(ctx context.Context, userID string, limit int) ([]domain.UserRecommendation, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	// Prvo friends-of-friends, zatim fallback na random korisnike koje korisnik ne prati
	query := `
		// Friends-of-friends preporuke
		MATCH (user:User {id: $userId})-[:FOLLOWS]->(following:User)-[:FOLLOWS]->(recommendation:User)
		WHERE recommendation.id <> $userId
		  AND NOT (user)-[:FOLLOWS]->(recommendation)
		WITH user, recommendation, count(following) as mutualFollowers
		ORDER BY mutualFollowers DESC
		LIMIT $limit
		RETURN recommendation.id as id, recommendation.username as username, recommendation.email as email, mutualFollowers
		UNION
		// Fallback: random korisnici koje korisnik ne prati i nisu on sam
		MATCH (user:User {id: $userId}), (other:User)
		WHERE other.id <> $userId AND NOT (user)-[:FOLLOWS]->(other)
		WITH other LIMIT $limit
		RETURN other.id as id, other.username as username, other.email as email, 0 as mutualFollowers
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"userId": userID,
		"limit":  limit,
	})

	if err != nil {
		log.Printf("Error getting recommendations: %v", err)
		return nil, err
	}

       var recommendations []domain.UserRecommendation
       for result.Next(ctx) {
	       record := result.Record()

	       idVal, _ := record.Get("id")
	       usernameVal, _ := record.Get("username")
	       emailVal, _ := record.Get("email")
	       mutualFollowersVal, _ := record.Get("mutualFollowers")

	       // Skip if any required field is nil
	       if idVal == nil || usernameVal == nil || emailVal == nil || mutualFollowersVal == nil {
		       continue
	       }

	       idStr, ok1 := idVal.(string)
	       usernameStr, ok2 := usernameVal.(string)
	       emailStr, ok3 := emailVal.(string)
	       mutualFollowersInt, ok4 := mutualFollowersVal.(int64)
	       if !ok1 || !ok2 || !ok3 || !ok4 {
		       continue
	       }

	       user := domain.User{
		       ID:       idStr,
		       Username: usernameStr,
		       Email:    emailStr,
	       }

	       recommendation := domain.UserRecommendation{
		       User:            user,
		       MutualFollowers: int(mutualFollowersInt),
		       Reason:          fmt.Sprintf("Followed by %d people you follow", mutualFollowersInt),
	       }
	       recommendations = append(recommendations, recommendation)
       }

       return recommendations, nil
}

// GetFollowedUserIds returns IDs of users that the given user follows (for blog filtering)
func (r *FollowerRepository) GetFollowedUserIds(ctx context.Context, userID string) ([]string, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (user:User {id: $userId})-[:FOLLOWS]->(following:User)
		RETURN following.id as id
	`

	result, err := session.Run(ctx, query, map[string]interface{}{
		"userId": userID,
	})

	if err != nil {
		log.Printf("Error getting followed user IDs: %v", err)
		return nil, err
	}

	var userIds []string
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		userIds = append(userIds, id.(string))
	}

	// Add the user's own ID so they can see their own blogs
	userIds = append(userIds, userID)

	return userIds, nil
}

// GetAllUsers returns all users from the database
func (r *FollowerRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	query := `
		MATCH (u:User)
		RETURN u.id as id, u.username as username, u.email as email, 
		       u.firstName as firstName, u.lastName as lastName
		ORDER BY u.username
	`

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var users []domain.User
	for result.Next(ctx) {
		record := result.Record()
		
		id, _ := record.Get("id")
		username, _ := record.Get("username")
		email, _ := record.Get("email")
		firstName, _ := record.Get("firstName")
		lastName, _ := record.Get("lastName")

		user := domain.User{
			ID:        toString(id),
			Username:  toString(username),
			Email:     toString(email),
			FirstName: toString(firstName),
			LastName:  toString(lastName),
		}
		
		// Ako nema ID, koristi username kao ID
		if user.ID == "" {
			user.ID = user.Username
		}
		
		users = append(users, user)
	}

	return users, result.Err()
}

// Helper function to convert interface{} to string
func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
