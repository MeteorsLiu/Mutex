package mutex

import (
	"sync"
	"sync/atomic"
)

const (
	UNLOCKED int32 = iota
	LOCKED
)

const (
	UNGRABBED int32 = iota
	GRABBED
)

type Mutex struct {
	m      sync.Mutex
	state  int32
	waiter int32
	grab   int32
}

func (m *Mutex) Lock() {
	atomic.AddInt32(&m.waiter, 1)
	m.m.Lock()
	atomic.CompareAndSwapInt32(&m.state, UNLOCKED, LOCKED)
}

func (m *Mutex) Unlock() {
	atomic.AddInt32(&m.waiter, -1)
	m.m.Unlock()
	atomic.CompareAndSwapInt32(&m.state, LOCKED, UNLOCKED)
}

func (m *Mutex) TryLock() bool {
	if m.m.TryLock() {
		atomic.AddInt32(&m.waiter, 1)
		atomic.CompareAndSwapInt32(&m.state, UNLOCKED, LOCKED)
		return true
	}
	return false
}

func (m *Mutex) TryUnlock() bool {
	if !m.IsLocked() {
		return false
	}
	if atomic.CompareAndSwapInt32(&m.grab, UNGRABBED, GRABBED) {
		m.Unlock()
		atomic.CompareAndSwapInt32(&m.grab, GRABBED, UNGRABBED)
		return true
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	return atomic.LoadInt32(&m.state) == LOCKED
}

func (m *Mutex) IsBusy() bool {
	return atomic.LoadInt32(&m.waiter) > 1
}
