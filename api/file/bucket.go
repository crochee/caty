// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package file

import (
	"net/http"
	"obs/pkg/e"
	"obs/pkg/log"
	"obs/pkg/resp"
	bucket2 "obs/pkg/service/business/bucket"
	tokenx2 "obs/pkg/service/business/tokenx"
	"obs/pkg/v"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
// @Failure 400 {object} e.Response
// @Failure 403 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/bucket [post]
func CreateBucket(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindBodyWith(&name, binding.JSON); err != nil {
		log.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		resp.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	token, err := tokenx2.QueryToken(ctx)
	if err != nil {
		log.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		resp.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx2.VerifyAuth(token.ActionMap, v.ServiceName, tokenx2.Write); err != nil {
		resp.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	if err = bucket2.CreateBucket(ctx.Request.Context(), token, name.BucketName); err != nil {
		resp.Errors(ctx, err)
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
// @Failure 400 {object} e.Response
// @Failure 401 {object} e.Response
// @Failure 403 {object} e.Response
// @Failure 404 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/bucket/{bucket_name} [get]
func GetBucket(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindUri(&name); err != nil {
		log.FromContext(ctx.Request.Context()).Errorf("bind uri failed.Error:%v", err)
		resp.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	token, err := tokenx2.QueryToken(ctx)
	if err != nil {
		log.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		resp.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx2.VerifyAuth(token.ActionMap, v.ServiceName, tokenx2.Read); err != nil {
		resp.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	var info *bucket2.Info
	if info, err = bucket2.HeadBucket(ctx.Request.Context(), token, name.BucketName); err != nil {
		resp.Errors(ctx, err)
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
// @Failure 400 {object} e.Response
// @Failure 403 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/bucket/{bucket_name} [delete]
func DeleteBucket(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindUri(&name); err != nil {
		log.FromContext(ctx.Request.Context()).Errorf("bind uri failed.Error:%v", err)
		resp.ErrorWith(ctx, e.ParseUrlFail, "bucket_name is nil")
		return
	}
	token, err := tokenx2.QueryToken(ctx)
	if err != nil {
		log.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		resp.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx2.VerifyAuth(token.ActionMap, v.ServiceName, tokenx2.Delete); err != nil {
		resp.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	if err = bucket2.DeleteBucket(ctx.Request.Context(), token, name.BucketName); err != nil {
		resp.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
