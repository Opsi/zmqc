package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelFatal
	LevelOff
)

var (
	DebugLog *log.Logger
	InfoLog  *log.Logger
	FatalLog *log.Logger
)

func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "fatal":
		return LevelFatal
	case "off":
		return LevelOff
	}
	fmt.Printf("Invalid log-level: %s\n", level)
	os.Exit(1)
	return LevelOff
}

func (level Level) ToPrefix() string {
	switch level {
	case LevelDebug:
		return "[DEBUG] "
	case LevelInfo:
		return "[INFO ] "
	case LevelFatal:
		return "[FATAL] "
	}
	return ""
}

func createLogger(writer *io.Writer, confLevel, loggerLevel Level, showTimestamps, showLogLevel bool) *log.Logger {
	if confLevel > loggerLevel {
		return log.New(ioutil.Discard, "", 0)
	}
	flags := 0
	if showTimestamps {
		flags |= log.Ldate | log.Ltime
	}
	if showLogLevel {
		return log.New(*writer, loggerLevel.ToPrefix(), flags)
	}
	return log.New(*writer, "", flags)
}

func InitLogger(level string, filePath string, showTimestamps, showLogLevel bool) {
	logLevel := ParseLevel(level)

	// log file or stdout
	var writer io.Writer = os.Stdout
	if filePath != "" {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Error opening log file: %s", err)
			os.Exit(1)
		}
		writer = file
	}

	DebugLog = createLogger(&writer, logLevel, LevelDebug, showTimestamps, showLogLevel)
	InfoLog = createLogger(&writer, logLevel, LevelInfo, showTimestamps, showLogLevel)
	FatalLog = createLogger(&writer, logLevel, LevelFatal, showTimestamps, showLogLevel)
}

func Debug(v ...interface{}) {
	DebugLog.Print(v...)
}

func Debugf(format string, v ...interface{}) {
	DebugLog.Printf(format, v...)
}

func Info(v ...interface{}) {
	InfoLog.Print(v...)
}

func Infof(format string, v ...interface{}) {
	InfoLog.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	FatalLog.Print(v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	FatalLog.Printf(format, v...)
	os.Exit(1)
}
