package main

import (
	"errors"
	"net/url"

	"github.com/hyqhyq3/through/common"
)

type ProxyCreateFunc func(name string, u *url.URL) (common.Dialer, error)

var proxySchemas map[string]ProxyCreateFunc

func RegisterProxyType(schema string, f ProxyCreateFunc) {
	if proxySchemas == nil {
		proxySchemas = make(map[string]ProxyCreateFunc)
	}
	proxySchemas[schema] = f
}

func ProxyFromURL(name, s string) (common.Dialer, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if proxySchemas != nil {
		if f, ok := proxySchemas[u.Scheme]; ok {
			return f(name, u)
		}
	}
	return nil, errors.New("proxy: unknown scheme: " + u.Scheme)
}
