package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 定义全局变量和互斥锁
var (
	db   *sql.DB
	once sync.Once
)

// InitDB 初始化数据库连接
func InitDB() {
	dataSourceName := "root:123456@tcp(192.168.3.74:3306)/test_db"
	// 使用 sync.Once 确保数据库连接池只初始化一次
	once.Do(func() {
		var err error
		// 打开数据库连接
		db, err = sql.Open("mysql", dataSourceName)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		// 设置数据库连接池参数
		db.SetMaxOpenConns(50)                 // 最大打开连接数
		db.SetMaxIdleConns(15)                 // 最大空闲连接数
		db.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期

		// 检查连接是否可用
		if err = db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		fmt.Println("Database connection pool initialized")
	})
}

// GetDB 返回数据库连接实例
func GetDB() *sql.DB {
	return db
}
