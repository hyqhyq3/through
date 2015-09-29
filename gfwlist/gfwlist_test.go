package gfwlist

import (
	"net/url"
	"testing"
)

// func TestGetListContent(t *testing.T) {
// 	content, err := GetListContent()
// 	print(content, err)
// }

// func TestParseList(t *testing.T) {
// 	content, err := GetListContent()
// 	rules, err := ParseList(content)
// 	fmt.Println(rules, err)
// }

func TestRulesTester_Test(t *testing.T) {
	tester, err := NewTester()
	if err != nil {
		t.Error(err)
	}
	u, _ := url.Parse("http://baidu.com")
	if tester.Test(u) {
		t.Error("百度没被墙")
	}
	u, _ = url.Parse("http://twitter.com")
	if !tester.Test(u) {
		t.Error("推特被墙了")
	}

	// 反向规则

	u, _ = url.Parse("http://qq.com")
	if tester.Test(u) {
		t.Error("qq.com 这个没被墙")
	}
}
