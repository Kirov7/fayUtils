package hash

import "testing"

type Item struct {
}

func TestS(t *testing.T) {
	items := []Item{}
	var all []*Item
	for _, item := range items {
		item := item
		all = append(all, &item)
	}
}
