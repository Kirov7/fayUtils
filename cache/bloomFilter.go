package cache

type Filter []byte

type BloomFilter struct {
	bitmap Filter
	k      uint8
}

func (bf *BloomFilter) MayContainKey(k []byte) bool {
	panic("implement me")
}

func (bf *BloomFilter) mayContain(h uint32) bool {
	panic("implement me")
}

func (bf BloomFilter) Len() int32 {
	panic("implement me")
}

func (bf *BloomFilter) InsertKey(k []byte) bool {
	panic("implement me")
}

func (bf *BloomFilter) insert(h uint32) bool {
	panic("implement me")
}

func (bf *BloomFilter) AllowKey(k []byte) bool {
	panic("implement me")
}

func (bf *BloomFilter) allow(h uint32) bool {
	panic("implement me")
}

func (bf *BloomFilter) reset() {
	panic("implement me")
}

func NewBloomFilter(numEntries int, falsePositive float64) *BloomFilter {
	panic("implement me")
}

func Hash(b []byte) uint32 {
	panic("implement me")
}
