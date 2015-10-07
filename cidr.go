package through

import (
	"log"
	"net"
)

type CIDRTester struct {
	*net.IPNet
}

func NewCIDRTester(s string) *CIDRTester {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		log.Fatal(err)
	}
	return &CIDRTester{ipnet}
}

func (c *CIDRTester) Test(host string) (bool, error) {
	ip := net.ParseIP(host)
	if ip == nil {
		return false, nil
	}
	if ip.Mask(c.Mask).Equal(c.IP) {
		return true, nil
	}
	return false, nil
}
