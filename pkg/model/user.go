// Date: 2021/10/9

// Package model
package model

import "github.com/crochee/lirity/db"

type User struct {
	ID             uint64 `json:"id,string" gorm:"primary_key:id"`
	AccountID      uint64 `json:"account_id" gorm:"column:account_id;not null;index:idx_account_id_name_primary_deleted,unique;comment:账号ID"`
	Name           string `json:"name" gorm:"column:name;type:varchar(255);not null;index:idx_account_id_name_primary_deleted,unique;comment:用户名"`
	Password       string `json:"password" gorm:"column:password;type:varchar(50);not null;comment:密码"`
	Email          string `json:"email" gorm:"column:email;type:varchar(50);not null;comment:邮箱"`
	Permission     string `json:"permission" gorm:"column:permission;type:json;not null;comment:权限文本"`
	Verify         uint8  `json:"verify" gorm:"column:verify;not null;comment:身份认证"`
	PrimaryAccount bool   `json:"primary_account" gorm:"column:primary_account;not null;index:idx_account_id_name_primary_deleted,unique,comment:是否主账号"`

	Desc string `json:"desc" gorm:"column:desc;type:json;not null;comment:详细描述"`

	Deleted db.Deleted `json:"deleted" gorm:"not null;index:idx_account_id_name_primary_deleted,unique;comment:软删除记录id"`
	db.Base
}

func (User) TableName() string {
	return "user"
}
