package main

import "strings"

type DomainTester struct {
	Domain            string
	IncludeSubDomains bool
}

func (rule *DomainTester) Test(host string) (bool, error) {
	if !rule.IncludeSubDomains {
		return rule.Domain == host, nil
	}
	for {
		if rule.Domain == host {
			return true, nil
		}
		if len(host) > len(rule.Domain) && strings.Contains(host, ".") {
			host = host[strings.Index(host, ".")+1:]
		} else {
			return false, nil
		}
	}
}

func NewDomainTester(domain string, includeSubDomain bool) *DomainTester {
	return &DomainTester{domain, includeSubDomain}
}
