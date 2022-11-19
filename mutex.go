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

	// if a goroutine is unlocking, the CAS may fail, however the lock state must be updated
	for !atomic.CompareAndSwapInt32(&m.grab, UNGRABBED, GRABBED) {
	}
	atomic.SwapInt32(&m.state, LOCKED)
}

func (m *Mutex) Unlock() {
	atomic.AddInt32(&m.waiter, -1)
	// only one can unlock
	if atomic.CompareAndSwapInt32(&m.grab, GRABBED, UNGRABBED) {
		m.m.Unlock()
		atomic.SwapInt32(&m.state, UNLOCKED)
	}

}

func (m *Mutex) TryLock() bool {
	if m.m.TryLock() {
		atomic.AddInt32(&m.waiter, 1)
		for !atomic.CompareAndSwapInt32(&m.grab, UNGRABBED, GRABBED) {
		}
		atomic.SwapInt32(&m.state, LOCKED)
		return true
	}
	return false
}

func (m *Mutex) TryUnlock() bool {
	if !m.IsLocked() {
		return false
	}
	m.Unlock()
	return !m.IsLocked()
}

func (m *Mutex) IsLocked() bool {
	return atomic.LoadInt32(&m.state) == LOCKED && atomic.LoadInt32(&m.grab) == GRABBED
}

func (m *Mutex) IsBusy() bool {
	return atomic.LoadInt32(&m.waiter) > 1
}
