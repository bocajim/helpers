package helpers

import (
	"net"
)

func inc(ip net.IP) {
	for j := len(ip)-1; j>=0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func NetCidrToIps(cidr string) ([]string,error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil,err
	}
	ret := make([]string,0,256)
	
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ret = append(ret, ip.String())
	}
	return ret, nil
}

