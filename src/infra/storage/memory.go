package storage

import (
	"context"
	"sync"
	"time"
)

type InMemory struct {
	*sync.Map
}

func (i InMemory) Store(ctx context.Context, key, value string, timeToLive time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (i InMemory) Fetch(ctx context.Context, key string) []byte {
	//TODO implement me
	panic("implement me")
}

func (i InMemory) Exists(ctx context.Context, key string) bool {
	//TODO implement me
	panic("implement me")
}

func (i InMemory) Delete(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}
