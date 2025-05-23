package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

var (
	_db *gorm.DB
)

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	dir, _ := os.Getwd()
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binDir := filepath.Dir(exePath)
	fmt.Println("[MHXYDB] shell work dir: ", dir)
	fmt.Println("[MHXYDB] bin work dir: ", binDir)

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/mhxyhelper.db", binDir)), &gorm.Config{
		//// 开启 WAL 模式
		//DSN: "mode=wal",
		//// 增加最大连接数为 100
		//MaxOpenConns: 100,
	})
	if err != nil {
		return nil, err
	}
	// 设置数据库连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	_, err = sqlDB.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	// TODO 开启SQL打印测试用
	db = db.Debug()

	_db = db

	return _db, nil
}

// LocalDB 获取数据库连接
func LocalDB() *gorm.DB {
	if _db != nil {
		return _db
	}
	panic("database connection is nil")
}
