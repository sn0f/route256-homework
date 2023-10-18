package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Group   string   `yaml:"group"`
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return fmt.Errorf("reading config file: %v", err)
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return fmt.Errorf("parsing yaml: %v", err)
	}

	return nil
}
