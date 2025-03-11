package litecache

import (
	"bytes"
	"testing"
	"time"
)

func TestLiteCache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		cache := NewLiteCache()
		cache.Set("key", []byte("value"), 0)

		val, ok := cache.Get("key")
		if !ok {
			t.Fatal("Expected to get value for key, but got not ok")
		}
		if !bytes.Equal(val, []byte("value")) {
			t.Fatalf("Expected value %q but got %q", "value", string(val))
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		cache := NewLiteCache()
		cache.Set("temp", []byte("value"), 50*time.Millisecond)
		cache.Set("perm", []byte("forever"), 0)

		// Initial values should be available
		_, ok := cache.Get("temp")
		if !ok {
			t.Fatal("Expected to get temp value initially")
		}

		// Wait for expiration
		time.Sleep(100 * time.Millisecond)

		// Temporary key should be gone
		_, ok = cache.Get("temp")
		if ok {
			t.Fatal("Expected temp key to be expired")
		}

		// Permanent key should still exist
		_, ok = cache.Get("perm")
		if !ok {
			t.Fatal("Expected perm key to still exist")
		}
	})

	t.Run("Overwrite", func(t *testing.T) {
		cache := NewLiteCache()
		cache.Set("key", []byte("value1"), time.Second)
		cache.Set("key", []byte("value2"), time.Second)

		val, ok := cache.Get("key")
		if !ok {
			t.Fatal("Expected to get value for key, but got not ok")
		}
		if !bytes.Equal(val, []byte("value2")) {
			t.Fatalf("Expected value %q but got %q", "value2", string(val))
		}
	})
}
