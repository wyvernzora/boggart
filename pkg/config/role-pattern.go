package config

import (
	"errors"
	"fmt"
	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
)

type RolePattern struct {
	glob glob.Glob
}

func (rp *RolePattern) Matches(arn string) bool {
	return rp.glob.Match(arn)
}

func (rp *RolePattern) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return errors.New("invalid role pattern")
	}

	g, err := glob.Compile(node.Value, '/', ':')
	if err != nil {
		return fmt.Errorf("failed to parse role pattern: %w", err)
	}
	rp.glob = g
	return nil
}
