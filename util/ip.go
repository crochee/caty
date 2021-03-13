// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/8

package util

import (
	"errors"
	"net"
)

// ExternalIP export native ip net.IP
func ExternalIP() (net.IP, error) {
	iFaceList, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iFace := range iFaceList {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrList, err := iFace.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrList {
			ip := getIPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("can't connected to the network")
}

func getIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	default:
		return nil
	}
	if ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}
