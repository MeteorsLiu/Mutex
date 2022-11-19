package mutex

import (
	"sync"
	"testing"
)

// Not stable
func _TestMutex(t *testing.T) {
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		lock.Lock()
		t.Log("Locked")
	}()
	go func() {
		defer wg.Done()
		lock.Lock()
		t.Log("Locked 2")
	}()
	go func() {
		defer wg.Done()
		lock.Unlock()
		t.Log("Unlocked 1")
	}()
	go func() {
		defer wg.Done()
		lock.Unlock()
		t.Log("Unlocked 1")
	}()

	wg.Wait()
}

func TestMuteX(t *testing.T) {
	var lock Mutex
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		t.Log("Locked 1 is unlocked: ", lock.TryUnlock())
		lock.Lock()
		t.Log("Locked")
	}()
	go func() {
		defer wg.Done()
		t.Log("Locked 2 is unlocked: ", lock.TryUnlock())
		lock.Lock()
		t.Log("Locked 2")
	}()
	go func() {
		defer wg.Done()
		t.Log("is unlocked: ", lock.TryUnlock())
	}()
	go func() {
		defer wg.Done()
		t.Log(lock.IsLocked())
	}()
	go func() {
		defer wg.Done()
		t.Log("is unlocked: ", lock.TryUnlock())
	}()
	wg.Wait()
	if lock.IsLocked() {
		t.Log("Unlock 2")
		lock.Unlock()
	}
}
