package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	Client   *mongo.Client
	Posts    *mongo.Collection
	Comments *mongo.Collection
}

func Connect(uri, dbName string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	postsCol := db.Collection("posts")
	commentsCol := db.Collection("comments")
	return &Storage{
		Client:   client,
		Posts:    postsCol,
		Comments: commentsCol,
	}, nil
}

func (db *Storage) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}
