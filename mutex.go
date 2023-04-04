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
	m sync.Mutex
	// atomic. indicate the lock waiter
	waiter int32
	// atomic. indicate the lock state
	state int32
}

func (m *Mutex) isUnlocked() bool {
	if atomic.CompareAndSwapInt32(&m.state, GRABBED, UNGRABBED) {
		atomic.AddInt32(&m.waiter, -1)
		return false
	}
	return true
}

func (m *Mutex) Lock() {
	atomic.AddInt32(&m.waiter, 1)
	m.m.Lock()
	atomic.StoreInt32(&m.state, GRABBED)
}

func (m *Mutex) Unlock() {
	atomic.AddInt32(&m.waiter, -1)
	// the state must be updated.
	atomic.StoreInt32(&m.state, UNGRABBED)
	m.m.Unlock()
}

func (m *Mutex) TryLock() bool {
	if m.m.TryLock() {
		atomic.AddInt32(&m.waiter, 1)
		atomic.StoreInt32(&m.state, GRABBED)
		return true
	}
	return false
}

// TryUnlock wouldn't promise "unlock", because the goroutines are all randomized.
// But this function is useful when you don't know the lock state.
// Be careful! YOU NEED TO Know what you are doing when you call this. It will break down "Synchronized Before".
func (m *Mutex) TryUnlock() bool {
	if !m.IsLocked() {
		return false
	}
	// no matter the lock is locked or not, try to lock it.
	// if the lock has beed locked, unlock as usual.
	// if not, still unlock, because we have grabbed the lock.
	couldBeLocked := m.m.TryLock()

	// the lock has been locked, then we have to unlock it, so update the unlock state in advance.
	if !couldBeLocked {
		// somebody has unlocked the lock.
		if m.isUnlocked() {
			// so there's nothing to do.
			return false
		}
	}

	m.m.Unlock()
	return !couldBeLocked
}

func (m *Mutex) IsLocked() bool {
	return atomic.LoadInt32(&m.state) == GRABBED
}

func (m *Mutex) IsBusy() bool {
	return atomic.LoadInt32(&m.waiter) > 1
}
