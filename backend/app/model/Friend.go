package model

import "time"

type Friend struct {
	Id            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	FromId        int       `json:"user1" `
	ToId          int       `json:"user2" `
	LastMessageId int       `json:"last_message_id"`
	Ctime         time.Time `json:"ctime" gorm:"autoCreateTime"`
	// 备注
	Content string `json:"content"`
}
