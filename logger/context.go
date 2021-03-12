// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/12

package logger

import "context"

type loggerKey struct{}

// With Adds fields.
func With(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// FromContext Gets the logger from context.
func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(loggerKey{}).(*Logger)
	if !ok {
		return nil
	}
	return l
}
