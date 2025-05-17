package groupService

import (
	"errors"
	"message/app/model"
	"message/config/database"
	"time"
)

func CreateGroup(group *model.Group) error {
	return database.DB.Create(group).Error
}

func GetGroupById(groupId int) (*model.Group, error) {
	var group model.Group
	if err := database.DB.First(&group, groupId).Error; err != nil {
		return nil, errors.New("群组不存在")
	}
	return &group, nil
}

func AddUserToGroup(userId, groupId, userType int) error {
	_, err := GetGroupById(groupId)
	if err != nil {
		return errors.New("群组不存在")
	}

	var groupUser model.GroupUser
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupId, userId).First(&groupUser).Error; err == nil {
		return errors.New("用户已在群组中")
	}

	groupUser = model.GroupUser{
		GroupId:  groupId,
		UserId:   userId,
		Ctime:    time.Now(),
		UserType: userType,
	}
	return database.DB.Create(&groupUser).Error
}

func UpdateGroup(group *model.Group) error {
	return database.DB.Model(&model.Group{}).Where("id = ?", group.ID).Updates(group).Error
}

func IsGroupOwner(userId, groupId int) bool {
	var groupUser model.GroupUser
	if err := database.DB.Where("group_id = ? AND user_id = ? AND user_type = 1", groupId, userId).First(&groupUser).Error; err != nil {
		return false
	}
	return true
}

func RemoveUserFromGroup(userId, groupId int) error {
	var groupUser model.GroupUser
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupId, userId).First(&groupUser).Error; err != nil {
		return errors.New("用户不在群组中")
	}
	return database.DB.Delete(&groupUser).Error
}
