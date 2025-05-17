package model

type Group struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupName   string `json:"group_name" gorm:"unique;not null"` // 群名称
	GroupAvatar string `json:"group_avatar"`                      // 群头像
	GroupDesc   string `json:"group_desc"`                        // 群描述
}
