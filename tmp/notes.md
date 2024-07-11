# Notes about Mutexes in Go

## sync.Mutex

`sync.Mutex` is a mutual exclusion lock. It's a basic locking mechanism used to protect shared resources from concurrent access by multiple goroutines. Only one goroutine can hold the lock at a time. Here's a simple breakdown:

- **m.Lock()**: When a goroutine calls `m.Lock()`, it attempts to acquire the lock. If no other goroutine holds the lock, it succeeds immediately. If another goroutine already holds the lock, the calling goroutine blocks (waits) until the lock becomes available.
- **m.Unlock()**: When a goroutine calls `m.Unlock()`, it releases the lock, allowing other waiting goroutines to acquire it. If no goroutine is waiting for the lock, it simply marks the lock as available.

## sync.RWMutex

`sync.RWMutex` stands for Read-Write Mutex. It's a more advanced locking mechanism that allows multiple readers or a single writer:

- **rw.RLock()**: When a goroutine calls `rw.RLock()`, it attempts to acquire a read lock. Multiple goroutines can hold read locks simultaneously, as long as no goroutine holds the write lock. If a write lock is already held, the calling goroutine blocks until the write lock is released.
- **rw.RUnlock()**: When a goroutine calls `rw.RUnlock()`, it releases a read lock. If it was the last read lock, waiting writers are allowed to proceed.
- **rw.Lock()**: When a goroutine calls `rw.Lock()`, it attempts to acquire a write lock. If no goroutine holds a read or write lock, it succeeds immediately. If other goroutines hold read locks or another goroutine holds the write lock, the calling goroutine blocks until all read locks and the write lock are released.
- **rw.Unlock()**: When a goroutine calls `rw.Unlock()`, it releases the write lock, allowing other waiting readers or writers to proceed.
