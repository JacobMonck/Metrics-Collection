package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigImpl struct {
	GuildID uint64 `yaml:"guild_id"`
}

var Config *ConfigImpl

func ParseConfig(configPath string) error {
	Config := &ConfigImpl{}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&Config); err != nil {
		return err
	}

	return nil
}
