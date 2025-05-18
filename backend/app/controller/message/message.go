package messageController

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	redisv9 "github.com/redis/go-redis/v9"
	"message/app/comm"
	"message/app/model"
	"message/app/service/messageService"
	"message/app/utils"
	"message/config/redis"
	"sort"
	"strconv"
	"time"
)

type SendData struct {
	SendTo   int              `json:"sendTo"`
	Content  string           `json:"content"`
	MsgType  comm.MessageType `json:"msgType"`
	SendType int              `json:"sendType"`
}

func SendMessage(c *gin.Context) {
	var data SendData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	if data.SendTo == 0 || data.Content == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	userid := c.MustGet("userId").(int)
	messageService.CreateMessage(&model.Message{
		SendType: data.SendType,
		SendTo:   data.SendTo,
		From:     userid,
		Content:  data.Content,
		MsgType:  data.MsgType,
	})
	key := "msg_" + strconv.Itoa(data.SendType) + "_" + strconv.Itoa(data.SendTo)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	_, err := redis.RedisClient.Get(context.Background(), key).Result()
	if errors.Is(err, redisv9.Nil) {
		redis.RedisClient.HSet(context.Background(), key, "time", currentTime, "count", 1)
	} else if err == nil {
		count, _ := redis.RedisClient.HGet(context.Background(), key, "count").Int()
		redis.RedisClient.HSet(context.Background(), key, "time", currentTime, "count", count+1)
	} else {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Redis错误",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "发送成功",
	})
}

type ListMessageResp struct {
	Total    int              `json:"total"`
	Messages []*model.Message `json:"messages"`
}

func ListMessage(c *gin.Context) {
	id, flag := c.Params.Get("friend_id")
	sendType, flag := c.Params.Get("send_type")
	if !flag {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	friendId, _ := strconv.Atoi(id)
	userid := c.MustGet("userId").(int)
	// Redis键
	key := "msg_" + sendType + "_" + strconv.Itoa(friendId)
	t, _ := strconv.Atoi(sendType)

	// 获取Redis中的消息数量
	count, err := redis.RedisClient.HGet(context.Background(), key, "count").Int()
	if err != nil && !errors.Is(err, redisv9.Nil) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "ok",
			"data":    nil,
		})
		return
	} else if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Redis错误",
		})
		return
	}

	messages, err := messageService.ListMessage(userid, friendId, count, t)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "服务器错误",
		})
		return
	}

	redis.RedisClient.Del(context.Background(), key)
	utils.JsonSuccessResponse(c, ListMessageResp{
		Total:    len(messages),
		Messages: messages,
	})
}

type GetAllNewMessageResp struct {
	Total int       `json:"total"`
	Ctime time.Time `json:"ctime"`
	Key   string    `json:"key"`
}

func GetAllNewMessage(c *gin.Context) {
	keys, err := redis.RedisClient.Keys(context.Background(), "msg_*").Result()
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Redis错误",
		})
		return
	}
	var responses []GetAllNewMessageResp

	for _, key := range keys {
		count, err := redis.RedisClient.HGet(context.Background(), key, "count").Int()
		if err != nil {
			continue
		}
		ctimeStr, err := redis.RedisClient.HGet(context.Background(), key, "time").Result()
		if err != nil {
			continue
		}

		ctime, err := time.Parse("2006-01-02 15:04:05", ctimeStr)
		if err != nil {
			continue
		}
		responses = append(responses, GetAllNewMessageResp{
			Total: count,
			Ctime: ctime,
			Key:   key,
		})
	}
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Ctime.After(responses[j].Ctime)
	})
	utils.JsonSuccessResponse(c, responses)
}

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
