package followers

import (
	"hornet/api/followers/handler"
	"hornet/api/followers/service"

	"github.com/gin-gonic/gin"
)

// Router sets up the Gin router with all the routes
func Router(followersService *service.FollowersService) *gin.Engine {
	r := gin.Default()

	// Create a new follow
	r.POST("/followers", handler.CreateFollow(followersService))

	// Delete a follow by ID
	r.DELETE("/followers/:follow_id", handler.DeleteFollow(followersService))

	// Get user followers
	r.GET("/followers/user/:user_id/followers", handler.GetUserFollowers(followersService))

	// Get user following
	r.GET("/followers/user/:user_id/following", handler.GetUserFollowing(followersService))

	// Get user followers count
	r.GET("/followers/user/:user_id/followers/count", handler.GetFollowersCount(followersService))

	// Get user following count
	r.GET("/followers/user/:user_id/following/count", handler.GetFollowingCount(followersService))

	return r
}
