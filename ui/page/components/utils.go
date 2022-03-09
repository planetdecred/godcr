package components

import (
	"context"
	"fmt"
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

// RetryFunc implements retry policy for processes that needs to be executed
// the after initial failure.
func RetryFunc(retryAttempts int, sleepDur time.Duration, funcDesc string, errFunc func() error) (int, error) {
	var err error
	for i := 0; i < retryAttempts; i++ {
		if i > 0 {
			if i > 1 {
				sleepDur *= 2
			}
			fmt.Printf("waiting %s to retry function %s after error: %v\n", sleepDur, funcDesc, err)
			time.Sleep(sleepDur)
		}
		err = errFunc()
		if err == nil {
			return i, nil
		}
	}

	return retryAttempts, fmt.Errorf("last error: %s", err)
}
