// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package auth

import "fmt"

const (
	Not    uint8 = 0
	Read   uint8 = 1
	Write  uint8 = 2
	Delete uint8 = 3
	Admin  uint8 = 4
)

var ActionString = map[uint8]string{
	Not:    "not",
	Read:   "read",
	Write:  "write",
	Delete: "delete",
	Admin:  "admin",
}

func VerifyAuth(actionMap map[string]uint8, serviceName string, action uint8) error {
	if tempAction, ok := actionMap[AllService]; ok {
		if tempAction >= action {
			return nil
		}
	}
	if actionMap[serviceName] >= action {
		return nil
	}
	return fmt.Errorf("must obtain %s access to the %s", ActionString[action], serviceName)
}
