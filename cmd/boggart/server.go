package main

import (
	"errors"
	"github.com/gliderlabs/ssh"
	"github.com/rs/zerolog/log"
	"github.com/wyvernzora/boggart/internal/format"
	boggart "github.com/wyvernzora/boggart/pkg"
	"github.com/wyvernzora/boggart/pkg/config"
	"github.com/wyvernzora/boggart/pkg/model"
)

type Server struct {
	*boggart.Boggart
	hostKeys   config.HostKeys
	format     format.Formatter
	listenAddr string
}

func (srv *Server) serve() error {
	sshServer := &ssh.Server{
		Addr: srv.listenAddr,
		Handler: func(s ssh.Session) {
			data, status := srv.handleRequest(s)
			s.Write(data)
			s.Exit(status)
		},
		HostSigners: srv.hostKeys.Signers,
	}

	sshAuth := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return srv.handleAuth(ctx, key)
	})
	if err := sshServer.SetOption(sshAuth); err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msgf("listening on %s", srv.listenAddr)
	return sshServer.ListenAndServe()
}

func (srv *Server) handleAuth(s ssh.Context, key ssh.PublicKey) bool {
	ctx := boggart.NewContext(s.RemoteAddr())
	return srv.Authenticate(ctx, key)
}

func (srv *Server) handleRequest(s ssh.Session) ([]byte, int) {
	ctx := boggart.NewContext(s.RemoteAddr())

	req, err := model.ParseRequest(s.RawCommand())
	if err != nil {
		ctx.Logger.Error().Err(err).Msg("failed to parse request")
		return []byte{}, 1
	}
	ctx.Request(req)

	fmt, err := format.ResolveFormat(req.Format, srv.format)
	if err != nil {
		ctx.Logger.Error().Err(err).Msg("failed to resolve formatter")
		return []byte{}, 1
	}

	if !srv.Authorize(ctx, s.PublicKey(), req) {
		err := errors.New("authorization failed")
		data, _ := fmt.Apply(model.NewErrorResponse(err))
		ctx.Logger.Error().Err(err).Msg("unauthorized")
		return data, 1
	}

	response := srv.AssumeRole(ctx, req)
	data, err := fmt.Apply(response)
	if err != nil {
		ctx.Logger.Error().Err(err).Msg("failed to assume role")
		return []byte{}, 1
	}
	if !response.Success {
		return data, 1
	}
	return data, 0
}
