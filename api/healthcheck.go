package api

import (
	"github.com/gin-gonic/gin"
)

func (h reqHandler) CheckHealth(c *gin.Context) {
	c.JSON(200, gin.H{"status": "serving"})
}
