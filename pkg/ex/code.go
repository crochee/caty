package ex

import "github.com/crochee/lib/e"

const (
	// 100~200为用户验证类

	ErrInvalidEmail e.Code = 40510100
)

func AddCode() error {
	return e.AddCode(map[e.Code]string{
		ErrInvalidEmail: "非法邮箱格式",
	})
}
