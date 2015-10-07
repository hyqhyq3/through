package through

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/hyqhyq3/through/common"
)

type Tester interface {
	Test(string) (bool, error)
}

type AlwaysTrueTeseter struct {
}

func (*AlwaysTrueTeseter) Test(string) (bool, error) {
	return true, nil
}

type DefaultDialer struct {
	name string
}

func (d *DefaultDialer) Dial(network, addr string) (c net.Conn, err error) {
	return net.Dial(network, addr)
}

func (d *DefaultDialer) Name() string {
	return d.name
}

var cfg *Config

func InitConfig(fileName string) {

	var err error
	cfg, err = ReadConfig(fileName)
	if err != nil {
		log.Fatal(err)
	}
	proxies := make(map[string]common.Dialer)

	proxies["Direct"] = &DefaultDialer{"Direct"}

	for _, proxy := range cfg.Proxies {
		server, err := ProxyFromURL(proxy.Name, proxy.Server)
		if err != nil {
			log.Fatal(err)
		}
		proxies[proxy.Name] = server
	}

	for _, rule := range cfg.Rules {
		var tester Tester
		switch rule.Type {
		case "cidr":
			tester = NewCIDRTester(rule.Net)
		case "domain":
			tester = NewDomainTester(rule.Domain, rule.IncludeSubDomain)
		case "geoip":
			tester = NewGeoIPTester(rule.Country, rule.Resolve)
		case "autoproxy":
			tester = NewAutoProxyTester()
		case "final":
			tester = &AlwaysTrueTeseter{}
		}
		if target, ok := proxies[rule.Target]; ok {
			AddRouteRule(tester, target)
		} else {
			log.Fatal("target not exitsts:", target)
		}
	}

}

type RouteDialer struct{}

func (r *RouteDialer) Dial(network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	d := Route(host)
	if d == nil {
		return nil, errors.New(fmt.Sprintf("No rule for %s", host))
	}

	return d.Dial(network, addr)
}

func (r *RouteDialer) Name() string {
	return "RouteDialer"
}
