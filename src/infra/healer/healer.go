package healer

import (
	"context"
	"fmt"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"time"
)

type Healer struct {
	Providers []contracts.IProvider
	Interval  time.Duration
}

func NewHealerService(providers []contracts.IProvider, interval time.Duration) *Healer {
	return &Healer{
		Providers: providers,
		Interval:  interval,
	}
}

func (h *Healer) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		defer func() {
			ticker.Stop()
			fmt.Println("Closing healer...")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, provider := range h.Providers {
					if err := provider.CheckHealth(ctx); err != nil {
						if err := provider.ResetConnection(ctx); err != nil {
							// todo: log the error
						}
					}
				}
			}
		}
	}()
}
