package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"mhelper_server/global"
	"mhelper_server/pkg/logger"
	"time"
)

// AccessLogWriter 访问日志 Writer
type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// AccessLog 访问日志中间件
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyW := &AccessLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = bodyW

		beginTime := time.Now().Unix()
		c.Next() // 记录本次请求的执行时间
		endTime := time.Now().Unix()
		fields := logger.Fields{ // 出参、入参
			"request":  c.Request.PostForm.Encode(),
			"response": bodyW.body.String(),
		}
		s := "access log: method: %s, status_code: %d, " +
			"begin_time: %d, end_time: %d"

		global.Logger.WithFields(fields).Infof(c, s,
			c.Request.Method, // 请求方法
			bodyW.Status(),   // 状态码
			beginTime,
			endTime,
		)
	}
}
