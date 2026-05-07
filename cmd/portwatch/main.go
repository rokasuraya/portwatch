// Command portwatch is a lightweight CLI daemon that monitors open ports
// and alerts on unexpected changes.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portwatch"
	"github.com/user/portwatch/internal/sighandler"
)

const version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		cfgPath  = flag.String("config", "", "path to config file (default: built-in defaults)")
		printVer = flag.Bool("version", false, "print version and exit")
		printCfg = flag.Bool("print-config", false, "print effective configuration and exit")
	)
	flag.Parse()

	if *printVer {
		fmt.Printf("portwatch %s\n", version)
		return nil
	}

	// Load configuration.
	var cfg *config.Config
	var err error
	if *cfgPath != "" {
		cfg, err = config.Load(*cfgPath)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
	} else {
		cfg = config.DefaultConfig()
	}

	if *printCfg {
		if err := cfg.Save(os.Stdout); err != nil {
			return fmt.Errorf("printing config: %w", err)
		}
		return nil
	}

	// Build the watcher.
	w, err := portwatch.New(cfg)
	if err != nil {
		return fmt.Errorf("initialising watcher: %w", err)
	}

	// Set up signal-based graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sh := sighandler.New()
	go func() {
		sig, _ := sh.Wait(ctx)
		if sig != nil {
			fmt.Fprintf(os.Stderr, "portwatch: received %s, shutting down\n", sig)
		}
		cancel()
	}()

	fmt.Fprintf(os.Stderr, "portwatch %s started\n", version)
	return w.Run(ctx)
}
