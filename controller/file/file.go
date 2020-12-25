// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/23

package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"obs/config"
	"obs/logger"
	"obs/model"
	"obs/response"
	"obs/util"
)

// UploadFile godoc
// @Summary UploadFile
// @Description upload file
// @Tags file
// @Accept multipart/form-data
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Param request formData model.FileInfo true "file"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file/bucket/{bucket_name} [post]
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
	ctx.Status(http.StatusOK)
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
// @Router /v1/file/bucket/{bucket_name} [delete]
func DeleteFile(ctx *gin.Context) {
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
	if err := os.RemoveAll(path); err != nil {
		logger.Errorf("delete path(%s) failed.Error:%v", path, err)
		response.ErrorWithMessage(ctx, "delete target failed")
		return
	}
	ctx.Status(http.StatusOK)
}

// HeadFile godoc
// @Summary HeadFile
// @Description head file sign info
// @Tags file
// @Accept application/json
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Param path query string true "target path"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file/bucket/{bucket_name} [head]
func SignFile(ctx *gin.Context) {

}

// DownloadFile godoc
// @Summary DownloadFile
// @Description head file sign info
// @Tags file
// @Accept application/json
// @Produce  application/json
// @Param bucket_name path string true "bucket name"
// @Param path query string true "target path"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file/bucket/{bucket_name} [get]
func DownloadFile(ctx *gin.Context) {

}
