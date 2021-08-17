// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/4

package db

import (
	"time"

	"gorm.io/gorm"
)

type Bucket struct {
	Bucket string `json:"bucket" gorm:"primary_key:bucket;type:varchar(50);not null"`
	Domain string `json:"domain" gorm:"column:domain;type:varchar(50);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Bucket) Delete() {
	NewDB().Unscoped()
	tx := NewDB().Begin()
	defer tx.Rollback()
	if err := tx.Unscoped().Where("`deleted_at` IS NOT NULL").Delete(b).Error; err != nil {
		return
	}
	tx.Commit()
}
