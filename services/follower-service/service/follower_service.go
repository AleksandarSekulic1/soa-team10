package service

import (
	"context"
	"errors"

	"follower-service/client"
	"follower-service/domain"
	"follower-service/repository"
)

type FollowerService struct {
	followerRepo     *repository.FollowerRepository
	sagaOrchestrator *SagaOrchestrator
}

func NewFollowerService(followerRepo *repository.FollowerRepository, blogClient *client.BlogClient) *FollowerService {
	sagaOrchestrator := NewSagaOrchestrator(followerRepo, blogClient)
	return &FollowerService{
		followerRepo:     followerRepo,
		sagaOrchestrator: sagaOrchestrator,
	}
}

// CreateUser creates or updates a user
func (s *FollowerService) CreateUser(ctx context.Context, user domain.User) error {
	if user.ID == "" {
		return errors.New("user ID is required")
	}
	if user.Username == "" {
		return errors.New("username is required")
	}

	return s.followerRepo.CreateUser(ctx, user)
}

// FollowUser allows one user to follow another
func (s *FollowerService) FollowUser(ctx context.Context, followerID, followingID string) error {
	if followerID == "" || followingID == "" {
		return errors.New("both follower and following IDs are required")
	}

	if followerID == followingID {
		return errors.New("users cannot follow themselves")
	}

	// Check if already following
	isFollowing, err := s.followerRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}

	if isFollowing {
		return errors.New("already following this user")
	}

	return s.followerRepo.FollowUser(ctx, followerID, followingID)
}

// UnfollowUser allows one user to unfollow another using saga pattern
func (s *FollowerService) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	if followerID == "" || followingID == "" {
		return errors.New("both follower and following IDs are required")
	}

	//Odkomentarisi kad se popravi follow service
	
	// Check if actually following
	// isFollowing, err := s.followerRepo.IsFollowing(ctx, followerID, followingID)
	// if err != nil {
	// 	return err
	// }

	// if !isFollowing {
	// 	return errors.New("not following this user")
	// }

	// Execute saga pattern for unfollow operation
	// This will:
	// 1. Remove follow relationship from Neo4j
	// 2. Remove all likes from the unfollowed user's blogs in blog service
	// If any step fails, compensations will be executed to rollback
	return s.sagaOrchestrator.ExecuteUnfollowSaga(ctx, followerID, followingID)
}

// GetFollowers returns all followers of a user
func (s *FollowerService) GetFollowers(ctx context.Context, userID string) (domain.FollowersResponse, error) {
	if userID == "" {
		return domain.FollowersResponse{}, errors.New("user ID is required")
	}

	followers, err := s.followerRepo.GetFollowers(ctx, userID)
	if err != nil {
		return domain.FollowersResponse{}, err
	}

	return domain.FollowersResponse{
		Followers: followers,
		Count:     len(followers),
	}, nil
}

// GetFollowing returns all users that a user follows
func (s *FollowerService) GetFollowing(ctx context.Context, userID string) (domain.FollowingResponse, error) {
	if userID == "" {
		return domain.FollowingResponse{}, errors.New("user ID is required")
	}

	following, err := s.followerRepo.GetFollowing(ctx, userID)
	if err != nil {
		return domain.FollowingResponse{}, err
	}

	return domain.FollowingResponse{
		Following: following,
		Count:     len(following),
	}, nil
}

// IsFollowing checks if one user follows another
func (s *FollowerService) IsFollowing(ctx context.Context, followerID, followingID string) (domain.IsFollowingResponse, error) {
	if followerID == "" || followingID == "" {
		return domain.IsFollowingResponse{}, errors.New("both follower and following IDs are required")
	}

	isFollowing, err := s.followerRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return domain.IsFollowingResponse{}, err
	}

	return domain.IsFollowingResponse{
		IsFollowing: isFollowing,
	}, nil
}

// GetRecommendations returns user recommendations based on mutual connections
func (s *FollowerService) GetRecommendations(ctx context.Context, userID string, limit int) (domain.RecommendationsResponse, error) {
	if userID == "" {
		return domain.RecommendationsResponse{}, errors.New("user ID is required")
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	recommendations, err := s.followerRepo.GetRecommendations(ctx, userID, limit)
	if err != nil {
		return domain.RecommendationsResponse{}, err
	}

	return domain.RecommendationsResponse{
		Recommendations: recommendations,
		Count:           len(recommendations),
	}, nil
}

// GetFollowedUserIds returns IDs of users that the given user follows (for blog filtering)
func (s *FollowerService) GetFollowedUserIds(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.followerRepo.GetFollowedUserIds(ctx, userID)
}

// GetAllUsers returns all users from the system
func (s *FollowerService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.followerRepo.GetAllUsers(ctx)
}
