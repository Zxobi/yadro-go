package sync

type Semaphore struct {
	C chan struct{}
}

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.C <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Release() {
	<-s.C
}

func NewSemaphore(size int) *Semaphore {
	return &Semaphore{C: make(chan struct{}, size)}
}
