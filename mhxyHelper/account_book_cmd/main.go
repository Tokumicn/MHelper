package main

import (
	"mhxyHelper/internal/database"
	"mhxyHelper/pkg/logger"
)

func main() {
	logger.NewLogger()

	// 初始化数据库连接
	_, err := database.InitDB()
	if err != nil {
		panic(err)
	}

}
