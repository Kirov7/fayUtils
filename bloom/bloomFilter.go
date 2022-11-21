package bloom

import (
	bitmap "github.com/Kirov7/fayUtils/bitmap"
	"github.com/Kirov7/fayUtils/hash"
	"math"
)

type BloomFilter struct {
	bitmap bitmap.Filter
	k      uint8 // hash Function nums
}

func (bf *BloomFilter) MayContainKey(k []byte) bool {
	return bf.mayContain(Hash(k))
}

func (bf *BloomFilter) mayContain(h uint32) bool {
	if bf.len() < 2 {
		return false
	}
	k := bf.k
	if k > 30 {
		return true
	}

	nBits := uint32(8 * (bf.len() - 1))
	delta := h>>17 | h<<15
	for j := uint8(0); j < k; j++ {
		bitPos := h & nBits
		if !bf.bitmap.Contains(bitPos) {
			return false
		}
		h += delta
	}
	return true
}

func (bf BloomFilter) len() int32 {
	return int32(len(bf.bitmap))
}

func (bf *BloomFilter) InsertKey(k []byte) bool {
	return bf.insert(Hash(k))
}

func (bf *BloomFilter) insert(h uint32) bool {
	k := bf.k
	if k > 30 {
		return true
	}
	nBits := uint32(8 * (bf.len() - 1))
	delta := h>>17 | h<<15
	for j := uint8(0); j < k; j++ {
		bitPos := h % uint32(nBits)
		bf.bitmap.Insert(bitPos)
		h += delta
	}
	return true
}

func (bf *BloomFilter) AllowKey(k []byte) bool {
	if bf == nil {
		return true
	}
	already := bf.MayContainKey(k)
	if !already {
		bf.InsertKey(k)
	}
	return already
}

func (bf *BloomFilter) Allow(h uint32) bool {
	if bf == nil {
		return true
	}
	already := bf.mayContain(h)
	if !already {
		bf.insert(h)
	}
	return already
}

func (bf *BloomFilter) Reset() {
	if bf == nil {
		return
	}
	for i := range bf.bitmap {
		bf.bitmap[i] = 0
	}
}

func NewBloomFilter(numEntries int, falsePositive float64) *BloomFilter {
	bitsPerKey := bitsPerKey(numEntries, falsePositive)
	return initFilter(numEntries, bitsPerKey)
}

func NewBloomFilterDefault(numEntries int) *BloomFilter {
	return initFilter(numEntries, 10)
}

func bitsPerKey(numEntries int, falsePositive float64) int {
	size := -1 * float64(numEntries) * math.Log(falsePositive) / math.Pow(math.Ln2, 2)
	locs := math.Ceil(size / float64(numEntries))
	return int(locs)
}

func initFilter(numEntries int, bitsPerKey int) *BloomFilter {
	bf := &BloomFilter{}
	if bitsPerKey < 0 {
		bitsPerKey = 0
	}
	k := uint32(float64(bitsPerKey) * math.Ln2)
	if k < 1 {
		k = 1
	}
	if k > 30 {
		k = 30
	}
	bf.k = uint8(k)

	nBits := numEntries * int(bitsPerKey)
	if nBits < 64 {
		nBits = 64
	}
	// bitmap`s []byte length
	nBytes := (nBits + 7) / 8
	nBits = nBytes * 8
	filter := bitmap.MakeBitmapWithByteSize(nBytes + 1)

	filter.SetHashNum(uint8(k))

	bf.bitmap = filter
	return bf
}

func Hash(b []byte) uint32 {
	return hash.MurmurHash(b)
}
