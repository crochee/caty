// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/20

// +build !darwin,!linux,!freebsd,!windows

package uuid

import "errors"

func readPlatformMachineID() (string, error) {
	return "", errors.New("not implemented")
}
