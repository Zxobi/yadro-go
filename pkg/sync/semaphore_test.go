package sync

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSemaphore(t *testing.T) {
	t.Parallel()

	t.Run("TryAcquireSize1", func(t *testing.T) {
		t.Parallel()
		doTestTryAcquire(t, 1)
	})
	t.Run("TryAcquireSize3", func(t *testing.T) {
		t.Parallel()
		doTestTryAcquire(t, 3)
	})
	t.Run("TryAcquireSize10", func(t *testing.T) {
		t.Parallel()
		doTestTryAcquire(t, 10)
	})

	t.Run("ParallelAcquireSize1", func(t *testing.T) {
		t.Parallel()
		doTestParallelAcquire(1)
	})
	t.Run("ParallelAcquireSize3", func(t *testing.T) {
		t.Parallel()
		doTestParallelAcquire(3)
	})
	t.Run("ParallelAcquireSize10", func(t *testing.T) {
		t.Parallel()
		doTestParallelAcquire(10)
	})
}

func doTestTryAcquire(t *testing.T, size int) {
	s := NewSemaphore(size)
	for i := 0; i < size; i++ {
		assert.True(t, s.TryAcquire(), "TryAcquire failed when not full")
	}
	assert.False(t, s.TryAcquire(), "TryAcquire succeeded when full")

	s.Release()
	assert.True(t, s.TryAcquire(), "TryAcquire failed when not full")
	assert.False(t, s.TryAcquire(), "TryAcquire succeeded when full")

	for i := 0; i < size; i++ {
		s.Release()
	}
	for i := 0; i < size; i++ {
		assert.True(t, s.TryAcquire(), "TryAcquire failed when not full")
	}
	assert.False(t, s.TryAcquire(), "TryAcquire succeeded when full")
}

func doTestParallelAcquire(size int) {
	s := NewSemaphore(size)
	acquire := make(chan bool)
	unlock := make(chan bool)
	done := make(chan bool)

	workers := size * 2
	for i := 0; i < workers; i++ {
		go parallelAcquier(s, acquire, unlock, done)
	}

	for i := 0; i < size; i++ {
		<-acquire
	}
	for i := 0; i < size; i++ {
		unlock <- true
	}
	for i := 0; i < size; i++ {
		<-acquire
	}
	for i := 0; i < size; i++ {
		unlock <- true
	}
	for i := 0; i < workers; i++ {
		<-done
	}
}

func parallelAcquier(s *Semaphore, acquire, unlock, done chan bool) {
	s.Acquire()
	acquire <- true
	<-unlock
	s.Release()
	done <- true
}
