package logger

import (
	"log/slog"
	"os"
)

var (
	Log *slog.Logger
	// logger      = logrus.New()
)

func NewLogger() {
	Log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// 设置日志级别为Debug
	//logger.SetLevel(logrus.DebugLevel) // logrus时使用
}
