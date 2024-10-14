package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mhelper_server/global"
	"mhelper_server/pkg/app"
	"mhelper_server/pkg/email"
	"mhelper_server/pkg/errcode"
	"time"
)

func Recovery() gin.HandlerFunc {
	defaultMailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				// 记录 panic 日志
				global.Logger.WithCallersFrames().Errorf(c, "[PANIC] panic recover err: %v", err)

				// 发送 Email 通知管理员
				defaultMailer.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("异常抛出，发生时间: %d", time.Now().Unix()),
					fmt.Sprintf("错误信息: %v", err),
				)
				if err != nil {
					global.Logger.Panicf(c, "mail.SendMail err: %v", err)
				}
				// 标准化错误日志
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort() // 终止本次请求
			}
		}()
		c.Next()
	}
}
