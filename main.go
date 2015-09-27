package main

import (
	"log"
	"net"
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

func initConfig() {
	//	ProxyA, err := ProxyFromURL("ProxyA", "https://hyq:pass@example.com/")
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	var err error
	cfg, err = ReadConfig("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	proxies := make(map[string]Dialer)

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

func serveHTTP(addr string, exit <-chan bool) {
	server := NewHttpProxyServer()
	if err := server.ListenAndServe(addr, exit); err != nil {
		log.Fatal(err)
	}
}

func startServer(exit <-chan bool) {

	for _, listenCfg := range cfg.Listen {
		switch listenCfg.Type {
		case "http":
			serveHTTP(listenCfg.Addr, exit)
		}
	}
}

func main() {

	initConfig()

	exit := make(chan bool)
	startServer(exit)
}
