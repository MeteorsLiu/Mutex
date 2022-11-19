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
		t.Log(lock.IsLocked())
	}()

	wg.Wait()
}
