package config

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/wyvernzora/boggart/internal/format"
	boggart "github.com/wyvernzora/boggart/pkg/model"
	"gopkg.in/yaml.v3"
)

type DefaultFormat struct {
	formatter format.Formatter
}

func (df *DefaultFormat) Name() string {
	return df.formatter.Name()
}

func (df *DefaultFormat) Apply(response *boggart.Response) ([]byte, error) {
	return df.formatter.Apply(response)
}

func (df *DefaultFormat) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return errors.New("invalid default format")
	}

	fmt, err := format.ResolveFormat(node.Value, nil)
	if err != nil {
		return err
	}
	log.Debug().Str("format", fmt.Name()).Msg("Loaded default format")
	df.formatter = fmt
	return nil
}
