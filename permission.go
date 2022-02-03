package boggart

import (
	gliderssh "github.com/gliderlabs/ssh"
	"github.com/gobwas/glob"
	"golang.org/x/crypto/ssh"
)

type Permission struct {
	Name          string
	Roles         []glob.Glob
	AuthorizedKey ssh.PublicKey
}

func NewPermission(name string, key ssh.PublicKey, roles []glob.Glob) *Permission {
	return &Permission{
		Name:          name,
		Roles:         roles,
		AuthorizedKey: key,
	}
}

func (p *Permission) AllowsAssumeRole(key ssh.PublicKey, roleArn string) bool {
	if !p.IsKeyAuthorized(key) {
		return false
	}
	for _, g := range p.Roles {
		if g.Match(roleArn) {
			return true
		}
	}
	return false
}

func (p *Permission) IsKeyAuthorized(key ssh.PublicKey) bool {
	return gliderssh.KeysEqual(p.AuthorizedKey, key)
}
