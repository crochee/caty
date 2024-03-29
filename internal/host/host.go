package host

import (
	"net"
)

var (
	privateAddress = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10", "fd00::/8"}
)

func isPrivateIP(addr string) bool {
	ipAddr := net.ParseIP(addr)
	for _, privateAddr := range privateAddress {
		if _, privy, err := net.ParseCIDR(privateAddr); err == nil {
			if privy.Contains(ipAddr) {
				return true
			}
		}
	}
	return false
}

// Extract returns a private addr and port.
func Extract(hostPort string) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", err
	}
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return net.JoinHostPort(addr, port), nil
	}
	var iFaceList []net.Interface
	if iFaceList, err = net.Interfaces(); err != nil {
		return "", err
	}
	for _, iFace := range iFaceList {
		addrList, err := iFace.Addrs()
		if err != nil {
			continue
		}
		for _, rawAddr := range addrList {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if isPrivateIP(ip.String()) {
				return net.JoinHostPort(ip.String(), port), nil
			}
		}
	}
	return "", nil
}
