package main

import "time"

func main() {
	cache := NewLiteCache()
	cache.Set("foo", []byte("bar"), 2*time.Second)
	cache.Set("baz", []byte("qux"))
	if val, ok := cache.Get("foo"); ok {
		println(string(val))
	}
	if val, ok := cache.Get("baz"); ok {
		println(string(val))
	}
	time.Sleep(3 * time.Second)
	if _, ok := cache.Get("foo"); !ok {
		println("foo expired!")
	}
	if val, ok := cache.Get("baz"); ok {
		println(string(val))
	}
	cache.Set("foo", []byte("new"), 2*time.Second)
	time.Sleep(1 * time.Second)
	if val, ok := cache.Get("foo"); ok {
		println(string(val))
	}
}
