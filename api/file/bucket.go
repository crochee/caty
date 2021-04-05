// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"obs/cmd"
	"obs/e"
	"obs/logger"
	"obs/service/business/bucket"
	"obs/service/business/tokenx"
)

// CreateBucket godoc
// @Summary CreateBucket
// @Description create bucket
// @Security ApiKeyAuth
// @Tags bucket
// @Accept application/json
// @Produce application/json
// @Param request body Name true "bucket name"
// @Success 201
// @Failure 400 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure 500 {object} e.ErrorResponse
// @Router /v1/bucket [post]
func CreateBucket(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindBodyWith(&name, binding.JSON); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		e.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, cmd.ServiceName, tokenx.Write); err != nil {
		e.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	if err = bucket.CreateBucket(ctx.Request.Context(), token, name.BucketName); err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

// GetBucket godoc
// @Summary GetBucket
// @Description get bucket
// @Security ApiKeyAuth
// @Tags bucket
// @Accept application/json
// @Produce application/json
// @Param bucket_name path string true "bucket name"
// @Success 200 {object} bucket.Info "bucket info"
// @Failure 400 {object} e.ErrorResponse
// @Failure 401 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure 404 {object} e.ErrorResponse
// @Failure 500 {object} e.ErrorResponse
// @Router /v1/bucket/{bucket_name} [get]
func GetBucket(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindUri(&name); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind uri failed.Error:%v", err)
		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		e.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, cmd.ServiceName, tokenx.Read); err != nil {
		e.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	var info *bucket.Info
	if info, err = bucket.HeadBucket(ctx.Request.Context(), token, name.BucketName); err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, info)
}

// DeleteBucket godoc
// @Summary DeleteBucket
// @Description delete bucket
// @Security ApiKeyAuth
// @Tags bucket
// @Accept application/json
// @Produce application/json
// @Param bucket_name path string true "bucket name"
// @Success 204
// @Failure 400 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure 500 {object} e.ErrorResponse
// @Router /v1/bucket/{bucket_name} [delete]
func DeleteBucket(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindUri(&name); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind uri failed.Error:%v", err)
		e.ErrorWith(ctx, e.ParseUrlFail, "bucket_name is nil")
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		e.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, cmd.ServiceName, tokenx.Delete); err != nil {
		e.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	if err = bucket.DeleteBucket(ctx.Request.Context(), token, name.BucketName); err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
