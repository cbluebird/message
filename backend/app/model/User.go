package model

import "time"

type User struct {
	Id          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	WxID        string    `json:"wx_id"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	SignOutTime time.Time `json:"sign_out_time"`
	Ctime       time.Time `json:"ctime" gorm:"autoCreateTime"`
}
