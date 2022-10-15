package container

import "github.com/Kirov7/fayUtils/container/skipLinkedList"

type Iterator interface {
	Next()
	Valid() bool
	Rewind()
	Item() Item
	Close() error
}

type Item interface {
	Entry() *skipLinkedList.Entry
}

type Options struct {
	Prefix []byte // 前缀
	IsAsc  bool   // 是否升序
}
