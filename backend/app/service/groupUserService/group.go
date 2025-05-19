package groupUserService

import (
	"message/app/model"
	"message/config/database"
)

func GetUsersInGroup(groupId int) ([]*model.GroupUser, error) {
	var groupUsers []*model.GroupUser
	if err := database.DB.Where("group_id = ?", groupId).Find(&groupUsers).Error; err != nil {
		return nil, err
	}
	return groupUsers, nil
}
