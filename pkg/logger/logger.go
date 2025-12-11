package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

// LogLevel представляет уровень логирования
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	info    *log.Logger
	warn    *log.Logger
	err     *log.Logger
	debug   *log.Logger
	logFile *os.File
	level   LogLevel
}

// parseLogLevel преобразует строку в LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo // default
	}
}

// New создает новый экземпляр логгера с записью в консоль и файл
func New(logFilePath string, level string) (*Logger, error) {
	// Открываем файл для логов
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Настраиваем мультиплексор для warning и error
	warnErrorWriter := io.MultiWriter(os.Stdout, logFile)

	return &Logger{
		debug:   log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
		info:    log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		warn:    log.New(warnErrorWriter, "[WARN] ", log.LstdFlags|log.Lshortfile),
		err:     log.New(warnErrorWriter, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		logFile: logFile,
		level:   parseLogLevel(level),
	}, nil
}

// Close закрывает файл логов
func (l *Logger) Close() error {
	if l != nil && l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Debug логирует отладочные сообщения (только консоль)
func (l *Logger) Debug(format string, v ...interface{}) {
	if l != nil && l.level <= LevelDebug {
		l.debug.Printf(format, v...)
	}
}

// Info логирует информационные сообщения (только консоль)
func (l *Logger) Info(format string, v ...interface{}) {
	if l != nil && l.level <= LevelInfo {
		l.info.Printf(format, v...)
	}
}

// Warn логирует предупреждения (консоль + файл)
func (l *Logger) Warn(format string, v ...interface{}) {
	if l != nil && l.level <= LevelWarn {
		l.warn.Printf(format, v...)
	}
}

// Error логирует ошибки (консоль + файл)
func (l *Logger) Error(format string, v ...interface{}) {
	if l != nil && l.level <= LevelError {
		l.err.Printf(format, v...)
	}
}

// Fatal логирует критическую ошибку и завершает программу
func (l *Logger) Fatal(format string, v ...interface{}) {
	if l != nil {
		l.err.Fatalf(format, v...)
	}
}
