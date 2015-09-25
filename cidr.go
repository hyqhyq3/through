package main

import "net"

type CIDRTester struct {
	*net.IPNet
}

func NewCIDRTester(s string) (*CIDRTester, error) {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &CIDRTester{ipnet}, nil
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
