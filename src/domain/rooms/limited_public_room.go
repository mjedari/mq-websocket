package rooms

import (
	"sync"
	"time"
)

type RateLimitedPublicRoom struct {
	*PublicRoom
	rate     int
	lastSent time.Time
	messages int
	mux      sync.Mutex
}

func NewRateLimitedPublicRoom(name string, ratePerMinute int) (*RateLimitedPublicRoom, error) {
	// you can set some rules to prevent get new rooms by returning err
	return &RateLimitedPublicRoom{
		PublicRoom: &PublicRoom{
			BaseRoom: &BaseRoom{
				Name: name,
			},
		},
		rate:     ratePerMinute,
		lastSent: time.Now(),
		messages: 0,
	}, nil
}

func (rlr *RateLimitedPublicRoom) Broadcast(message []byte) {
	rlr.mux.Lock()
	defer rlr.mux.Unlock()

	now := time.Now()
	if now.Sub(rlr.lastSent) >= time.Second {
		rlr.messages = 0
		rlr.lastSent = now
	}

	if rlr.messages < rlr.rate {
		rlr.PublicRoom.Broadcast(message)
		rlr.messages++
	}
}
