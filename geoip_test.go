package through

import (
	"log"
	"testing"
)

func TestGeoIPTester(t *testing.T) {
	var tester Tester
	tester = NewGeoIPTester("CN", true)
	ok, err := tester.Test("baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	if ok == false {
		t.Fail()
	}

	ok, err = tester.Test("114.114.114.114")
	if err != nil {
		log.Fatal(err)
	}
	if ok == false {
		t.Fail()
	}

	ok, err = tester.Test("8.8.8.8")
	if err != nil {
		log.Fatal(err)
	}
	if ok == true {
		t.Fail()
	}
}
