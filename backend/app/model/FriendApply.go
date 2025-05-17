package model

type FriendApply struct {
	ID          int    `json:"ID" gorm:"primary_key;AUTO_INCREMENT"`
	UserId      int    `json:"userId"`      // 申请人
	ApplyUserId int    `json:"applyUserId"` // 被申请人
	ApplyTime   string `json:"applyTime"`   // 申请时间
	Status      int    `json:"status"`      // 0:待处理 1:已同意 2:已拒绝
	Msg         string `json:"msg"`         // 申请信息
}
