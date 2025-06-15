package main

import (
	"mhxyHelper/pkg/database"
	"mhxyHelper/pkg/logger"
	"mhxyHelper/sever/cobra/cmd"
)

// TODO: 简单尝试暂时未用
func main() {

	logger.NewLogger()

	// 初始化数据库连接
	_, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	//ctx := context.Background()

	err = cmd.Execute()
	if err != nil {
		panic(err)
	}
}
