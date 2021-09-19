// Date: 2021/9/19

// Package e
package e

import (
	"fmt"
	"strings"
)

type Errors []error

func (s Errors) Error() string {
	var errs []string
	for i, e := range s {
		if e == nil {
			continue
		}
		errs = append(errs, fmt.Sprintf("[%d]:%s", i, e.Error()))
	}
	return strings.Join(errs, ";")
}
