package mutex

import (
	"sync"
	"sync/atomic"
)

const (
	UNGRABBED int32 = iota
	GRABBED
)

type Mutex struct {
	m      sync.Mutex
	waiter int32
	grab   int32
}

func (m *Mutex) Lock() {
	atomic.AddInt32(&m.waiter, 1)
	m.m.Lock()
	atomic.CompareAndSwapInt32(&m.grab, UNGRABBED, GRABBED)
}

func (m *Mutex) Unlock() {
	m.m.Unlock()
	atomic.AddInt32(&m.waiter, -1)
	atomic.CompareAndSwapInt32(&m.grab, GRABBED, UNGRABBED)
}

func (m *Mutex) TryLock() bool {
	if m.m.TryLock() {
		atomic.AddInt32(&m.waiter, 1)
		atomic.CompareAndSwapInt32(&m.grab, UNGRABBED, GRABBED)
		return true
	}
	return false
}

func (m *Mutex) TryUnlock() bool {
	if !m.IsLocked() {
		return false
	}
	m.Unlock()
	return true
}

func (m *Mutex) IsLocked() bool {
	return atomic.CompareAndSwapInt32(&m.grab, GRABBED, UNGRABBED)
}

func (m *Mutex) IsBusy() bool {
	return atomic.LoadInt32(&m.waiter) > 1
}
