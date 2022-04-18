package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ResolveFile(node *yaml.Node) ([]byte, error) {
	if node.Tag != "!file" {
		return []byte(node.Value), nil
	}

	// Host key file specified by reference instead of direct inclusion
	if node.Kind != yaml.ScalarNode {
		return nil, errors.New("!file on a non-scalar node")
	}
	data, err := ioutil.ReadFile(node.Value)
	if err != nil {
		return nil, err
	}
	return data, nil
}
