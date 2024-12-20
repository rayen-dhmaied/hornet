package service

import (
	"context"
	"fmt"
	"hornet/api/posts/model"
	"hornet/api/posts/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

// PostService defines the methods for handling post-related business logic
type PostService struct {
	postRepository *repository.PostRepository
}

// Declare a global variable for the singleton instance of PostService
var (
	postServiceInstance *PostService
	once                sync.Once
)

// NewPostService creates a new PostService instance if it doesn't exist
func NewPostService(postRepository *repository.PostRepository) *PostService {
	once.Do(func() {
		postServiceInstance = &PostService{
			postRepository: postRepository,
		}
	})
	return postServiceInstance
}

// GetPost retrieves a post by its ID
func (s *PostService) GetPost(ctx context.Context, postID uuid.UUID) (model.Post, error) {

	// Fetch the post from the repository
	post, err := s.postRepository.FindPostByID(ctx, postID)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

// GetReplies retrieves all replies for a given parent post
func (s *PostService) GetReplies(ctx context.Context, parentPostID uuid.UUID) ([]model.Post, error) {
	// Call the repository to fetch all posts with the given ParentPostID
	replies, err := s.postRepository.FindPostsByParentID(ctx, parentPostID)
	if err != nil {
		return nil, err
	}

	return replies, nil
}

// CreatePost handles the creation of a new post
func (s *PostService) CreatePost(ctx context.Context, req model.CreatePost) (model.Post, error) {
	// Generate a new Post ID
	postID := uuid.New()

	// Extract the content from the request
	var content string
	if req.Content != nil {
		content = *req.Content
	}
	// Create a new post instance
	post := model.Post{
		ID:             postID,
		AuthorID:       req.AuthorID,
		Content:        content,
		ParentPostID:   req.ParentPostID,   // Only set if it's a reply
		OriginalPostID: req.OriginalPostID, // Only set if it's shared
		RepliesCount:   0,                  // Initialize replies count to 0
		SharesCount:    0,                  // Initialize shares count to 0
		CreatedAt:      time.Now(),
	}

	// Insert the post into the database using the repository
	err := s.postRepository.SavePost(ctx, post)
	if err != nil {
		return model.Post{}, err
	}

	// If it's a reply, increment the replies count on the parent post
	if req.ParentPostID != nil {
		err = s.IncrementRepliesCount(ctx, *req.ParentPostID)
		if err != nil {
			// Log the error and continue with the successfully created post
			fmt.Printf("Failed to increment replies count for post %s: %v\n", req.ParentPostID, err)
		}
	}

	// If it's a shared post, increment the shares count on the original post
	if req.OriginalPostID != nil {
		err = s.IncrementSharesCount(ctx, *req.OriginalPostID)
		if err != nil {
			// Log the error and continue with the successfully created post
			fmt.Printf("Failed to increment shares count for post %s: %v\n", req.OriginalPostID, err)
		}
	}
	return post, nil
}

// DeletePost handles the deletion of a post
func (s *PostService) DeletePost(ctx context.Context, postID uuid.UUID, params ...string) error {
	// Fetch the post to be deleted
	post, err := s.postRepository.FindPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post with ID %s: %v", postID, err)
	}

	err = s.postRepository.DeletePost(ctx, post.ID)
	if err != nil {
		return fmt.Errorf("failed to delete post with ID %s: %v", postID, err)
	}

	if post.RepliesCount > 0 {
		// Fetch all replies for the post
		replies, err := s.GetReplies(ctx, post.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch replies for post with ID %s: %v", postID, err)
		}

		// Delete all replies
		for _, reply := range replies {
			err = s.DeletePost(ctx, reply.ID, "cascade")
			if err != nil {
				// Log the error and continue with the successfully deleted post
				fmt.Printf("Failed to delete reply %s: %v\n", reply.ID, err)
			}
		}
	}

	// If it's a reply, decrement the replies count on the parent post if not deleted with cascade
	if post.ParentPostID != nil && len(params) > 0 && params[0] != "cascade" {
		err = s.DecrementRepliesCount(ctx, *post.ParentPostID)
		if err != nil {
			// Log the error and continue with the successfully deleted post
			fmt.Printf("Failed to decrement replies count for post %s: %v\n", post.ParentPostID, err)
		}
	}

	// If it's a shared post, decrement the shares count on the original post
	if post.OriginalPostID != nil {
		err = s.DecrementSharesCount(ctx, *post.OriginalPostID)
		if err != nil {
			// Log the error and continue with the successfully deleted post
			fmt.Printf("Failed to decrement shares count for post %s: %v\n", post.OriginalPostID, err)
		}
	}
	return nil
}

// IncrementRepliesCount increments the replies count for a given post
func (s *PostService) IncrementRepliesCount(ctx context.Context, postID uuid.UUID) error {
	// Fetch the parent post
	post, err := s.postRepository.FindPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post with ID %s: %v", postID, err)
	}

	// Increment the replies count
	post.RepliesCount++

	// Save the updated post back to the repository
	err = s.postRepository.UpdatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to update replies count for post with ID %s: %v", postID, err)
	}

	return nil
}

// DecrementRepliesCount decrements the replies count for a given post
func (s *PostService) DecrementRepliesCount(ctx context.Context, postID uuid.UUID) error {
	// Fetch the parent post
	post, err := s.postRepository.FindPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post with ID %s: %v", postID, err)
	}

	// Decrement the replies count
	if post.RepliesCount > 0 {
		post.RepliesCount--
	}

	// Save the updated post back to the repository
	err = s.postRepository.UpdatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to update replies count for post with ID %s: %v", postID, err)
	}

	return nil
}

// IncrementRepliesCount increments the shares count for a given post
func (s *PostService) IncrementSharesCount(ctx context.Context, postID uuid.UUID) error {
	// Fetch the parent post
	post, err := s.postRepository.FindPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post with ID %s: %v", postID, err)
	}

	// Increment the shares count
	post.SharesCount++

	// Save the updated post back to the repository
	err = s.postRepository.UpdatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to update shares count for post with ID %s: %v", postID, err)
	}

	return nil
}

// DecrementSharesCount decrements the shares count for a given post
func (s *PostService) DecrementSharesCount(ctx context.Context, postID uuid.UUID) error {
	// Fetch the parent post
	post, err := s.postRepository.FindPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to find post with ID %s: %v", postID, err)
	}

	// Decrement the shares count
	if post.SharesCount > 0 {
		post.SharesCount--
	}

	// Save the updated post back to the repository
	err = s.postRepository.UpdatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to update shares count for post with ID %s: %v", postID, err)
	}

	return nil
}
