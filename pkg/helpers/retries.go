package helpers

import (
	"context"
	"errors"
	"time"
)

const retriesNumber = 3
const retriesInitialPauseDuration = time.Millisecond * 100
const requestTimeout = time.Second * 5

func WithRetries(ctx context.Context, action func(ctx context.Context) error) error {
	if action == nil {
		return errors.New("incorrect action")
	}

	var err error
	for idx := 1; idx <= retriesNumber; idx++ {
		ctx, cancel := context.WithTimeout(ctx, requestTimeout)
		defer cancel()

		err := action(ctx)
		if err == nil {
			return nil
		}

		time.Sleep(time.Duration(idx) * retriesInitialPauseDuration)
	}

	return err
}
