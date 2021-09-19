// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/3

package host

import (
	"net"
	"testing"
)

func TestGetIpByName(t *testing.T) {
	t.Log(net.InterfaceByName("eth0"))
	t.Log(GetIPByName("WLAN"))
	t.Log(ExternalIP())
}
