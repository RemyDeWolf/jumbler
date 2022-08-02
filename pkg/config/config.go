package config

import (
	_ "embed"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var EmbedConfig []byte

type Config struct {
	Path        string `yaml:"path"`
	Ext         string `yaml:"ext"`
	Password    string
	DryRun      bool
	AutoApprove bool
	Quiet       bool
}

func GetDefault() (Config, error) {
	var err error
	var config Config

	err = yaml.Unmarshal(EmbedConfig, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func ReadFile(file string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(file)
	if err != nil {
		return cfg, err
	}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
