package room

import (
	"sync"
	"time"
)

// ...

type RateLimitedRoom struct {
	*Room
	rate     int
	lastSent time.Time
	messages int
	mux      sync.Mutex
}

func NewRateLimitedRoom(name string, rate int) *RateLimitedRoom {
	return &RateLimitedRoom{
		Room: &Room{
			Name: name,
		},
		rate:     rate,
		lastSent: time.Now(),
		messages: 0,
	}
}

func (rlr *RateLimitedRoom) Broadcast(message []byte) {
	rlr.mux.Lock()
	defer rlr.mux.Unlock()

	now := time.Now()
	if now.Sub(rlr.lastSent) >= time.Second {
		rlr.messages = 0
		rlr.lastSent = now
	}

	if rlr.messages < rlr.rate {
		rlr.Room.Broadcast(message)
		rlr.messages++
	}
}
