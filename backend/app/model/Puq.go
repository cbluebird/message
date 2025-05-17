package model

type Puq struct {
	ID    int    `json:"ID" gorm:"primary_key;AUTO_INCREMENT"`
	Title string `json:"title"` // 标题
	From  int    `json:"from"`  // 发送人
	Ctime string `json:"ctime"` // 发送时间
}
