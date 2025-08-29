package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Debug *log.Logger
	Error *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),                 //время
		Debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile), //время +строка
		Error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile), //время +строка
	}
}

func (l *Logger) INFO(msg string) {
	l.Info.Println(msg)
}

func (l *Logger) DEBUG(msg string) {
	l.Debug.Println(msg)
}

func (l *Logger) ERROR(msg string) {
	l.Error.Println(msg)
}
