package code

import "github.com/crochee/lirity/e"

var (
	// 000~099 系统类

	ErrNoAccount = e.Froze("4001100000", "用户不存在")
	ErrNoUpdate  = e.Froze("3041100001", "数据无更新")

	// 100~199为账号)
	ErrRegisterAccount      = e.Froze("5001100100", "注册账号错误")
	ErrUpdateAccount        = e.Froze("5001100101", "编辑账号错误")
	ErrRetrieveAccount      = e.Froze("5001100102", "查询账号错误")
	ErrDeleteAccount        = e.Froze("5001100103", "删除账号错误")
	ErrExistAccount         = e.Froze("4001100104", "用户已存在")
	ErrLoginAccount         = e.Froze("5001100105", "用户登录错误")
	ErrWrongPasswordAccount = e.Froze("4001100106", "用户密码错误")

	// 200~299为权限类

	ErrCreateAuth  = e.Froze("5001100200", "生成token")
	ErrParseAuth   = e.Froze("5001100201", "解析token错误")
	ErrInvalidAuth = e.Froze("4001100202", "无效token")
	ErrExpireAuth  = e.Froze("4001100203", "过期token")
	ErrVerifyAuth  = e.Froze("4001100204", "错误token")
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
