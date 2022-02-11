package code

import "github.com/crochee/lirity/e"

var (
	// 000~099 系统类

	ErrNoAccount = e.Froze(40011000, "用户不存在")
	ErrNoUpdate  = e.Froze(3041101, "数据无更新")

	// 100~199为账号)
	ErrRegisterAccount      = e.Froze(50011100, "注册账号错误")
	ErrUpdateAccount        = e.Froze(50011101, "编辑账号错误")
	ErrRetrieveAccount      = e.Froze(50011102, "查询账号错误")
	ErrDeleteAccount        = e.Froze(50011103, "删除账号错误")
	ErrExistAccount         = e.Froze(40011104, "用户已存在")
	ErrLoginAccount         = e.Froze(50011105, "用户登录错误")
	ErrWrongPasswordAccount = e.Froze(40011106, "用户密码错误")

	// 200~299为权限类

	ErrCreateAuth  = e.Froze(50011200, "生成token")
	ErrParseAuth   = e.Froze(50011201, "解析token错误")
	ErrInvalidAuth = e.Froze(40011202, "无效token")
	ErrExpireAuth  = e.Froze(40011203, "过期token")
	ErrVerifyAuth  = e.Froze(40011204, "错误token")
)

func Loading() error {
	return e.AddCode(map[e.ErrorCode]struct{}{
		ErrNoAccount: {},
		ErrNoUpdate:  {},

		ErrRegisterAccount:      {},
		ErrUpdateAccount:        {},
		ErrRetrieveAccount:      {},
		ErrDeleteAccount:        {},
		ErrExistAccount:         {},
		ErrLoginAccount:         {},
		ErrWrongPasswordAccount: {},

		ErrCreateAuth:  {},
		ErrParseAuth:   {},
		ErrInvalidAuth: {},
		ErrExpireAuth:  {},
		ErrVerifyAuth:  {},
	})
}
