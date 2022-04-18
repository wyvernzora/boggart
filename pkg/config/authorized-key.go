package config

import (
	"errors"
	"fmt"
	gliderssh "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
	"strings"
)

type AuthorizedKey struct {
	key ssh.PublicKey
}

func (k *AuthorizedKey) Matches(key ssh.PublicKey) bool {
	return gliderssh.KeysEqual(k.key, key)
}

func (k *AuthorizedKey) String() string {
	return fmt.Sprintf("AuthorizedKey{ %s }", k.key)
}

func (k *AuthorizedKey) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return errors.New("invalid authorized key: unexpected node type")
	}

	data := []byte(strings.ReplaceAll(node.Value, "\n", ""))
	key, _, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return fmt.Errorf("invalid authorized key: %w", err)
	}

	k.key = key
	return nil
}
