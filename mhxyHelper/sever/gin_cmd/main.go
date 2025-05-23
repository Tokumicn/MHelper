package main

import (
	"github.com/gin-gonic/gin"
	"mhxyHelper/internal/handler"
	"mhxyHelper/pkg/database"
	"mhxyHelper/pkg/logger"
)

func main() {
	logger.NewLogger()
	// 初始化数据库连接
	_, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/ping", handler.Ping)

	// 字典构建
	dictRouter := r.Group("/dict")
	dictRouter.POST("/build", handler.BuildDict)

	// 物品相关接口
	stuffRouter := r.Group("/stuff")
	// 物品构建
	stuffRouter.POST("/build", handler.BuildStuff)
	// 物品查询
	stuffRouter.POST("/query", handler.QueryStuff)

	// 属性相关接口
	attrRouter := r.Group("/attribute")
	// 属性构建
	attrRouter.POST("/build", handler.BuildAttribute)
	// 属性查询
	attrRouter.POST("/query", handler.QueryAttribute)

	// 账本相关接口
	accountRouter := r.Group("/account")
	// 账单构建
	accountRouter.POST("/build", handler.BuildAccount)
	// 查询账单
	accountRouter.POST("/query", handler.QueryAccount)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
