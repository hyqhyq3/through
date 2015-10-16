package main

import (
	"net"
	"testing"
)

func TestGetLine(t *testing.T) {
	c, _ := net.Dial("tcp", "192.168.1.254:22")
	line := getLineFromPfctl(c.LocalAddr().String(), c.RemoteAddr().String())
	t.Log(line)
}
