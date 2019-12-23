package server

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rateLimiter struct {
	mu       sync.Mutex
	limiters map[string]limiter
	limit    rate.Limit
}

type limiter struct {
	lastAccessed time.Time
	internal     *rate.Limiter
}

func newRateLimiter() *rateLimiter {
	var mutex sync.Mutex

	limiter := rateLimiter{
		mu:       mutex,
		limiters: make(map[string]limiter),
	}

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for {
			<-ticker.C
			limiter.mu.Lock()
			for ip, l := range limiter.limiters {
				if time.Now().After(l.lastAccessed.Add(24 * time.Hour)) {
					delete(limiter.limiters, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()

	return &limiter
}

func (r *rateLimiter) add(ip string) *rate.Limiter {
	internal := rate.NewLimiter(r.limit, 1)
	limiter := limiter{
		lastAccessed: time.Now(),
		internal:     internal,
	}

	r.mu.Lock()
	r.limiters[ip] = limiter
	r.mu.Unlock()

	return internal
}

func (r *rateLimiter) get(ip string) *rate.Limiter {
	r.mu.Lock()
	limiter, exists := r.limiters[ip]

	if !exists {
		r.mu.Unlock()
		return r.add(ip)
	}
	limiter.lastAccessed = time.Now()

	r.mu.Unlock()

	return limiter.internal
}
