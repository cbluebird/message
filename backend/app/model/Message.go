package model

import (
	"message/app/comm"
	"time"
)

type Message struct {
	ID       int              `json:"ID" gorm:"primary_key;AUTO_INCREMENT"`
	SendType int              `json:"sendType"` // 1:单聊 2:群聊
	SendTo   int              `json:"sendTo"`
	From     int              `json:"sendFrom"`
	Ctime    time.Time        `json:"ctime" gorm:"autoCreateTime"`
	Content  string           `json:"content"`
	MsgType  comm.MessageType `json:"msgType"`
}
