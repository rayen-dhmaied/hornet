package service

import (
	"context"
	"fmt"
	"hornet/api/followers/model"
	"hornet/api/followers/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FollowersService defines the methods for handling followers-related business logic
type FollowersService struct {
	followersRepository *repository.FollowersRepository
}

var (
	followersServiceInstance *FollowersService
	followersOnce            sync.Once
)

// NewFollowersService creates a singleton instance of FollowersService.
func NewFollowersService(follRepository *repository.FollowersRepository) *FollowersService {
	if follRepository == nil {
		panic("repository cannot be nil")
	}

	followersOnce.Do(func() {
		followersServiceInstance = &FollowersService{
			followersRepository: follRepository,
		}
	})
	return followersServiceInstance
}

// CreateFollow creates a new follow relationship.
func (s *FollowersService) CreateFollow(ctx context.Context, senderID, receiverID uuid.UUID) (*model.Follow, error) {
	if senderID == receiverID {
		return nil, fmt.Errorf("sender and receiver IDs cannot be the same")
	}

	follow := &model.Follow{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		CreatedAt:  time.Now().UTC(),
	}

	savedFollow, err := s.followersRepository.CreateFollow(ctx, follow)
	if err != nil {
		return nil, fmt.Errorf("failed to create follow: %w", err)
	}

	return savedFollow, nil
}

// DeleteFollow deletes an existing follow relationship.
func (s *FollowersService) DeleteFollow(ctx context.Context, userID, followID uuid.UUID) error {
	if userID == uuid.Nil || followID == uuid.Nil {
		return fmt.Errorf("userID and followID cannot be nil")
	}

	err := s.followersRepository.DeleteFollow(ctx, followID)
	if err != nil {
		return fmt.Errorf("failed to delete follow with ID %s: %w", followID, err)
	}

	return nil
}

// GetFollowers retrieves a list of followers for a specific user.
func (s *FollowersService) GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("userID cannot be nil")
	}

	followers, err := s.followersRepository.GetFollowers(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers for user with ID %s: %w", userID, err)
	}

	return followers, nil
}

// GetFollowing retrieves a list of users a specific user is following.
func (s *FollowersService) GetFollowing(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("userID cannot be nil")
	}

	following, err := s.followersRepository.GetFollowing(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following for user with ID %s: %w", userID, err)
	}

	return following, nil
}
