package database

import (
	"fmt"
	"log"
	"message/config/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type RunOptions struct {
	User string
	Pass string
	Port string
	Name string
}

func NewRunOptions() *RunOptions {
	return &RunOptions{
		User: config.Config.GetString("db.user"),
		Pass: config.Config.GetString("db.password"),
		Port: config.Config.GetString("db.address"),
		Name: config.Config.GetString("db.name"),
	}
}

func (options *RunOptions) Init() { // 初始化数据库
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		options.User,
		options.Pass,
		options.Port,
		options.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln("Database Error: ", err)
	} else {
		fmt.Printf("database start")
	}
	err = autoMigrate(db)
	if err != nil {
		log.Fatal("DatabaseMigrateFailed", err)
	}
	DB = db
}
