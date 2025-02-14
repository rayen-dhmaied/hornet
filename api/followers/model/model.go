package model

import (
	"time"

	"github.com/google/uuid"
)

// Follow represents a :FOLLOW relationship in Neo4j.
// Typically, you'd store these fields as properties on the relationship.
type Follow struct {
	ID         uuid.UUID `json:"id"`          // Could be a UUID stored as a property on the relationship
	SenderID   uuid.UUID `json:"sender_id"`   // The ID of the user node that initiated the follow
	ReceiverID uuid.UUID `json:"receiver_id"` // The ID of the user node that received the follow
	CreatedAt  time.Time `json:"created_at"`  // When the relationship was created
}

// CreateFollow represents the request body for creating a new follow.
type CreateFollow struct {
	ReceiverID string `json:"receiver_id"` // The ID of the user node that received the Follow
}
