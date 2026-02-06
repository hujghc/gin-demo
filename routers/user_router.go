package routers

import (
	"gin-demo/handlers"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户相关的路由
func SetupUserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
	// 用户路由组
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)           // 创建用户
		userRoutes.GET("", userHandler.GetUsers)              // 获取用户列表
		userRoutes.GET("/:id", userHandler.GetUser)           // 获取单个用户（需要实现）
		userRoutes.PUT("/:id", userHandler.UpdateUser)       // 更新用户
		userRoutes.DELETE("/:id", userHandler.DeleteUser)    // 删除用户
		userRoutes.POST("/transfer", userHandler.TransferBalance) // 转账
	}
}
