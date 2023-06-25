package rooms

import (
	"sync"
	"time"
)

type RateLimitedRoom struct {
	*BaseRoom
	rate     int
	lastSent time.Time
	messages int
	mux      sync.Mutex
}

func NewRateLimitedRoom(name string, ratePerMinute int) (*RateLimitedRoom, error) {
	// you can set some rules to prevent get new rooms by returning err
	return &RateLimitedRoom{
		BaseRoom: &BaseRoom{
			Name: name,
		},
		rate:     ratePerMinute,
		lastSent: time.Now(),
		messages: 0,
	}, nil
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
		rlr.BaseRoom.Broadcast(message)
		rlr.messages++
	}
}
