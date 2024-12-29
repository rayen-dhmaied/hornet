package model

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a post document in MongoDB
type Post struct {
	ID             uuid.UUID  `bson:"_id,omitempty" json:"id"`
	Content        string     `bson:"content,omitempty" json:"content,omitempty" validate:"max=5000"`
	AuthorID       uuid.UUID  `bson:"author_id" json:"author_id" validate:"required"`
	ParentPostID   *uuid.UUID `bson:"parent_post_id,omitempty" json:"parent_post_id,omitempty"`     // For replies
	OriginalPostID *uuid.UUID `bson:"original_post_id,omitempty" json:"original_post_id,omitempty"` // For shared posts
	RepliesCount   int        `bson:"replies_count" json:"replies_count"`                           // For tracking nested replies
	SharesCount    int        `bson:"shares_count" json:"shares_count"`                             // For tracking shared posts
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
}

// CreatePost represents the structure of a new post creation request
type CreatePost struct {
	Content        *string    `json:"content,omitempty"`          // Content is optional, only when original_post_id is provided
	ParentPostID   *uuid.UUID `json:"parent_post_id,omitempty"`   // ID of the parent post if it's a reply, can be nil
	OriginalPostID *uuid.UUID `json:"original_post_id,omitempty"` // ID of the original post being shared, can be nil
	AuthorID       uuid.UUID  `json:"author_id"`                  // AuthorID is required and represents the user making the post
}
