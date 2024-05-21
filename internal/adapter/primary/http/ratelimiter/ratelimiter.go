package ratelimiter

import (
	"golang.org/x/time/rate"
	"sync"
)

type RateLimiter struct {
	rps     int
	idLimit map[string]*rate.Limiter
	mu      *sync.RWMutex
}

func NewRateLimiter(rps int) *RateLimiter {
	return &RateLimiter{
		rps:     rps,
		idLimit: make(map[string]*rate.Limiter),
		mu:      &sync.RWMutex{},
	}
}

func (rl *RateLimiter) Take(id string) bool {
	return rl.mustLimiter(id).Allow()
}

func (rl *RateLimiter) mustLimiter(id string) *rate.Limiter {
	if limiter := rl.limiter(id); limiter != nil {
		return limiter
	}
	return rl.addLimiter(id)
}

func (rl *RateLimiter) limiter(id string) *rate.Limiter {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.idLimit[id]
}

func (rl *RateLimiter) addLimiter(id string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter := rl.idLimit[id]
	if limiter != nil {
		return limiter
	}

	limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.rps)
	rl.idLimit[id] = limiter
	return limiter
}
