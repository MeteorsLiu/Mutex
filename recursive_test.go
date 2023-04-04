package mutex

import (
	"sync"
	"testing"
	"time"
)

func TestRecursive(t *testing.T) {
	var wg sync.WaitGroup
	var rm Recursive
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer rm.Unlock()
		rm.Lock()
		t.Log("1 Get the lock!")
		for i := 0; i < 5; i++ {
			if i == 3 {
				t.Logf("1 unlock")
				rm.Unlock()
			}
			rm.Lock()
			t.Logf("1 Get the lock: %d!", i+1)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			rm.Lock()
			t.Log("2 Get the lock!")
			rm.Unlock()
		}
		rm.Lock()
		time.Sleep(time.Second)
		t.Log("2 unlock")
		rm.Unlock()
	}()
	wg.Wait()
}
