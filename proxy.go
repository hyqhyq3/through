package main

import (
	"errors"
	"net/url"
)

type ProxyCreateFunc func(*url.URL) (Dialer, error)

var proxySchemas map[string]ProxyCreateFunc

func RegisterProxyType(schema string, f ProxyCreateFunc) {
	if proxySchemas == nil {
		proxySchemas = make(map[string]ProxyCreateFunc)
	}
	proxySchemas[schema] = f
}

func ProxyFromURL(s string) (Dialer, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if proxySchemas != nil {
		if f, ok := proxySchemas[u.Scheme]; ok {
			return f(u)
		}
	}
	return nil, errors.New("proxy: unknown scheme: " + u.Scheme)
}
