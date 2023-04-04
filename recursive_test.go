package mutex

import (
	"sync"
	"testing"
)

func TestRecursive(t *testing.T) {
	var wg sync.WaitGroup
	var rm Recursive
	incr := 0
	wg.Add(2)
	go func() {
		defer wg.Done()

		t.Log("1 Get the lock!")
		for i := 0; i < 5; i++ {
			rm.Lock()
			incr++
			t.Logf("1 Get the lock: %d!", incr)
		}
		for i := 0; i < 5; i++ {
			rm.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			rm.Lock()
			incr++
			t.Logf("2 Get the lock: %d!", incr)
		}
		for i := 0; i < 5; i++ {
			rm.Unlock()
		}
	}()

	/* Recursive Deadlock Example
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			rm.Lock()
			incr++
			t.Logf("2 Get the lock: %d!", incr)
		}

		for i := 0; i < 3; i++ {
			rm.Unlock()
		}
	}()

	*/
	wg.Wait()
}
