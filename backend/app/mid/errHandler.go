package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandleNotFound
//
//	404处理
func HandleNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, nil)
	return
}
