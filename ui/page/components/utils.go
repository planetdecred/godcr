package components

import (
	"context"
	"fmt"
	"log"
	"time"
)

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

// RetryFunc implement retry policy for processes that needs to be executed
// in the background after initial failure.
func RetryFunc(retryAttempts int, sleepDur time.Duration, errFunc func() error) (int, error) {
	var err error
	for i := 0; i < retryAttempts; i++ {
		if i > 0 {
			sleepDur *= 2
			log.Printf("waiting %v to retry after error: %v", sleepDur, err)
			time.Sleep(sleepDur)
		}
		err = errFunc()
		if err == nil {
			return i, nil
		}
	}

	return retryAttempts, fmt.Errorf("after %d attempts, last error: %s", retryAttempts, err)
}
