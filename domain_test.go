package main

import (
	"log"
	"testing"
)

func TestDomainTester(t *testing.T) {
	tester := NewDomainTester("baidu.com", false)
	ok, err := tester.Test("baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		t.Fail()

	}

	ok, err = tester.Test("google.com")
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		t.Fail()
	}

	ok, err = tester.Test("mp3.baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		t.Fail()
	}

	tester = NewDomainTester("baidu.com", true)
	ok, err = tester.Test("baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		t.Fail()
	}
	ok, err = tester.Test("mp3.baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		t.Fail()
	}
}
