package zzcache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func benchKey(i int) []byte {
	return []byte("key-" + strconv.Itoa(i%10000))
}

func benchValue() []byte {
	return make([]byte, 128)
}

func BenchmarkSet(b *testing.B) {
	c, _ := New(512, "")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.Set(benchKey(i), benchValue(), time.Hour); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetParallel(b *testing.B) {
	c, _ := New(512, "")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if err := c.Set(benchKey(i), benchValue(), time.Hour); err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

func BenchmarkGet(b *testing.B) {
	c, _ := New(512, "")
	for i := 0; i < 10000; i++ {
		if err := c.Set(benchKey(i), benchValue(), time.Hour); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := c.Get(benchKey(i)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetParallel(b *testing.B) {
	c, _ := New(512, "")
	for i := 0; i < 10000; i++ {
		if err := c.Set(benchKey(i), benchValue(), time.Hour); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if _, err := c.Get(benchKey(i)); err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

func BenchmarkSetGetMixed(b *testing.B) {
	c, _ := New(512, "")
	for i := 0; i < 10000; i++ {
		if err := c.Set(benchKey(i), benchValue(), time.Hour); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if err := c.Set(benchKey(i), benchValue(), time.Hour); err != nil {
				b.Fatal(err)
			}
			if _, err := c.Get(benchKey(i)); err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	c, _ := New(512, "")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("del-%d", i))
		if err := c.Set(key, benchValue(), time.Hour); err != nil {
			b.Fatal(err)
		}
		if err := c.Delete(key); err != nil {
			b.Fatal(err)
		}
	}
}
