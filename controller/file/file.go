// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/23

package file

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

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
// @Tags file
// @Accept multipart/form-data
// @Produce application/json
// @Param bucket_id path int true "bucket name"
// @Param file formData file true "file"
// @Success 201 int "file id"
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_id}/file [post]
func UploadFile(ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
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
		response.ErrorWith(ctx, response.Error(http.StatusInternalServerError, "Unauthorized"))
		return
	}
	if token.ActionMap["OBS"] < tokenx.Write { //权限不够
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "Insufficient permissions"))
		return
	}
	var fileId uint
	if fileId, err = bucket.UploadFile(ctx.Request.Context(), token, id.BucketId, fileInfo.File); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, fileId)
}

// DeleteFile godoc
// @Summary DeleteFile
// @Description delete file
// @Tags file
// @Accept application/json
// @Produce application/json
// @Param bucket_id path int true "bucket id"
// @Param file_id path int true "file id"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_id}/file/{file_id} [delete]
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
		response.ErrorWith(ctx, response.Error(http.StatusInternalServerError, "Unauthorized"))
		return
	}
	if token.ActionMap["OBS"] < tokenx.Delete { //权限不够
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "Insufficient permissions"))
		return
	}
	if err = bucket.DeleteFile(ctx.Request.Context(), token, target.BucketId, target.FileId); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// SignFile godoc
// @Summary SignFile
// @Description head file sign info
// @Tags file
// @Accept application/json
// @Produce application/json
// @Param bucket_id path int true "bucket id"
// @Param file_id query int true "file id"
// @Success 200 string "file link"
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_id}/file/{file_id} [head]
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
		response.ErrorWith(ctx, response.Error(http.StatusInternalServerError, "Unauthorized"))
		return
	}
	if token.ActionMap["OBS"] < tokenx.Read { //权限不够
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "Insufficient permissions"))
		return
	}
	var sign string
	if sign, err = bucket.SignFile(ctx.Request.Context(), token, target.BucketId, target.FileId); err != nil {
		response.ErrorWith(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("/v1/bucket/%d/file/%d?%s=%s",
		target.BucketId, target.FileId, middleware.Signature, sign))
}

// DownloadFile godoc
// @Summary DownloadFile
// @Description download File
// @Tags file
// @Accept application/json
// @Produce application/octet-stream
// @Param bucket_id path int true "bucket id"
// @Param file_id path int true "file id"
// @Param sign query string false "sign"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/bucket/{bucket_id}/file/{file_id} [get]
func DownloadFile(ctx *gin.Context) {
	bucketIdStr := ctx.Param("bucket_id")
	fileIdString := ctx.Param("file_id")
	if bucketIdStr == "" || fileIdString == "" {
		logger.FromContext(ctx.Request.Context()).Errorf("get url %s %s failed", bucketIdStr, fileIdString)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	bucketId, err := strconv.Atoi(bucketIdStr)
	if err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("strconv %s failed", bucketIdStr)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	var token *tokenx.Token
	if token, err = tokenx.QueryToken(ctx); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("query token failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusInternalServerError, "Unauthorized"))
		return
	}
	if token.ActionMap["OBS"] < tokenx.Read { //权限不够
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "Insufficient permissions"))
		return
	}

	conn := db.NewDB()
	b := new(db.Bucket)
	if err = conn.Model(b).Where("id =? AND domain= ?", bucketId, token.Domain).Find(b).Error; err != nil {
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		response.ErrorWith(ctx, response.Errors(http.StatusInternalServerError, err))
		return
	}
	bucketFile := &db.BucketFile{}
	if fileId, sErr := strconv.Atoi(fileIdString); sErr != nil {
		err = conn.Model(bucketFile).Where("file =? AND bucket_id= ?", fileIdString, bucketId).
			Find(bucketFile).Error
	} else {
		err = conn.Model(bucketFile).Where("id =? AND bucket_id= ?", fileId, bucketId).
			Find(bucketFile).Error
	}
	if err != nil {
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		response.ErrorWith(ctx, response.Errors(http.StatusInternalServerError, err))
		return
	}
	path := filepath.Clean(fmt.Sprintf("%s/%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath,
		b.Bucket, bucketFile.File))
	http.ServeFile(ctx.Writer, ctx.Request, path)
}
