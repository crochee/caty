// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"obs/e"
	"os"
	"path/filepath"
	"time"

	"obs/config"
	"obs/logger"
	"obs/model/db"
	"obs/response"
	"obs/service/tokenx"
)

// UploadFile 上传文件
func UploadFile(ctx context.Context, token *tokenx.Token, bucketName string, file *multipart.FileHeader) error {
	tx := db.NewDB().Begin()
	defer tx.Commit()
	bucket := &db.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?",
		bucketName, token.Domain).Find(bucket).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("insert db failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	path := filepath.Clean(fmt.Sprintf("%s/%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath,
		bucket.Bucket, file.Filename))
	dstFile, err := os.Create(path)
	if err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("create file %s failed.Error:%v", path, err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	defer dstFile.Close() // #nosec G307
	var srcFile multipart.File
	if srcFile, err = file.Open(); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("open %s file failed.Error:%v", file.Filename, err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	defer srcFile.Close() // #nosec G307

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("copy %s file failed.Error:%v", file.Filename, err)
		return e.Error(http.StatusInternalServerError, "copy failed")
	}
	var stat os.FileInfo
	if stat, err = dstFile.Stat(); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("get %s file stat failed.Error:%v", dstFile.Name(), err)
		return e.Error(http.StatusInternalServerError, "get file stat failed")
	}
	bucketFile := &db.BucketFile{
		File:    stat.Name(),
		Bucket:  bucket.Bucket,
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}
	if err = tx.Create(bucketFile).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("insert file failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	return nil
}

// DeleteFile 删除文件
func DeleteFile(ctx context.Context, token *tokenx.Token, bucketName, fileName string) error {
	tx := db.NewDB().Begin()
	defer tx.Commit()
	bucket := &db.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?", bucketName, token.Domain).Find(bucket).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	bucketFile := &db.BucketFile{}
	if err := tx.Model(bucketFile).Where("file =? AND bucket= ?", fileName, bucket.Bucket).
		Find(bucketFile).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	path := filepath.Clean(fmt.Sprintf("%s/%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath,
		bucket.Bucket, bucketFile.File))
	if err := os.Remove(path); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("remove file %s failed.Error:%v", path, err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	if err := tx.Delete(bucketFile).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("remove file %s failed.Error:%v", path, err)
		return response.Errors(http.StatusInternalServerError, err)
	}
	return nil
}

// SignFile 文件签名
func SignFile(ctx context.Context, token *tokenx.Token, bucketName, fileName string) (string, error) {
	tx := db.NewDB().Begin()
	defer tx.Commit()
	bucket := &db.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?",
		bucketName, token.Domain).Find(bucket).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return "", response.Errors(http.StatusInternalServerError, err)
	}
	bucketFile := &db.BucketFile{}
	if err := tx.Model(bucketFile).Where("file =? AND bucket= ?",
		fileName, bucket.Bucket).Find(bucketFile).Error; err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return "", response.Errors(http.StatusInternalServerError, err)
	}
	xToken := &tokenx.TokenClaims{
		Now: time.Now(),
		Token: &tokenx.Token{
			Domain: token.Domain,
			User:   token.User,
			ActionMap: map[string]tokenx.Action{
				"OBS": tokenx.Read,
			},
		},
	}
	signString, err := tokenx.CreateToken(xToken)
	if err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("create token failed.Error:%v", err)
		return "", response.Errors(http.StatusInternalServerError, err)
	}
	var (
		sign      string
		tokenSign = tokenx.Signature(signString)
	)
	if sign, err = tokenx.CreateSign(&tokenSign); err != nil {
		tx.Rollback()
		logger.FromContext(ctx).Errorf("create sian failed.Error:%v", err)
		return "", response.Errors(http.StatusInternalServerError, err)
	}
	return sign, nil
}
