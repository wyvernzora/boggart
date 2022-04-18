package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	boggart "github.com/wyvernzora/boggart/pkg"
	"github.com/wyvernzora/boggart/pkg/config"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.Load("/etc/boggart/config.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}

	server := &Server{
		Boggart:    boggart.New(cfg),
		hostKeys:   cfg.HostKeys,
		format:     &cfg.DefaultFormat,
		listenAddr: ":2222",
	}
	if err := server.serve(); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
