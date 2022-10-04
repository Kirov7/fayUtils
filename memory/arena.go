package memory

import (
	"log"
	"sync/atomic"
)

var (
	MaxElementSize int
)

type Arena struct {
	n          uint32
	shouldGrow bool
	buf        []byte
}

func NewArenaDefault(n int64) *Arena {
	return NewArenaWithMax(n, 5000)
}

func NewArenaWithMax(n int64, max int) *Arena {
	MaxElementSize = max
	return &Arena{
		n:   1,
		buf: make([]byte, n),
	}
}

func (a *Arena) Allocate(sz uint32) uint32 {
	offset := atomic.AddUint32(&a.n, sz)
	if !a.shouldGrow {
		if offset > uint32(len(a.buf)) {
			log.Fatal("state error")
		}
		return offset - sz
	}
	//for memory alignment
	if int(offset) > len(a.buf)-MaxElementSize {
		growBy := uint32(len(a.buf))
		if growBy > 1<<30 {
			growBy = 1 << 30
		}
		if growBy < sz {
			growBy = sz
		}
		newBuf := make([]byte, len(a.buf)+int(growBy))
		if len(a.buf) != copy(newBuf, a.buf) {
			log.Fatal("assert error")
		}
		a.buf = newBuf
	}
	return offset - sz
}

func (a *Arena) PutElement() uint32 {
	//todo implement me
	panic("need to be implement")
}

func (a *Arena) GetElement() uint32 {
	//todo implement me
	panic("need to be implement")
}
