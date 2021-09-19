// Package log
package log

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultLogSizeM int = 20
	DefaultMaxZip   int = 50
	MaxLogDays      int = 30
)

type Option struct {
	Path  string
	Level string
	Skip  int
}

// NewLogger 初始化日志对象
//
// @param: path 日志路径
// @param: level 日志等级
func NewLogger(opts ...func(*Option)) *Logger {
	l := &Logger{
		Option: Option{
			Path:  "",
			Level: "INFO",
			Skip:  1,
		},
	}
	for _, opt := range opts {
		opt(&l.Option)
	}
	l.Level = strings.ToUpper(l.Level)
	var encode encoder
	if l.Option.Path == "" {
		encode = zapcore.NewConsoleEncoder
	} else {
		encode = zapcore.NewJSONEncoder
	}
	l.Logger = newZap(l.Option.Level, encode, l.Option.Skip, SetLoggerWriter(l.Option.Path))
	l.LoggerSugar = l.Logger.Sugar()

	return l
}

type Logger struct {
	Logger      *zap.Logger
	LoggerSugar *zap.SugaredLogger
	Option
}

func (l *Logger) Opt() Option {
	return l.Option
}

func (l *Logger) With(key string, value interface{}) Builder {
	cpLog := l.Logger.With(zap.Any(key, value))
	return &Logger{
		Logger:      cpLog,
		LoggerSugar: cpLog.Sugar(),
		Option:      l.Option,
	}
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.LoggerSugar.Debugf(format, v...)
}

// Debug 打印Debug信息
//
// @param: message 格式信息
func (l *Logger) Debug(message string) {
	l.Logger.Debug(message)
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Infof(format string, v ...interface{}) {
	l.LoggerSugar.Infof(format, v...)
}

// Info 打印Info信息
//
// @param: message 格式信息
func (l *Logger) Info(message string) {
	l.Logger.Info(message)
}

// Warnf 打印Warn信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.LoggerSugar.Warnf(format, v...)
}

// Warn 打印Warn信息
//
// @param: message 信息
func (l *Logger) Warn(message string) {
	l.Logger.Warn(message)
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.LoggerSugar.Errorf(format, v...)
}

// Error 打印Error信息
//
// @param: message 信息
func (l *Logger) Error(message string) {
	l.Logger.Error(message)
}

// Fatalf 打印Fatalf信息
//
// @param: format 格式信息
// @param: v 参数信息
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.LoggerSugar.Errorf(format, v...)
}

// Fatal 打印Fatal信息
//
// @param: message 信息
func (l *Logger) Fatal(message string) {
	l.Logger.Error(message)
}

func (l *Logger) Sync() error {
	err := l.Logger.Sync()
	err1 := l.LoggerSugar.Sync()
	if err1 != nil {
		err = err1
	}
	return err
}
