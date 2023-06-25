package utils

import (
	"context"
	"log"
	"time"
)

const RetryTimes = 5
const RetryDelay = time.Second * 2

type Effector func(ctx context.Context) (any, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (any, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r > retries {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
