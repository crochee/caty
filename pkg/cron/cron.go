// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/28

package cron

import (
	"obs/pkg/log"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var timeCron *cron.Cron

// Cron return cron.Cron
func Cron() *cron.Cron {
	return timeCron
}

// Setup init a cron.Cron
func Setup() {
	timeCron = cron.New(
		cron.WithLocation(time.UTC),
		cron.WithSeconds(),
		cron.WithChain(cron.Recover(cronLogger{})),
		cron.WithLogger(cronLogger{}),
	)
	timeCron.Start()
}

type cronLogger struct{}

func (c cronLogger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	log.Infof(formatString(len(keysAndValues)), append([]interface{}{msg}, keysAndValues...)...)
}

func (c cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	log.Errorf(formatString(len(keysAndValues)+2),
		append([]interface{}{msg, "error", err}, keysAndValues...)...)
}

// formatString returns a logfmt-like format string for the number of
// key/values.
func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString("%v:%v")
	}
	return sb.String()
}

// formatTimes formats any time.Time values as RFC3339.
func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
