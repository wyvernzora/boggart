package format

import (
	"errors"
	"github.com/wyvernzora/boggart/pkg/model"
	"strings"
)

var formatters = map[string]Formatter{
	"json":  JsonFormatter{},
	"shell": ShellFormatter{},
}

type Formatter interface {
	Name() string
	Apply(response *model.Response) ([]byte, error)
}

func ResolveFormat(name string, fallback Formatter) (Formatter, error) {
	if name == "" {
		return fallback, nil
	}
	name = strings.ToLower(name)
	if fmt, ok := formatters[name]; ok {
		return fmt, nil
	}
	return nil, errors.New("unknown formatter " + name)
}
