package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Token               string `yaml:"token"`
	MetricsPort         string `yaml:"metrics_port"`
	ListCartWorkerCount int    `yaml:"list_cart_worker_count"`
	ProductsRpsLimit    uint64 `yaml:"products_rps_limit"`
	ProductsBurstLimit  uint64 `yaml:"products_burst_limit"`
	ProductsCacheTTL    uint64 `yaml:"products_cache_ttl"`
	ProductsCacheSize   int    `yaml:"products_cache_size"`
	Services            struct {
		Loms     string `yaml:"loms"`
		Products string `yaml:"products"`
	} `yaml:"services"`
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
