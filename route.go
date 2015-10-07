package through

import (
	"fmt"
	"log"

	"github.com/hyqhyq3/through/common"
)

type Rule struct {
	Tester
	Route common.Dialer
}

var rules []*Rule
var defaultDialer common.Dialer

func AddRouteRule(t Tester, d common.Dialer) {
	fmt.Println(d.Name())
	rules = append(rules, &Rule{t, d})
}

func SetDefaultRule(d common.Dialer) {
	defaultDialer = d
}

func Route(addr string) common.Dialer {
	for _, d := range rules {
		ok, err := d.Test(addr)
		if err != nil {
			log.Println(err)
			continue
		}
		if ok {
			log.Println("Route:", addr, d.Route.Name())
			return d.Route
		}
	}
	log.Println("Route:", addr, defaultDialer.Name())
	return defaultDialer
}
