// Package rotation provides log/state file rotation based on size or age.
package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Options configures rotation behaviour.
type Options struct {
	// MaxBytes rotates the file when it exceeds this size. Zero disables.
	MaxBytes int64
	// MaxAge rotates the file when it is older than this duration. Zero disables.
	MaxAge time.Duration
	// MaxBackups is the number of rotated files to keep. Zero keeps all.
	MaxBackups int
}

// Rotator manages rotation of a single file path.
type Rotator struct {
	path string
	opts Options
	now  func() time.Time
}

// New returns a Rotator for path with the given options.
func New(path string, opts Options) *Rotator {
	return &Rotator{path: path, opts: opts, now: time.Now}
}

// ShouldRotate reports whether the file at r.path needs rotation.
func (r *Rotator) ShouldRotate() (bool, error) {
	info, err := os.Stat(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if r.opts.MaxBytes > 0 && info.Size() >= r.opts.MaxBytes {
		return true, nil
	}
	if r.opts.MaxAge > 0 && r.now().Sub(info.ModTime()) >= r.opts.MaxAge {
		return true, nil
	}
	return false, nil
}

// Rotate renames the current file to a timestamped backup and prunes old backups.
func (r *Rotator) Rotate() error {
	timestamp := r.now().UTC().Format("20060102T150405Z")
	ext := filepath.Ext(r.path)
	base := r.path[:len(r.path)-len(ext)]
	dest := fmt.Sprintf("%s.%s%s", base, timestamp, ext)

	if err := os.Rename(r.path, dest); err != nil {
		return fmt.Errorf("rotation: rename %s -> %s: %w", r.path, dest, err)
	}

	if r.opts.MaxBackups > 0 {
		if err := r.pruneBackups(base, ext); err != nil {
			return err
		}
	}
	return nil
}

// pruneBackups removes the oldest backup files beyond MaxBackups.
func (r *Rotator) pruneBackups(base, ext string) error {
	pattern := base + ".*" + ext
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("rotation: glob %s: %w", pattern, err)
	}
	// matches are lexicographically ordered; oldest first.
	for len(matches) > r.opts.MaxBackups {
		if err := os.Remove(matches[0]); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rotation: remove %s: %w", matches[0], err)
		}
		matches = matches[1:]
	}
	return nil
}
