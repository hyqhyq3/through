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

func initConfig() {
	ProxyA, err := ProxyFromURL("https://hyq:pass@example.com/")
	if err != nil {
		log.Fatal(err)
	}
	ProxyB, err := ProxyFromURL("ss://method:pass@example.com:4000")
	if err != nil {
		log.Fatal(err)
	}
	AddRouteRule(NewDomainTester("google.com", true), ProxyA)
	AddRouteRule(NewDomainTester("google.co.jp", true), ProxyA)
	AddRouteRule(NewDomainTester("facebook.com", true), ProxyB)
	AddRouteRule(NewGeoIPTester("CN", true), &net.Dialer{})
	SetDefaultRule(ProxyA)
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
