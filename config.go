package through

import (
	"encoding/json"
	"io/ioutil"
)

type proxy struct {
	Name   string
	Server string
}

type rule struct {
	Type             string
	Net              string
	Domain           string
	IncludeSubDomain bool
	Resolve          bool
	Country          string
	Target           string
}

type listener struct {
	Type string
	Addr string
}

type Config struct {
	Proxies []*proxy
	Rules   []*rule
	Listen  []*listener
}

//	{
//		"listen": [
//			{
//				"type": "http",
//				"addr": ":4443"
//			}
//		],
//		"proxies":[
//			{
//				"name":"ProxyA",
//				"server":"https://hyq:k910407@s4.60in.com"
//			}
//		],
//		"rules":[
//			{
//				"type":"cidr",
//				"net":"10.0.0.0/8",
//				"target":"Direct"
//			},
//			{
//				"type":"cidr",
//				"net":"127.0.0.0/8",
//				"target":"Direct"
//			},
//			{
//				"type":"cidr",
//				"net":"192.168.1.0/24",
//				"target":"Direct"
//			},
//			{
//				"type":"geoip",
//				"country":"CN",
//				"target":"Direct"
//			},
//			{
//				"type":"domain",
//				"domain":"google.com",
//				"includeSubDomain":true,
//				"target":"ProxyA"
//			},
//			{
//				"type":"final",
//				"target":"Direct"
//			}
//		]
//	}
func ReadConfig(file string) (c *Config, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c = &Config{}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
