// File: github_activity/github_activity.go
package github_activity

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// GitHubEvent represents a GitHub event retrieved from the API
type GitHubEvent struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Repo      struct {
		Name string `json:"name"`
	} `json:"repo"`
}

// GitHubService interface represents operations to interact with GitHub
type GitHubService interface {
	GetRecentActivity(username string) ([]GitHubEvent, error)
}

// DefaultGitHubService implements GitHubService and fetches data from the GitHub API
type DefaultGitHubService struct {
	client *http.Client
}

// NewGitHubService creates a new instance of DefaultGitHubService with a given HTTP client
func NewGitHubService(client *http.Client) GitHubService {
	return &DefaultGitHubService{client: client}
}

// GetRecentActivity fetches recent GitHub events for the given username
func (s *DefaultGitHubService) GetRecentActivity(username string) ([]GitHubEvent, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var events []GitHubEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return events, nil
}

// PrintRecentActivity prints the recent GitHub activity of a user in a formatted way
func PrintRecentActivity(username string, service GitHubService) {
	events, err := service.GetRecentActivity(username)
	if err != nil {
		log.Fatalf("Error fetching activity: %v", err)
	}

	fmt.Printf("Recent activity for user %s:\n", username)
	for _, event := range events {
		fmt.Printf("- %s on %s at %s\n", event.Type, event.Repo.Name, event.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
