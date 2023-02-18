package bloom

import (
	"fmt"
	"github.com/Kirov7/fayUtils/binaryUtil"
	bitmap2 "github.com/Kirov7/fayUtils/bitmap"
	"testing"
)

func TestHash(t *testing.T) {
	t.Log(Hash([]byte("ssadafawghwqarfwASQP")))
}

func TestBitmap(t *testing.T) {
	bitmap := bitmap2.MakeBitmapWithByteSize(8)
	t.Log(bitmap)
	bitmap.Insert(3)
	bitmap.Insert(2)
	t.Log(bitmap.Contains(3))
	t.Log(bitmap)
	bitmap.Delete(3)
	t.Log(bitmap)
	t.Log(bitmap.Contains(3))
	bitmap.Reset()
	t.Log(bitmap)
	t.Log(bitmap.Contains(1))
}

func TestSec(t *testing.T) {
	fmt.Println(binaryUtil.ToBinaryString(9))
	bitmap := bitmap2.MakeBitmapWithBitSize(1)
	bitmap.Insert(2)
	str := binaryUtil.ToBinaryString([]byte(bitmap))
	println(bitmap.Contains(2))
	fmt.Println(str)
}
