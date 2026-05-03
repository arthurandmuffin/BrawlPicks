package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server *Server `yaml:"server"`
	UI     *UI     `yaml:"ui"`
}

type Server struct {
	BaseURL       string `yaml:"baseURL"`
	RecommendPath string `yaml:"recommendPath"`
}

type UI struct {
	DefaultMapName string `yaml:"defaultMapName"`
	DefaultMode    string `yaml:"defaultMode"`
	DefaultRank    int    `yaml:"defaultRank"`
	DefaultTopK    int    `yaml:"defaultTopK"`
}

func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
