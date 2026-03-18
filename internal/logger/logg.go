package logger

import (
	"fmt"
	"os"
	"time"
)
// наш логгер с уровнями debug info error
type Logger struct {
	debug bool
}

func (l *Logger) print(level string, format string, args ...any) {
	if l == nil {
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] [%s] \t %s\n", t, level, msg)
}
func NewLogger(debug bool) *Logger {
	return &Logger{debug: debug}
}
func (l *Logger) Info(format string, args ...any) {
	l.print("INFO", format, args...)
}
func (l *Logger) Debug(format string, args ...any) {
	if l.debug {
		l.print("DEBUG", format, args...)
	}
}
func (l *Logger) Error(format string, args ...any) {
	l.print("ERROR", format, args...)
}
