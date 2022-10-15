package memory

import (
	"log"
	"sync/atomic"
)

var (
	MaxElementSize int
)

type Arena struct {
	// The size of the has been allocated memory
	n          uint32
	shouldGrow bool
	// The memory space has been applied for
	buf []byte
}

func NewArenaDefault(n int64) *Arena {
	return NewArenaWithMax(n, 5000)
}

// NewArenaWithMax Initializes the amount of memory space requested and specifies element size limits
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
			log.Fatal("Error while copying")
		}
		a.buf = newBuf
	}
	return offset - sz
}

func (a *Arena) PutElement(e []byte) uint32 {
	keySz := uint32(len(e))
	offset := a.Allocate(keySz)
	buf := a.buf[offset : offset+keySz]
	if len(e) == copy(buf, e) {
		log.Fatal("Error while copying")
	}
	return offset
}

func (a *Arena) GetElement(offset uint32, size uint16) []byte {
	return a.buf[offset : offset+uint32(size)]
}
