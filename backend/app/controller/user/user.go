package userController

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"message/app/model"
	"message/app/service/userService"
	"message/app/utils"
	"message/config/database"
	"message/config/redis"
	"strconv"
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
	redisKey := "register_code:" + data.Email
	code, err := redis.RedisClient.Get(context.Background(), redisKey).Result()
	if err != nil {
		utils.JsonResponse(409, "验证码已过期", nil, c)
		return
	}
	if strconv.Itoa(data.Code) != code {
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

func FindUser(c *gin.Context) {
	var query struct {
		UserId int    `form:"userId"`
		Email  string `form:"email"`
		WxID   string `form:"wxid"`
	}

	// 绑定查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	if query.UserId == 0 && query.Email == "" && query.WxID == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "请提供至少一个查询条件",
		})
		return
	}
	var user model.User
	db := database.DB
	if query.UserId != 0 {
		db = db.Where("id = ?", query.UserId)
	}
	if query.Email != "" {
		db = db.Where("email = ?", query.Email)
	}
	if query.WxID != "" {
		db = db.Where("wx_id = ?", query.WxID)
	}
	if err := db.First(&user).Error; err != nil {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    user,
	})
}

func UpdateUserInfo(c *gin.Context) {
	userId := c.MustGet("userId").(int)

	var request struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 更新用户信息
	updates := map[string]interface{}{}
	if request.Username != "" {
		updates["username"] = request.Username
	}
	if request.Avatar != "" {
		updates["avatar"] = request.Avatar
	}
	if request.Email != "" {
		updates["email"] = request.Email
	}

	if len(updates) == 0 {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "没有提供更新内容",
		})
		return
	}

	if err := database.DB.Model(&model.User{}).Where("id = ?", userId).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "更新失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

func SendRegisterEmail(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	if err := utils.SendEmail(request.Email, code); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "发送邮件失败",
		})
		return
	}

	// 将验证码存储到 Redis，设置 5 分钟过期时间
	redisKey := "register_code:" + request.Email
	if err := redis.RedisClient.Set(context.Background(), redisKey, code, 5*time.Minute); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "存储验证码失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "验证码已发送",
	})
}
