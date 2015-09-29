package gfwlist

import "net/url"

type Tester interface {
	Test(u *url.URL) bool
}

type RulesTester struct {
	except []Rule
	rules  []Rule
}

func NewTester() (Tester, error) {
	content, err := GetListContent()
	if err != nil {
		return nil, err
	}
	except, rules, err := ParseList(content)
	if err != nil {
		return nil, err
	}
	return &RulesTester{except, rules}, nil
}

func (t RulesTester) Test(u *url.URL) bool {
	for _, rule := range t.except {
		if rule.Test(u) {
			return false
		}
	}
	for _, rule := range t.rules {
		if rule.Test(u) {
			return true
		}
	}
	return false
}
