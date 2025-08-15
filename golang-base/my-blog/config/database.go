package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"my-blog/model"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("MYSQL_DSN")
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // 日志输出
	//	logger.Config{
	//		SlowThreshold:             time.Second, // 慢 SQL 阈值
	//		LogLevel:                  logger.Info, // 日志级别
	//		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound错误
	//		Colorful:                  true,        // 彩色打印
	//	},
	//)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: newLogger, // 👈 这里设置
	})
	// 设置全局日志级别
	enableDebug := os.Getenv("ENABLE_PRINT_SQL")
	if enableDebug == "true" {
		DB = DB.Debug()
	}
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	log.Println("数据库连接成功")

	err = DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		log.Print("数据库迁移失败:", err)
		return
	}
}
