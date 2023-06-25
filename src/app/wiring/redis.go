package wiring

import (
	"repo.abanicon.com/abantheter-microservices/websocket/infra/storage"
)

func (w *Wire) SetNewRedisInstance() error {
	newInstance, err := storage.NewRedis(w.Configs.Redis)
	if err != nil {
		return err
	}
	w.Redis = newInstance
	return nil
}
