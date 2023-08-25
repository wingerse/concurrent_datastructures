package concurrentdatastructures

import (
	"sync"
	"testing"
)

func BenchmarkCounter(b *testing.B) {
	c := NewCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

func BenchmarkAtomicCounter(b *testing.B) {
	c := NewAtomicCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

func BenchmarkPartitionedAtomicCounter(b *testing.B) {
	c := NewPartitionedAtomicCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

func BenchmarkPaddedPartitionedAtomicCounter(b *testing.B) {
	c := NewPaddedPartitionedAtomicCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

func TestPaddedPartitionedAtomicCounter(t *testing.T) {
	c := NewPaddedPartitionedAtomicCounter()
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000000; j++ {
				c.Inc()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if c.Value() != 1000*1000000 {
		t.Fatal(c.Value())
	}
}
