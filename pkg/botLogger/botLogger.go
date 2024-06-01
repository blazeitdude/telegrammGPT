package botLogger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"strings"
	"sync"
)

type LoggerConfig struct {
	LogLevel string `yaml:"logLevel"`
}

type botLogger struct {
	Logger *zap.SugaredLogger
}

var (
	Logger *botLogger
	once   sync.Once
)

func (b botLogger) Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	b.Logger.Debug(message)
}

func (b botLogger) Printf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	b.Logger.Debug(message)
}

func getLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	default:
		log.Println("unknown logging level, the standard INFO mode is set")
		return zapcore.InfoLevel
	}
}

func InitLogger(conf LoggerConfig) {
	once.Do(func() {
		zapConfig := zap.NewProductionConfig()
		lvl := getLevel(conf.LogLevel)
		zapConfig.Level.SetLevel(lvl)
		prodlogger, err := zapConfig.Build()
		if err != nil {
			log.Fatal("failed to start logger")
		}
		defer prodlogger.Sync()

		suggaredLogger := prodlogger.Sugar()
		botLogger := botLogger{
			Logger: suggaredLogger,
		}
		Logger = &botLogger
	})
}

func GetLogger() *botLogger {
	return Logger
}
