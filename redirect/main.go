package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

func main() {
	syscall.Setreuid(1001, 1001)
	fmt.Println(syscall.Getuid())
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
	fmt.Println(c.LocalAddr())
	fmt.Println(c.RemoteAddr())
	orig, err := getOriginDst(c)
	fmt.Println(orig, err)
}
