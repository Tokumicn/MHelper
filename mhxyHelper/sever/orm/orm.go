package main

import (
	"context"
	"log"
	"mhxyHelper/pkg/database"
	"mhxyHelper/pkg/logger"
)

func main() {
	logger.NewLogger()
	ctx := context.Background()

	err := initDBWithAutoMigrate(ctx, true)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// 初始化数据库并创建表
func initDBWithAutoMigrate(ctx context.Context, needAutoMigrate bool) error {
	// 初始化数据库连接
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	if needAutoMigrate {
		// 用户表
		err = db.AutoMigrate(
		//data.User{},      // 测试用户表
		//data.StuffLog{},  // 物品信息表
		//data.Stuff{},     // 物品更新日志
		//data.Account{},   // 用户账单信息表
		//data.Attribute{}, //属性信息表
		//data.MHJLResponseLog{}, // 梦幻精灵回复日志表
		)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}
