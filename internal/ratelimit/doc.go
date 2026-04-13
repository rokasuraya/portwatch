// Package ratelimit provides a token-bucket rate limiter for controlling
// the frequency of outbound actions such as webhook dispatches or alert
// notifications.
//
// A Limiter is created with New, specifying the refill interval and optional
// burst capacity. Each call to Allow consumes one token; when the bucket is
// empty the call returns false immediately without blocking.
//
// The implementation is safe for concurrent use.
package ratelimit
