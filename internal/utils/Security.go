package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func Filter(c *gin.Context, callback func(c *gin.Context)) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	bearer := fmt.Sprint("Bearer ", os.Getenv("API_SECRET_KEY"))
	if token != bearer {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}
	callback(c)
}
