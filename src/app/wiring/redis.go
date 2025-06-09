package wiring

import (
	"github.com/mjedari/mq-websocket/infra/storage"
)

func (w *Wire) SetNewRedisInstance() error {
	newInstance, err := storage.NewRedis(w.Configs.Redis)
	if err != nil {
		return err
	}
	w.Redis = newInstance
	return nil
}
