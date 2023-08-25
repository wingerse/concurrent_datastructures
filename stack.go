package concurrentdatastructures

import (
	"sync"
	"sync/atomic"
)

type Node[T any] struct {
	next *Node[T]
	val  T
}

type Stack[T any] struct {
	head *Node[T]
	mu   *sync.Mutex
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{nil, &sync.Mutex{}}
}

func (s *Stack[T]) Push(val T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.head = &Node[T]{s.head, val}
}

func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	old := s.head
	if old == nil {
		var ret T
		return ret, false
	}
	s.head = s.head.next

	return old.val, true
}

type LockFreeNode[T any] struct {
	next atomic.Pointer[LockFreeNode[T]]
	val  T
}

type LockFreeStack[T any] struct {
	head atomic.Pointer[LockFreeNode[T]]
}

func NewLockFreeStack[T any]() *LockFreeStack[T] {
	return &LockFreeStack[T]{}
}

func (s *LockFreeStack[T]) Push(val T) {
	newNode := &LockFreeNode[T]{val: val}
	for {
		head := s.head.Load()
		newNode.next.Store(head)

		if s.head.CompareAndSwap(head, newNode) {
			break
		}
	}
}

func (s *LockFreeStack[T]) Pop() (T, bool) {
	for {
		old := s.head.Load()
		if old == nil {
			var ret T
			return ret, false
		}
		if s.head.CompareAndSwap(old, old.next.Load()) {
			return old.val, true
		}
	}
}
