package concurrentdatastructures

import (
	"testing"

	"github.com/petermattis/goid"
)

func TestLockFreeStack(t *testing.T) {
	s := NewLockFreeStack[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if v, ok := s.Pop(); !(ok && v == 3) {
		t.Fatal(v, ok)
	}
	if v, ok := s.Pop(); !(ok && v == 2) {
		t.Fatal(v, ok)
	}
	if v, ok := s.Pop(); !(ok && v == 1) {
		t.Fatal(v, ok)
	}
}

func BenchmarkStack(b *testing.B) {
	s := NewStack[int]()
	for i := 0; i < 10000; i++ {
		s.Push(1)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if goid.Get()%2 == 0 {
				s.Push(1)
			} else {
				s.Pop()
			}
		}
	})
}

func BenchmarkLockFreeStack(b *testing.B) {
	s := NewLockFreeStack[int]()
	for i := 0; i < 10000; i++ {
		s.Push(1)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if goid.Get()%2 == 0 {
				s.Push(1)
			} else {
				s.Pop()
			}
		}
	})
}
