// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/19

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"obs/model"
	"obs/service/verify"
	"obs/util"
)

// Verify 鉴权 (每一个请求还需要横向鉴权)
func Verify(ctx *gin.Context) {
	if SkipAuth(ctx) {
		return
	}
	var akSk model.AkSk
	if err := ctx.ShouldBindHeader(&akSk); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	t := &verify.Token{AkSecret: util.Slice(akSk.Ak)}
	if err := t.Verify(akSk.Sk); err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	if !verify.Authentication(t, MethodAction(ctx.Request.Method)) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Next()
}

var actionMap = map[string]verify.BucketAction{
	http.MethodGet:     verify.Read,
	http.MethodHead:    verify.Read,
	http.MethodOptions: verify.Read,
	http.MethodConnect: verify.Read,
	http.MethodTrace:   verify.Read,
	http.MethodPost:    verify.Write,
	http.MethodPut:     verify.Write,
	http.MethodPatch:   verify.Write,
	http.MethodDelete:  verify.Delete,
}

func MethodAction(method string) verify.BucketAction {
	if action, ok := actionMap[method]; ok {
		return action
	}
	return verify.Read
}

// 为所有不想加token到header的提供一个将token放入param key的途径
var notNeedVerifyApi = map[string]map[string]struct{}{
	"/v1/bucket/:bucket_name": {
		http.MethodPost: {},
	},
}

func SkipAuth(ctx *gin.Context) bool {
	if temp, ok := notNeedVerifyApi[ctx.FullPath()]; ok {
		if _, ok = temp[ctx.Request.Method]; ok {
			return true
		}
	}
	return false
}
