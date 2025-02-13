package handler

import (
	"hornet/api/followers/model"
	"hornet/api/followers/service"
	"hornet/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateFollow creates a new follow relationship.
func CreateFollow(service *service.FollowersService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			logger.WithContext(c).Warn("Missing X-User-ID header")
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid userID ", userIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
			return
		}

		var req model.CreateFollow
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.WithContext(c).Error("Invalid request body ", " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		receiverID, err := uuid.Parse(req.ReceiverID)
		if err != nil {
			logger.WithContext(c).Error("Invalid ReceiverID ", req.ReceiverID, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ReceiverID"})
			return
		}

		follow, err := service.CreateFollow(c.Request.Context(), userID, receiverID)
		if err != nil {
			logger.WithContext(c).Error("Error creating follow ", "error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Follow created successfully ", follow.ID)
		c.JSON(http.StatusCreated, follow)
	}
}

// DeleteFollow deletes an existing follow relationship.
func DeleteFollow(service *service.FollowersService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			logger.WithContext(c).Warn("Missing X-User-ID header")
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid userID ", userIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
			return
		}

		followIDStr := c.Param("follow_id")
		followID, err := uuid.Parse(followIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid followID ", followIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid FollowID"})
			return
		}

		if err := service.DeleteFollow(c.Request.Context(), userID, followID); err != nil {
			logger.WithContext(c).Error("Error deleting follow ", "error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Follow deleted successfully ", followID)
		c.JSON(http.StatusOK, gin.H{"message": "Follow deleted successfully"})
	}
}

// GetUserFollowers gets a list of followers for a user.
func GetUserFollowers(service *service.FollowersService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDStr := c.Param("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid targetUserID ", userIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
			return
		}

		followers, err := service.GetFollowers(c.Request.Context(), userID)
		if err != nil {
			logger.WithContext(c).Error("Error fetching followers ", "error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Followers fetched successfully for user ", userID)
		c.JSON(http.StatusOK, followers)
	}
}

// GetUserFollowing gets a list of users the given user is following.
func GetUserFollowing(service *service.FollowersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid userID ", userIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
			return
		}

		following, err := service.GetFollowing(c.Request.Context(), userID)
		if err != nil {
			logger.WithContext(c).Error("Error fetching following ", "error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Following fetched successfully for user ", userID)
		c.JSON(http.StatusOK, following)
	}
}

// GetFollowersCount gets the count of followers for a user.
func GetFollowersCount(service *service.FollowersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid userID ", userIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
			return
		}

		count, err := service.GetFollowersCount(c.Request.Context(), userID)
		if err != nil {
			logger.WithContext(c).Error("Error fetching followers count ", "error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Followers count fetched successfully for user ", userID)
		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}

// GetUserFollowingCount gets the count of users the given user is following.
func GetFollowingCount(service *service.FollowersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.WithContext(c).Error("Invalid userID ", userIDStr, " error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
			return
		}

		count, err := service.GetFollowingCount(c.Request.Context(), userID)
		if err != nil {
			logger.WithContext(c).Error("Error fetching following count ", "error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.WithContext(c).Info("Following count fetched successfully for user ", userID)
		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}
