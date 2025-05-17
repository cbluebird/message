package model

type Contact struct {
	Id     int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	FromId string `json:"user1" `
	ToId   string `json:"user2" `
	Ctime  int64  `json:"ctime" gorm:"autoCreateTime"`
	// 备注
	Content string `json:"content"`
}
