// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/23

package file

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"obs/config"
	"obs/logger"
	"obs/middleware"
	"obs/model/db"
	"obs/response"
	"obs/service/bucket"
	"obs/service/tokenx"
)

// UploadFile godoc
// @Summary UploadFile
// @Description upload file
// @Security ApiKeyAuth
// @Tags file
// @Accept multipart/form-data
// @Produce application/json
// @Param bucket_name path string true "bucket name"
// @Param file formData file true "file"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file [post]
func UploadFile(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindUri(&name); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	var fileInfo Info
	if err := ctx.ShouldBindWith(&fileInfo, binding.FormMultipart); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your payload"))
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, "OBS", tokenx.Write); err != nil {
		response.ErrorWith(ctx, response.Errors(http.StatusForbidden, err))
		return
	}
	if err = bucket.UploadFile(ctx.Request.Context(), token, name.BucketName, fileInfo.File); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// DeleteFile godoc
// @Summary DeleteFile
// @Description delete file
// @Security ApiKeyAuth
// @Tags file
// @Accept application/json
// @Produce application/json
// @Param bucket_name path string true "bucket name"
// @Param file_name path string true "file name"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file/{file_name} [delete]
func DeleteFile(ctx *gin.Context) {
	var target Target
	if err := ctx.ShouldBindUri(&target); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, "OBS", tokenx.Delete); err != nil {
		response.ErrorWith(ctx, response.Errors(http.StatusForbidden, err))
		return
	}
	if err = bucket.DeleteFile(ctx.Request.Context(), token, target.BucketName, target.FileName); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// SignFile godoc
// @Summary SignFile
// @Description get file sign info
// @Security ApiKeyAuth
// @Tags file
// @Accept application/json
// @Produce application/json
// @Param bucket_name path string true "bucket name"
// @Param file_name path string true "file name"
// @Success 200 {string} string "file link"
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file/{file_name}/sign [get]
func SignFile(ctx *gin.Context) {
	var target Target
	if err := ctx.ShouldBindUri(&target); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, "OBS", tokenx.Read); err != nil {
		response.ErrorWith(ctx, response.Errors(http.StatusForbidden, err))
		return
	}
	var sign string
	if sign, err = bucket.SignFile(ctx.Request.Context(), token, target.BucketName, target.FileName); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("/v1/bucket/%s/file/%s?%s=%s",
		target.BucketName, target.FileName, middleware.Signature, sign))
}

// DownloadFile godoc
// @Summary DownloadFile
// @Description download File
// @Security ApiKeyAuth
// @Tags file
// @Accept application/json
// @Produce application/octet-stream
// @Param bucket_name path string true "bucket name"
// @Param file_name path string true "file name"
// @Param sign query string false "sign"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file/{file_name} [get]
func DownloadFile(ctx *gin.Context) {
	var target Target
	if err := ctx.ShouldBindUri(&target); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}

	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, "OBS", tokenx.Read); err != nil {
		response.ErrorWith(ctx, response.Errors(http.StatusForbidden, err))
		return
	}

	conn := db.NewDB()
	b := new(db.Bucket)
	if err = conn.Model(b).Where("bucket =? AND domain= ?",
		target.BucketName, token.Domain).Find(b).Error; err != nil {
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		response.ErrorWith(ctx, response.Errors(http.StatusInternalServerError, err))
		return
	}
	bucketFile := &db.BucketFile{}
	if err = conn.Model(bucketFile).Where("file =? AND bucket= ?",
		target.FileName, b.Bucket).Find(bucketFile).Error; err != nil {
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		response.ErrorWith(ctx, response.Errors(http.StatusInternalServerError, err))
		return
	}
	path := filepath.Clean(fmt.Sprintf("%s/%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath,
		b.Bucket, bucketFile.File))
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, url.PathEscape(bucketFile.File)))
	http.ServeFile(ctx.Writer, ctx.Request, path)
}
