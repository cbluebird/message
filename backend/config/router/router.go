package router

import (
	"github.com/gin-gonic/gin"
	"message/app/controller/user"
)

func Init(r *gin.Engine) {
	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.Register)
		api.POST("/mail")
	}
}
