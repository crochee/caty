// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2020/12/8

// Package middleware
package middleware

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// Limit 限流器
func Limit(ctx *gin.Context) {
	if !Limiter().IsAvailable() {
		ctx.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	Limiter().Increase()
	ctx.Next()
}

var (
	limiter     LimiterInterface
	limiterOnce sync.Once
)

func Limiter() LimiterInterface {
	limiterOnce.Do(func() {
		limiter = NewLimiter(10*time.Millisecond, 10000)
	})
	return limiter
}

type requestLimit struct {
	interval time.Duration
	maxCount uint64
	reqCount uint64
	mux      *sync.Mutex
}

func NewLimiter(interval time.Duration, maxCnt uint64) *requestLimit {
	reqLimit := &requestLimit{
		interval: interval,
		maxCount: maxCnt,
		mux:      new(sync.Mutex),
	}
	go func() {
		t := time.NewTicker(interval)
		for range t.C {
			atomic.StoreUint64(&reqLimit.reqCount, 0)
		}
	}()
	return reqLimit
}

func (r *requestLimit) Increase() {
	atomic.AddUint64(&r.reqCount, 1)
}

func (r *requestLimit) IsAvailable() bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.reqCount < r.maxCount
}

type LimiterInterface interface {
	Increase()
	IsAvailable() bool
}
