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
	// if a goroutine is unlocking, the CAS may fail, however the lock state must be updated
	for !atomic.CompareAndSwapInt32(&m.grab, UNGRABBED, GRABBED) {
	}
}

func (m *Mutex) Unlock() {
	atomic.AddInt32(&m.waiter, -1)
	m.m.Unlock()
	atomic.CompareAndSwapInt32(&m.grab, GRABBED, UNGRABBED)
}

func (m *Mutex) TryLock() bool {
	if m.m.TryLock() {
		atomic.AddInt32(&m.waiter, 1)
		for !atomic.CompareAndSwapInt32(&m.grab, UNGRABBED, GRABBED) {
		}
		return true
	}
	return false
}

// TryUnlock wouldn't promise "unlock", because the goroutines are all randomized.
// But this function is useful when you don't know the lock state.
func (m *Mutex) TryUnlock() bool {
	if !m.IsLocked() {
		return false
	}
	if atomic.CompareAndSwapInt32(&m.grab, GRABBED, UNGRABBED) {
		atomic.AddInt32(&m.waiter, -1)
		m.m.Unlock()
		return true
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	return atomic.LoadInt32(&m.grab) == GRABBED
}

func (m *Mutex) IsBusy() bool {
	return atomic.LoadInt32(&m.waiter) > 1
}
