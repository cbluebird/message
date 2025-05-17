package router

import (
	"github.com/gin-gonic/gin"
	messageController "message/app/controller/message"
	"message/app/controller/user"
	middleware "message/app/mid"
)

func Init(r *gin.Engine) {
	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.Register)
		api.POST("/mail")
	}
	message := r.Group("/message", middleware.TokenAuth)
	{
		message.POST("/send", messageController.SendMessage)
		message.GET("/new", messageController.GetAllNewMessage)
		message.GET("/get", messageController.ListMessage)
	}
	group := api.Group("/group", middleware.TokenAuth)
	{

	}
}
