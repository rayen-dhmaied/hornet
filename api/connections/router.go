package connections

import (
	"hornet/api/connections/service"

	"github.com/gin-gonic/gin"
)

// Router sets up the Gin router with all the routes
func Router(postService *service.PostService) *gin.Engine {
	r := gin.Default()

	return r
}
