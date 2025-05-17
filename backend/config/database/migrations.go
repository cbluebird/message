package database

import (
	"gorm.io/gorm"
	"message/app/model"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Message{},
		&model.Group{},
		&model.GroupUser{},
		&model.InitReq{},
		&model.Friend{},
	)
}
