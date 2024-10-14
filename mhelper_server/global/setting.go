package global

import (
	"mhelper_server/pkg/logger"
	"mhelper_server/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS   // httpServer配置
	AppSetting      *setting.AppSettingS      // 应用配置
	EmailSetting    *setting.EmailSettingS    // Email配置
	JWTSetting      *setting.JWTSettingS      // JWTToken配置
	DatabaseSetting *setting.DatabaseSettingS // 数据库配置
	Logger          *logger.Logger            // 日志配置
)
