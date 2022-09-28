package bloom

import (
	"fmt"
	"github.com/imroc/biu"
)

type Filter []byte

func MakeBitmap(nBytes int) Filter {
	filter := make([]byte, nBytes)
	return filter
}

func (f Filter) SetHashNum(hashNum uint8) {
	f[len(f)-1] = byte(hashNum)
}

// Insert Change the bitPos th bit to 1
func (f Filter) Insert(bitPos uint32) {
	f[bitPos/8] |= 1 << (bitPos % 8)
}

// Delete Change the bitPos th bit to 0
func (f Filter) Delete(bitPos uint32) {
	f[bitPos/8] &= ^(1 << (bitPos % 8))
}

func (f Filter) Contains(bitPos uint32) bool {
	fmt.Println(biu.ByteToBinaryString(f[bitPos/8]))
	fmt.Println(biu.ByteToBinaryString(1 << (bitPos % 8)))
	fmt.Println(biu.ByteToBinaryString(f[bitPos/8] & (1 << (bitPos % 8))))
	bit := f[bitPos/8] & (1 << (bitPos % 8))
	fmt.Println("bit:", bit)
	return f[bitPos/8]&(1<<(bitPos%8)) != 0
}

func (f Filter) Reset() {
	for i := range f {
		f[i] = 0
	}
}

func (f Filter) BitSize(bitPos uint32) int {
	return len(f) * 8
}
