package groupController

import (
	"github.com/gin-gonic/gin"
	"message/app/model"
	"message/app/service/groupService"
	"strconv"
)

func CreateGroup(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userId := c.MustGet("userId").(int)

	if err := groupService.CreateGroup(&group); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "创建群组失败",
		})
		return
	}

	if err := groupService.AddUserToGroup(userId, group.ID, 1); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "加入群组失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "群组创建成功",
		"data":    group,
	})
}

func JoinGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userId := c.MustGet("userId").(int)
	if err := groupService.AddUserToGroup(userId, groupId, 3); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "加入群组失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "成功加入群组",
	})
}

func UpdateGroup(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	userId := c.MustGet("userId").(int)
	if !groupService.IsGroupOwner(userId, group.ID) {
		c.JSON(403, gin.H{
			"code":    403,
			"message": "无权限更新群组",
		})
		return
	}
	if err := groupService.UpdateGroup(&group); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "更新群组失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "群组更新成功",
	})
}

func LeaveGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userId := c.MustGet("userId").(int)
	if err := groupService.RemoveUserFromGroup(userId, groupId); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "退出群组失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "成功退出群组",
	})
}
