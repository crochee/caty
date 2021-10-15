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

	"cca/pkg/code"
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
	// 是否主账号
	PrimaryAccount bool `json:"primary_account"`
	// 是否认证
	Verify uint8 `json:"verify"`
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
		return nil, fmt.Errorf("register do transaction failed.Error:%v,%w", err, code.ErrRegisterAccount)
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

type UpdateRequest struct {
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

// Update 编辑账户
func Update(ctx context.Context, request *UpdateRequest) error {
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
		return fmt.Errorf("update failed.Error:%v,%w", err, code.ErrUpdateAccount)
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("update 0 rows affected,%w", code.ErrUpdateAccount)
	}
	return nil
}

type RetrieveRequest struct {
	// 账户ID
	// in: query
	// Required: true
	AccountID string `form:"account-id" binding:"omitempty,numeric"`
	// 用户
	// in: query
	// Required: true
	UserID string `form:"id" binding:"omitempty,numeric"`
	// 账户
	// in: query
	Account string `form:"account" binding:"omitempty"`
	// 邮箱
	// in: query
	Email string `form:"email" binding:"omitempty,email"`
}

type RetrieveResponses struct {
	// 结果集
	Result []*RetrieveResponse `json:"result"`
}

type RetrieveResponse struct {
	// 是否认证
	Verify uint8 `json:"verify"`
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
	// 描述
	Desc string `json:"desc"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// Retrieve 查询、获取账户信息
func Retrieve(ctx context.Context, request *RetrieveRequest) (*RetrieveResponses, error) {
	queryList := make(map[string]string)
	if request.AccountID != "" {
		queryList["account_id =?"] = request.AccountID
	}
	if request.UserID != "" {
		queryList["id"] = request.UserID
	}
	if request.Account != "" {
		queryList["name"] = request.Account
	}
	if request.Email != "" {
		queryList["email"] = request.Email
	}
	if len(queryList) == 0 {
		return nil, fmt.Errorf("retrieve has 0 conditions,%w", e.ErrInvalidParam)
	}
	query := db.With(ctx).Model(&model.User{})
	for k, v := range queryList {
		query = query.Where("? = ?", k, v)
	}
	var userList []*model.User
	if err := query.Find(userList); err != nil {
		return nil, fmt.Errorf("find user failed.Error:%v,%w", err, code.ErrRetrieveAccount)
	}
	responses := &RetrieveResponses{
		Result: make([]*RetrieveResponse, 0, len(userList)),
	}
	for _, v := range userList {
		responses.Result = append(responses.Result, &RetrieveResponse{
			AccountID:  strconv.FormatUint(v.AccountID, 10),
			Account:    v.Name,
			UserID:     strconv.FormatUint(v.ID, 10),
			Email:      v.Email,
			Permission: v.Permission,
			Verify:     v.Verify,
			Desc:       v.Desc,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		})
	}
	return responses, nil
}

type RetrieveSingleRequest struct {
	// 用户ID
	// in: path
	// Required: true
	UserID string `json:"user_id" uri:"id" binding:"required,numeric"`
}

// RetrieveSingle 查询、获取指定账户信息
func RetrieveSingle(ctx context.Context, request *RetrieveSingleRequest) (*RetrieveResponse, error) {
	user := &model.User{}
	if err := db.With(ctx).Model(user).Where("id =?", request.UserID).First(user).Error; err != nil {
		return nil, fmt.Errorf("%v.%w", err, code.ErrRetrieveAccount)
	}
	return &RetrieveResponse{
		AccountID:  strconv.FormatUint(user.AccountID, 10),
		Account:    user.Name,
		UserID:     strconv.FormatUint(user.ID, 10),
		Email:      user.Email,
		Permission: user.Permission,
		Verify:     user.Verify,
		Desc:       user.Desc,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

type DeleteRequest struct {
	// 用户ID
	// in: path
	// Required: true
	UserID string `json:"user_id" uri:"id" binding:"required,numeric"`
}

// Delete 删除账户
func Delete(ctx context.Context, request *DeleteRequest) error {
	user := &model.User{}
	if err := db.With(ctx).Model(user).Where("id =?", request.UserID).Delete(user).Error; err != nil {
		return fmt.Errorf("%v.%w", err, code.ErrDeleteAccount)
	}
	return nil
}
