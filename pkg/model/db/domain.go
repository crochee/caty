// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package db

import (
	"time"

	"gorm.io/gorm"
)

type Domain struct {
	Domain string `gorm:"primary_key:domain;type:varchar(50);not null"`

	Email    string `gorm:"type:varchar(50);not null;unique_index:email"`
	Nick     string `gorm:"type:varchar(50);not null;column:nick"`
	PassWord string `gorm:"type:varchar(20);not null;column:pass_word"`

	Permission string `gorm:"type:text;not null;column:permission"`

	Verify bool `gorm:"column:verify"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (d *Domain) Delete() {
	tx := NewDB().Begin()
	defer tx.Rollback()
	if err := tx.Unscoped().Where("`deleted_at` IS NOT NULL").Delete(d).Error; err != nil {
		return
	}
	tx.Commit()
}
