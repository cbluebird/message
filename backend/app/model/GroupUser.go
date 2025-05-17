package model

import "time"

type GroupUser struct {
	ID       int       `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupId  int       `json:"group_id"`  // 群组id
	UserId   int       `json:"user_id"`   // 用户id
	Nickname string    `json:"nickname"`  // 昵称
	UserType int       `json:"user_type"` //1:群主 2:管理员 3:普通成员
	Ctime    time.Time `json:"ctime" gorm:"autoCreateTime"`
}
