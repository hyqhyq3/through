// route.go
package main

import (
	"fmt"
	"log"
)

type Rule struct {
	Tester
	Route Dialer
}

var rules []*Rule
var defaultDialer Dialer

func AddRouteRule(t Tester, d Dialer) {
	fmt.Println(d.Name())
	rules = append(rules, &Rule{t, d})
}

func SetDefaultRule(d Dialer) {
	defaultDialer = d
}

func Route(addr string) Dialer {
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
