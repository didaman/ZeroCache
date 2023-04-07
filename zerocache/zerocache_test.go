package zerocache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "123",
	"Jack": "456",
	"Sam":  "789",
}

func TestGetter(t *testing.T) {
	// 测试GetterFunc类型转换， 将一个匿名回调函数转换成了接口Getter
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	except := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, except) {
		t.Errorf("callback failed")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	zero := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[TestDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := zero.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of xxx")
		}
		if _, err := zero.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}
	if view, err := zero.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
