package lru

import (
	"container/list"
)

// cache, LRU policy, not safe for concurrent access

// Value use Len to count how many bytes it takes
type Cache struct {
	maxBytes  int64 // 最大内存
	nbytes    int64 // 当前已经使用的内存
	ll        *list.List
	cache     map[string]*list.Element      // Element 是链表元素的抽象类型
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数
}

type entry struct {
	key   string
	value Value
}
type Value interface {
	Len() int
}

func New(maxBytes int64, onEivcted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEivcted,
	}
}

// Get a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok { // go 的多赋值语法，ele 即 map 中的对应元素，ok 表示 key 是否存在
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// remove the according LRU policy
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// add or mod (if exist)
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
