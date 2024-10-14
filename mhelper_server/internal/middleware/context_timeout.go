package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// ContextTimeout 统一的超时处理中间件
func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()

		// 将带有超时时间的Context注入给 Gin.Context 传递并用于超时处理
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
