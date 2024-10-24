package handler

import (
	"github.com/gin-gonic/gin"
	"mhxyHelper/internal/app"
	"mhxyHelper/internal/service"
	"mhxyHelper/pkg/errcode"
	"net/http"
)

type BuildAccountReq struct {
	AccountStrArr []string `json:"accountStrArr"`
}

// 通过识别文本建立账单信息
func BuildAccount(c *gin.Context) {
	var req BuildAccountReq
	response := app.NewResponse(c)

	err := c.BindJSON(&req)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if len(req.AccountStrArr) <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		// c.JSON(http.StatusBadRequest, gin.H{"error": "参数异常"})
		return
	}

	err = service.BuildAccountInfo(req.AccountStrArr)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorBuildAccountByStrFail.WithDetails(err.Error()))
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type QueryAccountReq struct {
	QueryStr string `json:"queryStr"` // 查询的字段
	UserId   int64  `json:"userId"`   // userId
}

// 查询账单信息
func QueryAccount(c *gin.Context) {
	var req QueryAccountReq
	response := app.NewResponse(c)

	err := c.BindJSON(&req)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	total, accounts, err := service.QueryUserAccountInfo(req.QueryStr, req.UserId)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorQueryAccountFail.WithDetails(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"list":  accounts,
	})
	return
}
