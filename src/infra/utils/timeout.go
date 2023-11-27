package utils

import (
	"context"
	"fmt"
)

type WithContext func(context.Context) error

type SlowFunction func() error

func Timeout(f SlowFunction) WithContext {
	return func(ctx context.Context) error {
		cherr := make(chan error)
		go func() {
			err := f()
			cherr <- err
		}()
		select {
		case <-ctx.Done():
			fmt.Println("function timeout")
			return ctx.Err()

		case err := <-cherr:
			fmt.Println("function processed")
			return err
		}
	}
}
