package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"message/app/utils"
	"net/http"
)

func TokenAuth(c *gin.Context) {
	key := c.Request.Header.Get("Authorization")
	if id, err := utils.ParseToken(key); err != nil {
		slog.Error("Failed to parse token", "Error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid token"})
		c.Abort()
		return
	} else {
		c.Set("userId", id)
	}
	c.Next()
}
