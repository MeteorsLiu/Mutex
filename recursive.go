package mutex

import (
	"sync"

	"github.com/v2pro/plz/gls"
)

type Recursive struct {
	mu  sync.Mutex
	id  int
	cnt int
}

func (r *Recursive) Lock() {
	currGoID := int(gls.GoID())
	// we can't get the goid
	if currGoID == 0 {
		panic("cannot get the current goid.")
	}
	if r.id != currGoID {
		r.mu.Lock()
		r.id = currGoID
	}
	// same goroutine, recursive lock
	// in pthread_mutex_lock, they increment the counter and return.
	// because we own the lock, so there's no race.
	r.cnt++
	// overflow
	if r.cnt == 0 {
		panic("too many recursive lock.")
	}
}

func (r *Recursive) TryLock() bool {
	currGoID := int(gls.GoID())
	if r.id != currGoID {
		if r.mu.TryLock() {
			r.id = currGoID
			r.cnt++
			return true
		}
		return false
	}
	r.cnt++
	return true
}

func (r *Recursive) Unlock() {
	// does it really get locked? or are this goroutine really hold the lock?
	currGoID := int(gls.GoID())
	if r.id != currGoID || r.cnt == 0 {
		return
	}
	r.cnt--
	if r.cnt == 0 {
		r.id = -1
		r.mu.Unlock()
	}
}
