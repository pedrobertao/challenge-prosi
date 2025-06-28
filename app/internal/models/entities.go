package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogPost struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	Comments  []Comment          `json:"comments,omitempty" bson:"-"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PostID    primitive.ObjectID `json:"post_id" bson:"post_id"`
	Author    string             `json:"author" bson:"author"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type BlogPostSummary struct {
	ID           primitive.ObjectID `json:"id"`
	Title        string             `json:"title"`
	CommentCount int64              `json:"comment_count"`
	CreatedAt    time.Time          `json:"createdAt"`
}
