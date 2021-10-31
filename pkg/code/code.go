package code

import "github.com/crochee/lirity/e"

var (
	// 000~099 系统类

	ErrNoAccount e.ErrorCode = &e.ErrCode{Err: 40011000, Msg: "用户不存在"}
	ErrNoUpdate  e.ErrorCode = &e.ErrCode{Err: 3041101, Msg: "数据无更新"}

	// 100~199为账号类

	ErrRegisterAccount      e.ErrorCode = &e.ErrCode{Err: 50011100, Msg: "注册账号错误"}
	ErrUpdateAccount        e.ErrorCode = &e.ErrCode{Err: 50011101, Msg: "编辑账号错误"}
	ErrRetrieveAccount      e.ErrorCode = &e.ErrCode{Err: 50011102, Msg: "查询账号错误"}
	ErrDeleteAccount        e.ErrorCode = &e.ErrCode{Err: 50011103, Msg: "删除账号错误"}
	ErrExistAccount         e.ErrorCode = &e.ErrCode{Err: 40011104, Msg: "用户已存在"}
	ErrLoginAccount         e.ErrorCode = &e.ErrCode{Err: 50011105, Msg: "用户登录错误"}
	ErrWrongPasswordAccount e.ErrorCode = &e.ErrCode{Err: 40011106, Msg: "用户密码错误"}

	// 200~299为权限类

	ErrCreateAuth  e.ErrorCode = &e.ErrCode{Err: 50011200, Msg: "生成token"}
	ErrParseAuth   e.ErrorCode = &e.ErrCode{Err: 50011201, Msg: "解析token错误"}
	ErrInvalidAuth e.ErrorCode = &e.ErrCode{Err: 40011202, Msg: "无效token"}
	ErrExpireAuth  e.ErrorCode = &e.ErrCode{Err: 40011203, Msg: "过期token"}
	ErrVerifyAuth  e.ErrorCode = &e.ErrCode{Err: 40011204, Msg: "错误token"}
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
