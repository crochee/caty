// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import (
	"obs/config"
	"testing"
)

func TestBucket_TableName(t *testing.T) {
	config.Cfg.ServiceConfig.List.Mysql = &config.SqlConfig{
		Type:     "mysql",
		User:     "root",
		Password: "123456",
		Host:     "192.168.31.62",
		Port:     3306,
		Database: "obs",
		Charset:  "utf8",
		Debug:    true,
	}
	Setup()
	test := &Bucket{
		Domain: "123",
		Bucket: "bucket",
		User:   "123-1",
		BucketFileList: []BucketFile{
			{
				File: "j.txt",
			},
		},
	}
	tx := NewDB().Begin()
	defer tx.Commit()
	if err := NewDB().Create(test).Error; err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}
}
