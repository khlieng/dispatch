package netutil

import "net"

var privateNets []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"::1/128",
		"fe80::/10",
		"fc00::/7",
	} {
		_, network, _ := net.ParseCIDR(cidr)
		privateNets = append(privateNets, network)
	}
}

func IsPrivate(host string) bool {
	if host == "localhost" {
		return true
	}
	return IsPrivateIP(net.ParseIP(host))
}

func IsPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	for _, privateNet := range privateNets {
		if privateNet.Contains(ip) {
			return true
		}
	}
	return false
}
