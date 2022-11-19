package mutex

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
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
	lock.Unlock()
}

func TestMuteX(t *testing.T) {
	var lock Mutex
	var wg sync.WaitGroup
	wg.Add(5)
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
		if lock.IsLocked() {
			t.Log("Unlock 1")
			lock.Unlock()
		}
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
