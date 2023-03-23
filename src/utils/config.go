package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GuildID uint64 `yaml:"guild_id"`
}

func ParseConfig(configPath string) (*Config, error) {
	var config *Config
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
