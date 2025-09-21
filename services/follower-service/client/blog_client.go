package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BlogClient struct {
	baseURL    string
	httpClient *http.Client
}

// Request structs
type RemoveLikesRequest struct {
	UserID   string `json:"userId"`
	AuthorID string `json:"authorId"`
}

// Response structs
type BlogServiceResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func NewBlogClient(baseURL string) *BlogClient {
	return &BlogClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RemoveLikesFromAuthorBlogs poziva blog servis da ukloni sve lajkove određenog korisnika sa blogova određenog autora
func (c *BlogClient) RemoveLikesFromAuthorBlogs(userID, authorID string) error {
	reqBody := RemoveLikesRequest{
		UserID:   userID,
		AuthorID: authorID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("%s/api/blogs/remove-likes", c.baseURL)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("blog service returned status %d", resp.StatusCode)
	}

	var response BlogServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("blog service operation failed: %s", response.Message)
	}

	return nil
}