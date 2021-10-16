// Date: 2021/10/12

// Package account
package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/crochee/lib"
	"github.com/crochee/lib/db"
	"github.com/crochee/lib/e"
	"gorm.io/gorm"

	"cca/pkg/code"
	"cca/pkg/dbx"
	"cca/pkg/model"
	"cca/pkg/service/auth"
)

type CreateRequest struct {
	// 用户名
	// Required: true
	Account string `json:"account" binding:"required"`
	// 账户ID
	AccountID string `json:"account_id" binding:"omitempty,numeric"`
	// 邮箱
	Email string `json:"email" binding:"omitempty,email"`
	// 密码
	// Required: true
	Password string `json:"password" binding:"required,alphanum"`
	// 描述信息
	Desc string `json:"desc" binding:"required,json"`
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
	actionMap := map[string]uint8{
		auth.AllService: auth.Admin,
	}
	if request.AccountID != "" {
		actionMap[auth.AllService] = auth.Read
	}
	permission, err := json.Marshal(actionMap)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed.Error:%v,%w", err, e.ErrInternalServerError)
	}
	userModel := &model.User{
		Name:       request.Account,
		Password:   request.Password,
		Email:      request.Email,
		Permission: lib.String(permission),
		Desc:       request.Desc,
	}
	err = dbx.With(ctx).Transaction(func(tx *gorm.DB) error {
		accountModel := &model.Account{}
		if request.AccountID != "" {
			if err = tx.Model(accountModel).Where("id =?", request.AccountID).
				First(accountModel).Error; err != nil {
				if errors.Is(err, db.NotFound) {
					return code.ErrNoAccount
				}
				return fmt.Errorf("first account failed;%v;%w", err, code.ErrRegisterAccount)
			}
		} else {
			accountModel.Name = request.Account
			if err = tx.Model(accountModel).Create(accountModel).Error; err != nil {
				if strings.Contains(err.Error(), db.ErrDuplicate) {
					return code.ErrExistAccount
				}
				return fmt.Errorf("insert account failed;%v;%w", err, code.ErrRegisterAccount)
			}
			userModel.PrimaryAccount = true
		}
		userModel.AccountID = accountModel.ID
		if err = tx.Model(userModel).Create(userModel).Error; err != nil {
			if strings.Contains(err.Error(), db.ErrDuplicate) {
				return code.ErrExistAccount
			}
			return fmt.Errorf("insert user failed;%v;%w", err, code.ErrRegisterAccount)
		}
		if err = tx.Model(userModel).First(userModel).Error; err != nil {
			return fmt.Errorf("first user failed;%v;%w", err, code.ErrRegisterAccount)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &CreateResponseResult{
		AccountID:      FormatUint(userModel.AccountID),
		Account:        userModel.Name,
		UserID:         FormatUint(userModel.ID),
		Email:          userModel.Email,
		Permission:     userModel.Permission,
		Verify:         userModel.Verify,
		PrimaryAccount: userModel.PrimaryAccount,
		Desc:           userModel.Desc,
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
	}, nil
}

type User struct {
	// 用户
	// Required: true
	// in: path
	ID string `json:"id" uri:"id" binding:"required,numeric"`
}

type UpdateRequest struct {
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
	Desc string `json:"desc" binding:"omitempty,json"`
}

// Update 编辑账户
func Update(ctx context.Context, user *User, request *UpdateRequest) error {
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
		return code.ErrNoUpdate
	}
	query := dbx.With(ctx).Model(&model.User{}).Where("id=? AND password=?",
		user.ID, request.OldPassword).Updates(updates)
	if err := query.Error; err != nil {
		return fmt.Errorf("update failed.Error:%v,%w", err, code.ErrUpdateAccount)
	}
	if query.RowsAffected == 0 {
		return code.ErrNoUpdate
	}
	return nil
}

type RetrievesRequest struct {
	// 账户ID
	// in: query
	AccountID string `json:"account-id" form:"account-id" binding:"omitempty,numeric"`
	// 用户
	// in: query
	ID string `json:"id" form:"id" binding:"omitempty,numeric"`
	// 账户
	// in: query
	Account string `json:"account" form:"account" binding:"omitempty"`
	// 邮箱
	// in: query
	Email string `json:"email" form:"email" binding:"omitempty,email"`
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

// Retrieves 查询、获取账户信息
func Retrieves(ctx context.Context, request *RetrievesRequest) (*RetrieveResponses, error) {
	query := dbx.With(ctx).Model(&model.User{})
	if request.AccountID != "" {
		query = query.Where("account_id = ?", request.AccountID)
	}
	if request.ID != "" {
		query = query.Where("id = ?", request.ID)
	}
	if request.Account != "" {
		query = query.Where("name = ?", request.Account)
	}
	if request.Email != "" {
		query = query.Where("email = ?", request.Email)
	}
	var userList []*model.User
	if err := query.Find(&userList).Error; err != nil {
		return nil, fmt.Errorf("find user failed.Error:%v,%w", err, code.ErrRetrieveAccount)
	}
	responses := &RetrieveResponses{
		Result: make([]*RetrieveResponse, 0, len(userList)),
	}
	for _, v := range userList {
		responses.Result = append(responses.Result, &RetrieveResponse{
			AccountID:  FormatUint(v.AccountID),
			Account:    v.Name,
			UserID:     FormatUint(v.ID),
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

// Retrieve 查询、获取指定账户信息
func Retrieve(ctx context.Context, request *User) (*RetrieveResponse, error) {
	user := &model.User{}
	if err := dbx.With(ctx).Model(user).Where("id =?", request.ID).First(user).Error; err != nil {
		if errors.Is(err, db.NotFound) {
			return nil, code.ErrNoAccount
		}
		return nil, fmt.Errorf("%v.%w", err, code.ErrRetrieveAccount)
	}
	return &RetrieveResponse{
		AccountID:  FormatUint(user.AccountID),
		Account:    user.Name,
		UserID:     FormatUint(user.ID),
		Email:      user.Email,
		Permission: user.Permission,
		Verify:     user.Verify,
		Desc:       user.Desc,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

// Delete 删除账户
func Delete(ctx context.Context, request *User) error {
	return dbx.With(ctx).Transaction(func(tx *gorm.DB) error {
		user := &model.User{}
		query := tx.Model(user).Where("id =?", request.ID)
		if err := query.First(user).Error; err != nil {
			if errors.Is(err, db.NotFound) {
				return code.ErrNoAccount
			}
			return fmt.Errorf("first user failed;%v;%w", err, code.ErrDeleteAccount)
		}
		if user.PrimaryAccount {
			accountModel := &model.Account{}
			queryAccountDel := tx.Model(accountModel).Where("id =?", user.AccountID).Delete(accountModel)
			if err := queryAccountDel.Error; err != nil {
				return fmt.Errorf("delete account failed;%v;%w", err, code.ErrDeleteAccount)
			}
			if queryAccountDel.RowsAffected == 0 {
				return code.ErrNoAccount
			}
		}
		queryDel := query.Delete(user)
		if err := queryDel.Error; err != nil {
			return fmt.Errorf("%v.%w", err, code.ErrDeleteAccount)
		}
		if queryDel.RowsAffected == 0 {
			return code.ErrNoAccount
		}
		return nil
	})
}

const PasswordMaxLength = 15

// ValidPassword 密码校验
func ValidPassword(password string) error {
	if len(password) < PasswordMaxLength {
		return fmt.Errorf("password's length is less than %d", PasswordMaxLength)
	}
	return nil
}

// ValidPermission 权限格式校验
func ValidPermission(permission string) error {
	data := make(map[string]uint8)
	return json.Unmarshal([]byte(permission), &data)
}

func FormatUint(data uint64) string {
	return strconv.FormatUint(data, 10)
}
