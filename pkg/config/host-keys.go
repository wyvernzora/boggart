package config

import (
	"errors"
	"fmt"
	gliderssh "github.com/gliderlabs/ssh"
	"github.com/rs/zerolog/log"
	"github.com/wyvernzora/boggart/internal/config"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type HostKeys struct {
	Signers []gliderssh.Signer
}

func (k *HostKeys) UnmarshalYAML(node *yaml.Node) error {
	var values []*yaml.Node
	switch node.Kind {
	case yaml.ScalarNode:
		values = []*yaml.Node{node}
		break
	case yaml.SequenceNode:
		values = node.Content
		break
	default:
		return errors.New("host keys must be a string or array of strings")
	}

	k.Signers = make([]gliderssh.Signer, 0)
	for _, key := range values {
		signer, err := resolveHostKey(key)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse host key")
			continue
		}
		k.Signers = append(k.Signers, signer)
	}

	return nil
}

func (k *HostKeys) String() string {
	return fmt.Sprintf("HostKeys{ ... }")
}

func resolveHostKey(node *yaml.Node) (ssh.Signer, error) {
	data, err := config.ResolveFile(node)
	log.Debug().Msg("loaded host key")
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host key: %w", err)
	}
	return signer, nil
}
