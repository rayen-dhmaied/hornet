package posts

import (
	"hornet/api/posts/handler"
	"hornet/api/posts/service"

	"github.com/gin-gonic/gin"
)

// Router sets up the Gin router with all the routes
func Router(postService *service.PostService) *gin.Engine {
	r := gin.Default()

	// Set the handler with the service layer
	r.POST("/posts", handler.CreatePost(postService))

	// Get a post by ID
	r.GET("/posts/:id", handler.GetPost(postService))

	// Get replies for a parent post
	r.GET("/posts/:id/replies", handler.GetReplies(postService))

	// Delete a post by ID
	r.DELETE("/posts/:id", handler.DeletePost(postService))

	return r
}
