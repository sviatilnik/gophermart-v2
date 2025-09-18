package accrual

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	mu              sync.RWMutex
	rateLimitActive bool
	retryAfter      time.Duration
	last429Time     time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

func (rl *RateLimiter) HandleResponse(resp *http.Response) bool {
	if resp.StatusCode == http.StatusTooManyRequests {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		if time.Since(rl.last429Time) < 5*time.Second {
			return true
		}

		rl.last429Time = time.Now()
		rl.rateLimitActive = true
		rl.retryAfter = 60 * time.Second

		// Парсим Retry-After заголовок
		retryAfterStr := resp.Header.Get("Retry-After")
		if retryAfterStr != "" {
			if seconds, err := strconv.Atoi(retryAfterStr); err == nil {
				rl.retryAfter = time.Duration(seconds) * time.Second
			} else if retryTime, err := time.Parse(time.RFC1123, retryAfterStr); err == nil {
				rl.retryAfter = time.Until(retryTime)
			}
		}

		return true
	}

	return false
}

func (rl *RateLimiter) ShouldStop() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.rateLimitActive
}

func (rl *RateLimiter) WaitIfNeeded() {
	rl.mu.Lock()
	if rl.rateLimitActive {
		retryAfter := rl.retryAfter
		rl.mu.Unlock()

		time.Sleep(retryAfter)

		rl.mu.Lock()
		rl.rateLimitActive = false
		rl.retryAfter = 0
	}
	rl.mu.Unlock()
}
