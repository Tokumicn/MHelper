package handler

import (
	"github.com/gin-gonic/gin"
	"mhxyHelper/internal/app"
	"mhxyHelper/internal/service"
	"mhxyHelper/pkg/errcode"
	"net/http"
)

type BuildAttributeReq struct {
	AttributeStrArr []string `json:"attributeStrArr"`
}

// 通过识别文本建立商品信息
func BuildAttribute(c *gin.Context) {
	var req BuildAttributeReq
	response := app.NewResponse(c)

	err := c.BindJSON(&req)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if len(req.AttributeStrArr) <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		// c.JSON(http.StatusBadRequest, gin.H{"error": "参数异常"})
		return
	}

	err = service.BuildAttributeByStr(req.AttributeStrArr)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorBuildAttributeByStrFail.WithDetails(err.Error()))
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type QueryAttributeReq struct {
	QueryStr string `json:"queryStr"` // 查询的字段
}

// 查询属性信息
func QueryAttribute(c *gin.Context) {
	var req QueryAttributeReq
	response := app.NewResponse(c)

	err := c.BindJSON(&req)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	total, attrs, err := service.QueryAttribute(req.QueryStr)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorQueryAttributeFail.WithDetails(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"list":  attrs,
	})
	return
}
