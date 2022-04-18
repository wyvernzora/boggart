package config

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type AwsConfig struct {
	Region      string
	Credentials *credentials.Credentials
}

func (c *AwsConfig) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("invalid credentials format")
	}

	var data struct {
		Region          string `yaml:"Region"`
		AccessKeyId     string `yaml:"AccessKeyId"`
		SecretAccessKey string `yaml:"SecretAccessKey"`
	}
	err := node.Decode(&data)
	if err != nil {
		return fmt.Errorf("invalid credentials format: %w", err)
	}

	log.Debug().Str("region", data.Region).Msg("loaded AWS region")
	c.Region = data.Region

	if data.AccessKeyId == "" || data.SecretAccessKey == "" {
		log.Debug().Msg("using AWS credentials from environment")
		c.Credentials = credentials.NewEnvCredentials()
	} else {
		log.Debug().Msg("loaded AWS credentials from config file")
		c.Credentials = credentials.NewStaticCredentials(data.AccessKeyId, data.SecretAccessKey, "")
	}

	return nil
}
