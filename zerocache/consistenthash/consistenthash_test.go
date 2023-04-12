package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	newhash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// 增加物理节点 2,4,6  对应添加虚拟节点是2,12,22,4,14,24,6,16,26
	newhash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	// 测试 testCases 中值会保存到哪个节点
	for k, v := range testCases {
		if v != newhash.Get(k) {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	// 增加一个节点, 虚拟节点会增加8,18,28
	newhash.Add("8")
	// 此时值27会保存到节点8
	testCases["27"] = "8"
	// 测试 testCases 中值会保存到哪个节点
	for k, v := range testCases {
		if v != newhash.Get(k) {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
