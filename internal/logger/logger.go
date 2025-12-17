package logger

import (
	"log"
	"os"
)

// Logger определяет интерфейс для логирования
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, err error, args ...interface{})
	Debug(msg string, args ...interface{})
}

// StandardLogger реализует Logger используя стандартный log пакет
type StandardLogger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
}

// New создает новый стандартный логгер
func New() Logger {
	return &StandardLogger{
		infoLog:  log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		errorLog: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debugLog: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *StandardLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLog.Printf(msg, args...)
	} else {
		l.infoLog.Println(msg)
	}
}

func (l *StandardLogger) Error(msg string, err error, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			l.errorLog.Printf(msg+": %v", append(args, err)...)
		} else {
			l.errorLog.Printf(msg+": %v", err)
		}
	} else {
		if len(args) > 0 {
			l.errorLog.Printf(msg, args...)
		} else {
			l.errorLog.Println(msg)
		}
	}
}

func (l *StandardLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.debugLog.Printf(msg, args...)
	} else {
		l.debugLog.Println(msg)
	}
}

// NoOpLogger реализует Logger но ничего не логирует (для тестов)
type NoOpLogger struct{}

func (n *NoOpLogger) Info(msg string, args ...interface{})   {}
func (n *NoOpLogger) Error(msg string, err error, args ...interface{}) {}
func (n *NoOpLogger) Debug(msg string, args ...interface{})  {}

