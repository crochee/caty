// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/23

package bucket

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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
// @Param request body model.FileInfo true "file"
// @Success 200
// @Failure 403
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/file [post]
func UploadFile(ctx *gin.Context) {
	var fileInfo model.FileInfo
	if err := ctx.ShouldBindWith(&fileInfo, binding.FormMultipart); err != nil {
		logger.Errorf("bind request failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusBadRequest, "check your payload"))
	}
	t := &verify.Token{
		AkSecret: util.Slice(fileInfo.Ak),
		Bucket:   fileInfo.BucketName.BucketName,
	}
	if err := t.Verify(fileInfo.Sk); err != nil {
		logger.Errorf("verify sk failed.Error:%v", err)
		response.ErrorWith(ctx, response.Error(http.StatusForbidden, "check your ak sk"))
	}
	if !verify.VerifyAuthentication(t, verify.Write) {
		ctx.Status(http.StatusForbidden)
		return
	}
	path := fmt.Sprintf("%s%s/%s", config.Cfg.YamlConfig.ServiceInformation.SaveRootPath,
		fileInfo.BucketName.BucketName, fileInfo.Target)
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
