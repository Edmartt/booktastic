package utils

import (
	"sync"

	"golang.org/x/time/rate"
)

var mu sync.Mutex

var visitors = make(map[string]*rate.Limiter)

func GetLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(5, 20)
		visitors[ip] = limiter
	}
	return limiter
}
