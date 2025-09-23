package utils

import (
	"net"
)

// ParseIPNet IPアドレス文字列をnet.IPNetに変換
func ParseIPNet(ipStr string) net.IPNet {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return net.IPNet{}
	}

	// IPアドレスのビット数を判定
	bits := 32
	if ip.To4() == nil {
		// IPv6の場合
		bits = 128
	}

	return net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(bits, bits),
	}
}
