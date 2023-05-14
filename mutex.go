package mutex

import (
	"sync"
	"sync/atomic"

	"github.com/v2pro/plz/gls"
)

type Mutex struct {
	m sync.Mutex
	// atomic. indicate the lock waiter
	waiter int32
	id     int
}

func (m *Mutex) Lock() {
	atomic.AddInt32(&m.waiter, 1)
	currGoID := int(gls.GoID())
	if m.id != currGoID {
		m.m.Lock()
		m.id = currGoID
	}
}

func (m *Mutex) Unlock() {
	currGoID := int(gls.GoID())
	if m.id == currGoID {
		m.id = 0
		m.m.Unlock()
		atomic.AddInt32(&m.waiter, -1)
	}
}

func (m *Mutex) TryLock() bool {
	if m.m.TryLock() {
		m.id = int(gls.GoID())
		atomic.AddInt32(&m.waiter, 1)
		return true
	}
	return false
}

func (m *Mutex) TryUnlock() bool {
	currGoID := int(gls.GoID())
	if m.id == currGoID {
		m.id = 0
		m.m.Unlock()
		atomic.AddInt32(&m.waiter, -1)
		return true
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	return m.id != 0
}

func (m *Mutex) IsBusy() bool {
	return atomic.LoadInt32(&m.waiter) > 1
}
