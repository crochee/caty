package ex

import "github.com/crochee/lib/e"

const (
	// 100~200为账号类

	ErrRegisterAccount e.Code = 40510100
	ErrUpdateAccount   e.Code = 40510101
	ErrRetrieveAccount e.Code = 40510102
	ErrDeleteAccount   e.Code = 40510103
)

func Loading() error {
	return e.AddCode(map[e.Code]string{
		ErrRegisterAccount: "注册账号错误",
		ErrUpdateAccount:   "编辑账号错误",
		ErrRetrieveAccount: "查询账号错误",
		ErrDeleteAccount:   "删除账号错误",
	})
}
