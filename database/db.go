package database

import (
	"fmt"
	"log"
	"time"

	"gin-demo/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	// 构建 DSN（数据源名称）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "123456", "localhost", 3306, "region")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用默认事务（提高性能）
		SkipDefaultTransaction: false,
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		// 日志配置
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	// 获取底层 SQL DB 对象以配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("获取 SQL DB 失败:", err)
	}
	// 配置连接池
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大生命周期

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("database ping failed:", err)
	}
	log.Println("数据库连接成功")
}

// 2.2 自动迁移最佳实践
// 自动迁移是 GORM 的强大功能，但需要谨慎使用：
func AutoMigrateModels() {
	err := DB.AutoMigrate(
		&models.User{},
		// &Product{}, // 其他模型
		// &Order{},   // 其他模型
	)
	if err != nil {
		log.Fatal("自动迁移失败:", err)
	}
	log.Println("自动迁移完成")
}
