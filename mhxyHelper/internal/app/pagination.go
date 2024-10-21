package app

import (
	"github.com/Tokumicn/theBookofChangesEveryDay/mhxyHelper/pkg/convert"
	"github.com/gin-gonic/gin"
)

const (
	DefaultPageSize = 25
	MaxPageSize     = 50
)

func GetPage(c *gin.Context) int {
	page := convert.StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := convert.StrTo(c.Query("pageSize")).MustInt()
	if pageSize <= 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}

	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}

	return result
}
