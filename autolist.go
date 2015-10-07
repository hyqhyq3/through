package through

import (
	"log"
	"net/url"

	"github.com/hyqhyq3/through/autoproxy"
)

type AutoProxyTester struct {
	autoproxy.Tester
}

func NewAutoProxyTester() *AutoProxyTester {
	tester, err := autoproxy.NewTester()
	if err != nil {
		log.Fatal(err)
	}
	return &AutoProxyTester{tester}
}

func (c *AutoProxyTester) Test(addr string) (bool, error) {
	var str string
	str = "http://" + addr + "/"
	u, err := url.Parse(str)
	if err != nil {
		return false, err
	}
	return c.Tester.Test(u), nil
}
