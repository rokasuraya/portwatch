package backoff

import (
	"context"
	"time"
)

// RetryFunc is a function that returns an error if the operation failed.
type RetryFunc func(ctx context.Context) error

// Retry calls fn up to maxAttempts times, waiting according to b between each
// attempt. It returns the last error if all attempts fail, or nil on success.
// The context is checked before each wait so callers can cancel early.
func Retry(ctx context.Context, b *Backoff, maxAttempts int, fn RetryFunc) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		if err = fn(ctx); err == nil {
			return nil
		}
		if i == maxAttempts-1 {
			break
		}
		wait := b.Next()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	return err
}
