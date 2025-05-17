package redis

import (
	"github.com/redis/go-redis/v9"
	"message/config/config"
)

var RedisClient *redis.Client

type RunOptions struct {
	Host     string
	Port     string
	DB       int
	Password string
}

func NewRunOptions() *RunOptions {
	Info := &RunOptions{
		Host:     "localhost",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
	if config.Config.IsSet("redis.host") {
		Info.Host = config.Config.GetString("redis.host")
	}
	if config.Config.IsSet("redis.port") {
		Info.Port = config.Config.GetString("redis.port")
	}
	if config.Config.IsSet("redis.db") {
		Info.DB = config.Config.GetInt("redis.db")
	}
	if config.Config.IsSet("redis.password") {
		Info.Password = config.Config.GetString("redis.password")
	}
	return Info
}

func (options *RunOptions) Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     options.Host + ":" + options.Port,
		Password: options.Password,
		DB:       options.DB,
	})
}
