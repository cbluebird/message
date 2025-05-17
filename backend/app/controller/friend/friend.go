package friendController

import (
	"github.com/gin-gonic/gin"
	"message/app/model"
	"message/config/database"
	"strconv"
	"time"
)

func SendFriendApply(c *gin.Context) {
	var apply model.FriendApply
	if err := c.ShouldBindJSON(&apply); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	apply.ApplyTime = time.Now().Format("2006-01-02 15:04:05")
	apply.Status = 0 // 默认状态为待处理

	if err := database.DB.Create(&apply).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "发送好友申请失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "好友申请已发送",
	})
}

func HandleFriendApply(c *gin.Context) {
	var request struct {
		ApplyId int `json:"applyId"`
		Status  int `json:"status"` // 1: 同意, 2: 拒绝
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var apply model.FriendApply
	if err := database.DB.First(&apply, request.ApplyId).Error; err != nil {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "好友申请不存在",
		})
		return
	}

	if apply.Status != 0 {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "好友申请已处理",
		})
		return
	}

	apply.Status = request.Status
	if err := database.DB.Save(&apply).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "处理好友申请失败",
		})
		return
	}

	if request.Status == 1 {
		friend := model.Friend{
			FromId: apply.UserId,
			ToId:   apply.ApplyUserId,
			Ctime:  time.Now(),
		}
		if err := database.DB.Create(&friend).Error; err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "添加好友失败",
			})
			return
		}
		friend.ToId, friend.FromId = friend.FromId, friend.ToId
		if err := database.DB.Create(&friend).Error; err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "添加好友失败",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "好友申请处理成功",
	})
}

func DeleteFriend(c *gin.Context) {
	var request struct {
		FromId string `json:"fromId"`
		ToId   string `json:"toId"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	if err := database.DB.Where("from_id = ? AND to_id = ?", request.FromId, request.ToId).Or("from_id = ? AND to_id = ?", request.ToId, request.FromId).Delete(&model.Friend{}).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "删除好友失败",
		})
		return
	}
	request.ToId, request.FromId = request.FromId, request.ToId
	if err := database.DB.Where("from_id = ? AND to_id = ?", request.FromId, request.ToId).Or("from_id = ? AND to_id = ?", request.ToId, request.FromId).Delete(&model.Friend{}).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "删除好友失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "好友已删除",
	})
}

func ListFriends(c *gin.Context) {
	userId := c.MustGet("userId").(int)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	// 查询好友关系
	var friendIds []int
	if err := database.DB.Model(&model.Friend{}).
		Where("from_id = ?", userId).
		Or("to_id = ?", userId).
		Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取好友列表失败",
		})
		return
	}

	// 查询好友详细信息
	var friends []model.User
	if err := database.DB.Where("id IN (?)", friendIds).
		Limit(limit).Offset(offset).
		Find(&friends).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取好友信息失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    friends,
	})
}
