package handler

import (
	"github.com/gin-gonic/gin"
	"mhxyHelper/internal/app"
	"mhxyHelper/internal/service/query/local_query"
	"mhxyHelper/pkg/errcode"
	"net/http"
)

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

	err = local_query.BuildStuffByStr(req.StuffStrArr)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorBuildStuffByStrFail.WithDetails(err.Error()))
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
	total, stuffs, err := local_query.QueryStuff(req.QueryStr)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorQueryStuffFail.WithDetails(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"list":  stuffs,
	})
	return
}
