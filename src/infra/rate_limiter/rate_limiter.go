package rate_limiter

import (
	"errors"
	"fmt"
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"sync"
	"time"
)

type visitor struct {
	*sync.Mutex
	Requests   int
	LastAccess time.Time
}

type RateLimiter struct {
	Visitors map[string]*visitor
	*sync.Mutex
	Rate   int
	Period time.Duration
}

func NewRateLimiter(conf configs.RateLimiter) *RateLimiter {
	return &RateLimiter{
		Visitors: make(map[string]*visitor),
		Mutex:    &sync.Mutex{},
		Rate:     conf.Rate,
		Period:   time.Duration(conf.Period) * time.Minute,
	}
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.Lock()
	defer rl.Unlock()

	v, exists := rl.Visitors[ip]
	if !exists {
		v = &visitor{
			Mutex:      &sync.Mutex{},
			LastAccess: time.Now(),
			Requests:   0,
		}
		rl.Visitors[ip] = v
	}

	return v
}

func (rl *RateLimiter) Handle(ip string) error {
	v := rl.getVisitor(ip)
	v.Lock()
	defer v.Unlock()

	// If too many requests have been made within the period, reject the request
	if time.Since(v.LastAccess) < rl.Period && v.Requests >= rl.Rate {
		//http.Error(w, "Too many requests", http.StatusTooManyRequests)
		fmt.Println("Got too many requests")
		return errors.New("too many requests")
	}

	// Reset the number of requests if the period has expired
	if time.Since(v.LastAccess) >= rl.Period {
		fmt.Println("reset the requests")
		v.Requests = 0
	}

	v.LastAccess = time.Now()
	v.Requests++

	return nil
}
