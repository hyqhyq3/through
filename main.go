package main

import (
	"log"
	"net"
)

type Tester interface {
	Test(string) (bool, error)
}

//type DirectDialer struct{}

//func NewDirectDialer() Dialer {
//	return &DirectDialer{}
//}

//func (d*DirectDialer) Dial(network,addr string) (net.con)

type DefaultDialer struct {
	name string
}

func (d *DefaultDialer) Dial(network, addr string) (c net.Conn, err error) {
	return net.Dial(network, addr)
}

func (d *DefaultDialer) Name() string {
	return d.name
}

func initConfig() {
	//	ProxyA, err := ProxyFromURL("ProxyA", "https://hyq:pass@example.com/")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	ProxyB, err := ProxyFromURL("ProxyB", "ss://method:pass@example.com:4000")
	if err != nil {
		log.Fatal(err)
	}
	Direct := &DefaultDialer{"Direct"}
	AddRouteRule(NewCIDRTester("10.0.0.0/8"), Direct)
	AddRouteRule(NewCIDRTester("127.0.0.0/8"), Direct)
	AddRouteRule(NewCIDRTester("192.168.1.0/24"), Direct)
	AddRouteRule(NewDomainTester("google.com", true), ProxyB)
	AddRouteRule(NewDomainTester("google.co.jp", true), ProxyB)
	AddRouteRule(NewDomainTester("facebook.com", true), ProxyB)
	AddRouteRule(NewGeoIPTester("CN", true), Direct)
	SetDefaultRule(Direct)
}

func startServer(exit <-chan bool) {

	server := NewHttpProxyServer()
	if err := server.ListenAndServe(":4443", exit); err != nil {
		log.Fatal(err)
	}
}

func main() {

	initConfig()

	exit := make(chan bool)
	startServer(exit)

}
