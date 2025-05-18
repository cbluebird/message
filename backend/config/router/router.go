package router

import (
	"github.com/gin-gonic/gin"
	"message/app/controller/friend"
	groupController "message/app/controller/group"
	messageController "message/app/controller/message"
	puqController "message/app/controller/puq"
	"message/app/controller/user"
	middleware "message/app/mid"
)

func Init(r *gin.Engine) {
	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.Register)
		api.POST("/mail", userController.SendRegisterEmail)
	}
	user := api.Group("/user", middleware.TokenAuth)
	{
		user.GET("/find", userController.FindUser)
		user.POST("/update", userController.UpdateUserInfo)
	}
	message := api.Group("/message", middleware.TokenAuth)
	{
		message.POST("/send", messageController.SendMessage)
		message.GET("/new", messageController.GetAllNewMessage)
		message.GET("/get", messageController.ListMessage)
		message.POST("/translate", messageController.Translate) // 查看某个用户的消息
	}
	group := api.Group("/group", middleware.TokenAuth)
	{
		group.POST("/create", groupController.CreateGroup)           // 创建群组
		group.POST("/join/:group_id", groupController.JoinGroup)     // 加入群组
		group.POST("/update", groupController.UpdateGroup)           // 更新群组
		group.DELETE("/leave/:group_id", groupController.LeaveGroup) // 退出群组
	}
	friend := api.Group("/friend", middleware.TokenAuth)
	{
		friend.GET("/list", friendController.ListFriends)          // 查看好友列表
		friend.POST("/apply", friendController.SendFriendApply)    // 发送好友申请
		friend.POST("/handle", friendController.HandleFriendApply) // 处理好友申请（同意或拒绝）
		friend.DELETE("/delete", friendController.DeleteFriend)    // 删除好友
	}
	puq := api.Group("/puq", middleware.TokenAuth)
	{
		puq.POST("/create", puqController.CreatePuq)         // 发布 puq
		puq.GET("/list", puqController.ListPuq)              // 查看好友的 puq
		puq.POST("/delete/:puq_id", puqController.DeletePuq) // 删除 puq
	}
}
