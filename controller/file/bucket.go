// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/13

package file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"obs/logger"
	"obs/response"
	"obs/service/bucket"
	"obs/service/tokenx"
)

// CreateBucket godoc
// @Summary CreateBucket
// @Description create bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param request body Name true "bucket name"
// @Success 201 {int} int "bucket_id"
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket [post]
func CreateBucket(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindBodyWith(&name, binding.JSON); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "Check your payload"))
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	if token.ActionMap["OBS"] < tokenx.Write { //权限不够
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "Insufficient permissions"))
		return
	}
	var id uint
	if id, err = bucket.CreateBucket(ctx.Request.Context(), token, name.BucketName); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, id)
}

// HeadBucket godoc
// @Summary HeadBucket
// @Description head bucket
// @Tags bucket
// @Accept application/json
// @Produce  application/json
// @Param bucket_id path int true "bucket id"
// @Success 200 {object} bucket.Info "bucket info"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_id} [head]
func HeadBucket(ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind uri failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "bucket_name is nil"))
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	if token.ActionMap["OBS"] < tokenx.Read { //权限不够
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "Insufficient permissions"))
		return
	}
	var info *bucket.Info
	if info, err = bucket.HeadBucket(ctx.Request.Context(), token, id.BucketId); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, info)
}

// DeleteBucket godoc
// @Summary DeleteBucket
// @Description delete bucket
// @Tags bucket
// @Accept application/json
// @Produce application/json
// @Param bucket_id path int true "bucket id"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_id} [delete]
func DeleteBucket(ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind uri failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "bucket_name is nil"))
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	if token.ActionMap["OBS"] < tokenx.Delete { //权限不够
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "Insufficient permissions"))
		return
	}
	if err = bucket.DeleteBucket(ctx.Request.Context(), token, id.BucketId); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
