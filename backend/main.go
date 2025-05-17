package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	midwares "message/app/mid"
	"message/config/config"
	"message/config/database"
	"message/config/redis"
	"message/config/router"
)

func main() {
	database.NewRunOptions().Init()
	redis.NewRunOptions().Init()
	r := gin.Default()
	r.Use(cors.Default())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	router.Init(r)
	err := r.Run(":" + config.Config.GetString("server.port"))
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}
}
