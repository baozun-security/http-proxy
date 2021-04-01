package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_logger *zap.Logger
)

const (
	// FormatText format log text
	FormatText = "text"
	// FormatJSON format log json
	FormatJSON = "json"
)

// type Level uint

// 日志配置
type Config struct {
	LogPath    string
	LogLevel   string
	Compress   bool
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Format     string
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func newLogWriter(logPath string, maxSize, maxBackups, maxAge int, compress bool) io.Writer {
	if logPath == "" || logPath == "-" {
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
}

func newZapEncoder() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return encoderConfig
}

func newLoggerCore(cfg *Config) zapcore.Core {
	hook := newLogWriter(cfg.LogPath, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)

	encoderConfig := newZapEncoder()

	atomLevel := zap.NewAtomicLevelAt(getZapLevel(cfg.LogLevel))

	var encoder zapcore.Encoder
	if cfg.Format == FormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook)),
		atomLevel,
	)
	return core
}

func newLoggerOptions() []zap.Option {
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	callerskip := zap.AddCallerSkip(1)
	// 开发者
	development := zap.Development()
	options := []zap.Option{
		caller,
		callerskip,
		development,
	}
	return options
}

// fill default config
func (c *Config) fillWithDefault() {
	if c.MaxSize <= 0 {
		c.MaxSize = 20
	}
	if c.MaxAge <= 0 {
		c.MaxAge = 7
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 7
	}
	if c.LogLevel == "" {
		c.LogLevel = "debug"
	}
	if c.Format == "" {
		c.Format = FormatText
	}
}

// InitLog conf
func Init(cfg *Config) {
	cfg.fillWithDefault()
	core := newLoggerCore(cfg)
	zapOpts := newLoggerOptions()
	_logger = zap.New(core, zapOpts...)
}

func Logger() *zap.Logger {
	return _logger
}

// Debug output log
func Debug(msg string, fields ...zap.Field) {
	_logger.Debug(msg, fields...)
}

// Info output log
func Info(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

// Warn output log
func Warn(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}

// Error output log
func Error(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}

// Panic output panic
func Panic(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}

// Fatal output log
func Fatal(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}
