package repository

import (
	"context"
	"hornet/api/followers/model"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// FollowersRepository defines the methods for interacting with the database for followers.
type FollowersRepository struct {
	driver neo4j.DriverWithContext
}

// Global variables for the singleton instance of FollowersRepository.
var (
	followersRepositoryInstance *FollowersRepository
	followersOnce               sync.Once
)

// NewFollowersRepository creates a new FollowersRepository instance if it doesn't exist.
func NewFollowersRepository(driver neo4j.DriverWithContext) *FollowersRepository {
	followersOnce.Do(func() {
		followersRepositoryInstance = &FollowersRepository{
			driver: driver,
		}
	})
	return followersRepositoryInstance
}

// CreateFollow saves a new follow relationship to the database.
func (r *FollowersRepository) CreateFollow(ctx context.Context, follow *model.Follow) (*model.Follow, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MERGE (sender:User {id: $senderID})
			MERGE (receiver:User {id: $receiverID})
			CREATE (sender)-[f:FOLLOW {id: $id, created_at: $createdAt}]->(receiver)
			RETURN f
		`
		params := map[string]interface{}{
			"id":         follow.ID.String(),
			"senderID":   follow.SenderID.String(),
			"receiverID": follow.ReceiverID.String(),
			"createdAt":  follow.CreatedAt,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	if err != nil {
		return nil, err
	}
	return follow, nil
}

// DeleteFollow deletes a follow relationship by its ID.
func (r *FollowersRepository) DeleteFollow(ctx context.Context, id uuid.UUID) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH ()-[f:FOLLOW {id: $id}]->()
			DELETE f
		`
		params := map[string]interface{}{
			"id": id.String(),
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

// GetFollowers retrieves a list of followers for a given user ID.
func (r *FollowersRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	var followers []model.Follow
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (receiver:User {id: $userID})<-[f:FOLLOW]-(sender:User)
			RETURN f.id AS id, sender.id AS senderID, receiver.id AS receiverID, f.created_at AS createdAt
		`
		params := map[string]interface{}{
			"userID": userID.String(),
		}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		// Iterate over the results.
		for result.Next(ctx) {
			record := result.Record()

			// Retrieve each value and ensure it exists.
			idVal, ok := record.Get("id")
			if !ok {
				continue
			}
			senderIDVal, ok := record.Get("senderID")
			if !ok {
				continue
			}
			receiverIDVal, ok := record.Get("receiverID")
			if !ok {
				continue
			}
			createdAtVal, ok := record.Get("createdAt")
			if !ok {
				continue
			}

			followers = append(followers, model.Follow{
				ID:         uuid.MustParse(idVal.(string)),
				SenderID:   uuid.MustParse(senderIDVal.(string)),
				ReceiverID: uuid.MustParse(receiverIDVal.(string)),
				CreatedAt:  createdAtVal.(time.Time),
			})
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}
	return followers, nil
}

// GetFollowing retrieves a list of users a given user is following.
func (r *FollowersRepository) GetFollowing(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	var following []model.Follow
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (sender:User {id: $userID})-[f:FOLLOW]->(receiver:User)
			RETURN f.id AS id, sender.id AS senderID, receiver.id AS receiverID, f.created_at AS createdAt
		`
		params := map[string]interface{}{
			"userID": userID.String(),
		}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		// Iterate over the results.
		for result.Next(ctx) {
			record := result.Record()

			idVal, ok := record.Get("id")
			if !ok {
				continue
			}
			senderIDVal, ok := record.Get("senderID")
			if !ok {
				continue
			}
			receiverIDVal, ok := record.Get("receiverID")
			if !ok {
				continue
			}
			createdAtVal, ok := record.Get("createdAt")
			if !ok {
				continue
			}

			following = append(following, model.Follow{
				ID:         uuid.MustParse(idVal.(string)),
				SenderID:   uuid.MustParse(senderIDVal.(string)),
				ReceiverID: uuid.MustParse(receiverIDVal.(string)),
				CreatedAt:  createdAtVal.(time.Time),
			})
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}
	return following, nil
}
