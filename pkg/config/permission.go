package config

import "golang.org/x/crypto/ssh"

type Permission struct {
	Name          string        `yaml:"Name"`
	AuthorizedKey AuthorizedKey `yaml:"AuthorizedKey"`
	Roles         []RolePattern `yaml:"Roles"`
}

func (p *Permission) IsKeyAuthorized(key ssh.PublicKey) bool {
	return p.AuthorizedKey.Matches(key)
}

func (p *Permission) CanAssumeRole(arn string) bool {
	for _, rp := range p.Roles {
		if rp.Matches(arn) {
			return true
		}
	}
	return false
}
