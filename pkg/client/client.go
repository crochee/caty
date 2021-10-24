// Date: 2021/10/24

// Package client
package client

import "fmt"

func New(service string) *Service {
	switch service {
	case AccountService:
		return &Service{Account: NewAccount()}
	case AuthService:
		return &Service{Auth: NewAuth()}
	default:
		panic(fmt.Sprintf("you must impl %s", service))
	}
}

type Service struct {
	Account
	Auth
}

const (
	AccountService = "account"
	AuthService    = "auth"
)
