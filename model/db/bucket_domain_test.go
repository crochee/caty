// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import (
	"testing"

	"obs/config"
)

func TestBucket_TableName(t *testing.T) {
	config.Cfg = config.Config{
		ServiceConfig: &config.ServiceConfig{
			List: config.Connection{
				Mysql: &config.SqlConfig{
					Type:     "mysql",
					User:     "root",
					Password: "1234567",
					Host:     "192.168.31.62",
					Port:     3307,
					Database: "obs",
					Charset:  "utf8",
					Debug:    true,
				},
				Mongo: nil,
			},
		},
	}
	Setup()
	test := &Bucket{
		Domain: "123",
		Bucket: "bucket",
		User:   "123-1",
	}
	tx := NewDB().Begin()
	defer tx.Commit()
	// 级联插入
	if err := tx.Create(test).Error; err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
}
