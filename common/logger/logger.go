package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	// Initialize the logger with development settings (you can change this to production or customize as needed)
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

// WithContext adds the Gin context to the logger to enrich log entries with request info
// This is a common pattern to use with Gin middleware
func WithContext(c *gin.Context) *zap.SugaredLogger {
	return logger.Sugar().With(
		zap.String("method", c.Request.Method),
		zap.String("uri", c.Request.RequestURI),
		zap.String("ip", c.ClientIP()),
	)
}
