# A "safe" mutex lock

Safe means that there's no deadlock.

Example
```
var mu sync.Mutex

mu.Lock()

More Code....

// deadlock! I don't know why Go runtime will not panic.
mu.Lock()

```

However, the project can prevent that.

```
var mu Mutex

mu.Lock()

More Code....

// safe to lock again. Actually it doesn't lock.
mu.Lock()

```

But this will cause the chaos of lock state.

So another function will be useful.

`IsLocked()`


```
var mu Mutex

mu.Lock()

More Code....

if !mu.Locked() {
    mu.Lock()
}

```

What about unlock?

To avoid the chaos of lock state, there's a `TryUnlock()`

If you can't ensure the lock is locked or not.

Use `TryUnlock()`

It's safe to call concurrently(Of course!).


# Recursive lock

Recursive lock acts like C++ STL Recursive lock.

It can lock multiple times, but if you need unlock, you need to unlock it multiple times.
Or it will cause deadlock.

```
var rm Recursive
for i := 0; i < 5; i++ {
    rm.Lock()
}

More Code ...

for i := 0; i < 5; i++ {
    rm.Unlock()
}

```