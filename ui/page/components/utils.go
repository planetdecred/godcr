package components

import (
	"context"
	"fmt"
	"log"
	"time"
)

var RetryAttempts = 7

// done returns whether the context's Done channel was closed due to
// cancellation or exceeded deadline.
func ContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func RetryFunc(attempts int, sleep time.Duration, errFunc func() error) (att int, err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}
		err = errFunc()
		if err == nil {
			return -1, nil
		}
	}
	return attempts, fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
