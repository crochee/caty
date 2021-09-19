// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import (
	"obs/pkg/model/db"
	"time"

	"gorm.io/gorm"
)

type User struct {
	User string `gorm:"primary_key:user;type:varchar(50);not null"`

	Domain string `gorm:"column:domain;type:varchar(50);not null"`

	Nick     string `gorm:"type:varchar(50);not null;column:nick"`
	PassWord string `gorm:"type:varchar(20);not null;column:pass_word"`

	Permission string `gorm:"type:text;not null;column:permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) Delete() {
	tx := db.NewDB().Begin()
	defer tx.Rollback()
	tx.Callback().Delete()
	if err := tx.Unscoped().Where("`deleted_at` IS NOT NULL").Delete(u).Error; err != nil {
		return
	}
	tx.Commit()
}
