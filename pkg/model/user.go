// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package model

import (
	db2 "cca/pkg/db"
	"github.com/crochee/lib/db"
	"github.com/crochee/lib/log"
)

type Account struct {
	ID         uint64 `json:"id" gorm:"primary_key:id"`
	AccountID  string `json:"account_id" gorm:"column:account_id;type:varchar(255);not null;comment:主账号ID"`
	Account    string `json:"account" gorm:"column:account;type:varchar(255);not null;comment:账号名"`
	UserID     string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;index:idx_user_id_deleted,unique;comment:用户ID"`
	Password   string `json:"password" gorm:"column:password;type:varchar(50);not null;comment:密码"`
	Email      string `json:"email" gorm:"column:email;type:varchar(50);not null;comment:邮箱"`
	Permission string `json:"permission" gorm:"column:permission;type:json;not null;comment:权限文本"`
	Verify     uint8  `json:"verify" gorm:"column:verify;not null;comment:身份认证"`
	Desc       string `json:"desc" gorm:"column:desc;type:text;not null;comment:详细描述"`

	Deleted db.Deleted `json:"deleted" gorm:"column:deleted;not null;index:idx_user_id_deleted,unique;comment:软删除标记"`
	db.Base
}

func DeleteUser() {
	u := new(Account)
	if err := db2.New().Model(u).Unscoped().Where("`deleted_at` IS NOT NULL").Delete(u).Error; err != nil {
		log.Warn(err.Error())
	}
}
