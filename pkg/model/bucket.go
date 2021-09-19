// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package model

import (
	"time"

	"gorm.io/gorm"

	"obs/pkg/db"
	"obs/pkg/log"
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
	if err := db.Client().Model(b).Unscoped().Where("`deleted_at` IS NOT NULL").Delete(b).Error; err != nil {
		log.Warn(err.Error())
	}
}
