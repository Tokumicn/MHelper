package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workHelper/internal/app"
	"workHelper/internal/service"
	"workHelper/pkg/errcode"
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

type BuildStuffReq struct {
	StuffStrArr []string `json:"stuffStrArr"`
}

// 通过识别文本建立商品信息
func BuildStuff(c *gin.Context) {
	var req BuildStuffReq
	response := app.NewResponse(c)

	err := c.BindJSON(&req)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if len(req.StuffStrArr) <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		// c.JSON(http.StatusBadRequest, gin.H{"error": "参数异常"})
		return
	}

	err = service.BuildStuffByStr(req.StuffStrArr)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorBuildHelperFail.WithDetails(err.Error()))
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type QueryStuffReq struct {
	QueryStr string `json:"queryStr"` // 查询的字段
}

// 查询物品信息
func QueryStuff(c *gin.Context) {
	var req QueryStuffReq
	response := app.NewResponse(c)

	err := c.BindJSON(&req)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	total, stuffs, err := service.QueryStuff(req.QueryStr)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGetHelperFail.WithDetails(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"list":  stuffs,
	})
	return
}
