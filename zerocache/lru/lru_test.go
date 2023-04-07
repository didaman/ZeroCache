package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Put("aaa", String("1111"))
	if v, ok := lru.Get("aaa"); !ok || string(v.(String)) != "1111" {
		t.Fatalf("cache hit key aaa =1111 failed!")
	}
	if _, ok := lru.Get("bb"); ok {
		t.Fatalf("cache miss key bb failed!")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Put(k1, String(v1))
	lru.Put(k2, String(v2))
	lru.Put(k3, String(v3))
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("Remove oldest key1 failed!")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	k1, k2, k3, k4 := "key1", "key2", "k3", "k4"
	v1, v2, v3, v4 := "v11", "v2", "v3", "444"
	caps := len(k1 + v1 + k2 + v2)
	lru := New(int64(caps), callback)
	lru.Put(k1, String(v1))
	lru.Put(k2, String(v2))
	lru.Put(k3, String(v3))
	lru.Put(k4, String(v4))
	lru.Get(k3)
	except := []string{"key1", "key2"}
	if !reflect.DeepEqual(except, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", except)
	}
}
