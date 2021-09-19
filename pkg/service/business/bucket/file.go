// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/13

package bucket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"obs/pkg/log"
	"obs/pkg/model"
	db2 "obs/pkg/model/db"
	"obs/pkg/service/business/tokenx"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"obs/config"
	"obs/e"
)

// UploadFile 上传文件
func UploadFile(ctx context.Context, token *tokenx.Token, bucketName string, file *multipart.FileHeader) error {
	tx := db2.NewDBWithContext(ctx).Begin()
	defer tx.Rollback()
	bucket := &model.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?",
		bucketName, token.Domain).First(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.New(e.NotFound, "not found record")
		}
		log.FromContext(ctx).Errorf("insert db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	path := filepath.Clean(fmt.Sprintf("%s/%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath,
		bucket.Bucket, file.Filename))
	dstFile, err := os.Create(path)
	if err != nil {
		log.FromContext(ctx).Errorf("create file %s failed.Error:%v", path, err)
		return e.New(e.GetAbsPathFail, "clear bucket failed")
	}
	defer dstFile.Close() // #nosec G307
	var srcFile multipart.File
	if srcFile, err = file.Open(); err != nil {
		log.FromContext(ctx).Errorf("open %s file failed.Error:%v", file.Filename, err)
		return e.New(e.OpenFileFail, err.Error())
	}
	defer srcFile.Close() // #nosec G307

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		log.FromContext(ctx).Errorf("copy %s file failed.Error:%v", file.Filename, err)
		return e.New(e.OpenFileFail, err.Error())
	}
	var stat os.FileInfo
	if stat, err = dstFile.Stat(); err != nil {
		log.FromContext(ctx).Errorf("get %s file stat failed.Error:%v", dstFile.Name(), err)
		return e.New(e.OpenFileFail, err.Error())
	}
	bucketFile := &model.BucketFile{
		File:    stat.Name(),
		Bucket:  bucket.Bucket,
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}
	if err = tx.Create(bucketFile).Error; err != nil {
		log.FromContext(ctx).Errorf("insert file failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	tx.Commit()
	return nil
}

// DeleteFile 删除文件
func DeleteFile(ctx context.Context, token *tokenx.Token, bucketName, fileName string) error {
	tx := db2.NewDBWithContext(ctx).Begin()
	defer tx.Rollback()
	bucket := &model.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?", bucketName, token.Domain).
		First(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.New(e.NotFound, "not found record")
		}
		log.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	bucketFile := &model.BucketFile{}
	if err := tx.Model(bucketFile).Where("file =? AND bucket= ?", fileName, bucket.Bucket).
		First(bucketFile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.New(e.NotFound, "not found record")
		}
		log.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return e.New(e.OperateDbFail, err.Error())
	}
	path := filepath.Clean(fmt.Sprintf("%s/%s/%s", config.Cfg.ServiceConfig.ServiceInfo.StoragePath,
		bucket.Bucket, bucketFile.File))
	if err := os.Remove(path); err != nil {
		log.FromContext(ctx).Errorf("remove file %s failed.Error:%v", path, err)
		return e.New(e.DeleteFileFail, err.Error())
	}
	if err := tx.Delete(bucketFile).Error; err != nil {
		log.FromContext(ctx).Errorf("remove file %s failed.Error:%v", path, err)
		return e.New(e.OperateDbFail, err.Error())
	}
	tx.Commit()
	return nil
}

// SignFile 文件签名
func SignFile(ctx context.Context, token *tokenx.Token, bucketName, fileName string) (string, error) {
	tx := db2.NewDBWithContext(ctx).Begin()
	defer tx.Rollback()
	bucket := &model.Bucket{}
	if err := tx.Model(bucket).Where("bucket =? AND domain= ?",
		bucketName, token.Domain).First(bucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", e.New(e.NotFound, "not found record")
		}
		log.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return "", e.New(e.OperateDbFail, err.Error())
	}
	bucketFile := &model.BucketFile{}
	if err := tx.Model(bucketFile).Where("file =? AND bucket= ?",
		fileName, bucket.Bucket).First(bucketFile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", e.New(e.NotFound, "not found record")
		}
		log.FromContext(ctx).Errorf("query db failed.Error:%v", err)
		return "", e.New(e.OperateDbFail, err.Error())
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
		log.FromContext(ctx).Errorf("create token failed.Error:%v", err)
		return "", e.New(e.GenerateTokenFail, err.Error())
	}
	var (
		sign      string
		tokenSign = &tokenx.Signature{signString}
	)
	if sign, err = tokenx.CreateSign(tokenSign); err != nil {
		log.FromContext(ctx).Errorf("create sian failed.Error:%v", err)
		return "", e.New(e.GenerateSignFail, err.Error())
	}
	tx.Commit()
	return sign, nil
}
