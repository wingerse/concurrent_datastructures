package concurrentdatastructures

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/cespare/xxhash"
)

const mapsize = 1 << 20

type entry struct {
	key   uint64
	value uint64
}

type Map struct {
	entries []entry
	mu      *sync.Mutex
}

func NewMap() *Map {
	return &Map{
		entries: make([]entry, mapsize),
		mu:      &sync.Mutex{},
	}
}

func hash(v uint64) uint64 {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	return xxhash.Sum64(buf[:])
}

func (m *Map) Get(key uint64) (uint64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i := hash(key)
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		if m.entries[i].key == key {
			return m.entries[i].value, true
		} else if m.entries[i].key == 0 {
			return 0, false
		}
		count++
		i++
	}

	return 0, false
}

func (m *Map) Set(key, value uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i := hash(key)
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		entry := &m.entries[i]
		if entry.key == key || entry.key == 0 {
			entry.key = key
			entry.value = value
			return
		}
		count++
		i++
	}

	panic("map full")
}

type MapRW struct {
	entries []entry
	mu      *sync.RWMutex
}

func NewMapRW() *MapRW {
	return &MapRW{
		entries: make([]entry, mapsize),
		mu:      &sync.RWMutex{},
	}
}

func (m *MapRW) Get(key uint64) (uint64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	i := hash(key)
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		if m.entries[i].key == key {
			return m.entries[i].value, true
		} else if m.entries[i].key == 0 {
			return 0, false
		}
		count++
		i++
	}

	return 0, false
}

func (m *MapRW) Set(key, value uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i := hash(key)
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		entry := &m.entries[i]
		if entry.key == key || entry.key == 0 {
			entry.key = key
			entry.value = value
			return
		}
		count++
		i++
	}

	panic("map full")
}

type PartitionedMap struct {
	entries []entry
	mus     []*sync.Mutex
}

func NewPartitionedMap() *PartitionedMap {
	m := &PartitionedMap{
		entries: make([]entry, mapsize),
	}
	for i := 0; i < 64; i++ {
		m.mus = append(m.mus, &sync.Mutex{})
	}
	return m
}

func (m *PartitionedMap) Get(key uint64) (uint64, bool) {
	i := hash(key)
	mui := i & 63
	m.mus[mui].Lock()
	defer m.mus[mui].Unlock()
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		if m.entries[i].key == key {
			return m.entries[i].value, true
		} else if m.entries[i].key == 0 {
			return 0, false
		}
		count++
		i++
	}

	return 0, false
}

func (m *PartitionedMap) Set(key, value uint64) {
	i := hash(key)
	mui := i & 63
	m.mus[mui].Lock()
	defer m.mus[mui].Unlock()
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		entry := &m.entries[i]
		if entry.key == key || entry.key == 0 {
			entry.key = key
			entry.value = value
			return
		}
		count++
		i++
	}

	panic("map full")
}

type atomicEntry struct {
	key   atomic.Uint64
	value atomic.Uint64
}

type LockFreeMap struct {
	entries []atomicEntry
}

func NewLockFreeMap() *LockFreeMap {
	return &LockFreeMap{
		entries: make([]atomicEntry, mapsize),
	}
}

func (m *LockFreeMap) Get(key uint64) (uint64, bool) {
	i := hash(key)
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		k := m.entries[i].key.Load()
		if k == key {
			return m.entries[i].value.Load(), true
		} else if k == 0 {
			return 0, false
		}
		count++
		i++
	}

	return 0, false
}

func (m *LockFreeMap) Set(key, value uint64) {
	i := hash(key)
	count := 0

	for count != len(m.entries) {
		i &= (uint64(len(m.entries)) - 1)
		entry := &m.entries[i]
		cmpSwap := entry.key.CompareAndSwap(0, key)
		if cmpSwap || !cmpSwap && entry.key.Load() == key {
			entry.key.Store(key)
			entry.value.Store(value)
			return
		}
		count++
		i++
	}

	for i := range m.entries {
		if m.entries[i].key.Load() == 0 {
			fmt.Printf("%v, %v\n", i, m.entries[i].value.Load())
		}
	}
	panic(fmt.Sprintf("map full: i=%v", i))
}
