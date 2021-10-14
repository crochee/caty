// Date: 2021/10/12

// Package account
package account

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/crochee/lib"
	"github.com/crochee/lib/e"
	"gorm.io/gorm"

	"cca/pkg/db"
	"cca/pkg/model"
	"cca/pkg/service/business/tokenx"
)

type CreateRequest struct {
	// 用户名
	// Required: true
	Name string `json:"account" binding:"required"`
	// 账户ID
	AccountID string `json:"account_id" binding:"omitempty,numeric"`
	// 邮箱
	Email string `json:"email" binding:"omitempty,email"`
	// 密码
	// Required: true
	Password string `json:"password" binding:"required,alphanum"`
	// 描述信息
	Desc string `json:"desc" binding:"json"`
}

type CreateResponseResult struct {
	// 账户ID
	AccountID string `json:"account_id"`
	// 账户
	Account string `json:"account"`
	// 用户
	UserID string `json:"user_id"`
	// 邮箱
	Email string `json:"email"`
	// 权限
	Permission string `json:"permission"`
	// 是否认证
	Verify uint8 `json:"verify"`
	// 是否主账号
	PrimaryAccount bool `json:"primary_account"`
	// 描述
	Desc string `json:"desc"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// Create 注册账户
func Create(ctx context.Context, request *CreateRequest) (*CreateResponseResult, error) {
	actionMap := map[string]tokenx.Action{
		tokenx.AllService: tokenx.Admin,
	}
	if request.AccountID != "" {
		actionMap[tokenx.AllService] = tokenx.Read
	}
	permission, err := json.Marshal(actionMap)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed.Error:%v,%w", err, e.ErrInternalServerError)
	}
	userModel := &model.User{
		Name:       request.Name,
		Password:   request.Password,
		Email:      request.Email,
		Permission: lib.String(permission),
		Desc:       request.Desc,
	}
	err = db.With(ctx).Transaction(func(tx *gorm.DB) error {
		accountModel := &model.Account{}
		if request.AccountID != "" {
			err = tx.Model(accountModel).Where("id =?", request.AccountID).First(accountModel).Error
			if err != nil {
				return err
			}
		} else {
			userModel.PrimaryAccount = true
			err = tx.Model(accountModel).Create(accountModel).Error
			if err != nil {
				return err
			}
		}
		userModel.AccountID = accountModel.ID
		return tx.Model(userModel).Create(userModel).Error
	})
	if err != nil {
		return nil, fmt.Errorf("register do transaction failed.Error:%v,%w", err, e.ErrOperateDB)
	}
	return &CreateResponseResult{
		AccountID:      strconv.FormatUint(userModel.AccountID, 10),
		Account:        userModel.Name,
		UserID:         strconv.FormatUint(userModel.ID, 10),
		Email:          userModel.Email,
		Permission:     userModel.Permission,
		Verify:         userModel.Verify,
		PrimaryAccount: userModel.PrimaryAccount,
		Desc:           userModel.Desc,
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
	}, nil
}

type ModifyRequest struct {
	// 账户ID
	// Required: true
	AccountID string `json:"account_id" binding:"required,numeric"`
	// 用户
	// Required: true
	UserID string `json:"user_id" binding:"required,numeric"`
	// 旧密码
	// Required: true
	OldPassword string `json:"old_password" binding:"required,alphanum"`
	// 账户
	Account string `json:"account" binding:"omitempty"`
	// 邮箱
	Email string `json:"email" binding:"omitempty,email"`
	// 新密码
	Password string `json:"password" binding:"omitempty,alphanum"`
	// 权限
	Permission string `json:"permission" binding:"omitempty,json"`
	// 描述信息
	Desc string `json:"desc"`
}

// Modify 编辑账户
func Modify(ctx context.Context, request *ModifyRequest) error {
	updates := make(map[string]interface{})
	if request.Account != "" {
		updates["name"] = request.Account
	}
	if request.Email != "" {
		updates["email"] = request.Email
	}
	if request.Password != "" {
		updates["password"] = request.Password
	}
	if request.Permission != "" {
		updates["permission"] = request.Permission
	}
	if request.Desc != "" {
		updates["desc"] = request.Desc
	}
	if len(updates) == 0 {
		return nil
	}
	query := db.With(ctx).Model(&model.User{}).Where("id=? AND account_id=? AND password=?",
		request.UserID, request.AccountID, request.OldPassword).Updates(updates)
	if err := query.Error; err != nil {
		return fmt.Errorf("update failed.Error:%v,%w", err, e.ErrOperateDB)
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("update 0 rows affected,%w", e.ErrOperateDB)
	}
	return nil
}

// Retrieve 查询、获取账户信息
func Retrieve(ctx context.Context, request *ModifyRequest) error {
	updates := make(map[string]interface{})
	if request.Account != "" {
		updates["name"] = request.Account
	}
	if request.Email != "" {
		updates["email"] = request.Email
	}
	if request.Password != "" {
		updates["password"] = request.Password
	}
	if request.Permission != "" {
		updates["permission"] = request.Permission
	}
	if request.Desc != "" {
		updates["desc"] = request.Desc
	}
	if len(updates) == 0 {
		return nil
	}
	query := db.With(ctx).Model(&model.User{}).Where("id=? AND account_id=? AND password=?",
		request.UserID, request.AccountID, request.OldPassword).Updates(updates)
	if err := query.Error; err != nil {
		return fmt.Errorf("update failed.Error:%v,%w", err, e.ErrOperateDB)
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("update 0 rows affected,%w", e.ErrOperateDB)
	}
	return nil
}
