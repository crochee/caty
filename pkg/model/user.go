// Date: 2021/10/9

// Package model
package model

import "cca/pkg/db"

type User struct {
	ID         uint64 `json:"id" gorm:"primary_key:id"`
	AccountID  string `json:"account" gorm:"column:account;type:varchar(255);not null;comment:账号ID"`
	Name       string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;comment:用户名"`
	Password   string `json:"password" gorm:"column:password;type:varchar(50);not null;comment:密码"`
	Email      string `json:"email" gorm:"column:email;type:varchar(50);not null;comment:邮箱"`
	Permission string `json:"permission" gorm:"column:permission;type:json;not null;comment:权限文本"`
	Verify     uint8  `json:"verify" gorm:"column:verify;not null;comment:身份认证"`

	Desc string `json:"desc" gorm:"column:desc;type:text;not null;comment:详细描述"`

	db.Base
}
