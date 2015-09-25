package main

import (
	"fmt"
	"log"
	"net"
)

type Dialer interface {
	Dial(network, addr string) (c net.Conn, err error)
}

type Tester interface {
	Test(string) (bool, error)
}

type Rule struct {
	Tester
	Route Dialer
}

func main() {
	ProxyA, err := ProxyFromURL("https://hyq:k910407@s4.60in.com/")
	if err != nil {
		log.Fatal(err)
	}
	rule := &Rule{NewDomainTester("google.com", true), ProxyA}
	rule.Test("hello")

	_, err = ProxyA.Dial("tcp", "google.com:80")
	if err != nil {
		fmt.Println(err)
	}
}
