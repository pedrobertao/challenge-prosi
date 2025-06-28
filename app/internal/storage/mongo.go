// Package storage provides MongoDB database connection and collection management
// for the blog application. It handles database initialization, connection pooling,
// and provides easy access to the required collections (posts and comments).
package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage struct encapsulates MongoDB client and collection references.
// It provides a centralized way to access database operations for the blog application.
// The struct maintains references to specific collections to avoid repeated lookups.
type Storage struct {
	Client   *mongo.Client     // MongoDB client for database operations
	Posts    *mongo.Collection // Collection for blog posts
	Comments *mongo.Collection // Collection for post comments
}

// Connect establishes a connection to MongoDB and initializes the Storage struct.
// It creates a new MongoDB client, tests the connection with a ping, and sets up
// collection references for posts and comments.
//
// Connection process:
//  1. Creates MongoDB client with provided URI
//  2. Tests connection with ping operation
//  3. Initializes database and collection references
//  4. Returns configured Storage instance
//
// Parameters:
//   - uri: MongoDB connection string (e.g., "mongodb://localhost:27017")
//   - dbName: name of the database to use (e.g., "blog")
//
// Returns:
//   - *Storage: configured storage instance with active connections
//   - error: connection error if any step fails
func Connect(uri, dbName string) (*Storage, error) {
	// Create context with timeout to prevent hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Establish connection to MongoDB server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Test the connection by pinging the server
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// Get database reference and collection handles
	db := client.Database(dbName)
	postsCol := db.Collection("posts")       // Collection for blog posts
	commentsCol := db.Collection("comments") // Collection for post comments

	// Return configured Storage instance with all references
	return &Storage{
		Client:   client,
		Posts:    postsCol,
		Comments: commentsCol,
	}, nil
}

// Close gracefully shuts down the MongoDB connection.
// This method should be called when the application is shutting down
// to ensure proper cleanup of database connections and resources.
//
// Parameters:
//   - ctx: context for controlling the disconnect timeout
//
// Returns error if the disconnect operation fails.
func (db *Storage) Close(ctx context.Context) error {
	// Disconnect the MongoDB client and clean up resources
	return db.Client.Disconnect(ctx)
}
