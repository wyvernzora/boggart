package boggart

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gobwas/glob"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	Region      *string
	Credentials *credentials.Credentials
	HostKey     *string
	Permissions []*Permission
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var rawConfig rawConfig
	err = yaml.Unmarshal(data, &rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	config := &Config{
		Permissions: make([]*Permission, 0),
	}

	if rawConfig.Region != nil {
		config.Region = rawConfig.Region
	}

	if rawConfig.Credentials != nil {
		config.Credentials = credentials.NewStaticCredentials(
			rawConfig.Credentials.AccessKeyId,
			rawConfig.Credentials.SecretAccessKey,
			"")
	}

	if rawConfig.HostKey != nil {
		config.HostKey = rawConfig.HostKey
	}

	for k, v := range rawConfig.Permissions {
		authorizedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.ReplaceAll(v.AuthorizedKey, "\n", "")))
		if err != nil {
			return nil, fmt.Errorf("failed to parse authorized key %s: %w", k, err)
		}

		roles := make([]glob.Glob, len(v.Roles))
		for i, role := range v.Roles {
			roles[i] = glob.MustCompile(role, '/', ':')
		}
		config.Permissions = append(config.Permissions, NewPermission(k, authorizedKey, roles))
	}

	return config, nil
}

type rawConfig struct {
	Region      *string                  `yaml:"Region"`
	Credentials *rawCredentials          `yaml:"Credentials"`
	HostKey     *string                  `yaml:"HostKey"`
	Permissions map[string]rawPermission `yaml:"Permissions"`
}

type rawCredentials struct {
	AccessKeyId     string `yaml:"AccessKeyId"`
	SecretAccessKey string `yaml:"SecretAccessKey"`
}

type rawPermission struct {
	AuthorizedKey string   `yaml:"AuthorizedKey"`
	Roles         []string `yaml:"Roles"`
}
