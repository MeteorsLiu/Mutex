package mutex

import (
	"testing"
	"sync"
)

func TestMutex(t *testing.T) {
	var lock Mutex
	var wg sync.WaitGroup
	wg.Add(4)
	go func(){
		defer wg.Done()
		lock.Lock()
		t.Log("Locked")
	}()
	go func{
		defer wg.Done()
		t.Log("is unlocked: ", lock.TryUnlock())
	}()
	go func{
		defer wg.Done()
		t.Log(lock.IsLocked())
	}()

	go func{
		defer wg.Done()
		t.Log("is unlocked: ", lock.TryUnlock())
	}()
	wg.Wait()
}