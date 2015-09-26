package main

import (
	"log"
	"net"
	"testing"
)

func TestRoute(t *testing.T) {
	ProxyA, err := ProxyFromURL("https://hyq:k910407@s4.60in.com/")
	if err != nil {
		log.Fatal(err)
	}
	ProxyB, err := ProxyFromURL("ss://aes-128-cfb:8vvfChCmQwTuh7@jp2.0bad.com:24832")
	if err != nil {
		log.Fatal(err)
	}
	Direct := &net.Dialer{}
	AddRouteRule(NewCIDRTester("10.0.0.0/8"), Direct)
	AddRouteRule(NewCIDRTester("127.0.0.0/8"), Direct)
	AddRouteRule(NewCIDRTester("192.168.1.0/24"), Direct)
	AddRouteRule(NewDomainTester("google.com", true), ProxyA)
	AddRouteRule(NewDomainTester("google.co.jp", true), ProxyA)
	AddRouteRule(NewDomainTester("facebook.com", true), ProxyB)
	AddRouteRule(NewGeoIPTester("CN", true), Direct)
	SetDefaultRule(ProxyA)

	if Route("127.0.0.1") != Direct {
		t.Fail()
	}
}
