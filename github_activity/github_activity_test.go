package github_activity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockGitHubService wraps around DefaultGitHubService for testing
// Takes a mock base URL to direct requests to the mock server
type MockGitHubService struct {
	baseURL string
	client  *http.Client
}

// NewMockGitHubService creates a new mock service with a mock base URL
func NewMockGitHubService(mockURL string) GitHubService {
	client := &http.Client{}
	return &MockGitHubService{baseURL: mockURL, client: client}
}

// GetRecentActivity fetches recent activity using the mock server's base URL
func (s *MockGitHubService) GetRecentActivity(username string) ([]GitHubEvent, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	url := fmt.Sprintf("%s/users/%s/events", s.baseURL, username)
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

// Helper function to create a mock server for testing
func createMockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))
}

// Test for successful retrieval of recent activity
func TestGitHubService_GetRecentActivity_Success(t *testing.T) {
	// Mock response from GitHub API
	mockResponse := `[
		{
			"type": "PushEvent",
			"created_at": "2023-10-10T14:12:15Z",
			"repo": { "name": "example/repo1" }
		},
		{
			"type": "IssueCommentEvent",
			"created_at": "2023-10-11T15:13:16Z",
			"repo": { "name": "example/repo2" }
		}
	]`

	// Start a mock server to return the mocked GitHub response
	mockServer := createMockServer(mockResponse, http.StatusOK)
	defer mockServer.Close()

	// Create a mock service that points to the mock server
	service := NewMockGitHubService(mockServer.URL)

	// Call the service
	events, err := service.GetRecentActivity("octocat")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate the response
	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	if events[0].Type != "PushEvent" || events[0].Repo.Name != "example/repo1" {
		t.Errorf("Unexpected event data for first event")
	}
}

// Test for GitHub API returning 404 error
func TestGitHubService_GetRecentActivity_NotFound(t *testing.T) {
	// Start a mock server to return a 404 status code
	mockServer := createMockServer("{}", http.StatusNotFound)
	defer mockServer.Close()

	// Create a mock service that points to the mock server
	service := NewMockGitHubService(mockServer.URL)

	// Call the service with an invalid username
	_, err := service.GetRecentActivity("nonexistent_user")
	if err == nil {
		t.Fatalf("Expected error for nonexistent user, got none")
	}

	expectedError := fmt.Sprintf("GitHub API returned status: %d", http.StatusNotFound)
	if err.Error() != expectedError {
		t.Errorf("Expected error: %v, got: %v", expectedError, err)
	}
}

// Test for handling invalid JSON response
func TestGitHubService_GetRecentActivity_InvalidJSON(t *testing.T) {
	// Start a mock server to return an invalid JSON
	mockServer := createMockServer(`{"invalid_json"}`, http.StatusOK)
	defer mockServer.Close()

	// Create a mock service that points to the mock server
	service := NewMockGitHubService(mockServer.URL)

	// Call the service
	_, err := service.GetRecentActivity("octocat")
	if err == nil {
		t.Fatalf("Expected error for invalid JSON, got none")
	}
}

// Test for handling empty activity (user has no events)
func TestGitHubService_GetRecentActivity_EmptyResponse(t *testing.T) {
	// Start a mock server to return an empty array
	mockServer := createMockServer(`[]`, http.StatusOK)
	defer mockServer.Close()

	// Create a mock service that points to the mock server
	service := NewMockGitHubService(mockServer.URL)

	// Call the service
	events, err := service.GetRecentActivity("octocat")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate that there are no events
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

// Test for empty username
func TestGitHubService_GetRecentActivity_EmptyUsername(t *testing.T) {
	client := &http.Client{}
	service := &DefaultGitHubService{client: client}

	// Call with empty username
	_, err := service.GetRecentActivity("")
	if err == nil {
		t.Fatalf("Expected error for empty username, got none")
	}

	expectedError := "username cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error: %v, got: %v", expectedError, err)
	}
}
