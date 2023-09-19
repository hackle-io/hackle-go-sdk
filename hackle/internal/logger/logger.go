package logger

import (
	"log"
	"os"
)

var globalLogger = NewLogger()

func Info(format string, v ...interface{}) {
	globalLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	globalLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	globalLogger.Error(format, v...)
}

type Logger interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

func NewLogger() Logger {
	return &HackleLogger{
		log: log.New(os.Stdout, "[Hackle] ", log.LstdFlags),
	}
}

type HackleLogger struct {
	log *log.Logger
}

func (l *HackleLogger) Info(format string, v ...interface{}) {
	l.log.Printf("INFO - "+format, v...)
}

func (l *HackleLogger) Warn(format string, v ...interface{}) {
	l.log.Printf("WARN - "+format, v...)
}

func (l *HackleLogger) Error(format string, v ...interface{}) {
	l.log.Printf("ERROR - "+format, v...)
}
