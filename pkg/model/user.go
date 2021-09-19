// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import (
	"obs/pkg/db"
	"obs/pkg/log"
)

type User struct {
	ID         uint64 `json:"id" gorm:"primary_key:id"`
	AccountID  string `json:"account_id" gorm:"column:account_id;type:varchar(255);not null;comment:主账号ID"`
	UserID     string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;index:idx_user_id_deleted,unique;comment:用户ID"`
	Nick       string `json:"nick" gorm:"column:nick;type:varchar(50);not null;comment:昵称"`
	PassWord   string `json:"pass_word" gorm:"column:pass_word;type:varchar(50);not null;comment:密码"`
	Email      string `json:"email" gorm:"column:email;type:varchar(50);not null;comment:邮箱"`
	Permission string `json:"permission" gorm:"column:permission;type:text;not null;comment:权限文本"`
	Verify     uint8  `json:"verify" gorm:"column:verify;not null;comment:身份认证"`

	Deleted db.Deleted `json:"deleted" gorm:"column:deleted;not null;index:idx_user_id_deleted,unique;comment:软删除标记"`
	db.Base
}

func DeleteUser() {
	u := new(User)
	if err := db.Client().Model(u).Unscoped().Where("`deleted_at` IS NOT NULL").Delete(u).Error; err != nil {
		log.Warn(err.Error())
	}
}
