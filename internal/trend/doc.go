// Package trend provides a lightweight sliding-window trend tracker for
// portwatch. It records successive open-port counts and exposes a Direction
// (Rising / Stable / Falling) that downstream components — such as the
// summarizer or notifier — can use to annotate reports with directional
// context without requiring a full time-series store.
package trend
