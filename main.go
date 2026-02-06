package main

import (
	"gin-demo/database"
	"gin-demo/handlers"
	"gin-demo/repositories"
	"gin-demo/routers"
	"gin-demo/services"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()
	// 自动迁移（仅在开发环境使用）
	if gin.Mode() == gin.DebugMode {
		database.AutoMigrateModels()
	}

	// 初始化 Repository
	userRepo := repositories.NewUserRepository(database.DB)

	// 初始化 Service
	userService := services.NewUserService(userRepo, database.DB)

	// 初始化 Handler
	userHandler := handlers.NewUserHandler(userService)

	// 初始化 Gin
	r := gin.Default()

	// 设置路由
	routers.SetupUserRoutes(r, userHandler)

	// 启动服务器
	log.Println("服务器启动在 :8080 端口")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
