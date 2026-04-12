// Package backoff implements an exponential back-off algorithm for use when
// retrying transient failures within portwatch.
//
// Basic usage:
//
//	b := backoff.New(100*time.Millisecond, 30*time.Second, 2.0)
//	err := backoff.Retry(ctx, b, 5, func(ctx context.Context) error {
//		return doSomething(ctx)
//	})
//
// The Backoff type is safe for concurrent use. Call Reset to reuse the same
// instance after a successful operation.
package backoff
