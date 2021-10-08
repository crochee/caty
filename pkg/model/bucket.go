// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package model

import (
	"time"

	"gorm.io/gorm"

	"cca/pkg/db"
	"cca/pkg/logx"
)

type Bucket struct {
	Bucket string `json:"bucket" gorm:"primary_key:bucket;type:varchar(50);not null"`
	Domain string `json:"domain" gorm:"column:domain;type:varchar(50);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func DeleteBucket() {
	b := new(Bucket)
	if err := db.New().Model(b).Unscoped().Where("`deleted_at` IS NOT NULL").Delete(b).Error; err != nil {
		logx.Warn(err.Error())
	}
}
