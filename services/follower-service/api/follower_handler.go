package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"follower-service/domain"
	"follower-service/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// JWT tajni ključ - treba da bude isti kao u stakeholders servisu
var jwtKey = []byte("super_secret_key")

// Claims struktura za JWT token
type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type FollowerHandler struct {
	followerService *service.FollowerService
}

func NewFollowerHandler(followerService *service.FollowerService) *FollowerHandler {
	return &FollowerHandler{
		followerService: followerService,
	}
}

// CreateUser creates or updates a user
func (h *FollowerHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.followerService.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// FollowUser allows a user to follow another user
func (h *FollowerHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID := h.getUserIDFromToken(r)
	if followerID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var followRequest domain.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&followRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.followerService.FollowUser(r.Context(), followerID, followRequest.FollowingID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully followed user"})
}

// UnfollowUser allows a user to unfollow another user
func (h *FollowerHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerID := h.getUserIDFromToken(r)
	if followerID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var unfollowRequest domain.UnfollowRequest
	if err := json.NewDecoder(r.Body).Decode(&unfollowRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.followerService.UnfollowUser(r.Context(), followerID, unfollowRequest.FollowingID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully unfollowed user"})
}

// GetFollowers returns all followers of a user
func (h *FollowerHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	response, err := h.followerService.GetFollowers(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFollowing returns all users that a user follows
func (h *FollowerHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	response, err := h.followerService.GetFollowing(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// IsFollowing checks if one user follows another
func (h *FollowerHandler) IsFollowing(w http.ResponseWriter, r *http.Request) {
	followerID := h.getUserIDFromToken(r)
	if followerID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followingID := vars["userId"]

	// Check if the authenticated user (followerID) follows the user in the URL (followingID)
	response, err := h.followerService.IsFollowing(r.Context(), followerID, followingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRecommendations returns user recommendations
func (h *FollowerHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromToken(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	limit := 10 // Default limit
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	response, err := h.followerService.GetRecommendations(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFollowedUserIds returns IDs of users that the current user follows (for blog filtering)
func (h *FollowerHandler) GetFollowedUserIds(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromToken(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIds, err := h.followerService.GetFollowedUserIds(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"userIds": userIds})
}

// getUserIDFromToken extracts user ID from JWT token
func (h *FollowerHandler) getUserIDFromToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return ""
	}

	tokenString := tokenParts[1]
	
	// Parse the JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return ""
	}

	return claims.ID
}

// GetAllUsers returns all users from the system
func (h *FollowerHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.followerService.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
