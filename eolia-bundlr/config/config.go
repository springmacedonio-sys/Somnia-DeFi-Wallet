package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ChainID          int64  `yaml:"chain_id"`
	RPCURL           string `yaml:"rpc_url"`
	EntryPoint       string `yaml:"entry_point"`
	Factory          string `yaml:"factory"`
	BundlrAddress    string `yaml:"bundlr_address"`
	BundlrPrivateKey string `yaml:"bundlr_private_key"`
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return &cfg
}
