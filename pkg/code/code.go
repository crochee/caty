package code

import "github.com/crochee/lib/e"

const (
	// 000~099 系统类

	ErrNoAccount e.Code = 40010000
	ErrNoUpdate  e.Code = 30410101

	// 100~200为账号类

	ErrRegisterAccount e.Code = 50010100
	ErrUpdateAccount   e.Code = 50010101
	ErrRetrieveAccount e.Code = 50010102
	ErrDeleteAccount   e.Code = 50010103
	ErrExistAccount    e.Code = 40010104
)

func Loading() error {
	return e.AddCode(map[e.Code]string{
		ErrNoAccount:       "用户不存在",
		ErrNoUpdate:        "数据无更新",
		ErrRegisterAccount: "注册账号错误",
		ErrUpdateAccount:   "编辑账号错误",
		ErrRetrieveAccount: "查询账号错误",
		ErrDeleteAccount:   "删除账号错误",
		ErrExistAccount:    "用户已存在",
	})
}
