package main

import (
	"bytes"
	"errors"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

func reverse(m []string) {
	l := len(m)
	for i := 0; i < l/2; i++ {
		m[i], m[l-1-i] = m[l-1-i], m[i]
	}
}

func getLineFromPfctl(local, remote string) string {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("pfctl", "-s", "state")
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Run()
	for {
		line, err := buf.ReadString(byte('\n'))
		if err != nil {
			break
		}
		if strings.Contains(line, local) && strings.Contains(line, remote) {
			return line
		}
	}
	return ""
}

func getOriginDst(c net.Conn) (addr net.Addr, err error) {
	line := getLineFromPfctl(c.LocalAddr().String(), c.RemoteAddr().String())
	re, _ := regexp.Compile(`([0-9]+\.){3}[0-9]+:\d+`)
	m := re.FindAllString(line, -1)
	if strings.Contains(line, "<-") {
		reverse(m)
	}
	if len(m) >= 2 {
		addr, err = net.ResolveTCPAddr("tcp", m[1])
		return
	}
	return nil, errors.New("cannot find dst")
}
