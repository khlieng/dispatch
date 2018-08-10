package server

import "net"

func addrToIPBytes(addr net.Addr) []byte {
	ip := addr.(*net.TCPAddr).IP

	if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4
	}

	return ip
}
