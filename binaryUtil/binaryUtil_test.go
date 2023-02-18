package binaryUtil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Kirov7/fayUtils/bitmap"
	"github.com/Kirov7/fayUtils/bloom"
	"math"
	"testing"
)

func TestUtil(t *testing.T) {
	str := ToBinaryString(10000)
	t.Log(str)
	var s int64
	err := ReadBinaryString(str, &s)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	t.Log(s)
}

func TestSec(t *testing.T) {
	bitmap := bitmap.MakeBitmapWithBitSize(math.MaxUint32)
	for i := uint32(200); i < 10000000; i++ {
		bitmap.Insert(bloom.Hash([]byte(fmt.Sprintf("%d", i))))
	}
	fmt.Println("插入成功")
	trueNum, falseNum := 0, 0
	for i := uint32(0); i < 10000000; i++ {
		if bitmap.Contains(bloom.Hash([]byte(fmt.Sprintf("%d", i)))) {
			trueNum++
		} else {
			falseNum++
		}
	}

	fmt.Println("ture: ", trueNum)
	fmt.Println("false: ", falseNum)
}

func Uint32ToBytes(n uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}
