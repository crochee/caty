// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/23

package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"obs/config"
	"obs/logger"
	"obs/model"
	"obs/response"
	"obs/service/verify"
	"obs/util"
)

// UploadFile godoc
// @Summary UploadFile
// @Description upload file
// @Tags file
// @Accept multipart/form-data
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Param file formData file true "file"
// @Param path formData string true "file"
// @Success 200 {object} model.FileTarget
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file/{bucket_name} [post]
func UploadFile(ctx *gin.Context) {
	var bucketName model.BucketName
	if err := ctx.ShouldBindUri(&bucketName); err != nil {
		logger.Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	var fileInfo model.FileInfo
	if err := ctx.ShouldBindWith(&fileInfo, binding.FormMultipart); err != nil {
		logger.Errorf("bind body failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your payload"))
		return
	}
	path := fmt.Sprintf("%s%s/%s/%s", config.Cfg.YamlConfig.ServiceInformation.SaveRootPath,
		bucketName.BucketName, fileInfo.Path, fileInfo.File.Filename)
	fileDir := filepath.Dir(path)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		logger.Errorf("mkdir %s failed.Error:%v", fileDir, err)
		response.ErrorWithMessage(ctx, "upload file failed")
		return
	}
	dstFile, err := os.Create(path)
	if err != nil {
		logger.Errorf("create %s file failed.Error:%v", path, err)
		response.ErrorWithMessage(ctx, "upload file failed")
		return
	}
	defer dstFile.Close()
	var srcFile multipart.File
	if srcFile, err = fileInfo.File.Open(); err != nil {
		logger.Errorf("open %s file failed.Error:%v", fileInfo.File.Filename, err)
		response.ErrorWithMessage(ctx, "open file failed")
		return
	}
	defer srcFile.Close()
	buf := util.AcquireBuf()
	defer util.ReleaseBuf(buf)
	var length int64
	if length, err = io.CopyBuffer(dstFile, srcFile, buf); err != nil {
		logger.Errorf("copy failed.Error:%v", err)
		response.ErrorWithMessage(ctx, "copy failed")
		return
	}
	if length != fileInfo.File.Size {
		logger.Errorf("file write %d,but need %d", length, fileInfo.File.Size)
		response.ErrorWithMessage(ctx, "write size wrong")
		return
	}
	newToken := &verify.Token{ExpiresAt: time.Now().Add(30 * time.Minute)}
	newToken.AddAction(verify.Read)
	ak, sk, err := newToken.Create()
	if err != nil {
		logger.Errorf("create token failed.Error:%v", err)
		response.ErrorWithMessage(ctx, "create sign failed")
		return
	}
	v := url.Values{}
	v.Set("path", fmt.Sprintf("%s/%s", fileInfo.Path, fileInfo.File.Filename))
	v.Add("ak", ak)
	v.Add("sk", sk)
	ctx.JSON(http.StatusOK, &model.FileTarget{Path: fmt.Sprintf("%s/v1/file/%s?%s",
		fmt.Sprintf("%s:%d", config.Cfg.IP.String(), config.Cfg.YamlConfig.ServiceInformation.Port),
		bucketName.BucketName,
		v.Encode(),
	)})
}

// DeleteFile godoc
// @Summary DeleteFile
// @Description delete file
// @Tags file
// @Accept application/json
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Param path query string true "target path"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file/{bucket_name} [delete]
func DeleteFile(ctx *gin.Context) {
	var bucketName model.BucketName
	if err := ctx.ShouldBindUri(&bucketName); err != nil {
		logger.Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	var fileTarget model.FileTarget
	if err := ctx.ShouldBindQuery(&fileTarget); err != nil {
		logger.Errorf("bind query failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your query"))
		return
	}
	path := fmt.Sprintf("%s%s/%s", config.Cfg.YamlConfig.ServiceInformation.SaveRootPath,
		bucketName.BucketName, fileTarget.Path)
	if err := os.RemoveAll(path); err != nil {
		logger.Errorf("delete path(%s) failed.Error:%v", path, err)
		response.ErrorWithMessage(ctx, "delete target failed")
		return
	}
	ctx.Status(http.StatusOK)
}

// SignFile godoc
// @Summary SignFile
// @Description head file sign info
// @Tags file
// @Accept application/json
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Param path query string true "target path"
// @Success 200 {object} model.FileTarget
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file/{bucket_name} [head]
func SignFile(ctx *gin.Context) {
	var bucketName model.BucketName
	if err := ctx.ShouldBindUri(&bucketName); err != nil {
		logger.Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	var fileTarget model.FileTarget
	if err := ctx.ShouldBindQuery(&fileTarget); err != nil {
		logger.Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	newToken := &verify.Token{ExpiresAt: time.Now().Add(30 * time.Minute)}
	newToken.AddAction(verify.Read)
	ak, sk, err := newToken.Create()
	if err != nil {
		logger.Errorf("create token failed.Error:%v", err)
		response.ErrorWithMessage(ctx, "create sign failed")
		return
	}
	v := url.Values{}
	v.Set("path", fileTarget.Path)
	v.Add("ak", ak)
	v.Add("sk", sk)
	ctx.JSON(http.StatusOK, &model.FileTarget{Path: fmt.Sprintf("%s/v1/file/%s?%s",
		fmt.Sprintf("%s:%d", config.Cfg.IP.String(), config.Cfg.YamlConfig.ServiceInformation.Port),
		bucketName.BucketName,
		v.Encode(),
	)})
}

// DownloadFile godoc
// @Summary DownloadFile
// @Description download File
// @Tags file
// @Accept application/json
// @Produce application/octet-stream
// @Param bucket_name path string true "bucket name"
// @Param path query string true "target path"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file/{bucket_name} [get]
func DownloadFile(ctx *gin.Context) {
	var bucketName model.BucketName
	if err := ctx.ShouldBindUri(&bucketName); err != nil {
		logger.Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	var fileTarget model.FileTarget
	if err := ctx.ShouldBindQuery(&fileTarget); err != nil {
		logger.Errorf("bind url failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your url"))
		return
	}
	path := fmt.Sprintf("%s%s/%s", config.Cfg.YamlConfig.ServiceInformation.SaveRootPath,
		bucketName.BucketName, fileTarget.Path)
	file, err := os.Open(path)
	if err != nil {
		logger.Errorf("open file failed.Error:%v", err)
		response.ErrorWithMessage(ctx, "open file failed")
		return
	}
	defer file.Close()
	buf := util.AcquireBuf()
	defer util.ReleaseBuf(buf)
	// 让浏览器自动下载
	ctx.Header("Content-Type", "application/octet-stream")
	// 让浏览器自动下载
	ctx.Header("Content-Disposition", "attachment; filename="+path)
	// 让浏览器自动打开
	//ctx.Header("Content-Disposition", "inline; filename="+path)
	if length, err := io.CopyBuffer(ctx.Writer, file, buf); err != nil {
		logger.Errorf("write file %d failed.Error:%v", length, err)
		response.ErrorWithMessage(ctx, "download file failed")
		return
	}
	ctx.Writer.Flush()
	ctx.Status(http.StatusOK)
}
