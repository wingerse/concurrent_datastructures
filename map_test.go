package concurrentdatastructures

import (
	"runtime"
	"sync"
	"testing"

	"github.com/petermattis/goid"
)

func TestMap(t *testing.T) {
	m := NewMap()
	for i := uint64(0); i < 100000; i++ {
		m.Set(i, i)
	}
	for i := uint64(0); i < 100000; i++ {
		if v, ok := m.Get(i); !(ok && v == i) {
			t.Fatal(v, ok)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	m := NewMap()
	for i := uint64(0); i < 1000; i++ {
		m.Set(i, i)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if goid.Get()%2 == 0 {
				m.Set(uint64(goid.Get()), 1)
			} else {
				m.Get(uint64(goid.Get()))
			}
		}
	})
}

func BenchmarkMapRW(b *testing.B) {
	m := NewMapRW()
	for i := uint64(0); i < 1000; i++ {
		m.Set(i, i)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get(uint64(goid.Get()))
		}
	})
}

func BenchmarkPartitionedMap(b *testing.B) {
	m := NewPartitionedMap()
	for i := uint64(0); i < 1000; i++ {
		m.Set(i, i)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if goid.Get()%2 == 0 {
				m.Set(uint64(goid.Get()), 1)
			} else {
				m.Get(uint64(goid.Get()))
			}
		}
	})
}

func TestLockFreeMap(t *testing.T) {
	m := NewLockFreeMap()
	for i := uint64(1); i <= mapsize; i++ {
		m.Set(i, i)
	}
	for i := uint64(1); i <= mapsize; i++ {
		if v, ok := m.Get(i); !(ok && v == i) {
			t.Fatal(v, ok)
		}
	}
}

func TestLockFreeMapParallel(t *testing.T) {
	m := NewLockFreeMap()
	wg := &sync.WaitGroup{}
	for g := 0; g < runtime.GOMAXPROCS(0); g++ {
		wg.Add(1)
		go func() {
			for i := uint64(1); i <= mapsize; i++ {
				m.Set(i, i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for i := uint64(1); i <= mapsize; i++ {
		if v, ok := m.Get(i); !(ok && v == i) {
			t.Fatal(v, ok)
		}
	}
}

func BenchmarkLockFreeMap(b *testing.B) {
	m := NewLockFreeMap()
	for i := uint64(0); i < 1000; i++ {
		m.Set(i, i)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if goid.Get()%2 == 0 {
				m.Set(uint64(goid.Get()), 1)
			} else {
				m.Get(uint64(goid.Get()))
			}
		}
	})
}
