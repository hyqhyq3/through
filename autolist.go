package main

import (
	"log"
	"net"
	"net/url"

	"github.com/hyqhyq3/through/autoproxy"
)

type AutoProxyTester struct {
	autoproxy.Tester
}

func NewAutoProxyTester(s string) *AutoProxyTester {
	tester, err := autoproxy.NewTester()
	if err != nil {
		log.Fatal(err)
	}
	return &AutoProxyTester{tester}
}

func (c *AutoProxyTester) Test(addr string) (bool, error) {
	host, port, _ := net.SplitHostPort(addr)
	var str string
	if port == "80" {
		str = "http://" + host + "/"
	} else if port == "443" {
		str = "https://" + host + "/"
	} else {
		return false, nil
	}
	u, err := url.Parse(str)
	if err != nil {
		return false, err
	}
	return c.Tester.Test(u), nil
}
