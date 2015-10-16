package main

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/hyqhyq3/through"
	"github.com/hyqhyq3/through/proxy"
)

var dialer = &through.RouteDialer{}

func main() {
	through.InitConfig("config.ini")

	syscall.Setregid(8347, 8347)
	fmt.Println(syscall.Getgid())
	l, err := net.Listen("tcp", ":8024")
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			break
		}
		go handle(c)
	}
}

func handle(c net.Conn) {
	defer c.Close()
	fmt.Println(c.LocalAddr())
	fmt.Println(c.RemoteAddr())
	orig, err := getOriginDst(c)
	if err != nil {
		log.Println(err)
		return
	}
	rc, err := dialer.Dial("tcp", orig.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer rc.Close()
	exit := make(chan bool, 1)
	go proxy.Splice(c, rc, exit)
	go proxy.Splice(rc, c, exit)
	<-exit
}
