package unit

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pedrobertao/challenge-prosi/app/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockCursor implements the necessary mongo.Cursor methods for testing.
// This struct simulates MongoDB cursor behavior by maintaining an internal
// slice of posts and an index to track iteration position.
type MockCursor struct {
	posts []models.BlogPost // Collection of blog posts to iterate through
	index int               // Current position in the posts slice
}

// Next simulates MongoDB cursor's Next() method.
// Returns true if there are more documents to process, false otherwise.
// This method is called by the handler to check if more posts are available.
func (m *MockCursor) Next(ctx context.Context) bool {
	return m.index < len(m.posts)
}

// Decode simulates MongoDB cursor's Decode() method.
// Copies the current post data into the provided interface.
// Increments the index to move to the next post after successful decode.
//
// Parameters:
//   - val: pointer to the struct where decoded data should be stored
//
// Returns error if no more documents are available.
func (m *MockCursor) Decode(val interface{}) error {
	// Check if we've reached the end of available posts
	if m.index >= len(m.posts) {
		return mongo.ErrNoDocuments
	}

	// Cast the interface to BlogPost pointer and copy current post data
	post := val.(*models.BlogPost)
	*post = m.posts[m.index]

	// Move to next post for subsequent calls
	m.index++
	return nil
}

// Close simulates MongoDB cursor's Close() method.
// In real MongoDB operations, this releases cursor resources.
// For testing purposes, this is a no-op that always succeeds.
func (m *MockCursor) Close(ctx context.Context) error {
	return nil
}

// MockCollection simulates MongoDB collection operations for testing.
// Uses testify/mock to track method calls and return configured responses.
type MockCollection struct {
	mock.Mock // Embedding mock.Mock provides expectation and assertion capabilities
}

// MockHandler simulates the actual handler structure for testing.
// Contains references to mock collections to control database behavior during tests.
type MockHandler struct {
	mock.Mock                 // Provides mock capabilities for handler methods
	Posts     *MockCollection // Mock posts collection for database operations
	Comments  *MockCollection // Mock comments collection for counting operations
}

