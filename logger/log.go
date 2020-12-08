// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package logger
package logger

import (
	"io"
	"obs/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

const (
	DefaultLogSizeM int = 20
	DefaultMaxZip   int = 50
	MaxLogDays      int = 30
)

func SetLoggerWriter(path string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    DefaultLogSizeM, //单个日志文件最大MaxSize*M大小 // megabytes
		MaxAge:     MaxLogDays,      //days
		MaxBackups: DefaultMaxZip,   //备份数量
		Compress:   false,           //不压缩
		LocalTime:  true,            //备份名采用本地时间
	}
}

// InitLogger 初始化日志组件
func InitLogger() {
	logger = NewZap(config.Cfg.ServiceInfo.Mode,
		zapcore.NewJSONEncoder, SetLoggerWriter(config.Cfg.ServiceInfo.Log))
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func Infof(format string, v ...interface{}) {
	if logger != nil {
		logger.Sugar().Infof(format, v...)
	}
}

func Info(message string) {
	if logger != nil {
		logger.Info(message)
	}
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func Debugf(format string, v ...interface{}) {
	if logger != nil {
		logger.Sugar().Debugf(format, v...)
	}
}

func Debug(message string) {
	if logger != nil {
		logger.Debug(message)
	}
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func Errorf(format string, v ...interface{}) {
	if logger != nil {
		logger.Sugar().Errorf(format, v...)
	}
}

func Error(message string) {
	if logger != nil {
		logger.Error(message)
	}
}

func Exit() {
	if logger != nil {
		logger.Info("Server exiting...")
		_ = logger.Sync()
	}
}
