package mutex

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add()
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

	wg.Wait()
	lock.Unlock()
}
