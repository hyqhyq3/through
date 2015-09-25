package main

import (
	"log"
	"testing"
)

func TestCIDRTester(t *testing.T) {
	tester, _ := NewCIDRTester("192.168.1.1/24")
	ok, err := tester.Test("192.168.1.254")
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		t.Fail()
	}

	ok, err = tester.Test("121.41.85.207")
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		t.Fail()
	}
}
