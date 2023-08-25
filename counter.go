package concurrentdatastructures

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

type Counter struct {
	val uint64
	mu  *sync.Mutex
}

func NewCounter() *Counter {
	return &Counter{
		val: 0,
		mu:  &sync.Mutex{},
	}
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val++
}

func (c *Counter) Value() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}

type AtomicCounter struct {
	val atomic.Uint64
}

func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{
		val: atomic.Uint64{},
	}
}

func (c *AtomicCounter) Inc() {
	c.val.Add(1)
}

func (c *AtomicCounter) Value() uint64 {
	return c.val.Load()
}

type PartitionedAtomicCounter struct {
	counters []atomic.Uint64
}

func NewPartitionedAtomicCounter() *PartitionedAtomicCounter {
	return &PartitionedAtomicCounter{
		counters: make([]atomic.Uint64, runtime.GOMAXPROCS(0)),
	}
}

func (c *PartitionedAtomicCounter) Inc() {
	c.counters[int(goid.Get())%len(c.counters)].Add(1)
}

func (c *PartitionedAtomicCounter) Value() uint64 {
	var val uint64 = 0
	for i := 0; i < len(c.counters); i++ {
		val += c.counters[i].Load()
	}
	return val
}

const cacheLineSize = 64

type cacheLinePaddedUint64 struct {
	v atomic.Uint64
	_ [cacheLineSize - 8]byte
}

type PaddedPartitionedAtomicCounter struct {
	counters []cacheLinePaddedUint64
}

func NewPaddedPartitionedAtomicCounter() *PaddedPartitionedAtomicCounter {
	return &PaddedPartitionedAtomicCounter{
		counters: make([]cacheLinePaddedUint64, runtime.GOMAXPROCS(0)),
	}
}

func (c *PaddedPartitionedAtomicCounter) Inc() {
	c.counters[int(goid.Get())%len(c.counters)].v.Add(1)
}

func (c *PaddedPartitionedAtomicCounter) Value() uint64 {
	var val uint64 = 0
	for i := 0; i < len(c.counters); i++ {
		val += c.counters[i].v.Load()
	}
	return val
}
