package setting

import "time"

// ServerSettingS http服务配置
type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// AppSettingS 应用配置
type AppSettingS struct {
	DefaultPageSize       int
	MaxPageSize           int
	DefaultContextTimeout time.Duration
	LogSavePath           string
	LogFileName           string
	LogFileExt            string
	UploadSavePath        string
	UploadServerUrl       string
	UploadImageMaxSize    int
	UploadImageAllowExts  []string
}

// EmailSettingS Email配置
type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

// JWTSettingS JWT Token 配置
type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

// DatabaseSettingS 数据库配置
type DatabaseSettingS struct {
	DBType                string
	UserName              string
	Password              string
	Host                  string
	DBName                string
	TablePrefix           string
	Charset               string
	ParseTime             bool
	MaxIdleConns          int
	MaxOpenConns          int
	SlowThreshold         int // 慢查询超时时间  单位: ms(毫秒)
	SlowThresholdDuration time.Duration
}

var sections = make(map[string]interface{})
