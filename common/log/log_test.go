package log

import (
	"go.uber.org/zap"
	"testing"
)

var (
	fakeConfig = &Config{
		MaxSize:    10,
		Compress:   true,
		LogPath:    "",
		MaxAge:     0,
		MaxBackups: 0,
		LogLevel:   "info",
	}
)

func TestLogger(t *testing.T) {
	Init(fakeConfig)
	Info("TestLog", zap.String("test", "eeyeyyeye"))
	Debug("debug", zap.String("debug", "debug"))
	Warn("warn", zap.String("warn", "warn"))
	Error("error", zap.String("error", "error"))
	//Panic("panic", zap.String("panic", "panic"))
	Fatal("fatal", zap.String("fatal", "fatal"))
}
