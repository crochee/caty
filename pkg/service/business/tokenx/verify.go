// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/14

package tokenx

import "fmt"

func VerifyAuth(actionMap map[string]Action, serviceName string, action Action) error {
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
