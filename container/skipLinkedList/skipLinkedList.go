package skipLinkedList

import (
	"bytes"
	"errors"
	"github.com/Kirov7/fayUtils/container"
	"time"

	"log"
	"math/rand"
	"sync"
)

const (
	defaultMaxHeight = 64
)

type SkipList struct {
	header *Element

	rand *rand.Rand

	maxLevel int
	length   int
	lock     sync.RWMutex
	size     int64
}
type Element struct {
	levels []*Element // levels[i] 第 i 层所指向的下一个节点
	entry  *Entry     // 存储的键值对
	score  float64    // 通过计算得出的分数,用于进行快速比较
}

type Entry struct {
	Key       []byte
	Value     []byte
	ExpiresAt uint64 // 过期时间
}

func NewEntry(key, value []byte) *Entry {
	return &Entry{
		Key:   key,
		Value: value,
	}
}

// WithTTL 为Entry添加自动过期时间
func (e *Entry) WithTTL(dur time.Duration) *Entry {
	e.ExpiresAt = uint64(time.Now().Add(dur).Unix())
	return e
}

func (e *Entry) Size() int64 {
	return int64(len(e.Key) + len(e.Value))
}

func (e *Entry) Entry() *Entry {
	return e
}

func newElement(score float64, entry *Entry, level int) *Element {
	return &Element{
		levels: make([]*Element, level+1),
		entry:  entry,
		score:  score,
	}
}

func (elem *Element) Entry() *Entry {
	return elem.entry
}

func NewSkipList() *SkipList {
	header := &Element{
		levels: make([]*Element, defaultMaxHeight),
	}

	return &SkipList{
		header:   header,
		maxLevel: defaultMaxHeight - 1,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (list *SkipList) Add(data *Entry) error {
	// 需要加锁
	list.lock.Lock()
	defer list.lock.Unlock()

	prevs := make([]*Element, list.maxLevel+1)

	key := data.Key
	keyScore := list.calcScore(key)
	header, maxLevel := list.header, list.maxLevel
	prev := header
	// 从最高层开始比较
	for i := maxLevel; i >= 0; i-- {
		for ne := prev.levels[i]; ne != nil; ne = prev.levels[i] {
			if comp := list.compare(keyScore, key, ne); comp <= 0 {
				if comp == 0 {
					// 如果kv对已存在,则直接更新数据
					ne.entry = data
					return nil
				} else {
					prev = ne
				}
			} else {
				// 如果同层下一个元素大于当前key的话则向下一层
				break
			}
		}
		prevs[i] = prev
	}
	// 找到插入的位置后,计算层数和key前8字节的摘要
	randLevel, keyScore := list.randLevel(), list.calcScore(key)
	e := newElement(keyScore, data, randLevel)

	for i := randLevel; i >= 0; i-- {
		ne := prevs[i].levels[i]
		prevs[i].levels[i] = e
		e.levels[i] = ne
	}
	return nil
}

func (list *SkipList) Search(key []byte) (e *Entry) {
	// 加读锁
	list.lock.RLock()
	defer list.lock.RUnlock()
	keyScore := list.calcScore(key)
	header, maxLevel := list.header, list.maxLevel
	prev := header
	for i := maxLevel; i >= 0; i-- {
		for ne := prev.levels[i]; ne != nil; ne = prev.levels[i] {
			if comp := list.compare(keyScore, key, ne); comp <= 0 {
				if comp == 0 {
					return ne.entry
				} else {
					prev = ne
				}
			} else {
				break
			}
		}
	}
	return nil
}

// Close 关闭skipList资源
func (list *SkipList) Close() error {
	return nil
}

// calcScore 计算 key 的摘要值
// 1byte 为 8bit, 一个uint64可以存储64bit,即可以存储8byte
// 将key的前8个字符存在uint64里面,就可以将逐位比较改为无符号整数比较,以优化比较速度
func (list *SkipList) calcScore(key []byte) (score float64) {
	var hash uint64
	l := len(key)

	if l > 8 {
		l = 8
	}

	// 将第i个byte存在到uint64中的搞 8 * i 位
	for i := 0; i < l; i++ {
		shift := uint(64 - 8 - i*8)
		hash |= uint64(key[i]) << shift
	}

	score = float64(hash)
	return
}

// compare 比较器函数
func (list *SkipList) compare(score float64, key []byte, next *Element) int {
	if score == next.score {
		return bytes.Compare(key, next.entry.Key)
	}

	if score < next.score {
		return -1
	} else {
		return 1
	}
}

// randLevel 添加新节点的时候,随机产生层数
func (list *SkipList) randLevel() int {
	// 2^(-i) 的概率返回 i
	for i := 0; i < list.maxLevel; i++ {
		if list.rand.Intn(2) == 0 {
			return i
		}
	}

	return list.maxLevel
}

func (list *SkipList) Size() int64 {
	return list.size
}

func (list *SkipList) NewSkipListIterator() container.Iterator {
	return &skipListIterator{skipList: list}
}

type skipListIterator struct {
	skipList *SkipList
	e        *Entry
}

func (s *skipListIterator) Next() {

	if !s.Valid() {
		log.Fatalf("%+v", errors.New("has no next"))
	}
	//TODO implement me
	panic("implement me")
}

func (s *skipListIterator) Valid() bool {
	return s.e != nil
}

func (s *skipListIterator) Rewind() {
	//TODO implement me
	panic("implement me")
}

func (s *skipListIterator) Item() container.Item {
	return &Entry{
		Key:       s.e.Key,
		Value:     s.e.Value,
		ExpiresAt: s.e.ExpiresAt,
	}
}

func (s *skipListIterator) Close() error {
	return nil
}
