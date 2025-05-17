package messageService

import (
	"message/app/model"
	"message/config/database"
)

func CreateMessage(message *model.Message) error {
	return database.DB.Create(message).Error
}

func ListMessage(userId, fromId, num, MessageType int) ([]*model.Message, error) {
	ansList := make([]*model.Message, 0)
	database.DB.Model(&model.Message{}).Where(&model.Message{
		SendType: MessageType,
		SendTo:   userId,
		From:     fromId,
	}).Order("ctime DESC").Limit(num).Find(&ansList)
	return ansList, nil
}
