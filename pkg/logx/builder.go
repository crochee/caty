// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/4

package logx

import "context"

type Builder interface {
	Opt() Option
	With(key string, value interface{}) Builder
	Debugf(format string, v ...interface{})
	Debug(message string)
	Infof(format string, v ...interface{})
	Info(message string)
	Warnf(format string, v ...interface{})
	Warn(message string)
	Errorf(format string, v ...interface{})
	Error(message string)
	Fatalf(format string, v ...interface{})
	Fatal(message string)
	Sync() error
}

type loggerKey struct{}

// WithContext Adds fields.
func WithContext(ctx context.Context, log Builder) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// FromContext Gets the logx from context.
func FromContext(ctx context.Context) Builder {
	l, ok := ctx.Value(loggerKey{}).(Builder)
	if !ok {
		return NoLogger{}
	}
	return l
}

type NoLogger struct {
}

func (n NoLogger) Opt() Option {
	return Option{}
}

func (n NoLogger) With(string, interface{}) Builder {
	return n
}

func (n NoLogger) Debugf(string, ...interface{}) {
}

func (n NoLogger) Debug(string) {
}

func (n NoLogger) Infof(string, ...interface{}) {
}

func (n NoLogger) Info(string) {
}

func (n NoLogger) Warnf(string, ...interface{}) {
}

func (n NoLogger) Warn(string) {
}

func (n NoLogger) Errorf(string, ...interface{}) {
}

func (n NoLogger) Error(string) {
}

func (n NoLogger) Fatalf(string, ...interface{}) {
}

func (n NoLogger) Fatal(string) {
}

func (n NoLogger) Sync() error {
	return nil
}
