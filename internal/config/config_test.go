package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ScanInterval != 30 {
		t.Errorf("expected scan interval 30, got %d", cfg.ScanInterval)
	}

	if len(cfg.Ports) != 3 {
		t.Errorf("expected 3 default ports, got %d", len(cfg.Ports))
	}

	if !cfg.AlertOnNew {
		t.Error("expected AlertOnNew to be true")
	}

	if !cfg.AlertOnClosed {
		t.Error("expected AlertOnClosed to be true")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid scan interval",
			config: &Config{
				ScanInterval: 0,
				Ports:        []int{80},
			},
			wantErr: true,
		},
		{
			name: "invalid port number",
			config: &Config{
				ScanInterval: 10,
				Ports:        []int{70000},
			},
			wantErr: true,
		},
		{
			name: "invalid port range",
			config: &Config{
				ScanInterval: 10,
				PortRanges:   []PortRange{{Start: 8080, End: 8000}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	original := DefaultConfig()
	original.ScanInterval = 60
	original.Ports = []int{8080, 9090}

	if err := original.Save(configPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.ScanInterval != original.ScanInterval {
		t.Errorf("scan interval mismatch: got %d, want %d", loaded.ScanInterval, original.ScanInterval)
	}

	if len(loaded.Ports) != len(original.Ports) {
		t.Errorf("ports length mismatch: got %d, want %d", len(loaded.Ports), len(original.Ports))
	}
}

func TestGetScanDuration(t *testing.T) {
	cfg := &Config{ScanInterval: 45}
	expected := 45 * time.Second

	if cfg.GetScanDuration() != expected {
		t.Errorf("expected duration %v, got %v", expected, cfg.GetScanDuration())
	}
}
