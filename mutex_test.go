package mutex

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
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
		lock.Unlock()
		t.Log("Unlocked 1")
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
		t.Log("Locked 4")
		lock.Unlock()
	}
}
