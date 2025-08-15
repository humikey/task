package main

import (
	"github.com/joho/godotenv"
	"log"
	"my-blog/config"
	"my-blog/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
		return
	}
	config.InitDB()           // 初始化数据库和数据迁移
	r := router.SetupRouter() // 设置路由
	err = r.Run(":8080")      // 启动服务
	if err != nil {
		log.Print("启动服务失败:", err)
		return
	}
}
