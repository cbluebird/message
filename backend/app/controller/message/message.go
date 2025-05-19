package messageController

import (
	"github.com/gin-gonic/gin"
	"log"
	"message/app/service/messageService"
)

func Translate(c *gin.Context) {
	// 获取请求参数
	var request struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 调用翻译服务
	result, err := messageService.Translate(request.Content)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "翻译失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "翻译成功",
		"data":    result,
	})

}

func WsServer(c *gin.Context) {
	userid := c.MustGet("userId").(int)
	if _, err := messageService.WsServer(c, messageService.Json, userid); err != nil {
		log.Println("ws conn error", err)
	}
}
