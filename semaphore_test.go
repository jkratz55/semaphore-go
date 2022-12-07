package semaphore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Panics(t, func() {
		New(0)
	})

	var sem *Semaphore
	assert.NotPanics(t, func() {
		sem = New(5)
	})
	assert.Equal(t, 5, cap(sem.sem))
	assert.Equal(t, 0, len(sem.sem))
}

func TestSemaphore_Acquire(t *testing.T) {
	sem := New(5)
	assert.NoError(t, sem.Acquire(context.Background()))
	assert.NoError(t, sem.Acquire(context.Background()))
	assert.NoError(t, sem.Acquire(context.Background()))
	assert.NoError(t, sem.Acquire(context.Background()))
	assert.NoError(t, sem.Acquire(context.Background()))

	assert.Equal(t, 5, len(sem.sem))

	// This will time out
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	assert.ErrorIs(t, sem.Acquire(ctx), context.DeadlineExceeded)
	cancel()
}

func TestSemaphore_TryAcquire(t *testing.T) {
	sem := New(5)
	assert.True(t, sem.TryAcquire())
	assert.True(t, sem.TryAcquire())
	assert.True(t, sem.TryAcquire())
	assert.True(t, sem.TryAcquire())
	assert.True(t, sem.TryAcquire())

	assert.Equal(t, 5, len(sem.sem))

	assert.False(t, sem.TryAcquire()) // This one will fail to acquire semaphore and return false
}

func TestSemaphore_Release(t *testing.T) {
	sem := New(5)
	sem.sem <- struct{}{}
	sem.sem <- struct{}{}
	sem.sem <- struct{}{}
	sem.sem <- struct{}{}
	sem.sem <- struct{}{}

	assert.Equal(t, 5, len(sem.sem))

	sem.Release()
	assert.Equal(t, 4, len(sem.sem))

	sem.Release()
	assert.Equal(t, 3, len(sem.sem))

	sem.Release()
	assert.Equal(t, 2, len(sem.sem))

	sem.Release()
	assert.Equal(t, 1, len(sem.sem))

	sem.Release()
	assert.Equal(t, 0, len(sem.sem))
}
