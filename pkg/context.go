package boggart

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wyvernzora/boggart/pkg/config"
	"github.com/wyvernzora/boggart/pkg/model"
	"net"
	"strings"
)

type Context struct {
	RemoteAddr string
	Permission *config.Permission
	Logger     zerolog.Logger
}

func (ctx *Context) Authenticated(permission *config.Permission) *Context {
	ctx.Permission = permission
	ctx.Logger = ctx.Logger.With().
		Bool("authorized", false).
		Str("key", permission.Name).
		Logger()
	return ctx
}

func (ctx *Context) Authorized(permission *config.Permission) *Context {
	ctx.Permission = permission
	ctx.Logger = ctx.Logger.With().
		Bool("authorized", true).
		Str("key", permission.Name).
		Logger()
	return ctx
}

func (ctx *Context) Request(req *model.Request) *Context {
	ctx.Logger = ctx.Logger.With().
		Str("role", req.RoleArn).
		Str("session", req.SessionName).
		Str("format", req.Format).
		Logger()
	return ctx
}

func NewContext(remote net.Addr) *Context {
	ip := strings.Split(remote.String(), ":")[0]

	return &Context{
		RemoteAddr: ip,
		Permission: nil,
		Logger:     log.With().Logger(),
	}
}
