package mutex

import (
	"sync"

	"github.com/v2pro/plz/gls"
)

type Recursive struct {
	state bool
	mu    sync.Mutex
	id    int
}

func (r *Recursive) Lock() {
	currGoID := int(gls.GoID())
	if r.id != currGoID {
		r.mu.Lock()
		r.id = currGoID
		r.state = true
		return
	}
	// same goroutine, recursive lock
	// in pthread_mutex_lock, they increment the counter and return.
	// however, we don't need a counter. so do nothing here.
}

func (r *Recursive) TryLock() bool {
	return r.mu.TryLock()
}

func (r *Recursive) Unlock() {
	// does it really get locked?
	if !r.state {
		return
	}
	r.state = false
	r.id = -1
	r.mu.Unlock()
}
