package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AWS   AWSConfig `yaml:"aws"`
	Theme string    `yaml:"theme"`
}

type AWSConfig struct {
	Profile string `yaml:"profile"`
	Region  string `yaml:"region"`
}

func DefaultConfig() *Config {
	return &Config{
		AWS: AWSConfig{
			Profile: "default",
			Region:  "ap-northeast-2",
		},
		Theme: "dark",
	}
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ecs9s", "config.yaml")
}

func Load() *Config {
	cfg := DefaultConfig()
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	_ = yaml.Unmarshal(data, cfg)
	return cfg
}

func (c *Config) Save() error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
