package handler

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 通过识别文本建立字典
func BuildDict(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "BuildDictByStr",
	})
}
