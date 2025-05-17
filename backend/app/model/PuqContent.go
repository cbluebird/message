package model

type PuqContent struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	NoticeId int    `json:"notice_id"` // 通知ID
	Content  string `json:"content"`
	Type     int    `json:"type"` // 1:文本 2:图片 3:视频
}
