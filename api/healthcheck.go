package api

import (
	"github.com/gin-gonic/gin"
)

// CheckHealth provides the implementation for a basic
// health endpoint.
func (h reqHandler) CheckHealth(c *gin.Context) {
	c.JSON(200, gin.H{"status": "serving"})
}
