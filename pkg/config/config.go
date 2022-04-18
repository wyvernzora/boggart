package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	HostKeys      HostKeys      `yaml:"HostKeys"`
	AwsConfig     AwsConfig     `yaml:"AwsConfig"`
	DefaultFormat DefaultFormat `yaml:"DefaultFormat"`
	Permissions   []Permission  `yaml:"Permissions"`
	ListenAddress string        `yaml:"ListenAddress"`
}

func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return &config, nil
}
