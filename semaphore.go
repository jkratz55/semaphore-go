package semaphore

import (
	"context"
)

// Semaphore provides functionality to limit and control concurrent access to a
// resource. The concurrency limit is controlled by the limited provided when
// a Semaphore is created using the New function.
//
// Setting the limit to one will result in the Semaphore behaving like a mutex,
// but if you need the behavior of a mutex you should use sync.Mutex instead.
//
// The zero-value is not usable and a Semaphore should be created using the New
// function.
type Semaphore struct {
	sem chan struct{}
}

// New creates a new Semaphore with the given concurrency limit.
//
// If a limit value of < 1 is provided this function will panic.
func New(limit int) *Semaphore {
	if limit < 1 {
		panic("cannot create a semaphore with a concurrency limit < 1")
	}
	return &Semaphore{
		sem: make(chan struct{}, limit),
	}
}

// Acquire acquires the semaphore blocking until the semaphore/resource is available
// or the context is done. On success a nil error value is returned. If the context
// is done before the semaphore is acquired ctx.Err() is returned and the semaphore
// is left unchanged.
//
// If ctx is already done Acquire might still succeed without blocking.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.sem <- struct{}{}:
		return nil
	}
}

// TryAcquire acquires the semaphore without blocking. On success returns true. On
// failure returns false and leaves the semaphore unchanged.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release releases a single acquirement of the Semaphore.
//
// Note: Improper use of this API can lead to a deadlock! Release should always be
// paired with either Acquire or TryAcquire, and only if they successfully acquired
// the Semaphore.
func (s *Semaphore) Release() {
	<-s.sem
}
