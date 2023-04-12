package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash
	replicas int            //虚拟节点倍数
	keys     []int          // 哈系环
	hashMap  map[int]string // 虚拟节点与真实节点的映射表 虚拟节点哈希指：真实节点名
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			nodehash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, nodehash)
			m.hashMap[nodehash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash_int := int(m.hash([]byte(key)))
	// binary search
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash_int
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
