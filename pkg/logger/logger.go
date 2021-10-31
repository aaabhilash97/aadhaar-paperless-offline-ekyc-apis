package logger

import (
	"log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel          string
	LogOutputPaths    string
	DisableStackTrace bool
}

var (
	// Log is global logger
	Log *zap.Logger
)

// Info - log to info level
var Info func(msg string, fields ...zap.Field)

// Error - log to error level
var Error func(msg string, fields ...zap.Field)

// Debug - log to debug level
var Debug func(msg string, fields ...zap.Field)

// Warn  - log to warn level
var Warn func(msg string, fields ...zap.Field)

// Fatal  - log to Fatal level
var Fatal func(msg string, fields ...zap.Field)

func init() {
	config := Config{}

	logger := InitLogger(config)
	Log = logger
	Info = logger.Info
	Error = logger.Error
	Debug = logger.Debug
	Warn = logger.Warn
	Fatal = logger.Fatal

}

func Init(config Config) {
	logger := InitLogger(config)
	Log = logger
	Info = logger.Info
	Error = logger.Error
	Debug = logger.Debug
	Warn = logger.Warn
	Fatal = logger.Fatal
}

// InitLogger initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
func InitLogger(config Config) *zap.Logger {
	if config.LogOutputPaths == "" {
		config.LogOutputPaths = "stdout"
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.DisableStacktrace = config.DisableStackTrace
	loggerConfig.OutputPaths = strings.Split(config.LogOutputPaths, ",")
	loggerConfig.Level = zap.NewAtomicLevelAt(
		zapcore.Level(LogLevels[config.LogLevel]),
	)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	return logger
}

// LogLevels is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
var LogLevels = map[string]int{
	"debug":  -1,
	"info":   0,
	"warn":   1,
	"error":  2,
	"dpanic": 3,
	"panic":  4,
	"fatal":  5,
}
