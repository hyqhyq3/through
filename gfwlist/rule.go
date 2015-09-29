package gfwlist

import (
	"net/url"
	"regexp"
	"strings"
)

type Rule interface {
	Test(u *url.URL) bool
}

type KeywordRule struct {
	Keyword string
}

func (r *KeywordRule) Test(u *url.URL) bool {
	return strings.Contains(u.String(), r.Keyword)
}

type DomainRule struct {
	Domain string
}

func (r *DomainRule) Test(u *url.URL) bool {
	str := u.String()
	if len(r.Domain) > len(str) {
		return false
	}
	if str[len(str)-len(r.Domain):] == r.Domain {
		return true
	}
	return false
}

type RegexRule struct {
	Reg *regexp.Regexp
}

func (r *RegexRule) Test(u *url.URL) bool {
	return r.Reg.FindString(u.String()) != ""
}

type PrefixRule struct {
	Keyword string
}

func (r *PrefixRule) Test(u *url.URL) bool {
	return strings.Index(r.Keyword, u.String()) == 0
}
