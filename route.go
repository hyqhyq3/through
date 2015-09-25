// route.go
package main

import "log"

type Rule struct {
	Tester
	Route Dialer
}

var rules []*Rule
var defaultDialer Dialer

func AddRouteRule(t Tester, d Dialer) {
	rules = append(rules, &Rule{t, d})
}

func SetDefaultRule(d Dialer) {
	defaultDialer = d
}

func Route(addr string) Dialer {
	for _, d := range rules {
		ok, err := d.Test(addr)
		log.Println(addr, ok, err)
		if err != nil {
			log.Println(err)
			continue
		}
		if ok {
			return d.Route
		}
	}
	return defaultDialer
}
