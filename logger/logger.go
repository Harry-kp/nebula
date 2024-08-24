package logger

import (
	"log"
	"os"
)

type Config struct {
	LogEnabled bool
}

type CustomLogger struct {
	logger     *log.Logger
	logEnabled bool
}

var globalLogger *CustomLogger

func Init(config Config) {
	globalLogger = &CustomLogger{
		logger:     log.New(os.Stdout, "", log.LstdFlags),
		logEnabled: config.LogEnabled,
	}
}

func Println(v ...interface{}) {
	if globalLogger.logEnabled {
		globalLogger.logger.Println(v...)
	}
}

func Printf(format string, v ...interface{}) {
	if globalLogger.logEnabled {
		globalLogger.logger.Printf(format, v...)
	}
}

func Fatal(v ...interface{}) {
	globalLogger.logger.Fatal(v...)
}
