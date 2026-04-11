// Package sampler wraps a port Scanner and drives periodic sampling with
// configurable interval and jitter.
//
// # Basic usage
//
//	sm := sampler.New(scanner, 30*time.Second, 5*time.Second, func(ports []string) {
//		// handle newly sampled ports
//	})
//	if err := sm.Run(ctx, 1, 65535); err != nil && err != context.Canceled {
//		log.Fatal(err)
//	}
//
// The jitter parameter adds a random duration in [0, jitter) to each wait
// period, which helps spread load when multiple instances start simultaneously.
package sampler
