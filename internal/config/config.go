package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the application configuration
type Config struct {
	// ScanInterval is how often to scan ports (in seconds)
	ScanInterval int `json:"scan_interval"`
	// Ports is the list of ports to monitor
	Ports []int `json:"ports"`
	// PortRanges defines ranges of ports to monitor (e.g., "8000-8100")
	PortRanges []PortRange `json:"port_ranges"`
	// AlertOnNew alerts when new ports are opened
	AlertOnNew bool `json:"alert_on_new"`
	// AlertOnClosed alerts when monitored ports are closed
	AlertOnClosed bool `json:"alert_on_closed"`
	// LogFile is the path to the log file
	LogFile string `json:"log_file"`
}

// PortRange represents a range of ports to monitor
type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		ScanInterval:  30,
		Ports:         []int{22, 80, 443},
		PortRanges:    []PortRange{},
		AlertOnNew:    true,
		AlertOnClosed: true,
		LogFile:       "/var/log/portwatch.log",
	}
}

// Load reads configuration from a JSON file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Save writes configuration to a JSON file
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ScanInterval < 1 {
		return fmt.Errorf("scan_interval must be at least 1 second")
	}

	for _, port := range c.Ports {
		if port < 1 || port > 65535 {
			return fmt.Errorf("invalid port number: %d", port)
		}
	}

	for _, pr := range c.PortRanges {
		if pr.Start < 1 || pr.Start > 65535 || pr.End < 1 || pr.End > 65535 {
			return fmt.Errorf("invalid port range: %d-%d", pr.Start, pr.End)
		}
		if pr.Start > pr.End {
			return fmt.Errorf("invalid port range: start (%d) > end (%d)", pr.Start, pr.End)
		}
	}

	return nil
}

// GetScanDuration returns the scan interval as a time.Duration
func (c *Config) GetScanDuration() time.Duration {
	return time.Duration(c.ScanInterval) * time.Second
}
