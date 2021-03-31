// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/23

package file

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"

	"obs/cmd"
	"obs/config"
	"obs/e"
	"obs/logger"
	"obs/model/db"
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
// @Failure 400 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure 500 {object} e.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file [post]
func UploadFile(ctx *gin.Context) {
	var name Name
	if err := ctx.ShouldBindUri(&name); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		e.Error(ctx, e.ParseUrlFail)
		return
	}
	var fileInfo Info
	if err := ctx.ShouldBindWith(&fileInfo, binding.FormMultipart); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind body failed.Error:%v", err)
		e.ErrorWith(ctx, e.ParsePayloadFailed, err.Error())
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		e.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, cmd.ServiceName, tokenx.Write); err != nil {
		e.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	if err = bucket.UploadFile(ctx.Request.Context(), token, name.BucketName, fileInfo.File); err != nil {
		e.Errors(ctx, err)
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
// @Failure 400 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure 500 {object} e.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file/{file_name} [delete]
func DeleteFile(ctx *gin.Context) {
	var target Target
	if err := ctx.ShouldBindUri(&target); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		e.Error(ctx, e.ParseUrlFail)
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		e.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, cmd.ServiceName, tokenx.Delete); err != nil {
		e.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	if err = bucket.DeleteFile(ctx.Request.Context(), token, target.BucketName, target.FileName); err != nil {
		e.Errors(ctx, err)
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
// @Failure 400 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure 500 {object} e.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file/{file_name}/sign [get]
func SignFile(ctx *gin.Context) {
	var target Target
	if err := ctx.ShouldBindUri(&target); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		e.Error(ctx, e.ParseUrlFail)
		return
	}
	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		e.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, cmd.ServiceName, tokenx.Read); err != nil {
		e.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}
	var sign string
	if sign, err = bucket.SignFile(ctx.Request.Context(), token, target.BucketName, target.FileName); err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("/v1/bucket/%s/file/%s?%s=%s",
		target.BucketName, target.FileName, "sign", sign))
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
// @Failure 400 {object} e.ErrorResponse
// @Failure 403 {object} e.ErrorResponse
// @Failure 500 {object} e.ErrorResponse
// @Router /v1/bucket/{bucket_name}/file/{file_name} [get]
func DownloadFile(ctx *gin.Context) {
	var target Target
	if err := ctx.ShouldBindUri(&target); err != nil {
		logger.FromContext(ctx.Request.Context()).Errorf("bind url failed.Error:%v", err)
		e.Error(ctx, e.ParseUrlFail)
		return
	}

	token, err := tokenx.QueryToken(ctx)
	if err != nil {
		e.ErrorWith(ctx, e.GetTokenFail, err.Error())
		return
	}
	if err = tokenx.VerifyAuth(token.ActionMap, cmd.ServiceName, tokenx.Read); err != nil {
		e.ErrorWith(ctx, e.Forbidden, err.Error())
		return
	}

	conn := db.NewDB()
	b := new(db.Bucket)
	if err = conn.Model(b).Where("bucket =? AND domain= ?",
		target.BucketName, token.Domain).First(b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Error(ctx, e.NotFound)
			return
		}
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		e.ErrorWith(ctx, e.OperateDbFail, err.Error())
		return
	}
	bucketFile := &db.BucketFile{}
	if err = conn.Model(bucketFile).Where("file =? AND bucket= ?",
		target.FileName, b.Bucket).First(bucketFile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Error(ctx, e.NotFound)
			return
		}
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		e.ErrorWith(ctx, e.OperateDbFail, err.Error())
		return
	}
	path := filepath.Clean(fmt.Sprintf("%s/%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath,
		b.Bucket, bucketFile.File))
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, url.PathEscape(bucketFile.File)))
	http.ServeFile(ctx.Writer, ctx.Request, path)
}
