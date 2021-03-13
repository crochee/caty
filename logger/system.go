// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package logger

import (
	"os"

	"go.uber.org/zap/zapcore"
)

var systemLogger *Logger

// InitSystemLogger 初始化系统级日志对象
//
// @param: path 日志路径
// @param: level 日志等级
func InitSystemLogger(path, level string) {
	systemLogger = &Logger{
		Logger: NewZap(level, zapcore.NewJSONEncoder, SetLoggerWriter(path)),
	}
	systemLogger.LoggerSugar = systemLogger.Logger.Sugar()
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func Debugf(format string, v ...interface{}) {
	if systemLogger != nil {
		systemLogger.LoggerSugar.Debugf(format, v...)
	}
}

// Debug 打印Debug信息
//
// @param: message 信息
func Debug(message string) {
	if systemLogger != nil {
		systemLogger.Logger.Debug(message)
	}
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func Infof(format string, v ...interface{}) {
	if systemLogger != nil {
		systemLogger.LoggerSugar.Infof(format, v...)
	}
}

// Info 打印Info信息
//
// @param: message 信息
func Info(message string) {
	if systemLogger != nil {
		systemLogger.Logger.Info(message)
	}
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func Errorf(format string, v ...interface{}) {
	if systemLogger != nil {
		systemLogger.LoggerSugar.Errorf(format, v...)
	}
}

// Error 打印Error信息
//
// @param: message 信息
func Error(message string) {
	if systemLogger != nil {
		systemLogger.Logger.Error(message)
	}
}

// Fatalf 打印Fatal信息
//
// @param: format 格式信息
// @param: v 参数信息
func Fatalf(format string, v ...interface{}) {
	if systemLogger != nil {
		systemLogger.LoggerSugar.Fatalf(format, v...)
	}
}

// Fatal 打印Fatal信息
//
// @param: message 信息
func Fatal(message string) {
	if systemLogger != nil {
		systemLogger.Logger.Fatal(message)
	}
}

// Exit 打印系统退出信息
//
// @param: message 信息
func Exit(message string) {
	if systemLogger != nil {
		systemLogger.Logger.Info(message)
		_ = systemLogger.Logger.Sync()
		_ = systemLogger.LoggerSugar.Sync()
	}
	os.Exit(1)
}
