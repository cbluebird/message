package model

import "time"

type GroupUser struct {
	ID       int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupId  int64     `json:"group_id"` // 群组id
	UserId   int64     `json:"user_id"`  // 用户id
	Nickname string    `json:"nickname"` // 昵称
	UserType int       `json:"user_type"`
	Ctime    time.Time `json:"ctime" gorm:"autoCreateTime"`
}
