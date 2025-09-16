package domain

import "time"

// User represents a user in the follower system
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Follow represents a follow relationship between users
type Follow struct {
	FollowerID  string    `json:"followerId"`
	FollowingID string    `json:"followingId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// FollowRequest represents a request to follow a user
type FollowRequest struct {
	FollowingID string `json:"followingId"`
}

// UnfollowRequest represents a request to unfollow a user
type UnfollowRequest struct {
	FollowingID string `json:"followingId"`
}

// FollowersResponse represents the response for getting followers
type FollowersResponse struct {
	Followers []User `json:"followers"`
	Count     int    `json:"count"`
}

// FollowingResponse represents the response for getting following
type FollowingResponse struct {
	Following []User `json:"following"`
	Count     int    `json:"count"`
}

// RecommendationsResponse represents the response for getting recommendations
type RecommendationsResponse struct {
	Recommendations []UserRecommendation `json:"recommendations"`
	Count           int                  `json:"count"`
}

// UserRecommendation represents a user recommendation with mutual connections
type UserRecommendation struct {
	User            User   `json:"user"`
	MutualFollowers int    `json:"mutualFollowers"`
	Reason          string `json:"reason"`
}

// IsFollowingResponse represents the response for checking if user follows another
type IsFollowingResponse struct {
	IsFollowing bool `json:"isFollowing"`
}
