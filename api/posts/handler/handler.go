package handler

import (
	"hornet/api/posts/model"
	"hornet/api/posts/service"
	"hornet/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetPost handles the retrieval of a post by its ID
func GetPost(postService *service.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")

		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid post ID ", postIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		post, err := postService.GetPost(c.Request.Context(), postID)
		if err != nil {
			logger.WithContext(c).Error("Error retrieving post ", postID, " error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if post.ID == uuid.Nil {
			logger.WithContext(c).Info("Post not found ", postID)
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}

		logger.WithContext(c).Info("Post retrieved successfully ", post.ID)
		c.JSON(http.StatusOK, post)
	}
}

// CreatePost handles the creation of a new post
func CreatePost(postService *service.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreatePost

		authorID := c.GetHeader("X-User-ID")
		if authorID == "" {
			logger.WithContext(c).Warn("Missing X-User-ID header")
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
			return
		}

		_, err := uuid.Parse(authorID)
		if err != nil {
			logger.WithContext(c).Error("Invalid AuthorID ", authorID, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid AuthorID"})
			return
		}

		req.AuthorID = authorID

		if err := c.ShouldBindJSON(&req); err != nil {
			logger.WithContext(c).Error("Invalid request body ", " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if req.OriginalPostID == nil && req.Content == nil {
			logger.WithContext(c).Warn("Missing content for new post ", authorID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required when creating a new post"})
			return
		}

		if req.Content != nil && (len(*req.Content) > 5000 || len(*req.Content) == 0) {
			logger.WithContext(c).Warn("Invalid content length ", authorID, " contentLength: ", len(*req.Content))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content length should be between 1 and 5000 characters"})
			return
		}

		post, err := postService.CreatePost(c.Request.Context(), req)
		if err != nil {
			logger.WithContext(c).Error("Error creating post ", "error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Post created successfully ", post.ID)
		c.JSON(http.StatusCreated, post)
	}
}

// DeletePost handles the deletion of a post by its ID
func DeletePost(postService *service.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")

		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid post ID ", postIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		err = postService.DeletePost(c.Request.Context(), postID)
		if err != nil {
			logger.WithContext(c).Error("Error deleting post ", postID, " error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Post deleted successfully ", postID)
		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
}

// GetReplies handles retrieval of replies for a post
func GetReplies(postService *service.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentPostIDStr := c.Param("id")
		parentPostID, err := uuid.Parse(parentPostIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid parent post ID ", parentPostIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent post ID"})
			return
		}

		replies, err := postService.GetReplies(c.Request.Context(), parentPostID)
		if err != nil {
			logger.WithContext(c).Error("Error fetching replies ", parentPostID, " error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(replies) == 0 {
			logger.WithContext(c).Info("No replies found ", parentPostID)
			c.JSON(http.StatusOK, gin.H{"message": "No replies found"})
			return
		}

		logger.WithContext(c).Info("Replies retrieved successfully ", parentPostID, " repliesCount: ", len(replies))
		c.JSON(http.StatusOK, replies)
	}
}
