package nc_init

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/gudn/iinit"
	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/internal/logging"
	"github.com/gudn/lesslog/internal/nc"
)

func InitNats() {
	if C.Nats == "" {
		return
	}
	conn, err := nats.Connect(C.Nats)
	if err != nil {
		log.Error().Err(err).Msg("failed connect to nats")
	} else {
		log.Info().Msg("success connect to nats")
		nc.Conn = conn
	}
}

func init() {
	iinit.Sequential(
		logging.Init,
		iinit.Static(InitNats),
	)
}
