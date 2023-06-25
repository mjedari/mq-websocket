package contracts

import "context"

type IProvider interface {
	CheckHealth(ctx context.Context) error
	ResetConnection(ctx context.Context) error
}
