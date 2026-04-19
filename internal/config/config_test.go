package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.AWS.Profile != "default" {
		t.Errorf("expected profile 'default', got %q", cfg.AWS.Profile)
	}
	if cfg.AWS.Region != "ap-northeast-2" {
		t.Errorf("expected region 'ap-northeast-2', got %q", cfg.AWS.Region)
	}
	if cfg.Theme != "dark" {
		t.Errorf("expected theme 'dark', got %q", cfg.Theme)
	}
}

func TestLoadMissingFile(t *testing.T) {
	// Load should return defaults when config file doesn't exist
	cfg := Load()
	if cfg == nil {
		t.Fatal("Load() returned nil")
	}
	if cfg.Theme != "dark" {
		t.Errorf("expected default theme 'dark', got %q", cfg.Theme)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temp dir to use as home
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := &Config{
		AWS: AWSConfig{
			Profile: "test-profile",
			Region:  "us-west-2",
		},
		Theme: "blue",
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify file exists
	path := filepath.Join(tmpDir, ".ecs9s", "config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("config file not created at %s", path)
	}

	// Load it back
	loaded := Load()
	if loaded.AWS.Profile != "test-profile" {
		t.Errorf("expected profile 'test-profile', got %q", loaded.AWS.Profile)
	}
	if loaded.AWS.Region != "us-west-2" {
		t.Errorf("expected region 'us-west-2', got %q", loaded.AWS.Region)
	}
	if loaded.Theme != "blue" {
		t.Errorf("expected theme 'blue', got %q", loaded.Theme)
	}
}