// Find simulates MongoDB collection's Find() method.
// Returns a cursor and potential error based on configured mock expectations.
//
// Parameters:
//   - ctx: context for the database operation
//   - filter: BSON filter criteria (typically bson.M{} for all documents)
//
// Returns configured mock cursor and error from test expectations.
func (m *MockCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

// CountDocuments simulates MongoDB collection's CountDocuments() method.
// Returns the number of documents matching the filter criteria.
//
// Parameters:
//   - ctx: context for the database operation
//   - filter: BSON filter criteria (typically contains post_id for comment counting)
//
// Returns configured mock count and error from test expectations.
func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// GetPosts simulates the actual GetPosts handler method for testing.
// Returns a hardcoded JSON response that matches the expected API format.
// This allows testing the HTTP layer without involving real database operations.
//
// Parameters:
//   - c: Fiber context containing request/response information
//
// Returns error if response sending fails (typically nil in tests).
func (m *MockHandler) GetPosts(c *fiber.Ctx) error {
	// Return a predefined JSON response that matches the expected API format
	// This response contains two blog posts with different comment counts
	return c.Status(200).Send([]byte(`{
		"success": true,
		"error": "",
		"data": [
			{
				"id": "686c3a82361beb165141b490",
				"title": "First Post",
				"comment_count": 3,
				"created_at": "2025-07-07T16:00:00Z"
			},
			{
				"id": "686c3a82361beb165141b491",
				"title": "Second Post",
				"comment_count": 1,
				"created_at": "2025-07-07T17:00:00Z"
			}
		]
	}`))
}

// TestGetPostsSuccess tests the successful retrieval of blog posts.
// This test verifies that the GetPosts handler correctly:
// 1. Fetches posts from the database
// 2. Counts comments for each post
// 3. Returns properly formatted API response
// 4. Handles the complete request-response cycle
func TestGetPostsSuccess(t *testing.T) {
	// === SETUP PHASE ===
	// Create specific ObjectIDs that match the expected response
	// These IDs are used to ensure consistency between test data and expected output
	id1, _ := primitive.ObjectIDFromHex("686c3a82361beb165141b490")
	id2, _ := primitive.ObjectIDFromHex("686c3a82361beb165141b491")

	// Create test blog posts that simulate real database records
	// Each post has realistic content and timestamps offset from current time
	testPosts := []models.BlogPost{
		{
			ID:        id1,
			Title:     "First Post",
			Content:   "This is the first post content", // Content won't appear in summary
			CreatedAt: time.Now().Add(-2 * time.Hour),   // Posted 2 hours ago
		},
		{
			ID:        id2,
			Title:     "Second Post",
			Content:   "This is the second post content", // Content won't appear in summary
			CreatedAt: time.Now().Add(-1 * time.Hour),    // Posted 1 hour ago
		},
	}

	// === MOCK SETUP PHASE ===
	// Create mock cursor that will iterate through our test posts
	// The cursor starts at index 0 and will return each post in sequence
	mockCursor := &MockCursor{
		posts: testPosts,
		index: 0, // Start at beginning of posts array
	}

	// Create mock collections to simulate database operations
	mockPosts := &MockCollection{}    // Will handle Posts.Find() operations
	mockComments := &MockCollection{} // Will handle Comments.CountDocuments() operations

	// === MOCK EXPECTATIONS SETUP ===
	// Configure the posts collection to return our mock cursor when Find() is called
	// The bson.M{} parameter represents an empty filter (get all posts)
	mockPosts.On("Find", mock.Anything, bson.M{}).Return((mockCursor), nil)

	// Configure comment counts for each post
	// First post should have 3 comments
	mockComments.On("CountDocuments", mock.Anything, bson.M{"post_id": testPosts[0].ID}).Return(int64(3), nil)
	// Second post should have 1 comment
	mockComments.On("CountDocuments", mock.Anything, bson.M{"post_id": testPosts[1].ID}).Return(int64(1), nil)

	// Create handler instance with mock collections
	handler := &MockHandler{
		Posts:    mockPosts,    // Mock posts collection
		Comments: mockComments, // Mock comments collection
	}

	// === HTTP TEST SETUP ===
	// Create Fiber application instance for testing HTTP endpoints
	app := fiber.New()
	// Register the GetPosts handler for the /api/posts route
	app.Get("/api/posts", handler.GetPosts)

	// Create HTTP GET request to test the endpoint
	req := httptest.NewRequest("GET", "/api/posts", nil)

	// === EXECUTION PHASE ===
	// Execute the HTTP request against our test application
	// This triggers the entire request-response cycle
	resp, err := app.Test(req)

	// === ASSERTION PHASE ===
	// Verify that the request executed without errors
	assert.NoError(t, err)
	// Verify that the response has HTTP 200 OK status
	assert.Equal(t, 200, resp.StatusCode)

	// === RESPONSE PARSING ===
	// Parse the JSON response body into our API response model
	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	// === RESPONSE STRUCTURE VERIFICATION ===
	// Verify that the API response indicates success
	assert.True(t, response.Success)
	// Verify that no error message is present
	assert.Empty(t, response.Error)
	// Verify that response contains data
	assert.NotNil(t, response.Data)

	// === DATA CONTENT VERIFICATION ===
	// Convert the generic data interface to a slice of interfaces
	// JSON unmarshaling creates []interface{} for arrays
	summariesData, ok := response.Data.([]interface{})
	assert.True(t, ok)
	// Verify that exactly 2 posts are returned
	assert.Len(t, summariesData, 2)

	// === FIRST POST VERIFICATION ===
	// Extract first post summary and verify its content
	firstSummary := summariesData[0].(map[string]interface{})
	// Verify the title matches our test data
	assert.Equal(t, testPosts[0].Title, firstSummary["title"])
	// Verify comment count (JSON numbers are unmarshaled as float64)
	assert.Equal(t, float64(3), firstSummary["comment_count"])

	// === SECOND POST VERIFICATION ===
	// Extract second post summary and verify its content
	secondSummary := summariesData[1].(map[string]interface{})
	// Verify the title matches our test data
	assert.Equal(t, testPosts[1].Title, secondSummary["title"])
	// Verify comment count matches expected value
	assert.Equal(t, float64(1), secondSummary["comment_count"])

}
