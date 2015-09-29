package gfwlist_test

import (
	"testing"

	"github.com/hyqhyq3/through/gfwlist"
)

func TestGetListContent(t *testing.T) {
	content, err := gfwlist.GetListContent()
	print(content, err)
}
