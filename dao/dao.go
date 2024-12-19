package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// 最大重试次数和重试间隔
const maxRetries = 5
const retryInterval = 2 * time.Second

func InitMysql() *gorm.DB {
	dsn := "root:kangqiao2006714@tcp(mysql-test:3306)/messagesboard?charset=utf8&parseTime=True&loc=Local"

	var mysqlLogger logger.Interface
	mysqlLogger = logger.Default.LogMode(logger.Error)

	var db *gorm.DB
	var err error

	// 尝试连接数据库，最多重试 maxRetries 次
	for attempts := 0; attempts < maxRetries; attempts++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: mysqlLogger,
		})
		if err == nil {
			// 如果连接成功，则设置连接池
			sqlDB, err := db.DB()
			if err != nil {
				fmt.Printf("获取数据库连接池失败: %v\n", err)
				return nil
			}
			// 配置连接池
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour * 4)

			// 连接成功，返回 db 实例
			fmt.Println("数据库连接成功")
			return db
		}

		// 如果连接失败，打印错误并等待重试
		fmt.Printf("数据库连接失败: %v，正在重试（尝试 %d/%d）...\n", err, attempts+1, maxRetries)
		time.Sleep(retryInterval)
	}

	// 如果重试了 maxRetries 次仍然失败，返回 nil
	fmt.Println("数据库连接失败，已达最大重试次数")
	return nil
}
