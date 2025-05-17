package model

import "time"

type InitReq struct {
	Id            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId        int       `json:"user_id"`
	LastMessageId int       `json:"last_message_id"`
	Ctime         time.Time `json:"ctime" gorm:"autoCreateTime"`
}
