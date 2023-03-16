package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigImpl struct {
	GuildID uint64 `yaml:"guild_id"`
}

var Config ConfigImpl

func ParseConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return err
	}

	return nil
}
