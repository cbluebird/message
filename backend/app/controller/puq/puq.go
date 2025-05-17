package puqController

import (
	"github.com/gin-gonic/gin"
	"message/app/model"
	"message/config/database"
	"strconv"
	"time"
)

func CreatePuq(c *gin.Context) {
	var request struct {
		Title   string             `json:"title"`
		Content []model.PuqContent `json:"content"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	userId := c.MustGet("userId").(int)
	puq := model.Puq{
		Title: request.Title,
		From:  userId,
		Ctime: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := database.DB.Create(&puq).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "发布失败"})
		return
	}

	for _, content := range request.Content {
		content.NoticeId = puq.ID
		if err := database.DB.Create(&content).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "message": "发布内容失败"})
			return
		}
	}
	c.JSON(200, gin.H{"code": 200, "message": "发布成功"})
}

type PuqResp struct {
	Puq      model.Puq          `json:"puq"`
	Contents []model.PuqContent `json:"contents"`
}

func ListPuq(c *gin.Context) {
	userId := c.MustGet("userId").(int)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 10

	var friendIds []int
	if err := database.DB.Model(&model.Friend{}).Where("from_id = ?", userId).Pluck("to_id", &friendIds).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取好友列表失败"})
		return
	}

	// 获取好友的 Puq
	var puqs []model.Puq
	if err := database.DB.Where("from IN (?)", friendIds).Order("ctime DESC").Limit(10).Offset(offset).Find(&puqs).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取失败"})
		return
	}
	ans := make([]PuqResp, 0)
	for i := range puqs {
		ans = append(ans, PuqResp{
			Puq:      puqs[i],
			Contents: make([]model.PuqContent, 0),
		})
		var contents []model.PuqContent
		if err := database.DB.Where("notice_id = ?", puqs[i].ID).Find(&contents).Error; err != nil {
			c.JSON(500, gin.H{"code": 500, "message": "获取内容失败"})
			return
		}
		ans[i].Contents = contents
	}

	c.JSON(200, gin.H{"code": 200, "message": "获取成功", "data": ans})
}

func DeletePuq(c *gin.Context) {
	puqId, err := strconv.Atoi(c.Param("puq_id"))
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	userId := c.MustGet("userId").(int)
	if err := database.DB.Where("id = ? AND from = ?", puqId, userId).Delete(&model.Puq{}).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	if err := database.DB.Where("notice_id = ?", puqId).Delete(&model.PuqContent{}).Error; err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "删除内容失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "删除成功"})
}
