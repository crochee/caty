// Date: 2021/10/24

// Package client
package client

type Auth interface {
	Sign() error
	Parse() error
}
