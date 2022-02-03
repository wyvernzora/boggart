package boggart

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configfile := "/etc/boggart/config.yml"
	if len(os.Args) > 1 {
		configfile = os.Args[1]
	}

	boggart, err := NewBoggart(configfile)
	if err != nil {
		log.Fatal().Err(err).Msgf("")
	}
	log.Fatal().Err(boggart.Serve(":2222")).Msgf("")
}
