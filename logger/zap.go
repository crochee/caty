// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package logger
package logger

import (
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Encoder func(zapcore.EncoderConfig) zapcore.Encoder

func NewZap(level string, encoderFunc Encoder, w io.Writer, fields ...zap.Field) *zap.Logger {
	core := zapcore.NewCore(
		encoderFunc(newEncoderConfig()),
		zap.CombineWriteSyncers(zapcore.AddSync(w)),
		newLevel(level),
	).With(fields) //自带node 信息
	//大于error增加堆栈信息
	return zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.DPanicLevel))
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "Message",
		LevelKey:       "Level",
		TimeKey:        "Time",
		NameKey:        "Logger",
		CallerKey:      "Caller",
		StacktraceKey:  "Stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func newLevel(level string) zapcore.Level {
	l := zap.InfoLevel
	if temp, ok := map[string]zapcore.Level{
		"DEBUG":  zap.DebugLevel,
		"INFO":   zap.InfoLevel,
		"WARN":   zap.WarnLevel,
		"ERROR":  zap.ErrorLevel,
		"DPANIC": zap.DPanicLevel,
		"PANIC":  zap.PanicLevel,
		"FATAL":  zap.FatalLevel,
	}[strings.ToUpper(level)]; ok {
		l = temp
	}
	return l
}
