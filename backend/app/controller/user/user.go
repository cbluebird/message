package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"message/app/model"
	"message/app/service/userService"
	"message/app/utils"
	"time"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utils.JsonResponseInternalServerError(c)
		return
	}
	user, err := userService.GetUserByEmail(data.Email)
	if err != nil {
		utils.JsonResponse(404, "该用户不存在", nil, c)
		return
	}
	flag := userService.Login(data.Email, data.Password)
	if !flag {
		utils.JsonResponse(409, "密码错误", nil, c)
		return
	}
	jwt, _ := utils.GenerateToken(user.Id)
	utils.JsonResponse(200, "OK", gin.H{
		"token": jwt,
	}, c)
}

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"name"`
	Avatar   string `json:"avatar"`
	Code     int    `json:"code"`
}

func Register(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utils.JsonResponseInternalServerError(c)
		return
	}
	if data.Code != 200 {
		utils.JsonResponse(409, "验证码错误", nil, c)
		return
	}

	if _, err = userService.GetUserByEmail(data.Email); err == nil {
		utils.JsonResponse(409, "该用户已存在", nil, c)
		return
	}

	if err = userService.CreateUser(&model.User{
		Email:       data.Email,
		Password:    data.Password,
		Username:    data.Username,
		Avatar:      data.Avatar,
		WxID:        utils.RandomString(12),
		SignOutTime: time.Now(),
	}); err != nil {
		log.Println(err)
		utils.JsonResponseInternalServerError(c)
		return
	}

	utils.JsonResponse(200, "注册成功", nil, c)
}
